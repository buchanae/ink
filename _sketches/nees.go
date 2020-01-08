package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/math"
	"github.com/buchanae/ink/rand"
)

const (
	Width  = 16
	Height = 30
)

func Ink(doc *Doc) {
	Clear(doc, White)

	box := RectCenter(XY{.5, .5}, XY{.5, .9})

	grid := NewGrid(Height+1, Width+1)
	for i, r := range grid.Rects() {

		r = Rect{
			A: box.Interp(r.A),
			B: box.Interp(r.B),
		}
		r = r.Shrink(0.002)

		row := i / Width
		row = Height - row - 5
		dr := float32(row) / Height
		dr = math.Clamp(dr, 0, 1)

		t := dr * 0.01
		r = r.Translate(rand.XYRange(-t, t))

		q := r.Quad()
		ang := rand.Range(-dr, dr)
		q = q.RotateAround(ang, r.Center())

		// TODO two types in a loop doesn't batch
		//Dot{XY: r.Center()}.Draw(doc)

		stk := q.Stroke()
		stk.Width = 0.001
		s := NewShader(stk)
		//s.Set("a_color", Red)
		s.Draw(doc)
	}
}
