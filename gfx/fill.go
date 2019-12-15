package gfx

import (
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
)

type Fill struct {
	Mesh  dd.Mesh
	Color color.RGBA
}

func (f Fill) Draw(doc *Doc) *Layer {
	return doc.Draw(&Shader{
		Name: "Fill",
		// TODO ! symbol relies on behavior defined
		//        outside this package
		Vert: "!default.vert",
		Frag: "!fill.frag",
		Mesh: f.Mesh,
		Uniforms: Uniforms{
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
