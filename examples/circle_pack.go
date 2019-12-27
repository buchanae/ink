package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
	"github.com/buchanae/ink/voronoi"
)

func Ink(doc *Doc) {
	rand.SeedNow()
	doc.Clear(White)

	noise := rand.BlueNoise(5000, 1, 1, 0.02)
	box := Rect{
		A: XY{.1, .1},
		B: XY{.9, .9},
	}

	var xys []XY
	for _, xy := range noise {
		if !box.Contains(xy) {
			continue
		}
		if rand.Bool(0.3) {
			continue
		}
		xys = append(xys, xy)
	}

	for _, xy := range xys {
		doc.Dot(xy, Red)
	}

	colors := rand.Palette()

	v := voronoi.New(xys, box)
	for _, cell := range v.Cells() {
		c := rand.Color(colors)
		c.A = 0.3
		for _, tri := range cell.Tris {
			s := doc.Shader(tri)
			s.Set("a_color", c)
		}

		for _, e := range cell.Edges {
			m := e.Stroke(0.002)
			doc.Shader(m)
		}
	}
}
