package main

import (
	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
	"github.com/buchanae/ink/voronoi"
)

const (
	N       = 10
	Padding = 0.01
	Margin  = Padding * 2
)

// TODO idea
// use MeshShade to shade an image of a face
// using image brightness to determine point
// density. kinda like pointilism, but connecting
// points into mesh

/*
TODO ideas:
- voronoi mesh of blue noise. draw lines of mesh edges.
  should create a mesh/web texture
- rand.BlueNoise needs to be cleaned up.
	- easily generate points in rect/shape
	- don't require N, just generate until full
- want easier interpolation, padding, etc.
- want easier variable curves (e.g. x^2)
  always find it unintuitive to define a model
	where a variable grows how I want
*/

func Ink(doc *app.Doc) {
	rand.SeedNow()

	for i := float32(0); i < N; i++ {
		p := i / N

		split := rand.Range(.3, .7)

		bot := i/N + (Padding / 2)
		top := (i+1)/N - (Padding / 2)

		ra := Rect{
			A: XY{Margin, bot},
			B: XY{split - (Padding / 2), top},
		}

		rb := Rect{
			A: XY{split + (Padding / 2), bot},
			B: XY{1 - Margin, top},
		}

		c := RGBA{p, p, p, 1}
		MeshShade(doc, ra, c, i)

		//gfx.Fill{ra, c}.Draw(doc)
		gfx.Fill{rb, c}.Draw(doc)
	}

}

// TODO be able to turn on tracing from within the doc/sketch

// TODO voronoi should definitely have a Mesh() method

// TODO this is a good example of a sketch where memoization
//      (smarter caching) would help. this takes time to generate
//      millions of blue noise points. would be nice to skip that
//      if all we need is to change the color

func MeshShade(doc gfx.Layer, r Rect, c RGBA, i float32) {
	space := Interp(0.002, 0.02, i/N)
	na := BlueNoiseInBox(100000, space, r)
	// TODO want voronoi without edges
	v := voronoi.New(na, r)

	seen := map[Line]struct{}{}

	for _, t := range v.Triangulate() {
		for _, e := range t.Edges() {

			if e.A.X < e.B.X || e.A.Y < e.B.Y {
				e.A, e.B = e.B, e.A
				// TODO cool mistake
				//e.A = e.B
			}
			if _, ok := seen[e]; ok {
				continue
			}
			seen[e] = struct{}{}

			gfx.Stroke{
				Target: e,
				Width:  0.0005,
				Color:  c,
			}.Draw(doc)
		}
	}
}

func CellShade(doc gfx.Layer, r Rect, c RGBA, i float32) {
	space := Interp(0.002, 0.01, i/N)
	na := BlueNoiseInBox(50000, space, r)
	// TODO want voronoi without edges
	v := voronoi.New(na, r)
	for _, e := range v.Edges() {

		gfx.Stroke{
			Target: e,
			Color:  c,
		}.Draw(doc)
	}
}

func DotShade(doc gfx.Layer, r Rect, c RGBA, i float32) {
	space := Interp(0.002, 0.009, i/N)
	na := BlueNoiseInBox(100000, space, r)
	for _, xy := range na {
		gfx.Dot{xy, c, 0.001}.Draw(doc)
	}
}

func Interp(from, to, p float32) float32 {
	return (to-from)*p + from
}

func BlueNoiseInBox(n int, space float32, box Rect) []XY {
	bn := rand.BlueNoise{
		Limit:   n,
		Rect:    box,
		Spacing: space,
	}
	return bn.Generate()
}
