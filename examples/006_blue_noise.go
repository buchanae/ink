package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(z *Layer) {

	xys := rand.BlueNoise(500, 1, 1, 0.02)

	for _, xy := range xys {
		z.Dot(xy, Red)
	}
}
