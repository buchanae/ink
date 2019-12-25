package gfx

import (
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
)

var currentID int

func NewLayer() *Layer {
	currentID++
	return &Layer{
		ID: currentID,
	}
}

type Layer struct {
	ID     int
	Values []interface{}
}

func (l *Layer) Clear(c color.RGBA) {
	l.Draw(Fill{
		Mesh:  Fullscreen,
		Color: c,
	})
}

func (l *Layer) Layer() *Layer {
	x := NewLayer()
	l.Values = append(l.Values, x)
	return x
}

func (l *Layer) Shader(m Meshable) *Shader {
	s := NewShader(m)
	l.Values = append(l.Values, s)
	return s
}

func (l *Layer) Dot(xy dd.XY, c color.RGBA) {
	l.Shader(dd.Circle{xy, 0.005, 8}).
		Set("a_color", c)
}

func (l *Layer) Draw(ds ...Drawable) {
	for _, d := range ds {
		d.Draw(l)
	}
}
