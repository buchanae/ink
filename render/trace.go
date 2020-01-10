package render

import (
	"bufio"
	"log"
	"os"
	"time"
)

var logger *log.Logger
var start time.Time

func init() {
	buf := bufio.NewWriterSize(os.Stderr, 50000)
	logger = log.New(buf, "", 0)
	go func() {
		for range time.Tick(100 * time.Millisecond) {
			buf.Flush()
		}
	}()
}

type tracer struct {
	start time.Time
}

func (t *tracer) StartTrace() {
	t.start = time.Now()
}

func (t *tracer) EndTrace() {
	t.start = time.Time{}
}

func (t *tracer) trace(msg string, args ...interface{}) {
	if t.start.IsZero() {
		return
	}
	d := time.Since(t.start)
	// for relative times
	//start = time.Now()
	fms := float64(d) / float64(time.Millisecond)
	nargs := append([]interface{}{}, fms)
	nargs = append(nargs, args...)
	logger.Printf("%5.1fms "+msg, nargs...)
}
