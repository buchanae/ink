package gfx

import (
	"encoding/gob"
	"log"
	"os"
	"time"
)

func init() {
	log.SetFlags(0)
}

type Doc struct {
	Layer
	OnFrame func(Frame)
}

type Frame struct {
	Time time.Time
}

func Run(f func(*Doc)) {

	doc := &Doc{}
	f(doc)
	send(doc.Layer)

	if doc.OnFrame == nil {
		return
	}

	for {
		time.Sleep(20 * time.Millisecond)
		frame := Frame{Time: time.Now()}
		doc.OnFrame(frame)
		send(doc.Layer)
	}
}

func send(layer Layer) {
	err := gob.NewEncoder(os.Stdout).Encode(layer)
	if err != nil {
		os.Stderr.Write([]byte("sending: "))
		os.Stderr.Write([]byte(err.Error()))
		os.Stderr.Write([]byte("\n"))
	}
}
