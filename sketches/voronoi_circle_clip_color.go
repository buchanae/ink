package main

import (
	"sort"

	"github.com/buchanae/ink/clip"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/math"
	"github.com/buchanae/ink/rand"
	"github.com/buchanae/ink/tess"
	"github.com/buchanae/ink/voronoi"
)

type PaletteIter struct {
	Palette []RGBA
	i       int
}

func (p *PaletteIter) Next() RGBA {
	c := p.Palette[p.i%len(p.Palette)]
	p.i++
	return c
}

func sortClockwise(xys []XY) {
	// TODO research best way to find center, best way to sort point
	min := xys[0].Copy()
	max := xys[0].Copy()

	for i := 1; i < len(xys); i++ {
		xy := xys[i]
		if xy.X > max.X {
			max.X = xy.X
		}
		if xy.Y > max.Y {
			max.Y = xy.Y
		}
		if xy.X < min.X {
			min.X = xy.X
		}
		if xy.Y < min.Y {
			min.Y = xy.Y
		}
	}

	center := max.Sub(min).DivScalar(0.5)
	// TODO find center
	center = XY{.5, .5}

	sort.Slice(xys, func(i, j int) bool {
		a := xys[i]
		b := xys[j]

		if a.X-center.X >= 0 && b.X-center.X < 0 {
			return true
		}
		if a.X-center.X < 0 && b.X-center.X >= 0 {
			return false
		}
		if a.X-center.X == 0 && b.X-center.X == 0 {
			if a.Y-center.Y >= 0 || b.Y-center.Y >= 0 {
				return a.Y > b.Y
			}
			return b.Y > a.Y
		}

		// compute the cross product of vectors (center -> a) x (center -> b)
		det := (a.X-center.X)*(b.Y-center.Y) - (b.X-center.X)*(a.Y-center.Y)
		if det < 0 {
			return true
		}
		if det > 0 {
			return false
		}

		// points a and b are on the same line from the center
		// check which point is closer to the center
		d1 := (a.X-center.X)*(a.X-center.X) + (a.Y-center.Y)*(a.Y-center.Y)
		d2 := (b.X-center.X)*(b.X-center.X) + (b.Y-center.Y)*(b.Y-center.Y)
		return d1 > d2
	})

}

func Ink(doc gfx.Doc) {
	rand.SeedNow()
	gfx.Clear(doc, White)
	pal := PaletteIter{Palette: rand.Palette()}

	circle := Circle{
		XY:       gfx.Center,
		Radius:   .3,
		Segments: 80,
	}
	cpts := []XY{}
	for _, curve := range circle.Path() {
		cpts = append(cpts, curve.Start())
	}

	rects := []Rect{
		RectCenter(XY{.5, .76}, XY{.9, .20}),
		RectCenter(XY{.5, .5}, XY{.9, .25}),
		RectCenter(XY{.5, .24}, XY{.9, .20}),
	}

	mask := doc.NewLayer()
	gfx.Clear(mask, RGBA{0, 0, 0, 0})

	type Clip struct {
		XYs   []XY
		Color RGBA
	}
	clips := []Clip{}

	for _, rect := range rects {
		pts := clip.Intersection(circle.Path(), rect.Path())
		clips = append(clips, Clip{
			XYs:   pts,
			Color: pal.Next(),
		})
	}

	for _, c := range clips {
		// TODO tesselation requires points to be in a specific order
		sortClockwise(c.XYs)
		tris := tess.Tesselate(c.XYs)

		gfx.Fill{
			Shape: Triangles(tris),
			Color: c.Color,
		}.Draw(mask)
	}

	s := doc.NewLayer()
	gfx.Fill{
		Shape: VoronoiCells{
			Rect:    gfx.Fullscreen,
			Spacing: math.Interp(0.001, 0.01, .9),
		},
		Color: Black,
	}.Draw(s)

	//gfx.Dots{cpts[:1], Red, 0.005}.Draw(doc)

	//gfx.Clear(doc, White)
	Mask{
		Rect:   gfx.Fullscreen,
		Source: s,
		Mask:   mask,
	}.Draw(doc)

	for _, c := range clips {
		gfx.Stroke{
			Shape: XYsToPath(c.XYs),
			Color: c.Color,
			Width: 0.005,
		}.Draw(doc)
	}

}

type Mask struct {
	Rect         Rect
	Source, Mask gfx.Layer
}

func (m Mask) Draw(l gfx.Layer) {
	l.AddShader(&gfx.Shader{
		Name: "Mask",
		Vert: gfx.DefaultVert,
		Frag: `
			#version 330 core

			uniform sampler2D u_source;
			uniform sampler2D u_mask;
			in vec2 v_uv;
			in vec2 v_vert;
			out vec4 color;

			void main() {
				vec4 m = texture(u_mask, v_vert);
				vec4 s = texture(u_source, v_vert);
				color = vec4(v_uv.xy, 0, 1);
				color = m;
				color = vec4(m.rgb, s.a * m.a);
			}
			`,
		Mesh: m.Rect.Fill(),
		Attrs: gfx.Attrs{
			"u_source": m.Source.LayerID(),
			"u_mask":   m.Mask.LayerID(),
			"a_uv": []float32{
				0, 0,
				0, 1,
				1, 1,
				1, 0,
			},
		},
	})
}

type VoronoiCells struct {
	Rect    Rect
	Spacing float32
}

// TODO fill? stroke? mesh?
func (vc VoronoiCells) Fill() Mesh {
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
