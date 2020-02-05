package main

import (
	"github.com/buchanae/ink/app"
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
)

func Ink(doc *app.Doc) {

	t := Triangle{
		XY{0.2, 0.2},
		XY{0.8, 0.2},
		XY{0.5, 0.8},
	}

	red := color.Hex(0xff0000)
	green := color.Hex(0x00ff00)
	blue := color.HexString("#0000ff")

	s := gfx.NewShader(t.Fill())
	s.Set("a_color", []color.RGBA{
		red, green, blue,
	})
	s.Draw(doc)
}
