package main

import (
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(doc gfx.Doc) {
	rand.SeedNow()

	grid := Grid{Rows: 10, Cols: 10}

	for _, cell := range grid.Cells() {
		r := cell.Rect
		bnd := r.Shrink(0.013)

		current := bnd.Interpolate(rand.XYRange(0.1, 0.9))
		pen := &Pen{}
		i := 0
		horizontal := false

		pen.MoveTo(current)

		for i < 20 {
			var add XY

			if horizontal {
				add.X = rand.Range(-0.2, 0.2)
			} else {
				add.Y = rand.Range(-0.2, 0.2)
			}

			next := current.Add(add)
			if !bnd.Contains(next) {
				continue
			}

			pen.LineTo(next)
			current = next
			horizontal = !horizontal
			i++
		}

		pen.Close()

		gfx.Stroke{
			Shape: pen,
			Width: 0.001,
			Color: color.Black,
		}.Draw(doc)
	}
}
