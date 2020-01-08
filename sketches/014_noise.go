package main

import (
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/gfx"
)

func Ink(doc Layer) {
	Clear(doc, color.White)
	n := DefaultNoise
	n.Size = 30
	n.Color = color.Red
	n.Draw(doc)
}
