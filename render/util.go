package render

import (
	loglib "log"
	"time"
)

var Trace bool

func log(msg string, args ...interface{}) {
	loglib.Printf(msg, args...)
}

func trace(msg string, args ...interface{}) {
	if Trace {
		loglib.Printf(msg, args...)
	}
}

func traceTime(msg string, start time.Time) {
	trace("%s took: %s", msg, time.Since(start))
}
