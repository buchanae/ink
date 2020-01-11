package main

import (
	"github.com/buchanae/ink/app"
)

func Ink(doc *app.Doc) {
	img := doc.LoadImage("toshiro.png")
	img.Draw(doc)
}
