package app

import (
	"encoding/gob"
	"os"
	"time"
)

var enc *gob.Encoder

func init() {
	enc = gob.NewEncoder(os.Stdout)
}

type Frame struct {
	Time time.Time
}

type Playback struct {
}

// TODO wants a better design with start/stop/pause
//      accurate timing, frame delta time, etc.
//      need feedback on how long the send+render took
func (pb Playback) Next() bool {
	time.Sleep(20 * time.Millisecond)
	return true
}

func Play() bool {
	time.Sleep(15 * time.Millisecond)
	return true
}

func Send(doc *Doc) {
	err := enc.Encode(doc)
	if err != nil {
		os.Stderr.Write([]byte("sending: "))
		os.Stderr.Write([]byte(err.Error()))
		os.Stderr.Write([]byte("\n"))
	}
}
