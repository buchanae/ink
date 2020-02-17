package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/math"
	"github.com/buchanae/ink/rand"
	"github.com/buchanae/ink/voronoi"
)

const (
	N       = 10
	Padding = 0.01
	Margin  = Padding * 2
)

func Ink(doc gfx.Doc) {
	rand.SeedNow()

	for i := float32(0); i < N; i++ {
		split := rand.Range(.3, .7)
		bot := i/N + (Padding / 2)
		top := (i+1)/N - (Padding / 2)

		p := i / N
		c := RGBA{p, p, p, 1}

		// left
		ra := Rect{
			A: XY{Margin, bot},
			B: XY{split - (Padding / 2), top},
		}
		// right
		rb := Rect{
			A: XY{split + (Padding / 2), bot},
			B: XY{1 - Margin, top},
		}

		ca := VoronoiCells{
			Rect:    ra,
			Spacing: math.Interp(0.003, 0.03, i/N),
		}

		gfx.Fill{
			Shape: ca.Mesh(),
			Color: c,
		}.Draw(doc)

		gfx.Fill{
			Shape: rb,
			Color: c,
		}.Draw(doc)
	}
}

type VoronoiCells struct {
	Rect
	Spacing float32
}

func (vc VoronoiCells) Mesh() Mesh {
	bn := rand.BlueNoise{
		Rect:    vc.Rect,
		Spacing: vc.Spacing,
	}
	noise := bn.Generate()
	v := voronoi.New(noise, vc.Rect)

	var meshes []Mesh
	for _, e := range v.Edges() {
		meshes = append(meshes, e.Stroke(StrokeOpt{}))
	}
	return Merge(meshes...)
}
