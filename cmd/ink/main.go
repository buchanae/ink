package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/buchanae/ink/app"
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, "usage: ink file.go")
		os.Exit(1)
	}

	gob.Register(gfx.Shader{})
	gob.Register(color.RGBA{})
	gob.Register(dd.XY{})
	gob.Register([]color.RGBA{})
	gob.Register([]dd.XY{})

	watch, err := newWatcher()
	if err != nil {
		log.Print(err)
		return
	}

	path := os.Args[1]
	watch.Watch(path)

	a, err := app.NewApp(app.DefaultConfig())
	if err != nil {
		panic(err)
	}

	refresh := func() {
		doc, err := run(path)
		if err != nil {
			log.Print(err)
			return
		}
		a.Render(doc)
	}

	go func() {
		refresh()
		for range watch.changes {
			refresh()
		}
	}()

	// Most access to the window must be done on a single OS thread,
	// so this code locks itself to the OS thread and handles all communication
	// via SDL queues and Go channels.
	a.Run()
}

func run(path string) (*gfx.Doc, error) {
	cmd := exec.Command("go", "run", "-tags", "sendonly", path)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, "INK_REMOTE=true")

	start := time.Now()
	err := cmd.Run()
	if err != nil {
		if ex, ok := err.(*exec.ExitError); ok {
			log.Println(string(ex.Stderr))
		}
		return nil, err
	}
	log.Printf("took: %s", time.Since(start))

	doc := &gfx.Doc{}
	dec := gob.NewDecoder(&stdout)
	err = dec.Decode(doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
