package main

import (
	"github.com/buchanae/ink/app"
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
)

func Ink(doc *app.Doc) {

	shapes := []gfx.Meshable{

		Rect{
			XY{0.1, 0.7},
			XY{0.3, 0.9},
		},

		Circle{
			XY:     XY{0.5, 0.5},
			Radius: 0.1,
		},

		Triangle{
			XY{.7, .1},
			XY{.8, .3},
			XY{.9, .1},
		},

		Ellipse{
			XY:   XY{.2, .2},
			Size: XY{.15, .1},
		},

		Quad{
			XY{0.65, 0.7},
			XY{0.9, 0.7},
			XY{0.85, 0.95},
			XY{0.7, 0.9},
		},
	}

	for _, s := range shapes {
		gfx.Fill{s, color.Blue}.Draw(doc)
	}
}
