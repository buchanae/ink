package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
)

func Ink(doc Layer) {
	// TODO blur looks different on black background. why is that?
	Clear(doc, White)

	Fill{
		Mesh: Rect{
			XY{0.2, 0.2},
			XY{0.8, 0.8},
		},
		Color: Blue,
	}.Draw(doc)

	Blur{
		Passes: 2,
		Source: doc,
	}.Draw(doc)
}
