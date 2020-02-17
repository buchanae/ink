package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(doc gfx.Doc) {

	grid := Grid{
		Rows: 15,
		Cols: 15,
		Rect: SquareCenter(gfx.Center, .8),
	}

	for _, cell := range grid.Cells() {
		const Z = 0.007
		const G = 0.005

		r := cell.Rect.Translate(XY{
			Y: rand.Range(-Z, Z),
		})
		grow := rand.Range(0, G)
		r.A.X -= grow
		r.B.X += grow

		col := Red
		col.A = 0.7
		gfx.Fill{r, col}.Draw(doc)
	}
}
