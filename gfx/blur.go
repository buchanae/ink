package gfx

type Blur struct {
	Passes int
	Source Layer
}

func (b Blur) Draw(out Layer) {

	if b.Passes == 0 {
		b.Passes = 10
	}
	if b.Source == nil {
		panic("blur source layer is nil")
	}

	// Each pass is really two passes: one horizontal, one vertical.
	// The output of each pass becomes the input to the next.
	// A temp layer is used to create this chain.
	horizontal := false
	in := b.Source

	for i := 0; i < b.Passes*2; i++ {
		out.AddShader(&Shader{
			Name: "Blur",
			Vert: DefaultVert,
			Frag: BlurFrag,
			Mesh: Fullscreen,
			Attrs: Attrs{
				"source":     in.LayerID(),
				"horizontal": horizontal,
				"a_uv": []float32{
					0, 0,
					0, 1,
					1, 1,
					1, 0,
				},
			},
		})
		horizontal = !horizontal
		in = out
	}
}

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

  if (horizontal) {
    for(int i = 1; i < 5; ++i) {
      result += texture(source, v_uv + vec2(tex_offset.x * i, 0.0)).rgba * weight[i];
      result += texture(source, v_uv - vec2(tex_offset.x * i, 0.0)).rgba * weight[i];
    }
  } else {
    for(int i = 1; i < 5; ++i) {
      result += texture(source, v_uv + vec2(0.0, tex_offset.y * i)).rgba * weight[i];
      result += texture(source, v_uv - vec2(0.0, tex_offset.y * i)).rgba * weight[i];
    }
  }
	color = result;
}
`
