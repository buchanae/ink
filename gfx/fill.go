package gfx

import (
	"github.com/buchanae/ink/color"
)

type Fill struct {
	Mesh  Meshable
	Color color.RGBA
}

func (f Fill) Draw(l *Layer) {
	l.Draw(&Shader{
		Name: "Fill",
		Vert: DefaultVert,
		Frag: FillFrag,
		Mesh: f.Mesh,
		Attrs: Attrs{
			"u_color": f.Color,
		},
	})
}

const FillFrag = `
#version 330 core

uniform vec4 u_color;
out vec4 f_color;

void main() {
  f_color = u_color;
}
`
