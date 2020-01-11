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

	grid := Grid{Rows: 10, Cols: 10}
	//p := rand.Palette()

	for _, cell := range grid.Cells() {
		r := cell.Rect
		bnd := r.Shrink(0.013)

		/*
			Fill{
				Mesh:  r,
				Color: rand.Color(p),
			}.Draw(doc)
		*/

		current := bnd.Interpolate(rand.XYRange(0.1, 0.9))
		path := Path{}
		i := 0
		horizontal := false

		path.MoveTo(current)

		for i < 20 {
			var add XY

			//if rand.Bool(0.5) {
			if horizontal {
				add.X = rand.Range(-0.2, 0.2)
			} else {
				add.Y = rand.Range(-0.2, 0.2)
			}

			next := current.Add(add)
			if !bnd.Contains(next) {
				continue
			}

			path.LineTo(next)
			current = next
			horizontal = !horizontal
			i++
		}

		path.Close()

		gfx.Stroke{
			Shape: &path,
			Width: 0.001,
			Color: color.Black,
		}.Draw(doc)
		//stk := path.Stroke()
		//stk.Width = 0.001
		//Fill{stk, color.Black}.Draw(doc)
	}
}
