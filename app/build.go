package app

import (
	"fmt"
	"image"
	"log"

	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/render"
)

type builder struct {
	renderer *render.Renderer
}

func (b *builder) build(nodes []gfx.Node) {
	for _, node := range nodes {
		switch z := node.Op.(type) {
		case *gfx.Shader:
			err := b.buildShader(node.LayerID, z)
			if err != nil {
				log.Printf("error: %v", err)
			}
		case gfx.Shader:
			err := b.buildShader(node.LayerID, &z)
			if err != nil {
				log.Printf("error: %v", err)
			}
		case image.Gray:
			b.buildImage(node.LayerID, &z)
		case image.RGBA:
			b.buildImage(node.LayerID, &z)
		default:
			log.Printf("error: unknown node type: %T", node.Op)
		}
	}
}

func (b *builder) buildImage(id int, img image.Image) {
	b.renderer.NewImage(id, img)
}

func (b *builder) buildShader(id int, shader *gfx.Shader) error {

	mesh := shader.Mesh.Mesh()
	verts := mesh.Verts

	if verts == nil {
		return fmt.Errorf("empty verts")
	}

	rl, err := b.renderer.NewLayer(render.Shader{
		ID:            id,
		Name:          shader.Name,
		Vert:          shader.Vert,
		Frag:          shader.Frag,
		VertexCount:   len(verts),
		InstanceCount: shader.InstanceCount,
	})
	if err != nil {
		return err
	}

	rl.SetAttr("a_vert", verts, len(verts)*2*4, 0)
	rl.SetAttr("a_uv", mesh.UV, len(mesh.UV)*2*4, 0)

	for _, name := range rl.AttrNames() {
		val, ok := shader.Attrs[name]
		if !ok {
			//log.Printf("missing attribute: %s", name)
			continue
		}

		divisor := shader.Divisors[name]

		switch z := val.(type) {

		case []float32:
			rl.SetAttr(name, z, len(z)*4, divisor)

		case float32:
			data := make([]float32, len(verts))
			for i := range data {
				data[i] = z
			}
			size := len(data) * 4
			rl.SetAttr(name, data, size, divisor)

		case float64:
			data := make([]float32, len(verts))
			for i := range data {
				data[i] = float32(z)
			}
			size := len(data) * 4
			rl.SetAttr(name, data, size, divisor)

		case []dd.XY:
			size := len(z) * 2 * 4
			rl.SetAttr(name, z, size, divisor)

		case []color.RGBA:
			size := len(z) * 4 * 4
			rl.SetAttr(name, z, size, divisor)

		case color.RGBA:
			data := make([]color.RGBA, len(verts))
			for i := range data {
				data[i] = z
			}
			size := len(data) * 4 * 4
			rl.SetAttr(name, data, size, divisor)

		case dd.XY:
			data := make([]dd.XY, len(verts))
			for i := range data {
				data[i] = z
			}
			size := len(data) * 2 * 4
			rl.SetAttr(name, data, size, divisor)

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
		default:
			rl.SetUniform(name, v)
		}
	}
	return nil
}
