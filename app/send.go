package app

import (
	"encoding/gob"
	"os"
	"time"
)

type Frame struct {
	Time time.Time
}

func Send(doc *Doc) {
	err := gob.NewEncoder(os.Stdout).Encode(doc)
	if err != nil {
		os.Stderr.Write([]byte("sending: "))
		os.Stderr.Write([]byte(err.Error()))
		os.Stderr.Write([]byte("\n"))
	}
}
