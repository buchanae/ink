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
