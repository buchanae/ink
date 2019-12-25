package render

import (
	"bufio"
	"log"
	"os"
	"time"
)

var Trace bool = false

var logger *log.Logger

func init() {
	buf := bufio.NewWriterSize(os.Stderr, 10000)
	logger = log.New(buf, "t: ", 0)
	go func() {
		for range time.Tick(time.Second) {
			buf.Flush()
		}
	}()
}

func trace(msg string, args ...interface{}) {
	if Trace {
		logger.Printf(msg, args...)
	}
}

func traceTime(name string) func() {
	start := time.Now()
	return func() {
		trace("%s took %s", name, time.Since(start))
	}
}
