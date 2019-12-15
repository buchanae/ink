package main

import (
	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

const (
	GridSize       = 15
	Boxes          = 5
	LineWidth      = 0.0015
	SkipCellChance = 0.05
	SkipBoxChance  = 0.20
	Shrink         = 0.0050
	TweakBox       = 0.0017
	Angle          = 0.2
)

func main() {
	doc := NewDoc()
	grid := NewGrid(GridSize, GridSize)

	for _, cell := range grid.Rects() {
		if rand.Bool(SkipCellChance) {
			continue
		}

		for i := 0; i < Boxes; i++ {
			if rand.Bool(SkipBoxChance) {
				continue
			}

			s := cell.Shrink(float32(i+1) * Shrink)
			q := QuadFromRect(s)
			q = rand.TweakQuad(q, TweakBox)
			m := q.Stroke(LineWidth)
			z := NewShader(m)
			z.SetColor(Red)
			z.Set("a_center", s.Center())
			z.Set("a_rot", rand.Range(-Angle, Angle))

			doc.Draw(z)
		}
	}

	app.Render(doc)
}
