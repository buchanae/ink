package main

import (
	"image"
	"image/draw"
	_ "image/png"
	"os"

	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
	"github.com/buchanae/ink/raster"
	"github.com/buchanae/ink/voronoi"
)

func Ink(doc *app.Doc) {
	rand.SeedNow()
	gfx.Clear(doc, Black)

	box := gfx.Fullscreen
	img := LoadImage("flower.png")
	//l := doc.NewImage(img)

	var initial []XY

	var i float32
	for ; i < 20; i++ {
		p := i / 20
		xy := box.Interpolate(XY{p, p})
		initial = append(initial, xy)
	}

	xys := rand.BlueNoise{
		Limit:   1000,
		Spacing: 0.03,
		Initial: initial,
	}.Generate()

	v := voronoi.New(xys, box)

	var colors []RGBA
	var mesh Mesh
	size := img.Bounds().Size()

	for _, cell := range v.Cells() {
		var rsum, bsum, gsum, count float32

		// TODO rename to Triangles
		lines := raster.Rasterize(cell.Tris, size.X, size.Y)
		for _, l := range lines {
			for x := l.A; x < l.B; x++ {
				i := img.PixOffset(x, l.Y)
				rsum += float32(img.Pix[i]) / 255
				gsum += float32(img.Pix[i+1]) / 255
				bsum += float32(img.Pix[i+2]) / 255
				count++
			}
		}

		avg := RGBA{
			R: rsum / count,
			G: gsum / count,
			B: bsum / count,
			A: 1,
		}

		var tris []Triangle
		// TODO use mesh faces instead
		moved := map[XY]XY{}

		for _, t := range cell.Tris {

			if xy, ok := moved[t.B]; ok {
				t.B = xy
			} else {
				mb := rand.Range(0.01, 0.2)
				vb := t.B.Sub(t.A).MulScalar(mb)
				m := t.B.Sub(vb)
				moved[t.B] = m
				t.B = m
			}

			if xy, ok := moved[t.C]; ok {
				t.C = xy
			} else {
				mc := rand.Range(0.01, 0.2)
				vc := t.C.Sub(t.A).MulScalar(mc)
				m := t.C.Sub(vc)
				moved[t.C] = m
				t.C = m
			}

			tris = append(tris, t)
			colors = append(colors,
				avg, avg, avg,
			)
		}
		mesh = mesh.AddTriangles(tris)
	}

	s := gfx.NewShader(mesh)
	s.Set("a_color", colors)
	s.Draw(doc)

	/*
				c := rand.Color(colors)
				c.A = 0.3
				for _, tri := range cell.Tris {
					s := doc.Shader(tri)
					s.Set("a_color", c)
				}

				for _, e := range cell.Edges {
					m := e.Stroke(0.002)
					s := doc.Shader(m)
				}
			}

		tris := v.Triangulate()

		for _, t := range tris {
			s := NewShader(t)
			c := rand.Color(colors)
			c.A = 0.3
			s.Set("a_color", c)
			s.Draw(doc)
		}
	*/
}

func LoadImage(path string) *image.RGBA {
	fh, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(fh)
	if err != nil {
		panic(err)
	}

	out := image.NewRGBA(img.Bounds())
	draw.Draw(out, out.Bounds(), img, image.ZP, draw.Src)
	return out
}
