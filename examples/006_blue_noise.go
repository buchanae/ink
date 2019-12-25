package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(z *Layer) {

	xys := rand.BlueNoise(500, 1, 1, 0.01)

	for _, xy := range xys {
		c := Dot(xy)
		s := z.Shader(c)
		s.Set("a_color", Red)
	}
}
