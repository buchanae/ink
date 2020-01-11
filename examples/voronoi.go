package main

import (
	"github.com/buchanae/ink/app"
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
	"github.com/buchanae/ink/voronoi"
)

func Ink(doc *app.Doc) {
	rand.SeedNow()
	gfx.Clear(doc, color.White)

	rect := Rect{
		A: XY{.1, .1},
		B: XY{.9, .9},
	}

	bn := rand.BlueNoise{
		Rect:    rect,
		Spacing: 0.05,
		Limit:   400,
	}

	palette := rand.Palette()
	xys := bn.Generate(rand.Default)
	v := voronoi.New(xys, rect)
	tris := v.Triangulate()

	for _, t := range tris {
		c := rand.Color(palette)
		f := gfx.Fill{t, c}
		f.Draw(doc)
	}
}
