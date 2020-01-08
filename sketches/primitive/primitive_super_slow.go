package main

import (
	"image"
	"log"
	"math"
	"os"
	"time"

	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
	"github.com/buchanae/ink/win"
)

func main() {
	rand.SeedNow()

	//log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.SetFlags(0)
	start := time.Now()

	conf := app.DefaultConfig()
	conf.Window.Width = 128
	conf.Window.Height = 128
	a, err := app.NewApp(conf)
	if err != nil {
		panic(err)
	}

	target := LoadImage("sketches/primitive/monalisa.png")
	palette := colors(target)
	avg := avgcolors(target)

	go func() {

		doc := NewDoc()
		targetLayer := doc.NewImage(target)
		Clear(doc, Black)
		Copy(doc, targetLayer)
		a.Render(doc)

		// Fill with average color.
		current := NewDoc()
		f := Fill{Fullscreen, avg}
		f.Draw(current)
		f.Draw(doc)

		a.Render(doc)
		a.RenderOffscreen(current)

		// Listen for "enter" keypresses
		next := make(chan struct{})
		go func() {
			for ev := range a.Events() {
				if ev == win.ReturnEvent {
					next <- struct{}{}
				}
			}
		}()
		//<-next

		diff := NewDoc()
		work := NewDoc()

		var bestScore = math.Inf(1)
		nodeCount := 0
		attempts := 0

		for i := 0; nodeCount < 1050; i++ {
			var best *Attempt
			var avgTime time.Duration
			count := 0

			start := RectCenter(
				rand.XYRange(0.0, 0.99),
				rand.XYRange(0.001, 1),
			).Quad()
			quad := start

			for j := 0; j < 50; j++ {
				//for j := 0; j < len(palette)*10; j++ {
				attempts++
				attemptStart := time.Now()

				quad = rand.TweakQuad(quad, 0.1)

				size := quad.Bounds().Size()
				if size.X < 0.003 || size.Y < 0.003 {
					log.Printf("BOUNDS: %v", size)
					break
				}

				col := rand.Color(palette)
				//col.A = rand.Range(.8, 1)
				col.A = 0.5
				attempt := Attempt{
					Mesh:  quad.Mesh(),
					Color: col,
				}

				diff.Clear()
				work.Clear()

				var score float64

				Copy(work, current)
				Fill{attempt.Mesh, attempt.Color}.Draw(work)

				a.RenderOffscreen(work)
				//r.RenderToScreen()

				Diff(diff, targetLayer, work)

				a.RenderOffscreen(diff)

				//diffStart := time.Now()
				pixels := r.CapturePixels(diff.LayerID(), 0, 0, 1, 1)
				//log.Printf("diff time: %s", time.Since(diffStart))
				score = sumdiff(pixels)
				//log.Printf("score: %v %v %s", score, bestScore)

				//r.RenderToScreen()

				//a.Swap()
				//<-next

				avgTime += time.Since(attemptStart)
				count++

				if score < bestScore {
					best = &attempt
					bestScore = score
					break
				}
			}
			if best == nil {
				log.Printf("FAIL: %v", bestScore)
				continue
			}

			nodeCount++
			log.Printf("NODES: %s %d %d %v %d", time.Since(start), i, nodeCount, bestScore, attempts)
			log.Printf("attempt time: %s", avgTime/time.Duration(count))

			fill := Fill{best.Mesh, best.Color}
			fill.Draw(current)
			fill.Draw(doc)

			a.Render(doc)
			a.RenderOffscreen(current)
		}

		a.Render(doc)
	}()

	a.Run()
}

func Diff(dest, target, attempt Layer) {
	dest.AddShader(&Shader{
		Vert: DefaultVert,
		Frag: `#version 330 core

			uniform sampler2D u_image;
			uniform sampler2D u_work;
			in vec2 v_vert;
			in vec2 v_uv;
			out vec4 color;

			void main() {
				vec4 target = texture(u_image, v_uv);
				vec4 work = texture(u_work, v_uv);
				vec3 d = work.rgb - target.rgb;
				d = d*d;
				color = vec4(d.r + d.g + d.b, 0, 0, 1);
			}
		`,
		Mesh: Fullscreen,
		Attrs: Attrs{
			"u_image": target.LayerID(),
			"u_work":  attempt.LayerID(),
		},
	})
}

type Attempt struct {
	Mesh  Mesh
	Color RGBA
}

func sumdiff(pixels []uint8) float64 {
	var total int
	var count = len(pixels)

	for i := 0; i < len(pixels); i++ {
		br := pixels[i]
		total += int(br)
	}
	//log.Printf("NUM PIX: %v %v", total, count)
	tot := float64(total) / 255
	score := math.Sqrt(tot / float64(count))
	if math.IsNaN(score) {
		log.Printf("NAN")
		return math.Inf(1)
	}
	return score
}

func colors(img image.Image) []RGBA {
	const N = 20
	size := img.Bounds().Size()
	gx := size.X / N
	gy := size.Y / N

	var out []RGBA

	for cx := 0; cx < gx; cx++ {
		for cy := 0; cy < gy; cy++ {
			col := RGBA{A: 1}
			count := 0

			for ix := 0; ix < N; ix++ {
				x := cx*N + ix
				if x >= size.X {
					break
				}

				for iy := 0; iy < N; iy++ {
					y := cy*N + iy
					if iy >= size.Y {
						break
					}
					r, g, b, _ := img.At(x, y).RGBA()
					col.R += float32(r) / 0xffff
					col.G += float32(g) / 0xffff
					col.B += float32(b) / 0xffff
					count++
				}
			}
			col.R = col.R / float32(count)
			col.G = col.G / float32(count)
			col.B = col.B / float32(count)
			out = append(out, col)
		}
	}

	return out
}

func avgcolors(img image.Image) RGBA {
	size := img.Bounds().Size()
	col := RGBA{A: 1}
	count := 0

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			col.R += float32(r) / 0xffff
			col.G += float32(g) / 0xffff
			col.B += float32(b) / 0xffff
			count++
		}
	}

	col.R = col.R / float32(count)
	col.G = col.G / float32(count)
	col.B = col.B / float32(count)

	return col
}

func LoadImage(path string) image.Image {
	fh, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(fh)
	if err != nil {
		panic(err)
	}
	return img
}

/*
	TODO cut/mask image using shape. could be useful in gfx.
		work.AddShader(&Shader{
			Offscreen: true,
			Vert:      DefaultVert,
			Frag: `
			#version 330 core

			uniform sampler2D u_image;
			in vec2 v_vert;
			out vec4 color;

			void main() {
				color = texture(u_image, v_vert);
			}
		`,
			Mesh: attempt,
			Attrs: Attrs{
				"u_image": targetLayer.LayerID(),
			},
		})
*/
