package midi

import (
	"log"
	"time"

	"github.com/rakyll/portmidi"
)

type MidiVal struct {
	Key int
	Val float32
}

type Midi struct {
	Debounce time.Duration

	data    map[int]int
	last    map[int]int
	version map[int]int
	change  chan MidiVal
}

func NewMidi() *Midi {
	return &Midi{
		Debounce: 100 * time.Millisecond,
		data:     map[int]int{},
		last:     map[int]int{},
		version:  map[int]int{},
		change:   make(chan MidiVal),
	}
}

func (m *Midi) Listen() <-chan MidiVal {
	return m.change
}

/*
func (m *Midi) OnKeyChange(key int, cb func()) {
}

func (m *Midi) Watch(key int) <-chan float32 {
	out := make(chan float32)
	go func() {
		// TODO this will conflict with Listen.
		//      can't have two consumers sharing
		for c := range m.change {
			if c.Key == key {
				out <- c.Val
			}
		}
	}()
	return out
}
*/

func (m *Midi) Run() {
	portmidi.Initialize()
	did := portmidi.DefaultInputDeviceID()
	in, err := portmidi.NewInputStream(did, 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	tick := time.Tick(m.Debounce)
	events := in.Listen()

	for {
		select {
		case <-tick:
			for key, ver := range m.version {
				last := m.last[key]
				if ver != last {

					mv := MidiVal{
						Key: key,
						Val: float32(m.data[key]) / 127,
					}

					select {
					case m.change <- mv:
					default:
					}

					m.last[key] = ver
				}
			}

		case ev := <-events:
			key := int(ev.Data1)
			val := int(ev.Data2)
			m.version[key] += 1
			m.data[key] = val
		}
	}
}
