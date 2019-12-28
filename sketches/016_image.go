package main

import (
	"image"
	_ "image/png"
	"os"

	"github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
)

func Ink(doc *Doc) {
	fh, err := os.Open("examples/toshiro.png")
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(fh)
	if err != nil {
		panic(err)
	}

	l := doc.NewImage(img)

	w, h := img.Bounds().Dx(), img.Bounds().Dy()

	rw := float32(w) / float32(800)
	rh := float32(h) / float32(800)

	doc.AddShader(&Shader{
		Vert: DefaultVert,
		Frag: CopyFrag,
		Mesh: dd.RectCenter(dd.XY{0.5, 0.5}, dd.XY{rw, rh}),
		Attrs: Attrs{
			"u_image": l.LayerID(),
			"a_uv": []float32{
				0, 0,
				0, 1,
				1, 1,
				1, 0,
			},
		},
	})
}
