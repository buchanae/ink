package glsl

import (
	"log"
	"testing"
)

func TestInspect(t *testing.T) {
	log.SetFlags(0)
	meta := Inspect(DefaultVert)
	ten := meta.Defs[10]
	if ten.Name != "v_vert" || ten.Type != "out" || ten.DataType != "vec2" {
		t.Error("incorrect meta")
	}
	log.Print(meta)
}

const DefaultVert = `
#version 330 core

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
  vec2 glspace = v * 2 - 1;
  gl_Position = vec4(glspace.x, glspace.y, 0, 1.0);
  v_color = a_color;
	v_uv = a_uv;
	v_normal = a_normal;
	v_vert = a_vert;
}
`

const DefaultFrag = `
#version 330 core

in vec4 v_color;
out vec4 color;

void main() {
  color = v_color;
}
`

const BlurFrag = `
#version 330 core

uniform bool horizontal;
uniform sampler2D source;

in vec2 v_uv;

out vec4 color;

void main() {             
	float weight[5] = float[] (0.227027, 0.1945946, 0.1216216, 0.054054, 0.016216);
  vec2 tex_offset = 1.0 / textureSize(source, 0); // gets size of single texel
  tex_offset = vec2(.004);
  vec4 pixel = texture(source, v_uv);
  vec4 result = pixel.rgba * weight[0]; // current fragment's contribution

	color = result;
}
`
