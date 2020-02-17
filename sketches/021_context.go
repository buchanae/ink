package main

import (
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
)

func Ink(doc gfx.Doc) {
	ctx := gfx.NewContext(doc)

	ctx.Clear(color.White)

	e := Ellipse{
		XY:       XY{.5, .5},
		Size:     XY{.3, .2},
		Segments: 100,
	}

	ctx.FillColor = color.Blue
	ctx.Fill(e)

	ctx.StrokeWidth = 0.005
	ctx.Stroke(e)
}
