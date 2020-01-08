// Package gfx contains types for drawing graphics to an Ink application.
package gfx

import (
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
