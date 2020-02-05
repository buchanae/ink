package gfx

import (
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
)

// TODO Dot doesn't batch
type Dot struct {
	XY     dd.XY
	Color  color.RGBA
	Radius float32
}

func (dot Dot) Draw(l Layer) {
	if dot.Color.IsZero() {
		dot.Color = color.Red
	}
	if dot.Radius == 0 {
		dot.Radius = 0.002
	}

	l.AddShader(&Shader{
		Vert: DotVert,
		Frag: DotFrag,
		Mesh: dd.Rect{
			A: dot.XY.SubScalar(dot.Radius),
			B: dot.XY.AddScalar(dot.Radius),
		}.Fill(),
		Attrs: Attrs{
			"u_color": dot.Color,
		},
	})
}

const DotVert = `
#version 330 core

in vec2 a_vert;
in vec2 a_uv;

out vec4 v_color;
out vec2 v_uv;
out vec2 v_vert;

void main() {
  vec2 v = a_vert;

  // OpenGL fragment coordinates are [-1, 1]
  vec2 glspace = v * 2 - 1;
  gl_Position = vec4(glspace.x, glspace.y, 0, 1.0);
	v_uv = a_uv;
	v_vert = a_vert;
}
`

const DotFrag = `
#version 330 core

uniform vec4 u_color;
in vec2 v_uv;
out vec4 color;

void main() {
	float d = distance(v_uv, vec2(0.5));
	float a = smoothstep(0.5, 0.49, d);
	color = vec4(u_color.rgb, u_color.a * a);
}
`
