package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/buchanae/ink/app"
	"github.com/buchanae/ink/gfx"
	. "github.com/buchanae/ink/trace"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, "usage: ink file.go")
		os.Exit(1)
	}

	watch := newWatcher()

	path, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}

	a, err := app.NewApp(app.DefaultConfig())
	if err != nil {
		panic(err)
	}

	watch.Watch(path)

	refresh := func() {
		run(a, path)
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

type firstByteReader struct {
	r     io.Reader
	total int
	done  bool
}

func (fbr *firstByteReader) Read(data []byte) (int, error) {
	n, err := fbr.r.Read(data)
	if !fbr.done {
		Trace("first byte")
		fbr.done = true
	}
	fbr.total += n
	return n, err
}

func run(app *app.App, path string) error {
	StartTrace()
	Trace("run")

	tempDir, err := ioutil.TempDir("", "ink-run-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	inkPath := filepath.Join(tempDir, "ink.go")

	Trace("copy")
	err = copyFile(inkPath, path)
	if err != nil {
		return err
	}

	Trace("write main")
	mainPath := filepath.Join(tempDir, "main.go")
	err = ioutil.WriteFile(mainPath, []byte(head), 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "run", inkPath, mainPath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = os.Stderr

	Trace("start")
	err = cmd.Start()
	if err != nil {
		return err
	}

	reader := &firstByteReader{r: stdout}

	for {
		Trace("decode")
		doc := &gfx.Layer{}
		dec := gob.NewDecoder(reader)
		err = dec.Decode(doc)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		Trace("render")
		app.Render(doc)
	}

	Trace("wait")
	defer Trace("done %d bytes", reader.total)
	return cmd.Wait()
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

import "os"
import "encoding/gob"
import "github.com/buchanae/ink/gfx"

func main() {
	layer := gfx.NewLayer()
	Ink(layer)
	err := gob.NewEncoder(os.Stdout).Encode(layer)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Stderr.Write([]byte("\n"))
	}
}
`
