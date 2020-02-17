package app

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/buchanae/ink/app/client"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/trac"
)

func (app *App) RunSketch(ctx context.Context, path string) error {

	wd, err := newWorkdir()
	defer wd.cleanup()
	if err != nil {
		return err
	}

	return run(ctx, app, wd, path)
}

func newDoc(app *App) *client.Doc {

	doc := client.NewDoc()
	c := &gfx.Config{}
	doc.Conf = c
	c.Width = app.conf.Window.Width
	c.Height = app.conf.Window.Height
	c.Snapshot.Width = app.conf.Snapshot.Width
	c.Snapshot.Height = app.conf.Snapshot.Height

	return doc
}

func build(wd workdir, path string) error {

	abspath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	inkPath := filepath.Join(wd.path, "ink.go")

	err = copyFile(inkPath, abspath, path)
	if err != nil {
		return err
	}

	cmd := exec.Command(
		"go", "build", "-o=inkbin", ".",
	)
	cmd.Env = []string{}
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, "GO111MODULE=on")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = wd.path
	return cmd.Run()
}

func run(ctx context.Context, app *App, wd workdir, sketchPath string) error {

	span := trac.Start("go build")
	err := build(wd, sketchPath)
	if err != nil {
		return err
	}
	span.End()

	binPath := filepath.Join(wd.path, "inkbin")
	cmd := exec.Command(binPath)
	cmd.Dir = filepath.Dir(sketchPath)
	cmd.Env = []string{}
	cmd.Env = append(cmd.Env, os.Environ()...)
	if trac.Enabled {
		cmd.Env = append(cmd.Env, "INK_TRACE=true")
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("getting stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("getting stdout pipe: %v", err)
	}
	cmd.Stderr = os.Stderr

	// CommandContext doesn't work as expected when used with subshells and stdout.
	// https://groups.google.com/forum/#!topic/golang-nuts/sbuYR7WpsZg
	// Close the stdout pipe explicitly when context is done.
	go func() {
		<-ctx.Done()
		stdout.Close()
	}()

	trac.Log("exec sketch")
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("starting: %v", err)
	}

	doc := newDoc(app)

	enc := gob.NewEncoder(stdin)
	err = enc.Encode(doc)
	if err != nil {
		return fmt.Errorf("sending initial doc: %v", err)
	}

	rdr := &dbgread{R: stdout}
	dec := gob.NewDecoder(rdr)

	for {
		msg := &client.RenderMessage{}
		err = dec.Decode(msg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("decoding: %v", err)
		}
		trac.Log("received")

		c := app.conf
		c.Window.Title = msg.Config.Title
		c.Window.Width = msg.Config.Width
		c.Window.Height = msg.Config.Height
		c.Snapshot.Width = msg.Config.Snapshot.Width
		c.Snapshot.Height = msg.Config.Snapshot.Height

		app.SetConfig(c)
		app.RenderPlan(msg.Plan)
		trac.Log("next loop")
		rdr.started = false
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("cmd: %v", err)
	}

	return nil
}

type dbgread struct {
	R       io.Reader
	started bool
}

func (d *dbgread) Read(p []byte) (int, error) {
	n, err := d.R.Read(p)
	if !d.started {
		trac.Log("first byte")
		d.started = true
	}
	return n, err
}

const head = `
package main
import "github.com/buchanae/ink/app/client"

func main() {
	client.Main(Ink)
}
`
