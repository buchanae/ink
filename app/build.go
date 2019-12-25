// +build !sendonly

package app

import (
	"fmt"
	"log"

	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/render"
)

type builder struct {
	layers   map[int]*render.Layer
	renderer *render.Renderer
}

func (b *builder) build(doc *gfx.Layer) {
	if b.layers == nil {
		b.layers = map[int]*render.Layer{}
	}

	for _, value := range doc.Values {
		switch z := value.(type) {
		case *gfx.Layer:
			b.build(z)
		case gfx.Layer:
			b.build(&z)
		case *gfx.Shader:
			err := b.buildShader(doc.ID, z)
			if err != nil {
				log.Printf("error: %v", err)
			}
		case gfx.Shader:
			err := b.buildShader(doc.ID, &z)
			if err != nil {
				log.Printf("error: %v", err)
			}
		default:
			log.Printf("error: unknown layer type: %T", value)
		}
	}
}

func (b *builder) buildShader(id int, shader *gfx.Shader) error {

	mesh := shader.Mesh.Mesh()
	verts := mesh.Verts

	if verts == nil {
		return fmt.Errorf("empty verts")
	}

	rl, err := b.renderer.NewLayer(render.Shader{
		ID:          id,
		Name:        shader.Name,
		Vert:        shader.Vert,
		Frag:        shader.Frag,
		VertexCount: len(verts),
	})
	if err != nil {
		return err
	}
	b.layers[id] = rl

	rl.SetAttr("a_vert", verts, len(verts)*2*4)

	for _, name := range rl.AttrNames() {
		val, ok := shader.Attrs[name]
		if !ok {
			//log.Printf("missing attribute: %s", name)
			continue
		}

		switch z := val.(type) {

		case []float32:
			rl.SetAttr(name, z, len(z)*4)

		case float32:
			data := make([]float32, len(verts))
			for i := range data {
				data[i] = z
			}
			size := len(data) * 4
			rl.SetAttr(name, data, size)

		case float64:
			data := make([]float32, len(verts))
			for i := range data {
				data[i] = float32(z)
			}
			size := len(data) * 4
			rl.SetAttr(name, data, size)

		case []dd.XY:
			size := len(z) * 2 * 4
			rl.SetAttr(name, z, size)

		case []color.RGBA:
			size := len(z) * 4 * 4
			rl.SetAttr(name, z, size)

		case color.RGBA:
			data := make([]color.RGBA, len(verts))
			for i := range data {
				data[i] = z
			}
			size := len(data) * 4 * 4
			rl.SetAttr(name, data, size)

		case dd.XY:
			data := make([]dd.XY, len(verts))
			for i := range data {
				data[i] = z
			}
			size := len(data) * 2 * 4
			rl.SetAttr(name, data, size)

		default:
			log.Printf("error: unsupported attribute value type %T: %v", z, z)
			continue
		}
	}

	faces := make([]uint32, 0, len(mesh.Faces))
	for _, f := range mesh.Faces {
		faces = append(faces,
			uint32(f[0]),
			uint32(f[1]),
			uint32(f[2]),
		)
	}
	rl.SetFaces(faces)

	for _, name := range rl.UniformNames() {
		v, ok := shader.Attrs[name]
		if !ok {
			log.Printf("missing uniform: %s", name)
			continue
		}

		switch z := v.(type) {
		case dd.XY:
			rl.SetUniform(name, [2]float32{z.X, z.Y})
		case color.RGBA:
			rl.SetUniform(name, [4]float32{z.R, z.G, z.B, z.A})
		case gfx.Layer:
			target, ok := b.layers[z.ID]
			if !ok {
				log.Printf("missing layer uniform: %d", z.ID)
				continue
			}
			rl.SetUniform(name, target)
		case *gfx.Layer:
			target, ok := b.layers[z.ID]
			if !ok {
				log.Printf("missing layer uniform: %d", z.ID)
				continue
			}
			rl.SetUniform(name, target)
		default:
			rl.SetUniform(name, v)
		}
	}
	return nil
}
