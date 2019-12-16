package gfx

import (
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
)

func Fullscreen() dd.Mesh {
	return dd.Rect{B: dd.XY{1, 1}}.Mesh()
}

type Shader struct {
	Name             string
	Vert, Frag, Geom string
	Offscreen        bool
	Mesh             dd.Mesh
	Uniforms         Uniforms
	Attrs
}

func NewShader(m dd.Mesh) *Shader {
	return &Shader{
		Vert: "!default.vert",
		Frag: "!default.frag",
		Mesh: m,
	}
}

func (s *Shader) Draw(doc *Doc) *Layer {
	return doc.NewLayer(s)
}

type Image struct {
	Name string
	Mesh dd.Mesh
}

type Uniforms map[string]interface{}

type Layer struct {
	ID    int
	Value interface{}
}

type Drawable interface {
	Draw(doc *Doc) *Layer
}

func NewDoc() *Doc {
	return &Doc{}
}

func Clear(c color.RGBA) Fill {
	return Fill{Fullscreen(), c}
}

type Doc struct {
	Layers    []*Layer
	currentID int
}

func (doc *Doc) NewLayer(val interface{}) *Layer {
	doc.currentID++
	layer := &Layer{
		ID:    doc.currentID,
		Value: val,
	}
	doc.Layers = append(doc.Layers, layer)
	return layer
}

func (doc *Doc) Draw(d Drawable) *Layer {
	return d.Draw(doc)
}
