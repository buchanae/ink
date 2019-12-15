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

	colors := []RGBA{
		Red, Blue, Green, Yellow,
	}

	grid := NewGrid(20, 20)
	for _, r := range grid.Rects() {
		m := NewShader(r.Mesh())
		c := rand.Color(colors)
		m.SetColor(c)
		doc.Draw(m)
	}

	app.Render(doc)
}
