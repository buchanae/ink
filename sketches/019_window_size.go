package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
)

func Ink(doc gfx.Doc) {

	conf := doc.Config()
	conf.Title = "example: set window size"
	conf.Width = 300
	conf.Height = 300
	conf.Snapshot.Width = 500
	conf.Snapshot.Height = 500

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
