package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

const (
	layers  = 15
	N       = 280
	opacity = 0.8
	blur    = 0.98
	S       = .3
)

func Ink(doc *gfx.Doc) {

	noiseLayer := doc.NewLayer()
	noiseLayer.AddShader(&Shader{
		Vert: DefaultVert,
		Frag: loadShader("sketches/polywater/snoise.frag"),
		Mesh: Fullscreen,
		Attrs: Attrs{
			"u_speed": 0,
			"u_time":  0,
			"a_uv": []float32{
				0, 0,
				0, 1,
				1, 1,
				1, 0,
			},
		},
	})

	polys := doc.NewLayer()

	for j := 0; j < layers; j++ {

		boxy := .69 - float32(j)*0.08

		// TODO clones (instancing)
		rec := gfx.Rect(0, 0, 1, 1)

		pos := make([]XY, N)
		rot := make(float32, N)
		size := make([]XY, N)

		for i := 0; i < N; i++ {
			size[i] = XY{
				X: rand.Range(S, S+.10) + 0.10*float32(j),
				Y: rand.Range(S, S+.00) + 0.00*float32(j),
			}
			pos[i] = XY{
				X: rand.Range(.00, 1.0),
				Y: rand.Range(.02, .1) + 0.02*float32(j),
			}
			rot[i] = rand.Range(0, 10.3)
		}

		polys.AddShader(&Shader{
			Name: "Polys",
			Vert: loadShader("sketches/polywater/paint.vert"),
			Frag: loadShader("sketches/polywater/paint.frag"),
			Mesh: rec,
			Attrs: Attrs{
				"a_pos":        pos,
				"a_rot":        rot,
				"a_size":       size,
				"u_turbulence": .24,
				"u_noise":      noise,
				"u_boxy":       boxy,
				"u_blur":       blur,
				"u_opacity":    opacity,
				"u_speed":      0,
				"u_time":       0,
				"u_color":      RGBA{1, 1, 1, 1},
			},
		})

		doc.AddShader(&gfx.Shader{
			Name: "PolyMask",
			Vert: DefaultVert,
			Frag: loadShader("sketches/polywater/mask.frag"),
			Mesh: Fullscreen,
			Attrs: Attrs{
				"u_color":   RGBA{0, .3, 0, 1},
				"u_opacity": 1 - float32(j)*0.25,
				"u_mask":    polys.LayerID,
				"a_uv": []float32{
					0, 0,
					0, 1,
					1, 1,
					1, 0,
				},
			},
		})
	}
}

func loadShader(path string) string {
	b, err := io.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(b)
}
