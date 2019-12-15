package main

import (
	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
)

func main() {
	doc := NewDoc()

	t := Triangle{
		XY{0.2, 0.2},
		XY{0.8, 0.2},
		XY{0.5, 0.8},
	}

	m := NewShader(t.Mesh())
	m.Set("a_color", []RGBA{
		Red, Green, Blue,
	})

	doc.Draw(m)
	app.Render(doc)
}
