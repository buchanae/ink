package main

import (
	"github.com/buchanae/ink/app"
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

func Ink(doc *app.Doc) {
	rand.SeedNow()

	gfx.Clear(doc, White)

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
			// TODO Fill is very common but doesn't
			//      merge well because of the uniform
			//      u_color.
			//f := Fill{t, rnd.Color(colors)}
			//l2.Draw(f)

			s := gfx.NewShader(t)
			s.Set("a_color", rand.Color(colors))
			s.Draw(l2)
		}

		if ShouldStroke {
			gfx.Stroke{
				Target: tris,
				Width:  StrokeWidth,
				Color:  White,
			}.Draw(l2)
		}

		gfx.Fill{
			Mesh: Circle{
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
