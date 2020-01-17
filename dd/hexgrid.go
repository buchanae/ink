package dd

import (
	"math"
)

// references
// https://www.redblobgames.com/grids/hexagons/

type HexGrid struct {
	Size int
}

type HexCell struct {
	Center   XY
	Radius   float32
	Row, Col int
}

func (hc HexCell) Mesh() Mesh {
	c := hc.Center
	rad := hc.Radius
	rot := float32((1.0 / 6.0) * math.Pi)
	opp := rad * sin(rot)
	adj := rad * cos(rot)

	return Mesh{
		Verts: []XY{
			c,
			{c.X + adj, c.Y + opp},
			{c.X, c.Y + rad},
			{c.X - adj, c.Y + opp},
			{c.X - adj, c.Y - opp},
			{c.X, c.Y - rad},
			{c.X + adj, c.Y - opp},
		},
		Faces: []Face{
			{0, 6, 1},
			{0, 1, 2},
			{0, 2, 3},
			{0, 3, 4},
			{0, 4, 5},
			{0, 5, 6},
		},
	}
}

func (hc HexCell) Stroke(opt StrokeOpt) Mesh {
	c := hc.Center
	rad := hc.Radius
	rot := float32((1.0 / 6.0) * math.Pi)
	opp := rad * sin(rot)
	adj := rad * cos(rot)
	verts := []XY{
		{c.X + adj, c.Y + opp},
		{c.X, c.Y + rad},
		{c.X - adj, c.Y + opp},
		{c.X - adj, c.Y - opp},
		{c.X, c.Y - rad},
		{c.X + adj, c.Y - opp},
	}
	path := Path{
		Line{verts[0], verts[1]},
		Line{verts[1], verts[2]},
		Line{verts[2], verts[3]},
		Line{verts[3], verts[4]},
		Line{verts[4], verts[5]},
		Line{verts[5], verts[0]},
	}
	return Stroke(path, opt)
}

func (g HexGrid) Cells() []HexCell {
	cells := make([]HexCell, 0, g.Size*2)

	w := 1 / float32(g.Size-2)
	radius := w / sqrt(3)
	vspace := radius * 1.5

	for r := 0; r < g.Size; r++ {
		for c := 0; c < g.Size; c++ {

			xy := XY{
				X: float32(c) * w,
				Y: float32(r) * vspace,
			}
			if r%2 != 0 {
				xy.X += w / 2
			}

			cells = append(cells, HexCell{
				Row:    r,
				Col:    c,
				Center: xy,
				Radius: radius,
			})
		}
	}
	return cells
}
