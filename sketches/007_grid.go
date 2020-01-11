package main

import (
	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(doc *app.Doc) {

	grid := Grid{Rows: 20, Cols: 20}
	palette := rand.Palette()

	for _, cell := range grid.Cells() {
		r := cell.Rect.Shrink(0.003)
		c := rand.Color(palette)
		gfx.Fill{r, c}.Draw(doc)
	}
}
