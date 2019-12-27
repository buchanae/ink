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
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/trace"
)

func init() {
	flag.BoolVar(&trace.On, "trace", trace.On, "Enable trace output")
}

func main() {
	log.SetFlags(0)
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		fmt.Fprint(os.Stderr, "usage: ink file.go")
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

	watch := newWatcher()
	watch.Watch(path)

	go func() {
		wg := sync.WaitGroup{}

		for {
			ctx, cancel := context.WithCancel(context.Background())

			wg.Add(1)
			go func() {
				run(ctx, a, path)
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
	a.Run()
}

type firstByteReader struct {
	r     io.Reader
	total int
	done  bool
}

func (fbr *firstByteReader) Read(data []byte) (int, error) {
	n, err := fbr.r.Read(data)
	if !fbr.done {
		trace.Log("first byte")
		fbr.done = true
	}
	fbr.total += n
	return n, err
}

func run(ctx context.Context, app *app.App, path string) error {
	trace.Start()
	trace.Log("run")

	tempDir, err := ioutil.TempDir("", "ink-run-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	inkPath := filepath.Join(tempDir, "ink.go")

	trace.Log("copy")
	err = copyFile(inkPath, path)
	if err != nil {
		return err
	}

	trace.Log("write main")
	mainPath := filepath.Join(tempDir, "main.go")
	err = ioutil.WriteFile(mainPath, []byte(head), 0644)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "go", "run", inkPath, mainPath)

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

	trace.Log("start")
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("starting: %v", err)
	}

	reader := &firstByteReader{r: stdout}

	for {
		trace.Log("decode")
		doc := &gfx.Layer{}
		dec := gob.NewDecoder(reader)
		err = dec.Decode(doc)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("decoding: %v", err)
		}

		trace.Log("render")
		app.Render(doc)
	}

	trace.Log("wait")
	defer trace.Log("done %d bytes", reader.total)
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("cmd: %v", err)
	}

	return nil
}

func copyFile(dstPath, srcPath string) error {
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}

	_, err = dst.Write([]byte("//line " + srcPath + ":1\n\n"))
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
import "github.com/buchanae/ink/gfx"

func main() {
	gfx.Run(Ink)
}
`
