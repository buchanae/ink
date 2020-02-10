package runner

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/render"
)

type App interface {
	// TODO
	//SetConfig(msg.Config)
	RenderPlan(render.Plan)
}

func RunSketch(ctx context.Context, app App, path string) error {

	wd, err := newWorkdir()
	log.Print(wd.path)
	//defer wd.cleanup()
	if err != nil {
		return err
	}

	//return run(ctx, app, wd, path)
	return runPlugin(ctx, app, wd, path)
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

func run(ctx context.Context, app App, wd workdir, sketchPath string) error {

	err := build(wd, sketchPath)
	if err != nil {
		return err
	}

	binPath := filepath.Join(wd.path, "inkbin")
	cmd := exec.Command(binPath)
	cmd.Dir = filepath.Dir(sketchPath)

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

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("starting: %v", err)
	}

	doc := gfx.NewDoc()
	// TODO
	//doc.Config = app.conf
	enc := gob.NewEncoder(stdin)
	err = enc.Encode(doc)
	if err != nil {
		return fmt.Errorf("sending initial doc: %v", err)
	}

	dec := gob.NewDecoder(stdout)

	for {
		msg := &RenderMessage{}
		err = dec.Decode(msg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("decoding: %v", err)
		}

		// TODO
		//app.SetConfig(msg.Config)
		app.RenderPlan(msg.Plan)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("cmd: %v", err)
	}

	return nil
}

const head = `
package main
import "log"
import "github.com/buchanae/ink/app/runner"

func main() {
	log.SetFlags(0)
	doc := runner.RecvDoc()
	Ink(doc)
	runner.Send(doc)
}
`
