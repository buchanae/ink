package main

import (
	"context"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/buchanae/ink/app"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	log.SetFlags(0)
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: ink file.go")
		os.Exit(1)
	}

	path, err := filepath.Abs(args[0])
	if err != nil {
		panic(err)
	}

	a, err := app.NewApp(app.DefaultConfig())
	if err != nil {
		panic(err)
	}

	wd, err := newWorkdir()
	defer wd.cleanup()
	if err != nil {
		panic(err)
	}

	watch, err := newWatcher()
	if err != nil {
		panic(err)
	}
	watch.Watch(path)

	a.AddKeyCallback(func(ev app.KeyEvent) {
		if ev.Pressed(glfw.KeyR) {
			watch.changes <- struct{}{}
		}
	})

	go func() {
		wg := sync.WaitGroup{}

		for {
			ctx, cancel := context.WithCancel(context.Background())

			wg.Add(1)
			go func() {
				run(ctx, a, wd, path, args[0])
				wg.Done()
			}()

			<-watch.changes
			log.Println("change")
			cancel()
			wg.Wait()
		}
	}()

	// Most access to the window must be done on a single OS thread,
	// so this code locks itself to the OS thread and handles all communication
	// via SDL queues and Go channels.
	err = a.Run()
	if err != nil {
		log.Printf("error: %v", err)
	}
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
	modContent := mod

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

func build(wd workdir, sketchPath, sketchName string) error {

	inkPath := filepath.Join(wd.path, "ink.go")
	err := copyFile(inkPath, sketchPath, sketchName)
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

func run(ctx context.Context, a *app.App, wd workdir, sketchPath, sketchName string) error {

	err := build(wd, sketchPath, sketchName)
	if err != nil {
		return err
	}
	log.Print("run")

	binPath := filepath.Join(wd.path, "inkbin")
	cmd := exec.Command(binPath)
	cmd.Dir = filepath.Dir(sketchPath)

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

	dec := gob.NewDecoder(stdout)

	for {
		start := time.Now()
		doc := &app.Doc{}
		err = dec.Decode(doc)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("decoding: %v", err)
		}

		renderStart := time.Now()
		a.Render(doc)
		log.Printf("render iter: %s", time.Since(renderStart))
		log.Printf("iter: %s", time.Since(start))
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("cmd: %v", err)
	}

	return nil
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
	doc := app.NewDoc()
	Ink(doc)
	app.Send(doc)
}
`

const mod = `
module temp
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
