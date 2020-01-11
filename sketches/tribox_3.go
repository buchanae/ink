package main

import (
	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

const (
	ShrinkRect = 0.019
	RandPoint  = 0.03
)

func Ink(doc *app.Doc) {
	rand.SeedNow()

	gfx.Clear(doc, White)

	grid := Grid{Rows: 10, Cols: 10}
	colors := rand.Palette()

	for _, cell := range grid.Cells() {
		r := cell.Rect.Shrink(ShrinkRect)
		q := r.Quad()
		p := rand.XYInRect(r.Shrink(RandPoint))
		p = r.Center()

		//a := q.A.Add(XY{0.009, 0})
		q = rand.TweakQuad(q, 0.005)

		tris := Triangles{
			{q.A, q.B, p},
			{q.B, q.C, p},
			{q.C, q.D, p},
			{q.D, q.A, p},
		}

		for _, t := range tris {
			gfx.Fill{t, rand.Color(colors)}.Draw(doc)
		}
	}
}
