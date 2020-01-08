package main

import (
	"github.com/buchanae/ink/app"
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/math"
)

func Ink(doc *app.Doc) {
	// TODO maybe this should be done automatically,
	//      otherwise textures can accumulate pixels from
	//      multiple renders
	Clear(doc, color.White)

	const N = 100
	inc := float32((2 * math.Pi) / N)
	t := float32(0)

	e := Ellipse{
		XY:   XY{.5, .5},
		Size: XY{.3, .2},
	}

	Fill{e, color.Blue}.Draw(doc)

	for i := 0; i < N; i++ {
		xy := e.Interpolate(t)
		Dot{XY: xy, Radius: 0.003}.Draw(doc)
		t += inc
	}
}
