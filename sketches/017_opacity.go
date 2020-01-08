package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
)

func Ink(doc Layer) {
	Clear(doc, White)

	fills := []Fill{
		{
			Rect{XY{0.3, 0.3}, XY{0.6, 0.6}},
			RGBA{1, 0, 0, 0.5},
		},
		{
			Rect{XY{0.4, 0.4}, XY{0.7, 0.7}},
			RGBA{1, 0, 0, 0.5},
		},

		{
			Rect{XY{0.4, 0.4}, XY{0.5, 0.5}},
			RGBA{0, 0, 0, 0},
		},
		{
			Rect{XY{0.2, 0.2}, XY{0.4, 0.4}},
			RGBA{0, 1, 0, 1},
		},
		{
			Rect{XY{0.2, 0.4}, XY{0.4, 0.6}},
			RGBA{0, 0, 1, 0.5},
		},
		{
			Rect{XY{0.1, 0.4}, XY{0.2, 0.6}},
			RGBA{1, 0, 0, 1},
		},
		{
			Rect{XY{0.6, 0.4}, XY{0.8, 0.6}},
			RGBA{1, 1, 0, 0.5},
		},
	}

	for _, f := range fills {
		f.Draw(doc)
	}
}
