package main

import (
	"log"
	"time"

	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/math"
	"github.com/buchanae/ink/rand"
)

const (
	// Number of _potential_ horizontal strokes
	Lines = 20

	// Actual horizontal strokes.
	// If this is greater than Lines, then the same line will
	// get multiple passes.
	// TODO maybe overcomplicated
	Strokes = 20
	N       = 5000
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

func random_float32s_in_range(r *rand.Rand, count int, min float32, max float32) []float32 {
	floats := make([]float32, count)
	for i := range floats {
		floats[i] = r.Range(min, max)
	}
	return floats
}

func Ink(doc gfx.Doc) {
	r := rand.New(time.Now().Unix())
	start := time.Now()
	palette := r.Palette()

	// random y-axis positions
	y_positions := random_float32s_in_range(r, Lines, MinY, MaxY)

	for j := 0; j < Strokes; j++ {

		y := y_positions[rand.Intn(len(y_positions))]
		h := rand.Range(0.01, 0.1)
		color := rand.Color(palette)

		for i := 0; i < N; i++ {
			x := float32(i) / N

			h += rand.Range(-D, D)
			h = math.Clamp(h, 0, MaxD)

			xy := XY{x, y}
			wh := XY{W, h}
			r := RectCenter(xy, wh)

			s := gfx.NewShader(r.Fill())
			sc := color
			sc.A = 1 - h/B - M
			s.Set("a_color", sc)
			s.Draw(doc)
		}
	}

	log.Printf("run time: %s", time.Since(start))
}
