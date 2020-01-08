package main

import (
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(doc *Doc) {
	Clear(doc, color.White)

	c := dd.Circle{
		XY:       dd.XY{0.5, 0.5},
		Radius:   0.3,
		Segments: 40,
	}
	m := c.Mesh()
	m = rand.TweakMesh(m, 0.03)

	Fill{m, color.Red}.Draw(doc)
}
