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
	"sync"

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

	modPath := filepath.Join(wd.path, "go.mod")
	err = ioutil.WriteFile(modPath, []byte(mod), 0644)
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

	reader := &firstByteReader{r: stdout}

	for {
		doc := &app.Doc{}
		dec := gob.NewDecoder(reader)
		err = dec.Decode(doc)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("decoding: %v", err)
		}

		a.Render(doc)
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
	app.Send(Ink)
}
`

const mod = `
module temp
`

type firstByteReader struct {
	r     io.Reader
	total int
	done  bool
}

func (fbr *firstByteReader) Read(data []byte) (int, error) {
	n, err := fbr.r.Read(data)
	if !fbr.done {
		fbr.done = true
	}
	fbr.total += n
	return n, err
}
