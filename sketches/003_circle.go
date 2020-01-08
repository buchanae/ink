package main

import (
	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
)

func Ink(doc *app.Doc) {
	Fill{
		Mesh: Circle{
			XY:     XY{0.5, 0.5},
			Radius: 0.2,
		},
		Color: Blue,
	}.Draw(doc)
}
