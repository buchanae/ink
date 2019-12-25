package gfx

import (
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
)

type Gradient struct {
	Rect dd.Rect
	A, B color.RGBA
}

func (g Gradient) Draw(l *Layer) {
	l.Draw(&Shader{
		Name: "Gradient",
		Vert: DefaultVert,
		Frag: GradientFrag,
		Mesh: g.Rect,
		Attrs: Attrs{
			"u_color_a": g.A,
			"u_color_b": g.B,
			"a_uv": []float32{
				0, 0,
				0, 1,
				1, 1,
				1, 0,
			},
		},
	})
}

const GradientFrag = `
#version 330 core

uniform vec4 u_color_a;
uniform vec4 u_color_b;
in vec2 v_uv;
out vec4 color;

void main() {
  color = mix(u_color_a, u_color_b, v_uv.x);
}
`
