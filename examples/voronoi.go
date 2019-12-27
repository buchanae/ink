package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
	"github.com/buchanae/ink/voronoi"
)

func Ink(doc *Doc) {
	rand.SeedNow()
	doc.Clear(White)

	box := Rect{
		A: XY{.1, .1},
		B: XY{.9, .9},
	}

	var xys []XY

	var initial []XY

	for i := 0; i < 20; i++ {
		p := float32(i) / 20
		xy := box.Interp(XY{p, p})
		initial = append(initial, xy)
	}

	noise := rand.BlueNoiseInitial(450, 1, 1, 0.050, initial)

	for _, xy := range noise {
		if !box.Contains(xy) {
			continue
		}
		xys = append(xys, xy)
	}

	//m := box.Stroke(0.001)
	//doc.Shader(m)

	v := voronoi.New(xys, box)

	colors := []RGBA{
		Blue, Yellow, Green, Black, Purple,
	}

	/*
		for _, cell := range v.Cells() {
			c := rand.Color(colors)
			c.A = 0.3
			for _, tri := range cell.Tris {
				s := doc.Shader(tri)
				s.Set("a_color", c)
			}

			for _, e := range cell.Edges {
				m := e.Stroke(0.002)
				s := doc.Shader(m)
			}
		}
	*/

	tris := v.Triangulate()

	for _, t := range tris {
		s := doc.Shader(t)
		c := rand.Color(colors)
		c.A = 0.3
		s.Set("a_color", c)
	}

	{
		m := Triangles(tris)
		stk := m.Stroke(0.001)
		s := doc.Shader(stk)
		s.Set("a_color", White)
	}

	/*
			for _, xy := range xys {
				doc.Dot(xy, Red)
			}

		for _, xy := range initial {
			doc.Dot(xy, Blue)
		}
	*/
}
