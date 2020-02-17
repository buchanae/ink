package main

import (
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(doc gfx.Doc) {

	r := Rect{
		A: XY{.1, .1},
		B: XY{.9, .9},
	}
	gfx.Fill{r, color.Lightgray}.Draw(doc)

	bn := rand.BlueNoise{
		Limit: 5050,
		Rect:  r,
	}
	xys := bn.Generate()

	for _, xy := range xys {
		gfx.Dot{XY: xy, Radius: 0.003}.Draw(doc)
	}
}
