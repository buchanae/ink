package main

import (
	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
)

func main() {
	doc := NewDoc()

	r := Rect{
		XY{0.2, 0.2},
		XY{0.8, 0.8},
	}

	doc.Draw(Fill{
		r.Mesh(),
		Red,
	})

	app.Render(doc)
}
