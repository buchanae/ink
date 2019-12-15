package app

import (
	"errors"
	"log"

	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/render"
)

var builtins = map[string][]byte{
	"!default.vert": []byte(gfx.DefaultVert),
	"!default.frag": []byte(gfx.DefaultFrag),
	"!fill.frag":    []byte(gfx.FillFrag),
	"!copy.frag":    []byte(gfx.CopyFrag),
	//"noise.frag":    gfx.NoiseFrag,
	//"noiselib.glsl": gfx.NoiseLib,
	//"blur.frag":     gfx.BlurFrag,
	//"gradient.frag": gfx.GradientFrag,
	//"mask.frag":     gfx.MaskFrag,
}

// TODO redo asset loader
// figure out a good way to allow gfx to define builtins on the doc
// without transfering on every doc render. need retained doc.
func asset(name string) ([]byte, error) {
	b, ok := builtins[name]
	if !ok {
		return nil, errors.New("unknown asset: " + name)
	}
	return b, nil
}

// build builds the doc into renderer layers
func build(doc *gfx.Doc, renderer *render.Renderer) error {

	for _, layer := range doc.Layers {
		shader, ok := layer.Value.(*gfx.Shader)
		if !ok {
			log.Printf("skipping non-shader layer: %T", layer.Value)
			// TODO
			continue
		}

		vert, err := asset(shader.Vert)
		if err != nil {
			return err
		}

		frag, err := asset(shader.Frag)
		if err != nil {
			return err
		}

		mesh := shader.Mesh
		verts := mesh.Verts
		faces := mesh.Faces

		if verts == nil {
			log.Println("empty verts")
			continue
		}

		rl := renderer.NewLayer(render.Shader{
			Vert: string(vert),
			Frag: string(frag),
		})
		rl.Name(shader.Name)

		rl.VertexCount(len(verts))
		rl.UnsafeAttr("a_vert", verts, len(verts)*2*4)

		for key, val := range shader.Attrs.Data {
			switch z := val.(type) {

			case []float32:
				rl.FloatAttr(key, z)

			case float32:
				data := make([]float32, len(verts))
				for i := range data {
					data[i] = z
				}
				rl.FloatAttr(key, data)

			case float64:
				data := make([]float32, len(verts))
				for i := range data {
					data[i] = float32(z)
				}
				rl.FloatAttr(key, data)

			case []dd.XY:
				size := len(z) * 2 * 4
				rl.UnsafeAttr(key, z, size)

			case []color.RGBA:
				size := len(z) * 4 * 4
				rl.UnsafeAttr(key, z, size)

			case color.RGBA:
				data := make([]color.RGBA, len(verts))
				for i := range data {
					data[i] = z
				}
				size := len(data) * 4 * 4
				rl.UnsafeAttr(key, data, size)

			case dd.XY:
				data := make([]dd.XY, len(verts))
				for i := range data {
					data[i] = z
				}
				size := len(data) * 2 * 4
				rl.UnsafeAttr(key, data, size)

			default:
				log.Printf("error: unsupported attribute value type %T: %v", z, z)
				continue
			}
		}

		faceData := make([]uint32, 0, len(faces))
		for _, f := range faces {
			faceData = append(faceData,
				uint32(f[0]),
				uint32(f[1]),
				uint32(f[2]),
			)
		}
		rl.Faces(faceData)

		for k, v := range shader.Uniforms {
			switch z := v.(type) {
			case dd.XY:
				rl.Uniform(k, [2]float32{z.X, z.Y})
			case color.RGBA:
				rl.Uniform(k, [4]float32{z.R, z.G, z.B, z.A})
			default:
				rl.Uniform(k, v)
			}
		}
	}
	return nil
}

// Determine the size (in bytes) of the data.
/*
	var size int
	switch z := val.(type) {
	}

	p.bindings = append(p.bindings, binding{
		attr:  attr,
		value: val,
		size:  size,
	})
*/
