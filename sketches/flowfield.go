package main

import (
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(doc *gfx.Doc) {

	const (
		Lines    = 10000
		Length   = 100
		GridSize = 50
		Opacity  = 0.5
		Decay    = 0.011
		Radius   = 0.002
		Segments = 5
		SegLen   = Radius * 5
		Jitter   = 0.9
		Tangent  = true
	)

	rand.SeedNow()

	grid := Grid{
		Rows: GridSize,
		Cols: GridSize,
		Rect: RectWH(1.0, 1.0),
	}

	vg := NewVectorGrid(grid)
	i := 0
	for row := 0; row < GridSize+1; row++ {
		for col := 0; col < GridSize+1; col++ {

			xy := XY{
				float32(col) / float32(GridSize),
				float32(row) / float32(GridSize),
			}
			xy = xy.Sub(Center)

			if Tangent {
				xy = XY{-xy.Y, xy.X}
			}

			if xy.IsZero() {
				xy = Unit()
			}
			xy = xy.Normalize()
			xy = xy.Rotate(rand.Range(-Jitter, Jitter))
			vg.Vecs[i] = xy

			i++
		}
	}

	palette := rand.Palette()
	pos := make([]XY, 0, Lines*Length)
	colors := make([]color.RGBA, 0, Lines*Length)

	for j := 0; j < Lines; j++ {

		//xy := rand.XYRange(0.1, 0.9)
		xy := rand.XYInCircle(Circle{
			XY:     Center,
			Radius: 0.4,
		})
		//xy := rand.XY()
		c := rand.Color(palette)
		c.A = Opacity

		for i := 0; i < Length; i++ {

			p := xy.Sub(grid.Rect.A).Div(grid.Size())

			vec, ok := vg.Interpolate(p)
			if !ok {
				continue
			}

			vec = vec.SetLength(SegLen)
			xy = xy.Add(vec)

			pos = append(pos, xy)
			cc := c
			cc.A += rand.Range(-0.5, 0.5)
			colors = append(colors, cc)
			c.A -= Decay
		}
	}

	shader := gfx.Fill{
		Shape: Circle{
			Radius:   Radius,
			Segments: Segments,
		},
	}.Shader()

	shader.Set("a_pos", pos)
	shader.Set("a_color", colors)
	shader.Instances = len(pos)
	shader.Divisors = map[string]int{
		"a_pos":   1,
		"a_color": 1,
	}
	shader.Draw(doc)
	//vg.DrawArrows(doc)
	//vg.DrawArrows2(doc)
}

type VectorGrid struct {
	Vecs []XY

	grid Grid
}

func NewVectorGrid(g Grid) VectorGrid {
	return VectorGrid{
		Vecs: make([]XY, (g.Rows+1)*(g.Cols+1)),
		grid: g,
	}
}

// TODO note "xy" is a percentage, range [0, 1]
func (vg VectorGrid) Interpolate(percent XY) (XY, bool) {
	nrow := vg.grid.Rows + 1
	ncol := vg.grid.Cols + 1

	// TODO should probably clamp to 0-1?
	row := int(percent.Y * float32(nrow-1))
	col := int(percent.X * float32(ncol-1))

	// TODO clamp row/col index to grid size
	ai := row*ncol + col
	bi := row*ncol + (col + 1)
	ci := (row+1)*ncol + (col + 1)
	di := (row+1)*ncol + col
	maxi := len(vg.Vecs) - 1

	// TODO hack. need better bounds logic.
	if ai < 0 || ai > maxi {
		return XY{}, false
	}
	if bi < 0 || bi > maxi {
		return XY{}, false
	}
	if ci < 0 || ci > maxi {
		return XY{}, false
	}
	if di < 0 || di > maxi {
		return XY{}, false
	}

	a := vg.Vecs[ai]
	b := vg.Vecs[bi]
	c := vg.Vecs[ci]
	d := vg.Vecs[di]

	xy := vg.grid.Rect.Interpolate(percent)
	cell := vg.grid.Cell(row, col)
	off := xy.Sub(cell.A)
	p := off.Div(cell.Size())

	e := a.InterpolateTo(b, p.X).Normalize()
	f := d.InterpolateTo(c, p.X).Normalize()

	return e.InterpolateTo(f, p.Y), true
}

func (vg VectorGrid) DrawArrows(doc *gfx.Doc) {
	GridSize := vg.grid.Rows + 1

	for row := 0; row < GridSize; row++ {
		for col := 0; col < GridSize; col++ {

			i := row*GridSize + col

			const size = 0.04

			vec := vg.Vecs[i].SetLength(size)

			xy := XY{
				float32(col) / float32(GridSize-1),
				float32(row) / float32(GridSize-1),
			}

			// create the vector arrow
			line := Line{
				A: xy,
				B: xy.Add(vec),
			}

			// draw the vector arrow
			gfx.Stroke{
				Shape: line,
				Color: color.Black,
				Width: 0.002,
			}.Draw(doc)

			gfx.Fill{
				Shape: Circle{
					XY:       line.A,
					Radius:   0.005,
					Segments: 6,
				},
				Color: color.Blue,
			}.Draw(doc)

		}
	}
}

func (vg VectorGrid) DrawArrows2(doc *gfx.Doc) {
	for _, cell := range vg.grid.Cells() {
		center := cell.Center()

		// figure out how long each vector arrow should be.
		// it should be half the size of a cell at most:
		size := cell.Size().Length() * 0.5
		// and then a bit smaller, just to add some padding:
		size *= 0.9

		// Determine the vector at this position.
		p := center.Sub(vg.grid.Rect.A).Div(vg.grid.Size())
		//angle, _ := vg.Interpolate(p)
		//vec := Unit().Rotate(angle)
		vec, _ := vg.Interpolate(p)

		// create the vector arrow
		line := Line{
			A: center,
			B: center.Add(vec.SetLength(size)),
		}

		// draw the vector arrow
		gfx.Stroke{
			Shape: line,
			Color: color.Red,
			Width: 0.002,
		}.Draw(doc)

		gfx.Fill{
			Shape: Circle{
				XY:       line.A,
				Radius:   0.003,
				Segments: 6,
			},
			Color: color.Red,
		}.Draw(doc)

		gc := color.Blue
		gc.A = 0.2
		gfx.Stroke{
			Shape: cell,
			Color: gc,
			Width: 0.0005,
		}.Draw(doc)
	}
}
