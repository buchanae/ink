package main

import (
	"time"

	"github.com/buchanae/ink/app/client"
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

const (
	GridSize       = 20
	Boxes          = 7
	LineWidth      = 0.0015
	SkipCellChance = 0.05
	SkipBoxChance  = 0.20
	Shrink         = 0.0050
	Angle          = 0.2
)

func Ink(idoc gfx.Doc) {
	doc := idoc.(*client.Doc)

	for {
		doc.Clear()
		gfx.Clear(doc, color.White)
		rand.SeedNow()

		grid := Grid{Rows: GridSize, Cols: GridSize}

		p := rand.Palette()

		for _, cell := range grid.Cells() {
			rect := cell.Rect

			if rand.Bool(SkipCellChance) {
				continue
			}

			for i := 0; i < Boxes; i++ {
				if rand.Bool(SkipBoxChance) {
					continue
				}

				r := rect.Shrink(float32(i+1) * Shrink)
				q := r.Quad()
				m := q.Stroke(StrokeOpt{
					Width: LineWidth,
				})

				s := gfx.NewShader(m.Fill())
				s.Set("a_color", rand.Color(p))
				s.Set("a_pivot", r.Center())
				s.Draw(doc)

				if cell.Row > 5 {
					s.Set("a_rot", rand.Range(-Angle, Angle))
				}
			}

		}
		client.Send(doc)
		time.Sleep(1 * time.Second)
	}
}
