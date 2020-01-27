package main

import (
	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
	"github.com/buchanae/ink/voronoi"
)

func Ink(doc *app.Doc) {
	rand.SeedNow()

	box := Rect{
		A: XY{.1, .1},
		B: XY{.9, .9},
	}

	var initial []XY

	for i := 0; i < 20; i++ {
		p := float32(i) / 20
		xy := box.Interpolate(XY{p, p})
		initial = append(initial, xy)
	}

	noise := rand.BlueNoise{
		Limit:   450,
		Spacing: 0.05,
		Initial: initial,
		Rect:    box,
	}.Generate()

	//m := box.Stroke(0.001)
	//doc.Shader(m)

	v := voronoi.New(noise, box)

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
		c := rand.Color(colors)
		c.A = 0.3
		gfx.Fill{
			Mesh:  t,
			Color: c,
		}.Draw(doc)
	}

	gfx.Stroke{
		Target: Triangles(tris),
		Width:  0.001,
		Color:  Black,
	}.Draw(doc)
}
