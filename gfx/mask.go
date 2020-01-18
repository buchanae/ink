package gfx

import "github.com/buchanae/ink/dd"

type Mask struct {
	Rect         dd.Rect
	Source, Mask Layer
}

func (m Mask) Draw(l Layer) {
	l.AddShader(&Shader{
		Name: "Mask",
		Vert: DefaultVert,
		Frag: MaskFrag,
		Mesh: m.Rect,
		Attrs: Attrs{
			"u_source": m.Source.LayerID(),
			"u_mask":   m.Mask.LayerID(),
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
