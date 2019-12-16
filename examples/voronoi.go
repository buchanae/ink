package main

import (
	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
	"github.com/buchanae/ink/tess"
)

func main() {
	doc := NewDoc()

	bg := Fill{Fullscreen(), White}
	doc.Draw(bg)

	xys := rand.BlueNoise(150, 1, 1, 0.05)

	for _, xy := range xys {
		// TODO four lines just to draw a small point is too much
		// doc.Circle(c)
		c := Circle{xy, 0.005}
		s := NewShader(c.Mesh(10))
		s.SetColor(Red)
		doc.Draw(s)
	}

	tris := tess.Tesselate(xys)
	m := StrokeTriangles(tris, 0.005)
	s := NewShader(m)
	doc.Draw(s)

	app.Render(doc)
}
