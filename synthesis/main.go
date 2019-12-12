package main

/*
Algorithm based on:
http://graphics.stanford.edu/papers/texture-synthesis-sig00/texture.pdf

nearest neighbor data structure based on:
http://www1.cs.columbia.edu/CAVE/publications/pdfs/Kumar_ECCV08_2.pdf
*/

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {

	var conf = struct {
		source    string
		output    string
		size      int
		patchSize int
	}{
		source:    "source.png",
		output:    "output.png",
		size:      1024,
		patchSize: 9,
	}

	flag.StringVar(&conf.source, "source", conf.source, "Source PNG")
	flag.Parse()

	if conf.source == "" {
		flag.Usage()
		os.Exit(1)
	}

	log.SetFlags(0)
	rand.Seed(time.Now().Unix())

	src, err := readImageFile(conf.source)
	if err != nil {
		log.Printf("error: loading source image: %v", err)
		os.Exit(1)
	}

	out := synthesize(src, conf.size, conf.patchSize)

	err = writePNG(out, conf.output)
	if err != nil {
		log.Printf("error: writing output image: %v", err)
		os.Exit(1)
	}
}

func synthesize(source image.Image, size, patchSize int) image.Image {

	patches := patches(source, patchSize)
	tree := newVPTree(patches)
	rect := image.Rect(0, 0, size, size)
	out := image.NewRGBA(rect)
	noise(out)
	start := time.Now()

	for y := 0; y < size; y++ {

		progress := float64(y) / float64(size) * 100
		fmt.Fprintf(os.Stderr, "progress: %.1f%%\r", progress)

		for x := 0; x < size; x++ {
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

	log.Printf("finished in %s", time.Since(start))
	return out
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

func writePNG(img image.Image, path string) error {
	outFh, err := os.Create(path)
	defer outFh.Close()
	if err != nil {
		return err
	}
	return png.Encode(outFh, img)
}

func readImageFile(path string) (image.Image, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// rgb holds an RGB color using an int for each channel
// so that it's easy to sum/average pixels during distance calculation
// (stdlib RBGA uses uint8, which is too limited).
type rgb struct {
	r, g, b int
}

// patch holds a set of pixels used during synthesis.
//
// A patch has a specific L-shape (it's not rectangular).
// For example, patch of size 5 around the point C looks like:
//
//   xxxxx
//   xxxxx
//   xxC__
//   _____
//   _____
//
// ...where "x" is the patch. The patch is shaped this way because
// during synthesis the algorithm fills in pixels from top-to-bottom,
// left-to-right, so only the pixels above and to the left of C are useful
// (the rest are random noise).
type patch []rgb

func patchLen(patchSize int) int {
	ns := patchSize / 2
	return ns*patchSize + ns + 1
}

func newPatch(size, cx, cy int, img image.Image) patch {
	ns := size / 2
	n := make([]rgb, patchLen(size))
	i := 0

	for y := cy - ns; y < cy; y++ {
		for x := cx - ns; x < cx+ns+1; x++ {
			coli := wrap(x, y, img)
			col, ok := coli.(color.RGBA)
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
		coli := wrap(x, cy, img)
		col, ok := coli.(color.RGBA)
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

// distance returns the distance between two patches.
// the distance is sum of square difference for each pixel.
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

// patches returns all the patches in the given source image.
func patches(img image.Image, size int) []patch {
	bounds := img.Bounds()
	// don't create patches around the edges of the image,
	// in order to avoid wrapping/edge issues.
	ns := size / 2
	minY := bounds.Min.Y + ns
	maxY := bounds.Max.Y
	minX := bounds.Min.X + ns
	maxX := bounds.Max.X - ns
	patches := make([]patch, 0, (maxY-minY)*(maxX-minX))

	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			n := newPatch(size, x, y, img)
			patches = append(patches, n)
		}
	}
	return patches
}
