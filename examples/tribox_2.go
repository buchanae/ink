package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

const (
	ShrinkRect   = 0.003
	RandPoint    = 0.03
	CircleRadius = 0.05
	Stroke       = true
	StrokeWidth  = 0.001
)

func Ink(doc *Doc) {
	rand.SeedNow()

	doc.Clear(White)

	grid := NewGrid(10, 10)
	colors := rand.Palette()

	l2 := doc.Layer()
	mask := doc.Layer()

	for _, r := range grid.Rects() {
		r = r.Shrink(ShrinkRect)
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

			l2.Shader(t).
				Set("a_color", rand.Color(colors))
		}

		if Stroke {
			stk := tris.Stroke(StrokeWidth)
			l2.Shader(stk).
				Set("a_color", White)
		}

		mask.Draw(Fill{
			Mesh: Circle{
				XY:       r.Center(),
				Radius:   CircleRadius,
				Segments: 40,
			},
			Color: Black,
		})
	}

	doc.Draw(Mask{
		Rect:   Fullscreen,
		Source: l2,
		Mask:   mask,
	})
}
