package gfx

const CopyFrag = `
#version 330 core

uniform sampler2D u_image;
in vec2 v_uv;
out vec4 color;

void main() {
	color = texture(u_image, v_uv);
}
`
