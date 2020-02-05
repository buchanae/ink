package main

import (
	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
)

func Ink(doc *app.Doc) {
	doc.Config.Window.Width = 300
	doc.Config.Window.Height = 300

	t := Triangle{
		XY{0.2, 0.2},
		XY{0.8, 0.2},
		XY{0.5, 0.8},
	}

	s := gfx.NewShader(t.Fill())
	s.Set("a_color", []RGBA{
		Red, Green, Blue,
	})
	s.Draw(doc)
}
