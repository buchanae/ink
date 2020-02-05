package main

import (
	"log"
	"time"

	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/math"
	"github.com/buchanae/ink/rand"
)

const (
	Lines   = 20
	Strokes = 20
	N       = 1000
	// percentage of W
	// only works for small N
	Padding = 0.00
	W       = (1 / float32(N)) * (1 - Padding)

	// min/max y-axis starting position
	MinY = 0.1
	MaxY = 0.9

	D    = 0.004
	B    = 0.3
	MaxD = 0.2
	M    = 0.4
)

func Ink(doc *app.Doc) {
	rand.SeedNow()
	palette := rand.Palette()

	start := time.Now()
	ys := make([]float32, Lines)
	for i := range ys {
		ys[i] = rand.Range(MinY, MaxY)
	}

	for j := 0; j < Strokes; j++ {

		y := ys[rand.Intn(len(ys))]
		dy := rand.Range(0.01, 0.1)
		color := rand.Color(palette)

		for i := 0; i < N; i++ {
			x := float32(i) / N

			dy += rand.Range(-D, D)
			dy = math.Clamp(dy, 0, MaxD)

			xy := XY{x, y}
			wh := XY{W, dy}
			r := RectCenter(xy, wh)

			s := gfx.NewShader(r.Fill())
			sc := color
			sc.A = 1 - dy/B - M
			s.Set("a_color", sc)
			s.Draw(doc)
		}
	}

	log.Printf("run time: %s", time.Since(start))
}
