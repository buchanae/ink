package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(doc gfx.Doc) {

	for i := 0; i < 100; i++ {
		gfx.Fill{
			Shape: Circle{
				XY:       rand.XY(),
				Radius:   0.01,
				Segments: 5,
			},
			Color: Red,
		}.Draw(doc)
	}
}
