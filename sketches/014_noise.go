package main

import (
	"github.com/buchanae/ink/app"
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/gfx"
)

func Ink(doc *app.Doc) {
	n := gfx.DefaultNoise
	n.Size = 30
	n.Color = color.Red
	n.Draw(doc)
}
