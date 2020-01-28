package main

import (
	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/math"
	"github.com/buchanae/ink/rand"
	"github.com/buchanae/ink/voronoi"
)

/*
TODO ideas:
- want easier padding/spacing of blocks etc.
- want easier variable curves (e.g. x^2)
  always find it unintuitive to define a model
	where a variable grows how I want
*/

const (
	N       = 10
	Padding = 0.01
	Margin  = Padding * 2
)

func Ink(doc *app.Doc) {
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

		/*
			ma := VoronoiMesh{
				Rect:    ra,
				Spacing: math.Interp(0.002, 0.02, i/N),
			}
		*/

		ca := VoronoiCells{
			Rect:    ra,
			Spacing: math.Interp(0.003, 0.03, i/N),
		}

		gfx.Fill{
			Mesh:  ca.Mesh(),
			Color: c,
		}.Draw(doc)

		gfx.Fill{
			Color: c,
			Mesh:  rb,
		}.Draw(doc)
	}

}

// TODO voronoi should have a Mesh() method

// TODO this is a good example of a sketch where memoization
//      (smarter caching) would help. this takes time to generate
//      millions of blue noise points. would be nice to skip that
//      if all we need is to change the color

type VoronoiMesh struct {
	Rect
	Spacing float32
}

func (vm VoronoiMesh) Mesh() Mesh {

	bn := rand.BlueNoise{
		Rect:    vm.Rect,
		Spacing: vm.Spacing,
	}
	noise := bn.Generate()
	v := voronoi.New(noise, vm.Rect)

	// voronoi generates triangles that share edges,
	// so track which lines have already been drawn
	// to avoid double-drawing shared edges.
	seen := map[Line]struct{}{}

	var meshes []Mesh

	for _, t := range v.Triangulate() {
		for _, e := range t.Edges() {

			if e.A.X < e.B.X || e.A.Y < e.B.Y {
				e.A, e.B = e.B, e.A
				// cool mistake
				//e.A = e.B
			}
			if _, ok := seen[e]; ok {
				continue
			}
			seen[e] = struct{}{}

			stk := e.Stroke(StrokeOpt{
				Width: 0.0005,
			})
			meshes = append(meshes, stk.Mesh())
		}
	}

	return Merge(meshes...)
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

func DotShade(doc gfx.Layer, r Rect, c RGBA, i float32) {
	bn := rand.BlueNoise{
		Rect:    r,
		Spacing: math.Interp(0.002, 0.009, i/N),
	}
	noise := bn.Generate()
	for _, xy := range noise {
		gfx.Dot{xy, c, 0.001}.Draw(doc)
	}
}
