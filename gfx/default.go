package gfx

const DefaultVert = `
precision mediump float;

in vec2 a_vert;
in vec2 a_uv;
in vec2 a_pivot;
in vec2 a_pos;
in float a_rot;
in vec4 a_color;
in vec2 a_normal;

out vec4 v_color;
out vec2 v_uv;
out vec2 v_normal;
out vec2 v_vert;

void main() {
  vec2 v = a_vert;

  v -= a_pivot;
  v *= mat2(
     cos(a_rot), sin(a_rot),
    -sin(a_rot), cos(a_rot)
  );
  v += a_pivot;
  v += a_pos;

  // OpenGL fragment coordinates are [-1, 1]
  vec2 glspace = v * 2.0 - 1.0;
  gl_Position = vec4(glspace.x, glspace.y, 0, 1.0);
  v_color = a_color;
	v_uv = a_uv;
	v_normal = a_normal;
	v_vert = a_vert;
}
`

const DefaultFrag = `
precision mediump float;

in vec4 v_color;

void main() {
  gl_FragColor = v_color;
}
`
