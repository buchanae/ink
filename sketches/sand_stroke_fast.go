package main

import (
	"image"
	colorlib "image/color"

	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/math"
	"github.com/buchanae/ink/rand"
)

const (
	Lines   = 3
	Strokes = 10
	N       = 300
	W       = 1 / float32(N)

	Amp = 1000

	// min/max y-axis starting position
	MinY = 0.1
	MaxY = 0.9

	D           = 0.0040
	MaxD        = 0.2
	AlphaOffset = -0.3
)

func Ink(doc *Doc) {
	rand.SeedNow()
	palette := rand.Palette()

	Clear(doc, White)

	lines := make([]float32, Lines)
	for i := range lines {
		lines[i] = rand.Range(MinY, MaxY)
	}

	for j := 0; j < Strokes; j++ {

		y := lines[rand.Intn(len(lines))]

		heights := make([]float32, N)
		h := rand.Range(0.01, 0.1)
		for i := range heights {
			h += rand.Range(-D, D)
			h = math.Clamp(h, 0, MaxD)
			heights[i] = h
		}

		img := makeHeightMap(heights)
		heightmap := doc.NewImage(img)

		s := &Shader{
			Name: "Stroke",
			Vert: DefaultVert,
			Frag: Frag,
			Mesh: Rect{
				XY{0, y - MaxD},
				XY{1, y + MaxD},
			},
			Attrs: Attrs{
				"u_heightmap":    heightmap,
				"u_color":        rand.Color(palette),
				"u_alpha_offset": float32(AlphaOffset),
			},
		}
		s.Draw(doc)
	}
}

const Frag = `
#version 330 core

uniform vec4 u_color;
uniform sampler2D u_heightmap;
uniform float u_alpha_offset;

in vec2 v_uv;
out vec4 color;

void main() {
	float h = texture(u_heightmap, v_uv).r;
	float d = abs(v_uv.y - 0.5);
	float a = step(d, h) * (1-h*2) + u_alpha_offset;
	color = vec4(u_color.rgb, a);
}
`

// TODO want to easily create a texture without
//      involving the "image" library
func makeHeightMap(heights []float32) *image.Gray {
	r := image.Rect(0, 0, len(heights), 1)
	img := image.NewGray(r)

	for i, h := range heights {
		img.SetGray(i, 0, colorlib.Gray{
			Y: uint8(h * Amp),
		})
	}
	return img
}
