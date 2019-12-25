package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(doc *Layer) {
	rand.SeedNow()
	doc.Clear(White)

	// TODO grid size is confusing
	//      because the grid is actually a point grid, not a rectangle grid
	grid := NewGrid(17, 33)
	box := RectCenter(XY{.5, .5}, XY{1, .5})

	var strokes []Mesh

	for i, rect := range grid.Rects() {

		r := rect.Shrink(0.000)
		//r = r.Rotate(-math.Pi/2, XY{.5, .5})
		r = Rect{
			A: box.Interp(r.A),
			B: box.Interp(r.B),
		}
		stk := r.Stroke(0.0005)
		strokes = append(strokes, stk)

		sub := NewGrid(4, 4)
		for j, sr := range sub.Rects() {

			a := r.Interp(sr.A)
			b := r.Interp(sr.B)
			xr := Rect{a, b}

			//stk := xr.Stroke(0.0005)
			//strokes = append(strokes, stk)

			// TODO interleaving a stroke
			//      causes all the batching to fail
			// TODO move these things to an examples
			//      folder demonstrating performance
			//      issues
			//doc.Shader(stk)

			mask := 1 << j
			if i&mask == mask {
				doc.Shader(xr)
			}
		}
	}

	for _, stk := range strokes {
		shd := doc.Shader(stk)
		shd.Set("a_color", White)
	}
}
