// +build darwin

package main

import (
	"log"
	"time"

	"github.com/fsnotify/fsevents"
)

// newWatcher returns a new file watcher.
func newWatcher() (*watcher, error) {
	return &watcher{
		changes: make(chan struct{}, 100),
	}, nil
}

// watcher watches files for changes.
// changes are signaled via the "changes" channel.
// errors are logged.
type watcher struct {
	watcher *fsevents.EventStream
	changes chan struct{}
}

func (w *watcher) init(path string) {
	dev, err := fsevents.DeviceForPath("/")
	if err != nil {
		log.Print(err)
		return
	}

	es := &fsevents.EventStream{
		Latency: 500 * time.Millisecond,
		Flags:   fsevents.FileEvents | fsevents.NoDefer,
		Paths:   []string{path},
		Device:  dev,
	}
	es.Start()
	w.watcher = es

	go func() {
		for msg := range es.Events {
			for _, event := range msg {
				if event.Flags&fsevents.ItemModified == fsevents.ItemModified {

					w.changes <- struct{}{}
				}
			}
		}
	}()
}

// Watch watches a file for changes.
func (w *watcher) Watch(file string) {
	if w.watcher == nil {
		w.init(file)
	} else {
		w.watcher.Paths = append(w.watcher.Paths, file)
		w.watcher.Restart()
	}
}

// StopWatchingAll stops watching all files.
func (w *watcher) StopWatchingAll() {
	panic("unimplemented")
}
