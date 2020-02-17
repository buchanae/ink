package main

import (
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(doc gfx.Doc) {
	rand.SeedNow()

	center := XY{0.5, 0.5}
	grid := Grid{
		Rows: 32,
		Cols: 16,
		Rect: RectCenter(center, XY{.5, .97}),
	}
	sub := Grid{Rows: 3, Cols: 3}

	var bold []Strokeable
	var strokes []Strokeable

	for i, cell := range grid.Cells() {
		r := cell.Rect.Shrink(0.003)
		bold = append(bold, r)

		for j, sc := range sub.Cells() {
			sr := sc.Rect

			xr := Rect{
				A: r.Interpolate(sr.A),
				B: r.Interpolate(sr.B),
			}

			strokes = append(strokes, xr)

			// TODO interleaving a stroke
			//      causes all the batching to fail
			// TODO move these things to an examples
			//      folder demonstrating performance
			//      issues
			//doc.Shader(stk)

			mask := 1 << j
			if i&mask == mask {
				gfx.Fill{xr, color.Black}.Draw(doc)
			}
		}
	}

	for _, stk := range strokes {
		gfx.Stroke{
			Shape: stk,
			Width: 0.0002,
			Color: color.Black,
		}.Draw(doc)
	}

	for _, stk := range bold {
		gfx.Stroke{
			Shape: stk,
			Width: 0.0009,
			Color: color.Black,
		}.Draw(doc)
	}
}
