#version 330 core

uniform vec4 u_color;
uniform float u_opacity;
uniform float u_blur;
uniform float u_boxy;
uniform float u_time;
in vec2 v_uv;
out vec4 color;

void main() {
  float base = 1 - u_boxy;
  float d = smoothstep(
      base,
      base + u_blur,
      1 - distance(v_uv, vec2(0.5))
  );
  color = vec4(u_color.rgb, d * u_opacity);
  //color = vec4(vec2(vUV), 1, 1);
  //color = vec4(1, 1, 1, d * u_opacity);
}
