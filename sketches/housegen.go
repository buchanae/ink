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
	wh := rand.XYRange(.3, .4)
	shape := Rect{
		A: XY{
			X: .5 - wh.X,
			Y: .1,
		},
		B: XY{
			X: .5 + wh.X,
			Y: .1 + wh.Y,
		},
	}

	gfx.Fill{
		Mesh:  shape,
		Color: randGray(),
	}.Draw(doc)
	gfx.Stroke{
		Target: shape,
		Color:  randGray(),
	}.Draw(doc)

	Roof(doc, shape.B.Y)
	Windows(doc, shape)
}

func randGray() RGBA {
	g := rand.Range(0, 1)
	return RGBA{g, g, g, 1}
}

func Roof(doc *app.Doc, baseY float32) {
	wh := rand.XYRange(.2, .4)
	shape := Rect{
		A: XY{.5 - wh.X, baseY},
		B: XY{.5 + wh.X, baseY + wh.Y},
	}

	quad := shape.Quad()
	if rand.Bool(0.3) {
		amt := rand.Range(0.01, 0.1)
		quad.B.X -= amt
		quad.C.X += amt
	}

	gfx.Fill{
		Mesh:  quad,
		Color: randGray(),
	}.Draw(doc)

	gfx.Stroke{
		Target: quad,
		Color:  randGray(),
	}.Draw(doc)

	/*
		gfx.Fill{
			Mesh: Rect{
				A: XY{.055, .51},
				B: XY{.945, .52},
			},
			Color: randGray(),
		}.Draw(doc)
	*/

	Chimney(doc, RectAWH(
		shape.Interpolate(XY{
			X: rand.Range(.1, .9),
			Y: rand.Range(.5, .9),
		}),
		XY{
			X: rand.Range(.02, .2),
			Y: rand.Range(.05, .2),
		},
	))
}

func Chimney(doc *app.Doc, rect Rect) {
	gfx.Fill{
		Mesh:  rect,
		Color: randGray(),
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

	col := randGray()
	for _, c := range circles {
		gfx.Fill{c, col}.Draw(doc)
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
		Color: randGray(),
	}.Draw(doc)

	// door
	amt := rand.Range(0.01, 0.03)
	door := Rect{
		A: XY{base.A.X + amt, base.A.Y},
		B: XY{base.B.X - amt, base.B.Y - amt},
	}
	gfx.Fill{
		Mesh:  door,
		Color: randGray(),
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
		Color: randGray(),
	}.Draw(doc)
}

func Ground(doc *app.Doc) {
	gfx.Stroke{
		Target: Line{A: XY{0, .1}, B: XY{1, .1}},
		Color:  randGray(),
	}.Draw(doc)
}

func Windows(doc *app.Doc, base Rect) {

	win := Window{
		Rows:       rand.IntRange(1, 5),
		Cols:       rand.IntRange(1, 5),
		ShrinkBase: rand.Range(0.004, 0.007),
		ShrinkPane: rand.Range(0.001, 0.003),
		PaneColor:  randGray(),
		TrimColor:  randGray(),
	}

	win.Gen(doc, base.SubRect(Rect{
		A: XY{.1, .15},
		B: XY{.3, .5},
	}))

	win.Gen(doc, base.SubRect(Rect{
		A: XY{.7, .15},
		B: XY{.9, .5},
	}))

	second := base.SubRect(Rect{
		A: XY{0, .6},
		B: XY{1, 1},
	})

	SecondFloor(doc, win, second)
}

func SecondFloor(doc *app.Doc, win Window, base Rect) {
	win.Rows = rand.IntRange(1, 5)
	win.Cols = rand.IntRange(1, 5)
	win.ShrinkBase = rand.Range(0.004, 0.007)
	win.ShrinkPane = rand.Range(0.001, 0.003)

	win.Gen(doc, base.SubRect(Rect{
		A: XY{.1, .15},
		B: XY{.3, .85},
	}))

	win.Gen(doc, base.SubRect(Rect{
		A: XY{.4, .15},
		B: XY{.6, .85},
	}))

	win.Gen(doc, base.SubRect(Rect{
		A: XY{.7, .15},
		B: XY{.9, .85},
	}))
}

type Window struct {
	Rows       int
	Cols       int
	ShrinkBase float32
	ShrinkPane float32
	PaneColor  RGBA
	TrimColor  RGBA
}

func (w Window) Gen(doc *app.Doc, base Rect) {
	gfx.Fill{
		Mesh:  base,
		Color: w.TrimColor,
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
			Color: w.PaneColor,
		}.Draw(doc)
	}
}

func Ink(doc *app.Doc) {
	rand.SeedNow()
	gfx.Clear(doc, White)

	Base(doc)
	Door(doc)
	Ground(doc)
	Bushes(doc)

}
