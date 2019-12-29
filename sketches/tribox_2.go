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
	Stroke       = false
	StrokeWidth  = 0.001
)

func Ink(doc *Doc) {
	rand.SeedNow()

	Clear(doc, White)

	grid := NewGrid(10, 10)
	colors := rand.Palette()

	l2 := doc.NewLayer()
	mask := doc.NewLayer()

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

			s := NewShader(t)
			s.Set("a_color", rand.Color(colors))
			s.Draw(l2)
		}

		if Stroke {
			stk := tris.Stroke()
			stk.Width = StrokeWidth
			s := NewShader(stk)
			s.Set("a_color", White)
			s.Draw(l2)
		}

		Fill{
			Mesh: Circle{
				XY:       r.Center(),
				Radius:   CircleRadius,
				Segments: 40,
			},
			Color: Black,
		}.Draw(mask)
	}

	Mask{
		Rect:   Fullscreen,
		Source: l2,
		Mask:   mask,
	}.Draw(doc)
}
