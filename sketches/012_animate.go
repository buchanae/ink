package main

import (
	"math"

	"github.com/buchanae/ink/app"
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
)

func Ink(doc *app.Doc) {

	a := 0.6
	h := (math.Sqrt(3) / 2) * a

	// TODO add helper for equilateral
	t := Triangle{
		XY{0.2, 0.3},
		XY{0.8, 0.3},
		XY{0.5, float32(0.3 + h)},
	}
	center := XY{0.5, float32(h/3) + 0.3}

	m := gfx.Fill{Shape: t}.Shader()
	m.Set("a_color", []color.RGBA{
		color.Red, color.Green, color.Blue,
	})
	m.Set("a_pivot", center)
	m.Draw(doc)

	var rot float32
	for app.Play() {
		// TODO want animation system to handle clearing
		//      for you. But also like the ability to easily
		//      accumulate changes over frames by not clearing.
		doc.Ops = nil
		//gfx.Clear(doc, color.Black)

		m.Set("a_rot", rot)
		m.Draw(doc)

		rot += 0.01
		app.Send(doc)
	}
}
