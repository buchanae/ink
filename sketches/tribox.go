package main

import (
	"log"
	"time"

	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
	"github.com/rakyll/portmidi"
)

func main() {
	portmidi.Initialize()
	did := portmidi.DefaultInputDeviceID()
	in, err := portmidi.NewInputStream(did, 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	data := map[int64]int64{}

	draw(0.003)

	events := in.Listen()
	tick := time.Tick(500 * time.Millisecond)
	var last, version int

	for {
		select {

		case <-tick:
			if version != last {
				last = version
				draw(0.003)
			}

		case event := <-events:
			log.Print(event)
			version++
			data[event.Data1] = event.Data2
		}
	}
}

func draw(shrink float32) {
	doc := NewDoc()

	doc.Draw(Clear(White))

	rnd := rand.New(time.Now().Unix())
	grid := NewGrid(10, 10)
	colors := rnd.Palette()

	for _, r := range grid.Rects() {
		r = r.Shrink(shrink)
		q := r.Quad()
		p := rnd.XYInRect(r.Shrink(0.01))

		/*
			point(doc, q.A, Black)
			point(doc, q.B, Black)
			point(doc, q.C, Black)
			point(doc, q.D, Black)
			point(doc, p, Red)
		*/

		tris := []Triangle{
			{q.A, q.B, p},
			{q.B, q.C, p},
			{q.C, q.D, p},
			{q.D, q.A, p},
		}

		for _, t := range tris {
			m := t.Mesh()
			doc.Draw(Fill{m, rnd.Color(colors)})
		}

		stk := StrokeTriangles(tris, 0.001)
		s := NewShader(stk)
		s.SetColor(White)
		doc.Draw(s)
	}

	app.Render(doc)
}

func point(doc *Doc, xy XY, col RGBA) {
	c := Circle{xy, 0.002}
	m := c.Mesh(10)
	s := NewShader(m)
	s.SetColor(col)
	doc.Draw(s)
}
