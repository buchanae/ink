package gfx

import (
	"github.com/buchanae/ink/color"
)

type Fill struct {
	Mesh  Meshable
	Color color.RGBA
}

func (f Fill) Draw(l Layer) {
	l.AddShader(&Shader{
		Name: "Fill",
		Vert: DefaultVert,
		Frag: DefaultFrag,
		Mesh: f.Mesh,
		Attrs: Attrs{
			"a_color": f.Color,
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
