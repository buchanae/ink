package main

import (
	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
	"github.com/buchanae/ink/voronoi"
)

func Ink(doc *app.Doc) {
	rand.SeedNow()

	box := Rect{
		A: XY{.1, .1},
		B: XY{.9, .9},
	}
	bn := rand.BlueNoise{
		Limit:   5000,
		Rect:    box,
		Spacing: 0.02,
	}

	var xys []XY
	for _, xy := range bn.Generate() {
		if rand.Bool(0.3) {
			continue
		}
		xys = append(xys, xy)
	}

	for _, xy := range xys {
		gfx.Dot{XY: xy}.Draw(doc)
	}

	colors := rand.Palette()

	v := voronoi.New(xys, box)
	for _, cell := range v.Cells() {
		c := rand.Color(colors)
		c.A = 0.3

		for _, tri := range cell.Tris {
			s := gfx.NewShader(tri)
			s.Set("a_color", c)
			s.Draw(doc)
		}

		for _, e := range cell.Edges {
			gfx.Stroke{
				Target: e,
				Width:  0.002,
			}.Draw(doc)
		}
	}
}
