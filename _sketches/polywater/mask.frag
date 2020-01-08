#version 330 core

uniform sampler2D u_mask;
uniform vec4 u_color;
uniform float u_opacity;
in vec2 v_uv;
out vec4 color;

void main() {
  vec4 m = texture(u_mask, v_uv.xy);
  color = vec4(u_color.rgb, m.r * u_opacity);
}
