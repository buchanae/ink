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

	doc := NewDoc()
	f(doc)
	send(doc)
}

func send(doc *Doc) {
	err := gob.NewEncoder(os.Stdout).Encode(doc)
	if err != nil {
		os.Stderr.Write([]byte("sending: "))
		os.Stderr.Write([]byte(err.Error()))
		os.Stderr.Write([]byte("\n"))
	}
}
