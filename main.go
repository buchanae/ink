package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
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

	abspath, err := filepath.Abs(args[0])
	if err != nil {
		panic(err)
	}

	a, err := app.NewApp(app.DefaultConfig())
	if err != nil {
		panic(err)
	}

	watch, err := newWatcher()
	if err != nil {
		panic(err)
	}
	watch.Watch(abspath)

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
				err := a.RunSketch(ctx, args[0])
				if err != nil {
					log.Printf("error: %v", err)
				}
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
