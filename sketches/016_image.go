package main

import "github.com/buchanae/ink/gfx"

func Ink(doc gfx.Doc) {
	img := doc.LoadImage("toshiro.png")
	img.Draw(doc)
}
