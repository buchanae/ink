package main

import (
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

/*
Floors
Eves
- Thick

Gutters

Windows
- Shutters
Doors
Symmetry
Siding
Colors
Roof lines
- Even
- Multi step but flat
- Pointy

Arcs
Circles
Stained glass

garage

front yard

mail box

pillars


wide vs tall vs both vs neither

lighting and shadows, depth

Porch

Chimney

Stairs

Bushes

Basement / Raised first floor


Decorative:
- Eves
- Windows
*/

func Base(doc gfx.Layer) Rect {
	origin := XY{.5, .1}
	wh := rand.XYRange(.2, .4)
	shape := Rect{
		A: XY{
			X: origin.X - wh.X,
			Y: origin.Y,
		},
		B: XY{
			X: origin.X + wh.X,
			Y: origin.Y + wh.Y,
		},
	}

	gfx.Stroke{
		Shape: shape,
		Color: color.Black,
	}.Draw(doc)
	return shape
}

type Op func(gfx.Layer, Rect)

var ops []Op

func window(doc gfx.Layer, rect Rect) {
	amt := rand.XYRange(0.5, 0.9)
	size := rect.Size()
	s := rect.ShrinkXY(XY{
		X: size.X * amt.X,
		Y: size.Y * amt.Y,
	})

	gfx.Stroke{
		Shape: s,
		Color: color.Blue,
	}.Draw(doc)
}

func door(doc gfx.Layer, rect Rect) {
	amt := rand.XYRange(0.5, 0.9)
	wh := rect.Size().Mul(amt).MulScalar(.5)
	origin := XY{
		X: rect.Center().X,
		Y: rect.A.Y,
	}

	gfx.Stroke{
		Shape: Rect{
			A: XY{
				X: origin.X - wh.X,
				Y: origin.Y,
			},
			B: XY{
				X: origin.X + wh.X,
				Y: origin.Y + wh.Y,
			},
		},
		Color: color.Red,
	}.Draw(doc)
}

func divide(doc gfx.Layer, rect Rect) {

	grid := dd.Grid{
		Rect: rect,
		Rows: rand.IntRange(0, 4),
		Cols: rand.IntRange(0, 4),
	}

	for _, cell := range grid.Cells() {
		if cell.Rect.Area() > .01 {
			println("area", cell.Rect.Area())

			i := rand.IntRange(0, len(ops))
			ops[i](doc, cell.Rect)
		}
		gfx.Stroke{
			Shape: cell.Rect,
			Color: color.Black,
		}.Draw(doc)
	}
}

func Ink(doc gfx.Doc) {
	rand.SeedNow()
	gfx.Clear(doc, White)

	ops = []Op{divide, window, door}

	base := Base(doc)
	divide(doc, base)

}
