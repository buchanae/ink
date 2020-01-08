// +build !sendonly

package app

import (
	"image"
	"log"

	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/internal/glsl"
	"github.com/buchanae/ink/render"
)

func buildPlan(doc *Doc) render.Plan {

	plan := render.Plan{
		Images: map[int]image.Image{},
	}

	for k, v := range doc.Images {
		plan.Images[k] = v
	}

	for _, op := range doc.Ops {

		s := op.Shader
		if s.Mesh == nil {
			continue
		}

		mesh := s.Mesh.Mesh()

		rs := &render.Shader{
			Name:      s.Name,
			Vert:      s.Vert,
			Frag:      s.Frag,
			Layer:     op.LayerID,
			Vertices:  len(mesh.Verts),
			Instances: s.Instances,
			Faces:     make([]uint32, 0, len(mesh.Faces)*3),
			Uniforms:  map[string]interface{}{},
			Attrs: map[string]render.Attr{
				"a_vert": {
					Value: mesh.Verts,
					Size:  len(mesh.Verts) * 2 * 4,
				},
				"a_uv": {
					Value: mesh.UV,
					Size:  len(mesh.UV) * 2 * 4,
				},
			},
		}

		for _, f := range mesh.Faces {
			rs.Faces = append(rs.Faces,
				uint32(f[0]), uint32(f[1]), uint32(f[2]),
			)
		}

		meta := glsl.Inspect(s.Vert, s.Frag)

		for _, uni := range meta.Uniforms {
			val := convertUniform(s.Attrs[uni.Name])
			rs.Uniforms[uni.Name] = val
		}

		for _, attr := range meta.Attributes {
			if _, ok := rs.Attrs[attr.Name]; ok {
				// skip attrs like a_vert and a_uv
				// which are set above.
				continue
			}
			val := convertAttr(s, attr.Name, len(mesh.Verts))
			rs.Attrs[attr.Name] = val
		}

		plan.Shaders = append(plan.Shaders, rs)
	}

	return plan
}

func convertAttr(s *gfx.Shader, name string, verts int) render.Attr {
	attr := render.Attr{
		Divisor: s.Divisors[name],
	}

	val, ok := s.Attrs[name]
	if !ok {
		return attr
	}

	switch z := val.(type) {

	case []float32:
		attr.Value = z
		attr.Size = len(z) * 4

	case float32:
		data := make([]float32, verts)
		for i := range data {
			data[i] = z
		}
		attr.Value = data
		attr.Size = len(data) * 4

	case float64:
		data := make([]float32, verts)
		for i := range data {
			data[i] = float32(z)
		}
		attr.Value = data
		attr.Size = len(data) * 4

	case []dd.XY:
		attr.Value = z
		attr.Size = len(z) * 2 * 4

	case []color.RGBA:
		attr.Value = z
		attr.Size = len(z) * 4 * 4

	case color.RGBA:
		data := make([]color.RGBA, verts)
		for i := range data {
			data[i] = z
		}
		attr.Value = data
		attr.Size = len(data) * 4 * 4

	case dd.XY:
		data := make([]dd.XY, verts)
		for i := range data {
			data[i] = z
		}
		attr.Value = data
		attr.Size = len(data) * 2 * 4

	default:
		log.Printf("unsupported attribute type %T for %s", z, name)
	}

	return attr
}

func convertUniform(v interface{}) interface{} {
	switch z := v.(type) {
	case dd.XY:
		return [2]float32{z.X, z.Y}
	case color.RGBA:
		return [4]float32{z.R, z.G, z.B, z.A}
	case gfx.Image:
		return z.ID
	default:
		return v
	}
}
