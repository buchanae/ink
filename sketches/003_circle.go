package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
)

func Ink(doc Layer) {
	Fill{
		Mesh: Circle{
			XY:       XY{0.5, 0.5},
			Radius:   0.2,
			Segments: 40,
		},
		Color: Blue,
	}.Draw(doc)
}
