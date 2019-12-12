package main

/*
Algorithm based on:
http://graphics.stanford.edu/papers/texture-synthesis-sig00/texture.pdf

nearest neighbor data structure based on:
http://www1.cs.columbia.edu/CAVE/publications/pdfs/Kumar_ECCV08_2.pdf
*/

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	patchSize  = 9
	outputSize = 1024
)

func main() {
	rand.Seed(time.Now().Unix())

	patches, err := loadSource("source.png")
	if err != nil {
		panic(err)
	}
	tree := newVPTree(patches)

	bounds := image.Rect(0, 0, outputSize, outputSize)
	out := image.NewRGBA(bounds)
	noise(out)

	start := time.Now()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			q := newPatch(patchSize, x, y, out)
			_, best := tree.findClosest(q)
			b := best[len(best)-1]
			out.Set(x, y, color.RGBA{
				R: uint8(b.r),
				G: uint8(b.g),
				B: uint8(b.b),
				A: 255,
			})
		}
	}

	log.Printf("took %s", time.Since(start))
	writePNG(out, "output.png")
}

// noise fills the given image with random colors.
func noise(img *image.RGBA) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(rand.Intn(0xff)),
				G: uint8(rand.Intn(0xff)),
				B: uint8(rand.Intn(0xff)),
				A: 0xff,
			})
		}
	}
}

// distance returns the distance between two patches.
// the distance metric is sum of square difference for each
// pixel.
//
// the last pixel is not included because that pixel is usually
// a random color in the image being synthesized.
func distance(a, b patch) int {
	d := 0
	for i := range a[:len(a)-1] {
		av := a[i]
		bv := b[i]
		r := av.r - bv.r
		g := av.g - bv.g
		b := av.b - bv.b
		d += r*r + g*g + b*b
	}
	return d
}

func writePNG(img image.Image, path string) {
	outFh, err := os.Create(path)
	defer outFh.Close()
	if err != nil {
		panic(err)
	}

	err = png.Encode(outFh, img)
	if err != nil {
		panic(err)
	}
}

// loadSource loads a source image from the given "path"
// into a list of patches.
func loadSource(path string) ([]patch, error) {
	f, err := os.Open("source.png")
	defer f.Close()
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	var patches []patch
	ns := patchSize / 2

	for y := bounds.Min.Y + ns; y < bounds.Max.Y-ns; y++ {
		for x := bounds.Min.X + ns; x < bounds.Max.X-ns; x++ {
			n := newPatch(patchSize, x, y, img)
			patches = append(patches, n)
		}
	}
	return patches, nil
}

type rgb struct {
	r, g, b int
}

type patch []rgb

func patchLen() int {
	ns := patchSize / 2
	return ns*patchSize + ns + 1
}

func newPatch(size, cx, cy int, img image.Image) patch {
	bounds := img.Bounds()
	ns := size / 2
	n := make([]rgb, patchLen())
	i := 0

	for y := cy - ns; y < cy; y++ {
		for x := cx - ns; x < cx+ns+1; x++ {
			col, ok := wrap(x, y, img).(color.RGBA)
			if !ok {
				col = color.RGBAModel.Convert(coli).(color.RGBA)
			}

			n[i] = rgb{
				r: int(col.R),
				g: int(col.G),
				b: int(col.B),
			}
			i++
		}
	}
	for x := cx - ns; x < cx+1; x++ {
		col, ok := wrap(x, y, img).(color.RGBA)
		if !ok {
			col = color.RGBAModel.Convert(coli).(color.RGBA)
		}

		n[i] = rgb{
			r: int(col.R),
			g: int(col.G),
			b: int(col.B),
		}
		i++
	}

	return n
}

// wrap wraps the x,y coordinate around the edges of the image
// and returns the pixel.
func wrap(x, y int, img image.Image) color.Color {
	bounds := img.Bounds()
	if x < bounds.Min.X {
		x = bounds.Max.X - bounds.Min.X - x
	}
	if y < bounds.Min.Y {
		y = bounds.Max.Y - bounds.Min.Y - y
	}
	if x > bounds.Max.X {
		x = x - bounds.Max.X
	}
	if y > bounds.Max.Y {
		y = y - bounds.Max.Y
	}
	return img.At(x, y)
}
