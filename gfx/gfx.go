package gfx

import (
	"log"

	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
)

func init() {
	log.SetFlags(0)
}

// TODO use opengl Clear command
func Clear(l Layer, c color.RGBA) {
	Fill{Fullscreen, c}.Draw(l)
}

var Fullscreen = dd.Rect{B: dd.XY{1, 1}}

type Meshable interface {
	Mesh() dd.Mesh
}

type Layer interface {
	LayerID() int
	NewLayer() Layer
	AddShader(*Shader)
}
