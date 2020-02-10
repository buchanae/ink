package runner

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"

	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/render"
)

func runPlugin(ctx context.Context, app App, wd workdir, sketchPath string) error {

	err := buildPlugin(wd, sketchPath)
	if err != nil {
		return err
	}
	log.Print("load plugin")

	binPath := filepath.Join(wd.path, "inkbin")
	plug, err := plugin.Open(binPath)
	if err != nil {
		return err
	}

	sym, err := plug.Lookup("Ink")
	if err != nil {
		return err
	}

	inkfunc, ok := sym.(func(*gfx.Doc))
	if !ok {
		return fmt.Errorf("type of Ink is not func(*gfx.Doc)")
	}

	log.Print("run plugin")

	// TODO plugin runner won't be able to change directories?
	//cmd := exec.Command(binPath)
	//cmd.Dir = filepath.Dir(sketchPath)

	doc := gfx.NewDoc()
	inkfunc(doc)
	plan := render.BuildPlan(doc)

	// TODO
	//doc.Config = app.conf

	// TODO plugin runner needs animation/sending updates rethink

	// TODO
	//app.SetConfig(msg.Config)
	app.RenderPlan(plan)
	log.Print("after render")

	return nil
}

func buildPlugin(wd workdir, path string) error {

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
		"go", "build", "-tags=sendonly", "-buildmode=plugin", "-o=inkbin", ".",
	)
	cmd.Env = []string{}
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, "GO111MODULE=on")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = wd.path
	return cmd.Run()
}
