package main

import (
	"github.com/buchanae/ink/app"
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(doc *app.Doc) {
	Clear(doc, color.White)

	r := dd.Rect{
		A: dd.XY{.1, .1},
		B: dd.XY{.9, .9},
	}
	Fill{r, color.Gray}.Draw(doc)

	bn := rand.BlueNoise{
		Limit: 5050,
		Rect:  r,
	}
	xys := bn.Generate(rand.Default)

	for _, xy := range xys {
		Dot{XY: xy, Radius: 0.003}.Draw(doc)
	}
}
