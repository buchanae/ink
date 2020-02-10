package main

import (
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/math"
	"github.com/buchanae/ink/rand"
)

func Ink(doc *gfx.Doc) {
	rand.SeedNow()

	grid := Grid{
		Rows: 10,
		Cols: 10,
		Rect: RectWH(1.5, 1.5),
	}

	vg := NewVectorGrid(grid)
	for i := range vg.Angles {
		vg.Angles[i] = rand.Angle() * 8
	}

	palette := rand.Palette()
	pos := []XY{}
	colors := []color.RGBA{}

	for j := 0; j < 100; j++ {

		xy := rand.XY()
		c := rand.Color(palette)
		c.A = 0.20

		for i := 0; i < 100; i++ {

			p := xy.Sub(grid.Rect.A).Div(grid.Size())
			angle, ok := vg.Interpolate(p)
			if !ok {
				continue
			}

			vec := Unit().Rotate(angle)
			vec = vec.SetLength(0.008)

			xy = xy.Add(vec)
			//xy = xy.Clamp(XY{}, XY{1, 1})

			pos = append(pos, xy)
			colors = append(colors, c)
		}
	}

	shader := gfx.Fill{
		Shape: Circle{
			Radius:   0.002,
			Segments: 5,
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
}

type VectorGrid struct {
	// TODO used angles because they're easy to interpolate
	Angles []float32

	grid Grid
}

func NewVectorGrid(g Grid) VectorGrid {
	return VectorGrid{
		Angles: make([]float32, (g.Rows+1)*(g.Cols+1)),
		grid:   g,
	}
}

// TODO note "xy" is a percentage, range [0, 1]
func (vg VectorGrid) Interpolate(percent XY) (float32, bool) {
	// TODO should probably clamp to 0-1?
	row := int(percent.Y * float32(vg.grid.Rows))
	col := int(percent.X * float32(vg.grid.Cols))

	// TODO clamp row/col index to grid size
	ai := row*vg.grid.Cols + col
	bi := row*vg.grid.Cols + (col + 1)
	ci := (row+1)*vg.grid.Cols + (col + 1)
	di := (row+1)*vg.grid.Cols + col
	maxi := len(vg.Angles) - 1

	// TODO hack. need better bounds logic.
	if ai < 0 || ai > maxi {
		return 0, false
	}
	if bi < 0 || bi > maxi {
		return 0, false
	}
	if ci < 0 || ci > maxi {
		return 0, false
	}
	if di < 0 || di > maxi {
		return 0, false
	}

	a := vg.Angles[ai]
	b := vg.Angles[bi]
	c := vg.Angles[ci]
	d := vg.Angles[di]

	xy := vg.grid.Rect.Interpolate(percent)
	cell := vg.grid.Cell(row, col)
	off := xy.Sub(cell.A)
	p := off.Div(cell.Size())

	e := math.Interp(a, b, p.X)
	f := math.Interp(d, c, p.X)
	return math.Interp(e, f, p.Y), true
}

func (vg VectorGrid) DrawArrows(doc *gfx.Doc) {

	for _, cell := range vg.grid.Cells() {
		center := cell.Center()

		// pick a random angle
		//angle := rand.Angle()

		// figure out how long each vector arrow should be.
		// it should be half the size of a cell at most:
		size := cell.Size().Length() * 0.5
		// and then a bit smaller, just to add some padding:
		size *= 0.25

		// Determine the vector at this position.
		p := center.Sub(vg.grid.Rect.A).Div(vg.grid.Size())
		angle, _ := vg.Interpolate(p)
		vec := Unit().Rotate(angle)

		// create the vector arrow
		line := Line{
			A: center,
			B: center.Add(vec.MulScalar(size)),
		}

		// draw the vector arrow
		gfx.Stroke{
			Shape: line,
			Color: color.Black,
		}.Draw(doc)

		gfx.Dot{
			XY: center,
		}.Draw(doc)
	}
}
