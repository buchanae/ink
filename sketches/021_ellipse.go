package main

import (
	"github.com/buchanae/ink/app"
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
)

func Ink(doc *app.Doc) {
	// TODO maybe this should be done automatically,
	//      otherwise textures can accumulate pixels from
	//      multiple renders
	Clear(doc, color.White)

	const N = 100

	e := Ellipse{
		XY:       XY{.5, .5},
		Size:     XY{.3, .2},
		Segments: 100,
	}

	Fill{e, color.Blue}.Draw(doc)
	s := e.Stroke()
	s.Width = 0.005
	Fill{s, color.Black}.Draw(doc)

	/*
		for i := 0; i < N; i++ {
			xy := e.Interpolate(float32(i) / N)
			Dot{XY: xy, Radius: 0.003}.Draw(doc)
		}
	*/
}
