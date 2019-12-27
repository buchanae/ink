package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/tess"
)

func Ink(doc *Doc) {
	doc.Clear(White)

	xys := []XY{
		{0.2, 0.2},
		{0.2, 0.6},
		{0.4, 0.7},
		{0.9, 0.7},
		{0.3, 0.5},
		{0.5, 0.4},
		{0.4, 0.3},
	}

	tris := tess.Tesselate(xys)
	m := Triangles(tris)
	doc.Shader(m)

	for _, xy := range xys {
		doc.Dot(xy, Red)
	}
}
