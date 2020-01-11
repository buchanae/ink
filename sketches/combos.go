package main

import (
	"github.com/buchanae/ink/app"
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(doc *app.Doc) {
	rand.SeedNow()
	gfx.Clear(doc, color.White)

	center := XY{0.5, 0.5}
	grid := Grid{
		Rows: 32,
		Cols: 16,
		Rect: RectCenter(center, XY{.5, .97}),
	}
	sub := Grid{Rows: 3, Cols: 3}

	var bold []Stroke
	var strokes []Stroke

	for i, cell := range grid.Cells() {
		r := cell.Rect.Shrink(0.003)
		bold = append(bold, r.Stroke())

		for j, sc := range sub.Cells() {
			sr := sc.Rect

			xr := Rect{
				A: r.Interpolate(sr.A),
				B: r.Interpolate(sr.B),
			}

			stk := xr.Stroke()
			stk.Width = 0.0005
			strokes = append(strokes, stk)

			// TODO interleaving a stroke
			//      causes all the batching to fail
			// TODO move these things to an examples
			//      folder demonstrating performance
			//      issues
			//doc.Shader(stk)

			mask := 1 << j
			if i&mask == mask {
				gfx.NewShader(xr).Draw(doc)
			}
		}
	}

	for _, stk := range strokes {
		stk.Width = 0.0002
		shd := gfx.NewShader(stk)
		shd.Set("a_color", color.Black)
		shd.Draw(doc)
	}

	for _, stk := range bold {
		stk.Width = 0.0009
		shd := gfx.NewShader(stk)
		shd.Set("a_color", color.Black)
		shd.Draw(doc)
	}
}
