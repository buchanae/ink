package main

import (
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

const (
	N       = 1000
	Scale   = 1.5
	Octaves = 40
	Shift   = 0.55
)

func Ink(doc gfx.Doc) {

	for i := 0; i < N; i++ {
		x := float32(i) / N

		h := octaves(x, 6)
		h -= 0.7
		h *= .2

		xy := XY{x, 0.5 + h}
		c := Circle{xy, 0.001, 10}
		s := gfx.NewShader(c.Fill())
		s.Draw(doc)
	}
}

func octaves(x float32, N int) float32 {
	var n float32
	var z float32 = 10
	var amp float32 = 1

	for j := 0; j < N; j++ {
		n += rand.Noise1(x*z) * amp
		amp *= 0.5
		z = z * 2
	}
	return n
}
