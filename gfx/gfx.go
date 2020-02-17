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
var Center = dd.XY{.5, .5}

type Doc interface {
	Layer

	// TODO should be on Layer?
	LoadImage(name string) Image

	/*
		TODO
		- load image
		- create image
		- get/set config
			- window title, size
			- doc size
			- snapshot size
		- animation
			- get animation frame
			- send updates
			- clear doc?
		- get events
			- key press, mouse move
	*/
}

// TODO be able to turn layers on/off easily
type Layer interface {
	Clear()
	LayerID() int
	NewLayer() Layer
	NewImage(image.Image) Image
	AddShader(*Shader)
}
