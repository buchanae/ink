// +build darwin

package main

import (
	"time"

	"github.com/fsnotify/fsevents"
)

// newWatcher returns a new file watcher.
func newWatcher() (*watcher, error) {

	es := &fsevents.EventStream{
		Latency: 500 * time.Millisecond,
		Flags:   fsevents.FileEvents,
	}
	es.Start()

	w := &watcher{
		watcher: es,
		changes: make(chan struct{}, 100),
	}

	go func() {
		for msg := range es.Events {
			for _, event := range msg {
				if event.Flags&fsevents.ItemModified == fsevents.ItemModified {

					w.changes <- struct{}{}
				}
			}
		}
	}()

	return w, nil
}

// watcher watches files for changes.
// changes are signaled via the "changes" channel.
// errors are logged.
type watcher struct {
	watcher *fsevents.EventStream
	changes chan struct{}
}

// Watch watches a file for changes.
func (w *watcher) Watch(file string) {
	w.watcher.Paths = append(w.watcher.Paths, file)
	w.watcher.Restart()
}

// StopWatchingAll stops watching all files.
func (w *watcher) StopWatchingAll() {
	panic("unimplemented")
}
