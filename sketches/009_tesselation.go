package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/tess"
)

func Ink(doc gfx.Doc) {

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
	gfx.Fill{Shape: m, Color: Black}.Draw(doc)

	for _, xy := range xys {
		d := gfx.Dot{XY: xy, Color: Red, Radius: 0.005}
		d.Draw(doc)
	}
}
