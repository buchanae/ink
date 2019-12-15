package main

import (
	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func main() {
	doc := NewDoc()

	for _, xy := range rand.BlueNoise(100, 1, 1, 0.01) {
		c := Circle{xy, 0.005}
		s := NewShader(c.Mesh(10))
		s.SetColor(Red)
		doc.Draw(s)
	}

	app.Render(doc)
}
