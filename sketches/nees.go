package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/math"
	"github.com/buchanae/ink/rand"
)

const (
	Width  = 16
	Height = 30
)

func Ink(doc gfx.Doc) {

	center := XY{.5, .5}
	grid := Grid{
		Rows: Height,
		Cols: Width,
		Rect: RectCenter(center, XY{.5, .9}),
	}

	for i, cell := range grid.Cells() {
		r := cell.Rect
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

		gfx.Stroke{
			Shape: q,
			Width: 0.001,
			Color: Black,
		}.Draw(doc)
	}
}
