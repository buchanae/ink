package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(doc gfx.Doc) {
	gfx.Clear(doc, White)

	grid := Grid{
		Rows: 40,
		Cols: 40,
		Rect: RectCenter(XY{.5, .5}, XY{.95, .95}),
	}
	col := Blue
	col.A = 0.7

	for _, cell := range grid.Cells() {
		r := cell.Rect.Shrink(0.001)
		size := r.Size()
		center := r.Center()
		a := XY{
			X: r.A.X + size.X/2,
			Y: r.A.Y,
		}
		b := XY{
			X: a.X,
			Y: r.B.Y,
		}
		rot := rand.Angle()

		gfx.Stroke{
			Shape: Line{
				a.RotateAround(rot, center),
				b.RotateAround(rot, center),
			},
			Color: col,
			Width: 0.0025,
		}.Draw(doc)
	}
}
