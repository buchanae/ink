// +build !darwin

package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

// newWatcher returns a new file watcher.
func newWatcher() (*watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("creating FS watcher: %s", err)
	}

	w := &watcher{
		watcher: fsw,
		changes: make(chan struct{}, 100),
	}

	go func() {
		for err := range fsw.Errors {
			log.Println("fswatch error:", err)
		}
	}()

	go func() {
		for event := range fsw.Events {
			if event.Op&fsnotify.Write == fsnotify.Write {
				w.changes <- struct{}{}
			}
		}
	}()

	return w, nil
}

// watcher watches files for changes.
// changes are signaled via the "changes" channel.
// errors are logged.
type watcher struct {
	watcher  *fsnotify.Watcher
	changes  chan struct{}
	watching []string
}

// Watch watches a file for changes.
func (w *watcher) Watch(file string) {
	w.watching = append(w.watching, file)
	w.watcher.Add(file)
}

// StopWatchingAll stops watching all files.
func (w *watcher) StopWatchingAll() {
	for _, f := range w.watching {
		w.watcher.Remove(f)
	}
	w.watching = nil
}
