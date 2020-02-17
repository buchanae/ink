package main

import (
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(doc gfx.Doc) {

	c := Circle{
		XY:       XY{0.5, 0.5},
		Radius:   0.3,
		Segments: 40,
	}
	m := c.Fill()
	m = rand.TweakMesh(m, 0.03)

	gfx.Fill{m, color.Red}.Draw(doc)
}
