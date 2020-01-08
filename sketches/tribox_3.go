package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

const (
	ShrinkRect = 0.019
	RandPoint  = 0.03
)

func Ink(doc Layer) {
	rand.SeedNow()

	Clear(doc, White)

	grid := NewGrid(10, 10)
	colors := rand.Palette()

	for _, r := range grid.Rects() {
		r = r.Shrink(ShrinkRect)
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
			// TODO Fill is very common but doesn't
			//      merge well because of the uniform
			//      u_color.
			//f := Fill{t, rnd.Color(colors)}
			//l2.Draw(f)

			s := NewShader(t)
			s.Set("a_color", rand.Color(colors))
			s.Draw(doc)
		}
	}
}
