#version 330 core

uniform float u_time;
uniform float u_speed;
uniform float u_turbulence;
uniform sampler2D u_noise;
in vec2 a_vert;
in vec2 a_uv;
in vec2 a_pos;
in float a_rot;
in vec2 a_size;
out vec2 v_uv;

void main() {
  float amp = .05;
  vec2 v = a_vert;
  vec2 s = a_size;
  float z = texture(u_noise, a_vert).r * u_turbulence;
  s += vec2(0, z);

  mat2 scaleMat = mat2(
    s.x, 0,
    0, s.y
  );

  mat2 rotmat = mat2(
     cos(a_rot), sin(a_rot),
    -sin(a_rot), cos(a_rot)
  );

  v *= scaleMat;
  // position at center of rect/ellipse
  v -= vec2(s.xy / 2);
  v *= rotmat;

  // Fragment coordinates are [-1, 1]
  v += a_pos * 2 - 1;

  gl_Position = vec4(v, 0, 1.0);

  v_uv = a_uv;
}

