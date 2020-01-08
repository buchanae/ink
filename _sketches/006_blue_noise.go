package main

import (
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(doc *Doc) {
	Clear(doc, color.White)

	xys := rand.BlueNoise(250, 1, 1, 0.02)

	for _, xy := range xys {
		Dot{XY: xy, Radius: 0.003}.Draw(doc)
	}
}
