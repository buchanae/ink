package gfx

import (
	"encoding/gob"
	"os"
	"time"
)

type Frame struct {
	Time time.Time
}

func Run(f func(*Doc)) {

	doc := &Doc{}
	f(doc)
	send(doc)

	if doc.OnFrame == nil {
		return
	}

	for {
		time.Sleep(20 * time.Millisecond)
		frame := Frame{Time: time.Now()}
		doc.OnFrame(frame)
		send(doc)
	}
}

func send(doc *Doc) {
	err := gob.NewEncoder(os.Stdout).Encode(doc.nodes)
	if err != nil {
		os.Stderr.Write([]byte("sending: "))
		os.Stderr.Write([]byte(err.Error()))
		os.Stderr.Write([]byte("\n"))
	}
}
