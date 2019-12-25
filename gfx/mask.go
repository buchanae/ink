package gfx

type Mask struct {
	Mesh         Meshable
	Source, Mask *Layer
}

func (m Mask) Draw(l *Layer) {
	l.Draw(&Shader{
		Name: "Mask",
		Vert: DefaultVert,
		Frag: MaskFrag,
		Mesh: m.Mesh,
		Attrs: Attrs{
			"u_source": m.Source,
			"u_mask":   m.Mask,
			"a_uv": []float32{
				0, 0,
				0, 1,
				1, 1,
				1, 0,
			},
		},
	})
}

const MaskFrag = `
#version 330 core

uniform sampler2D u_source;
uniform sampler2D u_mask;
in vec2 v_uv;
out vec4 color;

void main() {
  vec4 m = texture(u_mask, v_uv);
  vec4 s = texture(u_source, v_uv);
	color = vec4(s.rgb, m.a * s.a);
}
`
