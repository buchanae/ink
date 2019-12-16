package main

import (
	"time"

	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

const (
	N        = 1050
	Passes   = 10
	Padding  = 0.000
	MinRange = 0.005
	MaxRange = 0.1
	RangeInc = 0.005
)

func main() {
	doc := NewDoc()

	doc.Draw(Clear(White))

	r := rand.New(time.Now().Unix())
	p := r.Palette()

	for j := 0; j < Passes; j++ {

		h := r.Range(MinRange, MaxRange)
		l := r.Range(0.1, 0.9)
		c := r.Color(p)

		for i := 0; i < N; i++ {
			x := float32(i) / N
			w := float32(1) / N

			h += r.Range(-RangeInc, RangeInc)
			if h > MaxRange {
				h = MaxRange
			}
			if h < MinRange {
				h = MinRange
			}

			r := Rect{
				A: XY{x + Padding, l - h},
				B: XY{x + w - Padding, l + h},
			}

			alpha := 1 - h/MaxRange
			s := NewShader(r.Mesh())
			s.SetColor(RGBA{c.R, c.G, c.B, alpha})
			doc.Draw(s)
		}
	}

	app.Render(doc)
}
