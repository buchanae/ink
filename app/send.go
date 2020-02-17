package app

import (
	"encoding/gob"
	"io"
	"log"
	"os"
	"time"

	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/render"
)

var enc *gob.Encoder
var dec *gob.Decoder

type writeDebug struct {
	w     io.Writer
	count int
}

func (wd *writeDebug) Write(b []byte) (int, error) {
	n, err := wd.w.Write(b)
	wd.count += n
	return n, err
}

var wdbg = &writeDebug{w: os.Stdout}

func init() {
	enc = gob.NewEncoder(wdbg)
	dec = gob.NewDecoder(os.Stdin)
}

type Frame struct {
	Time time.Time
}

type Playback struct{}

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

func RecvDoc() *Doc {
	doc := &Doc{}
	err := dec.Decode(doc)
	if err != nil {
		os.Stderr.Write([]byte("RecvDoc: "))
		os.Stderr.Write([]byte(err.Error()))
		os.Stderr.Write([]byte("\n"))
	}
	return doc
}

type RenderMessage struct {
	Config gfx.Config
	Plan   render.Plan
}

func Main(inkfunc func(gfx.Doc)) {
	log.SetFlags(0)
	doc := RecvDoc()
	inkfunc(doc)
	Send(doc)
}

func Send(doc *Doc) {
	wdbg.count = 0
	//start := time.Now()

	// TODO move plan optimization out of opengl to client-side
	plan := buildPlan(doc)

	err := enc.Encode(RenderMessage{
		Config: *doc.Conf,
		Plan:   plan,
	})

	// TODO ongoing debugging
	//fmt.Fprintf(os.Stderr, "sent %d bytes\n", wdbg.count)
	//fmt.Fprintf(os.Stderr, "send took: %s\n", time.Since(start))

	if err != nil {
		os.Stderr.Write([]byte("sending: "))
		os.Stderr.Write([]byte(err.Error()))
		os.Stderr.Write([]byte("\n"))
	}
}
