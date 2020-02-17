package main

import (
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/math"
	"github.com/buchanae/ink/rand"
)

func Ink(doc gfx.Doc) {
	rand.SeedNow()

	ctx := gfx.NewContext(doc)

	center := XY{0.5, 0.5}
	grid := Grid{
		Rows: 25, Cols: 25,
		Rect: SquareCenter(center, 0.95),
	}

	lines := []Line{
		{XY{0, 0}, XY{1, 1}},
		{XY{0, 1}, XY{1, 0}},
		{XY{0, 0.5}, XY{1, 0.5}},
		{XY{0.5, 0}, XY{0.5, 1}},
		{XY{0, 0.5}, XY{0.5, 1}},
		{XY{0, 0.5}, XY{0.5, 0}},
		{XY{0.5, 0}, XY{1, 0.5}},
		{XY{0.5, 1}, XY{1, 0.5}},
	}

	for _, cell := range grid.Cells() {
		r := cell.Rect
		r = r.Shrink(0.005)

		p := float32(cell.Row) / float32(grid.Rows)
		n := int(math.Interp(1, 15, p))
		i := 0
		for i < n {
			l := lines[rand.Intn(len(lines))]
			c := r.Interpolate(l.A)
			d := r.Interpolate(l.B)

			ctx.Stroke(Line{c, d})

			i++
		}

		ctx.Stroke(r)
	}
}
