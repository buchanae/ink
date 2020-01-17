package main

import (
	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
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

func Base(doc *app.Doc) {
	shape := Rect{
		A: XY{.1, .1},
		B: XY{.9, .5},
	}

	gfx.Fill{
		Mesh:  shape,
		Color: Lightgray,
	}.Draw(doc)
	gfx.Stroke{
		Target: shape,
		Color:  Slategray,
	}.Draw(doc)
}

func Roof(doc *app.Doc) {
	shape := Rect{
		A: XY{.05, .5},
		B: XY{.95, .75},
	}

	gfx.Fill{
		Mesh:  shape,
		Color: Gray,
	}.Draw(doc)

	gfx.Stroke{
		Target: shape,
		Color:  Black,
	}.Draw(doc)

	gfx.Fill{
		Mesh: Rect{
			A: XY{.055, .51},
			B: XY{.945, .52},
		},
		Color: HexString("#dddddd"),
	}.Draw(doc)

	Chimney(doc, RectAWH(
		XY{
			X: rand.Range(.1, .7),
			Y: rand.Range(.6, .7),
		},
		XY{
			X: rand.Range(.02, .2),
			Y: rand.Range(.05, .2),
		},
	))
}

func Chimney(doc *app.Doc, rect Rect) {
	gfx.Fill{
		Mesh:  rect,
		Color: Crimson,
	}.Draw(doc)
}

func Bushes(doc *app.Doc) {
	circles := []Circle{
		{
			XY:     XY{0.1, 0.13},
			Radius: 0.03,
		},
		{
			XY:     XY{0.2, 0.13},
			Radius: 0.03,
		},
		{
			XY:     XY{0.3, 0.13},
			Radius: 0.03,
		},
		{
			XY:     XY{0.7, 0.13},
			Radius: 0.03,
		},
		{
			XY:     XY{0.8, 0.13},
			Radius: 0.03,
		},
		{
			XY:     XY{0.9, 0.13},
			Radius: 0.03,
		},
	}

	for _, c := range circles {
		gfx.Fill{c, Green}.Draw(doc)
	}
}

func Door(doc *app.Doc) {
	base := Rect{
		A: XY{.42, .1},
		B: XY{.58, .3},
	}

	// Door trim
	gfx.Fill{
		Mesh:  base,
		Color: White,
	}.Draw(doc)

	// door
	amt := rand.Range(0.01, 0.03)
	door := Rect{
		A: XY{base.A.X + amt, base.A.Y},
		B: XY{base.B.X - amt, base.B.Y - amt},
	}
	gfx.Fill{
		Mesh:  door,
		Color: Red,
	}.Draw(doc)

	// handle
	DoorHandle(doc, door)
}

func DoorHandle(doc *app.Doc, door Rect) {
	handle := door.Interpolate(XY{
		X: rand.Range(0.8, .95),
		Y: rand.Range(0.4, 0.6),
	})
	gfx.Fill{
		Mesh: Circle{
			XY:     handle,
			Radius: rand.Range(0.002, 0.008),
		},
		Color: Black,
	}.Draw(doc)
}

func Ground(doc *app.Doc) {
	gfx.Stroke{
		Target: Line{A: XY{0, .1}, B: XY{1, .1}},
		Color:  Black,
	}.Draw(doc)
}

func Windows(doc *app.Doc) {
	gen := Window{
		Rows:       rand.IntRange(1, 5),
		Cols:       rand.IntRange(1, 5),
		ShrinkBase: rand.Range(0.004, 0.007),
		ShrinkPane: rand.Range(0.001, 0.003),
	}

	gen.Gen(doc, Rect{
		A: XY{.2, .15},
		B: XY{.3, .3},
	})

	gen.Gen(doc, Rect{
		A: XY{.7, .15},
		B: XY{.8, .3},
	})
}

func SecondFloor(doc *app.Doc) {
	gen := Window{
		Rows:       rand.IntRange(1, 5),
		Cols:       rand.IntRange(1, 5),
		ShrinkBase: rand.Range(0.004, 0.007),
		ShrinkPane: rand.Range(0.001, 0.003),
	}

	gen.Gen(doc, Rect{
		A: XY{.2, .35},
		B: XY{.3, .45},
	})

	gen.Gen(doc, Rect{
		A: XY{.4, .35},
		B: XY{.6, .45},
	})

	gen.Gen(doc, Rect{
		A: XY{.7, .35},
		B: XY{.8, .45},
	})
}

type Window struct {
	Rows       int
	Cols       int
	ShrinkBase float32
	ShrinkPane float32
}

func (w Window) Gen(doc *app.Doc, base Rect) {
	gfx.Fill{
		Mesh:  base,
		Color: White,
	}.Draw(doc)

	grid := Grid{
		Rows: w.Rows,
		Cols: w.Cols,
		Rect: base.Shrink(w.ShrinkBase),
	}
	for _, cell := range grid.Cells() {
		pane := cell.Rect.Shrink(w.ShrinkPane)
		gfx.Fill{
			Mesh:  pane,
			Color: Gray,
		}.Draw(doc)
	}
}

func Ink(doc *app.Doc) {
	rand.SeedNow()
	gfx.Clear(doc, White)

	Base(doc)
	Roof(doc)
	Door(doc)
	Windows(doc)
	SecondFloor(doc)
	Ground(doc)
	Bushes(doc)

}
