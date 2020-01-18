package gfx

func Copy(dst, src Layer) {
	dst.AddShader(&Shader{
		Vert: DefaultVert,
		Frag: CopyFrag,
		Mesh: Fullscreen,
		Attrs: Attrs{
			"u_image": src.LayerID(),
		},
	})
}

const CopyFrag = `
uniform sampler2D u_image;
in vec2 v_uv;
out vec4 color;

void main() {
	color = texture(u_image, v_uv);
}
`
