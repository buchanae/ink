package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

const (
	ShrinkRect   = 0.003
	RandPoint    = 0.03
	CircleRadius = 0.05
	ShouldStroke = false
	StrokeWidth  = 0.001
)

func Ink(doc gfx.Doc) {
	rand.SeedNow()

	grid := Grid{Rows: 10, Cols: 10}
	colors := rand.Palette()

	l2 := doc.NewLayer()
	mask := doc.NewLayer()

	for _, cell := range grid.Cells() {
		r := cell.Rect.Shrink(ShrinkRect)
		q := r.Quad()
		p := rand.XYInRect(r.Shrink(RandPoint))

		tris := Triangles{
			{q.A, q.B, p},
			{q.B, q.C, p},
			{q.C, q.D, p},
			{q.D, q.A, p},
		}

		for _, t := range tris {
			s := gfx.NewShader(t.Fill())
			s.Set("a_color", rand.Color(colors))
			s.Draw(l2)
		}

		if ShouldStroke {
			gfx.Stroke{
				Shape: tris,
				Width: StrokeWidth,
				Color: White,
			}.Draw(l2)
		}

		gfx.Fill{
			Shape: Circle{
				XY:       r.Center(),
				Radius:   CircleRadius,
				Segments: 40,
			},
			Color: Black,
		}.Draw(mask)
	}

	gfx.Mask{
		Rect:   gfx.Fullscreen,
		Source: l2,
		Mask:   mask,
	}.Draw(doc)
}
