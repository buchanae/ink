package main

import (
	"image"
	"os"

	"github.com/buchanae/ink/app"
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/math"
	"github.com/buchanae/ink/midi"
	"github.com/buchanae/ink/rand"
	"github.com/buchanae/ink/voronoi"
)

func Ink(doc *app.Doc) {

	midin := midi.NewMidi()
	go midin.Run()

	gfx.Clear(doc, color.White)

	fh, err := os.Open("toshiro_copy.png")
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	img, _, err := image.Decode(fh)
	if err != nil {
		panic(err)
	}

	imgn := doc.NewImage(img)
	//imgn.Draw(doc)

	noise := rand.BlueNoise{
		Spacing: 0.0025,
		Limit:   70000,
	}.Generate()

	var (
		alpha   float32 = 0.3
		lumBase float32 = 130
		size    float32 = 0.3
	)

	pos, col := buildPoints(
		img, imgn.Rect, noise, alpha, lumBase,
	)
	build(doc, pos, col, size)
	app.Send(doc)

	for val := range midin.Listen() {

		switch val.Key {
		case 22:
			alpha = val.Val
			pos, col = buildPoints(
				img, imgn.Rect, noise, alpha, lumBase,
			)
		case 21:
			size = val.Val
		default:
			continue
		}

		build(doc, pos, col, size)
		app.Send(doc)
	}

	/*
		vm := VoronoiMesh{
			XYs:  pos,
			Rect: imgn.Rect,
		}
		gfx.Fill{vm.Mesh(), color.Black}.Draw(doc)
	*/
}

func buildPoints(img image.Image, rect dd.Rect, noise []dd.XY, alpha, lumBase float32) ([]dd.XY, []color.RGBA) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()

	var pos []dd.XY
	var col []color.RGBA

	for _, xy := range noise {
		if !rect.Contains(xy) {
			continue
		}

		ix := int(xy.X * float32(w))
		iy := int((1 - xy.Y) * float32(h))
		r, g, b, _ := img.At(ix, iy).RGBA()
		rf := float32(r) / 255
		gf := float32(g) / 255
		bf := float32(b) / 255
		lum := 0.2126*rf + 0.7152*gf + 0.0722*bf
		lumf := lum / lumBase

		if rand.Bool(lumf) {
			continue
		}

		pos = append(pos, xy)
		col = append(col, color.RGBA{
			R: lumf,
			G: lumf,
			B: lumf,
			A: alpha,
		})
	}
	return pos, col
}

func build(doc *app.Doc, pos []dd.XY, col []color.RGBA, sizeInput float32) {
	doc.Ops = nil
	gfx.Clear(doc, color.White)
	circle := dd.Circle{
		Radius:   math.Interp(0.0018, 0.008, sizeInput),
		Segments: 20,
	}

	s := gfx.NewShader(circle)
	s.Set("a_color", col)
	s.Set("a_pos", pos)
	s.Instances = len(pos)
	s.Divisors = map[string]int{
		"a_color": 1,
		"a_pos":   1,
	}
	s.Draw(doc)
}

type VoronoiMesh struct {
	XYs  []dd.XY
	Rect dd.Rect
}

func (vm VoronoiMesh) Mesh() dd.Mesh {

	v := voronoi.New(vm.XYs, vm.Rect)

	// voronoi generates triangles that share edges,
	// so track which lines have already been drawn
	// to avoid double-drawing shared edges.
	seen := map[dd.Line]struct{}{}

	var meshes []dd.Mesh

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

			stk := e.Stroke()
			stk.Width = 0.0005
			meshes = append(meshes, stk.Mesh())
		}
	}

	return dd.Merge(meshes...)
}
