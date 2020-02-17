package trac

import (
	"fmt"
	"log"
	"os"
	"time"
)

var logger *log.Logger
var Enabled bool
var start time.Time

func init() {
	logger = log.New(os.Stderr, "", 0)
	//log.Ltime|log.Lmicroseconds)
	Enabled = os.Getenv("INK_TRACE") == "true"
	Reset()
}

func Reset() {
	start = time.Now()
}

type Span struct {
	msg string
	t   time.Time
}

func (s Span) End() {
	dur := time.Since(s.t)
	Log("%s : done in %s", s.msg, dur)
}

func Start(msg string) Span {
	Log(msg)
	return Span{msg, time.Now()}
}

func Log(msg string, args ...interface{}) {
	if !Enabled {
		return
	}
	//prefix := time.Now().Format("04:05.000")
	//prefix := time.Now().Format("05.000")
	d := time.Since(start)
	prefix := fmt.Sprintf("%.3f", d.Seconds())
	logger.Printf(prefix+" "+msg, args...)
}
