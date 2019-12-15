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

	m := NewShader(r.Mesh())
	m.Set("a_center", r.Center())
	m.Set("a_rot", 0.4)
	m.SetColor(Red)

	doc.Draw(m)
	app.Render(doc)
}
