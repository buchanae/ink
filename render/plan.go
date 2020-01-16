package render

import (
	"image"
)

// TODO is a plan resolution independent?
//      i.e. vertex data only, no texture specs.
type Plan struct {
	Shaders []*Shader
	Images  map[int]image.Image
}

type Shader struct {
	Name             string
	Vert, Frag, Geom string
	Output           string
	Layer            int
	Vertices         int
	Instances        int
	Faces            []uint32
	Uniforms         map[string]interface{}
	Attrs            map[string]Attr
	Blend            Blend
}

type Attr struct {
	Value   interface{}
	Size    int
	Divisor int
}

type Blend int

const (
	Normal Blend = iota
	Darken
	Multiply
)
