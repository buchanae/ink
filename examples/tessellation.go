package main

import (
	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/tess"
)

func main() {
	doc := NewDoc()

	bg := Fill{Fullscreen(), White}
	doc.Draw(bg)

	xys := []XY{
		{0.2, 0.2},
		{0.2, 0.6},
		{0.4, 0.7},
		{0.9, 0.7},
		{0.3, 0.5},
		{0.5, 0.4},
		{0.4, 0.3},
	}

	tris := tess.Tesselate(xys)
	m := Triangles(tris)
	s := NewShader(m)
	doc.Draw(s)

	for _, xy := range xys {
		// TODO four lines just to draw a small point is too much
		// doc.Circle(c)
		c := Circle{xy, 0.005}
		s := NewShader(c.Mesh(10))
		s.SetColor(Red)
		doc.Draw(s)
	}

	app.Render(doc)
}
