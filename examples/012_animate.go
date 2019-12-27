package main

import (
	"log"

	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
)

func Ink(z *Doc) {

	// TODO add helper for equilateral
	t := Triangle{
		XY{0.2, 0.3},
		XY{0.8, 0.3},
		XY{0.5, 0.8},
	}

	m := z.Shader(t)

	m.Set("a_color", []RGBA{
		Red, Green, Blue,
	})
	m.Set("a_pivot", t.Centroid())

	z.Dot(t.Centroid(), White)

	var rot float32
	z.OnFrame = func(f Frame) {
		log.Println("frame")
		m.Set("a_rot", rot)
		rot += 0.01
	}
}
