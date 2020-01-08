package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

const (
	GridSize       = 20
	Boxes          = 7
	LineWidth      = 0.0015
	SkipCellChance = 0.05
	SkipBoxChance  = 0.20
	Shrink         = 0.0050
	TweakBox       = 0.0017
	Angle          = 0.2
)

func Ink(doc *Doc) {
	rand.SeedNow()
	Clear(doc, White)
	grid := NewGrid(GridSize, GridSize)

	p := rand.Palette()

	for i, rect := range grid.Rects() {
		cell := grid.Cells[i]

		if rand.Bool(SkipCellChance) {
			continue
		}

		for i := 0; i < Boxes; i++ {
			if rand.Bool(SkipBoxChance) {
				continue
			}

			r := rect.Shrink(float32(i+1) * Shrink)
			q := r.Quad()
			//q = rand.TweakQuad(q, TweakBox)
			m := q.Stroke()
			m.Width = LineWidth

			s := NewShader(m)
			s.Set("a_color", rand.Color(p))
			s.Set("a_pivot", r.Center())
			s.Draw(doc)

			if cell.Row > 5 {
				s.Set("a_rot", rand.Range(-Angle, Angle))
			}
		}
	}
}
