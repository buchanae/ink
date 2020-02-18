package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
)

func Ink(doc gfx.Doc) {

	// The "hello, world" of graphics:
	// a triangle with different colored vertices.
	for i := 0; i < 100; i++ {
		t := Triangle{
			XY{0.2, 0.2},
			XY{0.8, 0.2},
			XY{0.5, 0.8},
		}

		s := gfx.Fill{Shape: t, Color: Red}.Shader()
		s.Draw(doc)
	}
}
