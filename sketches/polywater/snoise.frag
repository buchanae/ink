#version 330 core

in vec2 v_uv;
uniform float u_time;
uniform float u_speed;
out vec4 color;

include /snoise.glsl

void main() {
  vec2 v2 = v_uv * vec2(05) + vec2(u_time * u_speed, 0);
  float z2 = noise2(v2.xx).y;

  vec2 v = v_uv * vec2(40) + vec2(u_time * u_speed, 0);
  float z = noise2(v.xx).y;
  z += z2 + .3;
  z += 0.1;
  //z = z + z * z;
  //z = abs(z);
  color = vec4(z, z, z, 1);
}
