// +build !sendonly

package app

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/buchanae/ink/color"
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

func (app *App) NewDoc() *Doc {
	doc := &Doc{
		ID: nextID(),
	}
	// TODO this shouldn't be here
	gfx.Clear(doc, color.White)

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
		"go", "build", "-tags=sendonly", "-o=inkbin", ".",
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

	doc := app.NewDoc()

	enc := gob.NewEncoder(stdin)
	err = enc.Encode(doc)
	if err != nil {
		return fmt.Errorf("sending initial doc: %v", err)
	}

	rdr := &dbgread{R: stdout}
	dec := gob.NewDecoder(rdr)

	for {
		msg := &RenderMessage{}
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

func copyFile(dstPath, srcPath, name string) error {
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}

	_, err = dst.Write([]byte("//line " + name + ":1\n"))
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	defer src.Close()
	return nil
}

const head = `
package main
import "github.com/buchanae/ink/app"

func main() {
	app.Main(Ink)
}
`

// search directory tree for a "go.mod" file for the ink module
// in order to decide if ink should build against a local codebase.
func findInkCode() (string, error) {
	if v := os.Getenv("INK_PATH"); v != "" {
		return v, nil
	}

	const inkMod = "module github.com/buchanae/ink"

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting current directory: %w", err)
	}

	current := wd
	for {
		mod := filepath.Join(current, "go.mod")

		info, err := os.Stat(mod)
		if err != nil && !os.IsNotExist(err) {
			return "", fmt.Errorf("getting file info: %w", err)
		}

		if info != nil && info.Mode().IsRegular() {
			b, err := ioutil.ReadFile(mod)
			if err != nil {
				return "", fmt.Errorf("reading file %q: %w", mod, err)
			}
			if strings.HasPrefix(string(b), inkMod) {
				return current, nil
			}
		}

		dir := filepath.Dir(current)
		base := filepath.Base(current)
		if dir == base {
			break
		}
		current = dir
	}
	return "", nil
}

type workdir struct {
	path string
}

func newWorkdir() (wd workdir, err error) {

	wd.path, err = ioutil.TempDir("", "ink-run-")
	if err != nil {
		return
	}

	mainPath := filepath.Join(wd.path, "ink_main_wrapper.go")
	err = ioutil.WriteFile(mainPath, []byte(head), 0644)
	if err != nil {
		return
	}

	// TODO should look for an existing go.mod file in the sketch directory
	modPath := filepath.Join(wd.path, "go.mod")
	modContent := "module temp\n"

	inkCode, err := findInkCode()
	if err != nil {
		log.Printf("error: finding ink code: %v", err)
	}
	if inkCode != "" {
		modContent += "\n\nreplace github.com/buchanae/ink => " + inkCode
	}

	err = ioutil.WriteFile(modPath, []byte(modContent), 0644)
	if err != nil {
		return
	}

	return
}

func (w workdir) cleanup() {
	if w.path != "" {
		os.RemoveAll(w.path)
	}
}
