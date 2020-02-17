package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/buchanae/ink/app"
	inkcolor "github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
	"github.com/buchanae/ink/raster"
)

func main() {
	log.SetFlags(0)
	rand.SeedNow()

	conf := app.DefaultConfig()
	a, err := app.NewApp(conf)
	if err != nil {
		panic(err)
	}

	res := make(chan *Attempt, 100)
	doc := a.NewDoc()

	go run(res)
	go func() {
		ticker := time.Tick(time.Second)
		stale := false
		for {
			select {
			case <-ticker:
				if stale {
					a.Render(doc)
					stale = false
				}
			case x := <-res:
				stale = true
				col := inkcolor.FromGo(x.Color)
				col.A = 0.8
				doc.AddShader(&gfx.Shader{
					// Flip coorindates since Go coordinate system
					// is upside down compared to Ink
					Vert: `#version 330 core
					in vec2 a_vert;
					in vec4 a_color;
					out vec4 v_color;
					void main() {
						vec2 v = a_vert * 2 - 1;
						gl_Position = vec4(v.x, -v.y, 0, 1);
						v_color = a_color;
					}`,
					Frag: gfx.DefaultFrag,
					Mesh: x.Mesh,
					Attrs: gfx.Attrs{
						"a_color": col,
					},
				})
			}
		}
	}()

	a.Run()
}

func run(output chan *Attempt) {

	target := LoadImage("monalisa_256.png")
	palette := colors(target)
	avg := avgcolors(target)

	output <- &Attempt{
		Mesh:  gfx.Fullscreen.Mesh(),
		Color: avg,
	}

	current := image.NewRGBA(target.Bounds())
	Fill(current, avg)

	var bestScore = math.Inf(1)
	nodeCount := 0
	attempts := 0

	for i := 0; nodeCount < 4000; i++ {

		workers := make([]*worker, 1)
		for j := range workers {
			workers[j] = &worker{
				palette:   palette,
				target:    target,
				current:   current,
				bestScore: math.Inf(1),
			}
		}

		wg := sync.WaitGroup{}
		for j := range workers {
			wg.Add(1)
			go func(j int) {
				workers[j].generate()
				wg.Done()
			}(j)
		}
		wg.Wait()

		for j := range workers {
			w := workers[j]
			attempts += w.attempts
			//log.Printf("attempt time: %s", w.avgTime/time.Duration(w.attempts))

			if w.bestScore < bestScore {
				bestScore = w.bestScore
				output <- w.best
				rasterFill(w.best.Mesh, current, w.best.Color)
				nodeCount++
				log.Printf("NODES: %d %d %v %d", i, nodeCount, bestScore, attempts)
			}
		}
	}
}

type worker struct {
	attempts  int
	avgTime   time.Duration
	palette   []color.RGBA
	current   *image.RGBA
	target    *image.RGBA
	best      *Attempt
	bestScore float64
}

func (w *worker) generate() {

	shape := dd.Circle{
		rand.XYRange(0.0, 0.99),
		rand.Range(0.001, 0.04),
		30,
	}
	/*
		quad := dd.RectCenter(
			rand.XYRange(0.0, 0.99),
			rand.XYRange(0.001, .4),
		).Quad()
	*/
	mesh := shape.Mesh()

	sb := shape.Bounds()
	tb := w.target.Bounds()
	ts := tb.Size()

	r := image.Rect(
		tb.Min.X+int(sb.A.X*float32(ts.X)),
		tb.Min.Y+int(sb.A.Y*float32(ts.Y)),
		tb.Min.X+int(sb.B.X*float32(ts.X)),
		tb.Min.Y+int(sb.B.Y*float32(ts.Y)),
	)
	sub := w.target.SubImage(r).(*image.RGBA)
	col := avgcolors(sub)
	//col := w.palette[rand.Intn(len(w.palette))]
	//col.A = 255

	for j := 0; j < 10; j++ {
		w.attempts++
		attemptStart := time.Now()

		//tweaked := rand.TweakQuad(quad, 0.01)
		tweaked := rand.TweakMesh(mesh, 0.002)

		/*
			size := tweaked.Bounds().Size()
			if size.X < 0.01 || size.Y < 0.01 {
				log.Printf("BOUNDS: %v", size)
				continue
			}
		*/

		try := Copy(w.current)
		rasterFill(tweaked, try, col)
		score := diff(w.target, try)

		w.avgTime += time.Since(attemptStart)

		if score < w.bestScore {
			w.best = &Attempt{
				Mesh:  tweaked,
				Color: col,
			}
			w.bestScore = score
			mesh = tweaked
		}
	}
}

func rasterFill(mesh dd.Mesh, img *image.RGBA, c color.RGBA) {
	size := img.Bounds().Size()
	tris := mesh.Triangles()
	lines := raster.Rasterize(tris, size.X, size.Y)

	for _, l := range lines {
		for x := l.A; x < l.B; x++ {
			i := img.PixOffset(x, l.Y)
			img.Pix[i] = c.R
			img.Pix[i+1] = c.G
			img.Pix[i+2] = c.B
		}
	}
}

func diff(a, b *image.RGBA) float64 {
	var sum, count uint64
	bnd := a.Bounds()

	for x := bnd.Min.X; x < bnd.Max.X; x++ {
		for y := bnd.Min.Y; y < bnd.Max.Y; y++ {
			i := a.PixOffset(x, y)
			r := uint64(a.Pix[i]) - uint64(b.Pix[i])
			g := uint64(a.Pix[i+1]) - uint64(b.Pix[i+1])
			b := uint64(a.Pix[i+2]) - uint64(b.Pix[i+2])
			sum += r*r + g*g + b*b
			count++
		}
	}
	return math.Sqrt(float64(sum)/float64(count)) / 255
}

type Attempt struct {
	Mesh  dd.Mesh
	Color color.RGBA
}

func colors(img *image.RGBA) []color.RGBA {
	const N = 20
	size := img.Bounds().Size()
	gx := size.X / N
	gy := size.Y / N

	var out []color.RGBA

	for cx := 0; cx < gx; cx++ {
		for cy := 0; cy < gy; cy++ {
			r := image.Rect(cx*N, cy*N, (cx+1)*N, (cy+1)*N)
			s := img.SubImage(r).(*image.RGBA)
			c := avgcolors(s)
			out = append(out, c)
		}
	}

	return out
}

func avgcolors(img *image.RGBA) color.RGBA {
	b := img.Bounds()
	var rsum, gsum, bsum, count uint64

	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			i := img.PixOffset(x, y)
			rsum += uint64(img.Pix[i])
			gsum += uint64(img.Pix[i+1])
			bsum += uint64(img.Pix[i+2])
			count++
		}
	}

	if count == 0 {
		return color.RGBA{}
	}

	return color.RGBA{
		R: uint8(rsum / count),
		G: uint8(gsum / count),
		B: uint8(bsum / count),
		A: 255,
	}
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

func Copy(img *image.RGBA) *image.RGBA {
	out := image.NewRGBA(img.Bounds())
	draw.Draw(out, out.Bounds(), img, image.ZP, draw.Src)
	return out
}

func Fill(img *image.RGBA, c color.Color) {
	b := img.Bounds()
	draw.Draw(img, b, &image.Uniform{c}, b.Min, draw.Src)
}
