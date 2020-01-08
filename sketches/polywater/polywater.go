package main

import (
	"io/ioutil"

	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

const (
	layers             = 1
	N                  = 400
	opacity    float32 = 0.5
	turbulence float32 = .14
	blur       float32 = 0.28
	speed      float32 = 0
	time       float32 = 0
	S          float32 = .3
)

func Ink(doc *Doc) {
	Clear(doc, White)

	/*
			noiseLayer := doc.NewLayer()
			noiseLayer.AddShader(&Shader{
				Vert: DefaultVert,
				Frag: loadShader("sketches/polywater/snoise.frag"),
				Mesh: Fullscreen,
				Attrs: Attrs{
					"u_speed": speed,
					"u_time":  time,
				},
			})
		log.Printf("noise layer ID: %d", noiseLayer.LayerID())
	*/

	//polys := doc.NewLayer()

	for j := 0; j < layers; j++ {

		boxy := .69 - float32(j)*0.01

		rec := RectWH(0.5, 0.5)

		pos := make([]XY, N)
		rot := make([]float32, N)
		size := make([]XY, N)

		for i := 0; i < N; i++ {
			size[i] = XY{
				X: rand.Range(S, S+.30) + 0.10*float32(j),
				Y: rand.Range(S, S+.00) + 0.00*float32(j),
			}
			pos[i] = XY{
				X: rand.Range(.00, 1.0),
				Y: rand.Range(.00, .005) + 0.02*float32(j),
			}
			rot[i] = rand.Angle()
		}

		doc.AddShader(&Shader{
			Name:      "Polys",
			Instances: N,
			Vert:      loadShader("paint.vert"),
			Frag:      loadShader("paint.frag"),
			Mesh:      rec,
			Divisors: map[string]int{
				"a_pos":  1,
				"a_rot":  1,
				"a_size": 1,
			},
			Attrs: Attrs{
				"a_pos":        pos,
				"a_rot":        rot,
				"a_size":       size,
				"u_turbulence": turbulence,
				//"u_noise":      noiseLayer.LayerID(),
				"u_boxy":    boxy,
				"u_blur":    blur,
				"u_opacity": opacity,
				"u_speed":   speed,
				"u_time":    time,
				"u_color":   RGBA{0, 0, 1, 1},
			},
		})

		/*
			doc.AddShader(&gfx.Shader{
				Name: "PolyMask",
				Vert: DefaultVert,
				Frag: loadShader("sketches/polywater/mask.frag"),
				Mesh: Fullscreen,
				Attrs: Attrs{
					"u_color":   RGBA{0, .3, 0, 1},
					"u_opacity": 1 - float32(j)*0.25,
					"u_mask":    polys.LayerID(),
				},
			})
		*/
	}
}

func loadShader(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(b)
}
