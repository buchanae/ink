package gfx

import (
	"image"

	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
)

// TODO use opengl Clear command
func Clear(l Layer, c color.RGBA) {
	Fill{Fullscreen, c}.Draw(l)
}

var Fullscreen = dd.Rect{B: dd.XY{1, 1}}

type Meshable interface {
	Mesh() dd.Mesh
}

type Strokeable interface {
	Stroke() dd.Stroke
}

type Layer interface {
	LayerID() int
	NewLayer() Layer
	NewImage(image.Image) Image
	AddShader(*Shader)
}

type Image struct {
	ID            int
	Width, Height float32
}

func (img Image) Draw(out Layer) {
	out.AddShader(&Shader{
		Vert: DefaultVert,
		Frag: CopyFrag,
		Mesh: dd.RectCenter(
			dd.XY{0.5, 0.5},
			dd.XY{img.Width, img.Height},
		),
		Attrs: Attrs{
			"u_image": img,
		},
	})
}
