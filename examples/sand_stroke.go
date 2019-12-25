package main

import (
	"time"

	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/math"
	"github.com/buchanae/ink/rand"
)

const (
	N       = 850
	Strokes = 1
	Padding = 0.000

	MinRange = 0.001
	MaxRange = 0.3
	RangeInc = 0.005

	MaxAlpha = 1
	MinAlpha = 0

	MinY = 0.1
	MaxY = 0.9
)

func main() {
	doc := NewDoc()

	doc.Draw(Clear(White))

	r := rand.New(time.Now().Unix())
	r = rand.New(6)

	p := r.Palette()

	for j := 0; j < Strokes; j++ {

		l := r.Range(MinY, MaxY)
		l = 0.5
		c := r.Color(p)
		h := float32(0)
		inc := float32(0)

		for i := 0; i < N; i++ {
			x := float32(i) / N
			w := float32(1) / N

			if i%10 == 0 {
				inc = r.Range(-RangeInc, RangeInc)
			}

			h += inc
			h = math.Clamp(h, MinRange, MaxRange)

			r := Rect{
				A: XY{x + Padding, l - h},
				B: XY{x + w - Padding, l + h},
			}

			//	alpha := Interp(h/MaxRange, MaxAlpha, MinAlpha)
			//log.Println(h, h2, alpha, h2/0.5)
			//alpha = h2 / 0.5
			//alpha := float32(1.0)
			alpha := float32(1)

			s := NewShader(r.Mesh())
			s.SetColor(RGBA{c.R, c.G, c.B, alpha})
			doc.Draw(s)
		}
	}

	app.Render(doc)
}

func octaves(x, y float32, N int) float32 {
	var n float32
	var z float32 = 50
	var amp float32 = 1

	for j := 0; j < N; j++ {
		n += rand.Noise2(x*z, y) * amp
		amp *= 0.010
		z = z * 2
	}
	return n
}
