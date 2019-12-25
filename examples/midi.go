package main

import (
	"log"

	"github.com/rakyll/portmidi"
)

func main() {

	portmidi.Initialize()
	//log.Print(portmidi.CountDevices())
	//log.Print(portmidi.Info(0))
	did := portmidi.DefaultInputDeviceID()
	in, err := portmidi.NewInputStream(did, 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	log.Print(did)

	for event := range in.Listen() {
		log.Print(event)
	}
}
