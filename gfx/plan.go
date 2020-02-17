package gfx

import (
	"log"

	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/internal/glsl"
	"github.com/buchanae/ink/render"
)

func (doc *Doc) Plan() render.Plan {

	plan := render.Plan{
		RootLayer: doc.LayerID(),
		Shaders:   map[int]render.Shader{},
		//Images:    map[int]image.Image{},
	}

	nextSID := 1
	uniqShaders := map[render.Shader]int{}

	/*
		for k, v := range doc.Images {
			plan.Images[k] = v
		}
	*/

	for _, op := range doc.Ops {

		s := op.Shader

		shader := render.Shader{
			Vert: s.Vert,
			Frag: s.Frag,
		}
		sid, ok := uniqShaders[shader]
		if !ok {
			sid = nextSID
			uniqShaders[shader] = sid
			plan.Shaders[sid] = shader
			nextSID++
		}

		rs := render.Pass{
			Name:      s.Name,
			Shader:    sid,
			Layer:     op.LayerID,
			Instances: s.Instances,
			Vertices:  len(s.Mesh.Verts),
			Faces: render.Faces{
				Offset: len(plan.FaceData),
				Count:  len(s.Mesh.Faces),
			},
			Uniforms: map[string]interface{}{},
		}

		if rs.Vertices == 0 {
			// TODO move to a "warnings" or "errors" list
			log.Print("empty verts")
		}

		rs.Attrs = append(rs.Attrs, render.Attr{
			Name:   "a_vert",
			Offset: len(plan.AttrData),
			Count:  len(s.Mesh.Verts) * 2,
		})
		for _, v := range s.Mesh.Verts {
			plan.AttrData = append(plan.AttrData, v.X, v.Y)
		}

		rs.Attrs = append(rs.Attrs, render.Attr{
			Name:   "a_uv",
			Offset: len(plan.AttrData),
			Count:  len(s.Mesh.UV) * 2,
		})
		for _, v := range s.Mesh.UV {
			plan.AttrData = append(plan.AttrData, v.X, v.Y)
		}

		for _, f := range s.Mesh.Faces {
			plan.FaceData = append(plan.FaceData,
				uint32(f[0]), uint32(f[1]), uint32(f[2]),
			)
		}

		meta := glsl.Inspect(s.Vert, s.Frag)

		for _, uni := range meta.Uniforms {
			val := convertUniform(s.Attrs[uni.Name])
			if val == nil {
				// TODO move to a "warnings" or "errors" list
				log.Printf("missing uniform: %s", uni.Name)
				continue
			}
			rs.Uniforms[uni.Name] = val
		}

		for _, attr := range meta.Attributes {

			if attr.Name == "a_vert" || attr.Name == "a_uv" {
				// These special cases are set above.
				continue
			}

			val, ok := s.Attrs[attr.Name]
			if !ok {
				//log.Printf("missing attribute: %s", attr.Name)
				continue
			}

			data := convertAttr(val, len(s.Mesh.Verts))

			if len(data) == 0 {
				log.Printf("empty attribute: %s", attr.Name)
				continue
			}

			rs.Attrs = append(rs.Attrs, render.Attr{
				Name:    attr.Name,
				Offset:  len(plan.AttrData),
				Count:   len(data),
				Divisor: s.Divisors[attr.Name],
			})
			plan.AttrData = append(plan.AttrData, data...)
		}

		plan.Passes = append(plan.Passes, rs)
	}

	return plan
}

func convertAttr(val interface{}, verts int) []float32 {

	if val == nil {
		return nil
	}

	switch z := val.(type) {

	case []float32:
		return z

	case float32:
		data := make([]float32, verts)
		for i := range data {
			data[i] = z
		}
		return data

	case float64:
		data := make([]float32, verts)
		for i := range data {
			data[i] = float32(z)
		}
		return data

	case []color.RGBA:
		data := make([]float32, 0, len(z)*4)
		for _, v := range z {
			data = append(data, v.R, v.G, v.B, v.A)
		}
		return data

	case color.RGBA:
		data := make([]float32, 0, verts*4)
		for i := 0; i < verts; i++ {
			data = append(data, z.R, z.G, z.B, z.A)
		}
		return data

	case []dd.XY:
		data := make([]float32, 0, len(z)*2)
		for _, v := range z {
			data = append(data, v.X, v.Y)
		}
		return data

	case dd.XY:
		data := make([]float32, 0, verts*2)
		for i := 0; i < verts; i++ {
			data = append(data, z.X, z.Y)
		}
		return data

	default:
		log.Printf("unsupported attribute type %T", z)
		return nil
	}
}

func convertUniform(v interface{}) interface{} {
	switch z := v.(type) {
	case dd.XY:
		return [2]float32{z.X, z.Y}
	case color.RGBA:
		return [4]float32{z.R, z.G, z.B, z.A}
	case Image:
		return z.ID
	default:
		return v
	}
}
