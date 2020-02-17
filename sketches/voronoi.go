package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
	"github.com/buchanae/ink/voronoi"
)

func Ink(doc gfx.Doc) {
	rand.SeedNow()

	box := Rect{
		A: XY{.1, .1},
		B: XY{.9, .9},
	}

	var initial []XY

	for i := 0; i < 20; i++ {
		p := float32(i) / 20
		xy := box.Interpolate(XY{p, p})
		initial = append(initial, xy)
	}

	noise := rand.BlueNoise{
		Limit:   450,
		Spacing: 0.05,
		Initial: initial,
		Rect:    box,
	}.Generate()

	v := voronoi.New(noise, box)

	colors := []RGBA{
		Blue, Yellow, Green, Black, Purple,
	}

	tris := v.Triangulate()

	for _, t := range tris {
		c := rand.Color(colors)
		c.A = 0.3
		gfx.Fill{
			Shape: t,
			Color: c,
		}.Draw(doc)
	}

	gfx.Stroke{
		Shape: Triangles(tris),
		Width: 0.001,
		Color: Black,
	}.Draw(doc)
}
