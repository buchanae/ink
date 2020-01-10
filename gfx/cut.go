package gfx

type Cut struct {
	Shape  Meshable
	Source Layer
}

func (cut Cut) Draw(out Layer) {
	out.AddShader(&Shader{
		Vert: DefaultVert,
		Frag: `
			#version 330 core

			uniform sampler2D u_image;
			in vec2 v_vert;
			out vec4 color;

			void main() {
				color = texture(u_image, v_vert);
			}
		`,
		Mesh: cut.Shape,
		Attrs: Attrs{
			"u_image": cut.Source.LayerID(),
		},
	})
}
