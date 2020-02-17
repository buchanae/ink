package gfx

import (
	"bytes"
	"encoding/gob"
	"io"
	"log"
	"os"
	"time"

	"github.com/buchanae/ink/render"
	"github.com/buchanae/ink/trac"
)

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

func init() {
	dec = gob.NewDecoder(os.Stdin)
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
	Plan render.Plan
}

func Main(inkfunc func(*Doc)) {
	trac.Reset()
	log.SetFlags(0)
	doc := RecvDoc()

	span := trac.Start("sketch: Ink")
	inkfunc(doc)
	span.End()

	Send(doc)
}

func Send(doc *Doc) {

	// TODO move plan optimization out of opengl to client-side
	span := trac.Start("sketch: build plan")
	plan := doc.Plan()
	span.End()
	trac.Log("sketch: len(plan.FaceData): %d", len(plan.FaceData))
	trac.Log("sketch: len(plan.AttrData): %d", len(plan.AttrData))

	span = trac.Start("sketch: Encode")

	var enc *gob.Encoder
	/*
		var out io.Writer
		out = os.Stdout
		out = bufio.NewWriterSize(out, 4096)
	*/
	out := &bytes.Buffer{}

	var wdbg = &writeDebug{w: out}
	enc = gob.NewEncoder(wdbg)

	err := enc.Encode(RenderMessage{
		//Config: doc.Config,
		Plan: plan,
	})

	span.End()

	// TODO ongoing debugging
	if wdbg.count > 1_000_000 {
		trac.Log("sent %d MB\n", wdbg.count/1000000)
	} else {
		trac.Log("sent %d bytes\n", wdbg.count)
	}

	span = trac.Start("sketch: Send")
	io.Copy(os.Stdout, out)
	span.End()

	if err != nil {
		os.Stderr.Write([]byte("sending: "))
		os.Stderr.Write([]byte(err.Error()))
		os.Stderr.Write([]byte("\n"))
	}
}
