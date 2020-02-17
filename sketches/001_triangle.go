package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
)

func Ink(doc *gfx.Doc) {

	// The "hello, world" of graphics:
	// a triangle with different colored vertices.
	t := Triangle{
		XY{0.2, 0.2},
		XY{0.8, 0.2},
		XY{0.5, 0.8},
	}

	s := gfx.Fill{Shape: t}.Shader()
	s.Set("a_color", []RGBA{
		Red, Green, Blue,
	})
	s.Draw(doc)
}
