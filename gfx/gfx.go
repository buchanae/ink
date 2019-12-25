package gfx

import (
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
)

func Clear(c color.RGBA) Fill {
	return Fill{Fullscreen, c}
}

var Fullscreen = dd.Rect{B: dd.XY{1, 1}}

type Image struct {
	Name string
	Mesh dd.Mesh
}

type Drawable interface {
	Draw(*Layer)
}

type Meshable interface {
	Mesh() dd.Mesh
}
