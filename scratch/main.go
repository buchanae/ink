package main

import (
	"log"
	"runtime"
	"time"
)

func init() {
	runtime.LockOSThread()
	go func() {
		for range time.Tick(time.Second) {
			log.Print("main")
		}
	}()
	runtime.UnlockOSThread()
}

func main() {
	go func() {
		for range time.Tick(time.Second) {
			log.Print("tick")
		}
	}()

	block := make(chan struct{})
	<-block
	log.Print("done")
}
