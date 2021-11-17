package client

import (
	"image"
	"log"

	"github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/internal/glsl"
	"github.com/buchanae/ink/render"
	"github.com/buchanae/ink/trac"
)

type batcher struct {
	shaders *uniqueShaders
	batches []*batch
	faces   int
	attrs   int
}

type batch struct {
	pass   render.Pass
	meshes []dd.Mesh
	attrs  map[string]*batchAttr
}

type batchAttr struct {
	divisor int
	vals    [][]float32
}

func (b *batcher) Add(op Op) {

	shader := op.Shader
	mesh := shader.Mesh

	if len(mesh.Verts) == 0 {
		// TODO move to a "warnings" or "errors" list
		log.Print("empty verts")
		return
	}

	pass := render.Pass{
		Name: shader.Name,
		Shader: b.shaders.Add(render.Shader{
			Vert: shader.Vert,
			Frag: shader.Frag,
		}),
		Layer:     op.LayerID,
		Instances: shader.Instances,
		Vertices:  len(mesh.Verts),
		Faces: render.Faces{
			Offset: b.faces,
			Count:  len(mesh.Faces),
		},
		Uniforms: map[string]interface{}{},
	}
	b2 := &batch{
		pass:   pass,
		meshes: []dd.Mesh{mesh},
		attrs:  map[string]*batchAttr{},
	}

	b.faces += len(mesh.Faces) * 3
	b.attrs += len(mesh.Verts) * 2
	b.attrs += len(mesh.UV) * 2

	meta := glsl.Inspect(shader.Vert, shader.Frag)

	for _, uni := range meta.Uniforms {
		val := convertUniform(shader.Attrs[uni.Name])
		if val == nil {
			// TODO move to a "warnings" or "errors" list
			log.Printf("missing uniform: %s", uni.Name)
			continue
		}
		pass.Uniforms[uni.Name] = val
	}

	for _, attr := range meta.Attributes {

		if attr.Name == "a_vert" || attr.Name == "a_uv" {
			// These special cases are set above.
			continue
		}

		val, ok := shader.Attrs[attr.Name]
		if !ok {
			//log.Printf("missing attribute: %s", attr.Name)
			continue
		}

		data := convertAttr(val, len(mesh.Verts))

		if len(data) == 0 {
			log.Printf("empty attribute: %s", attr.Name)
			continue
		}

		b.attrs += len(data)

		b2.attrs[attr.Name] = &batchAttr{
			divisor: shader.Divisors[attr.Name],
			vals:    [][]float32{data},
		}
	}

	b.batches = append(b.batches, b2)
}

func buildPlan(doc *Doc) render.Plan {

	plan := render.Plan{
		RootLayer: doc.LayerID(),
		Shaders:   map[int]render.Shader{},
		Images:    map[int]image.Image{},
	}

	for k, v := range doc.Images {
		plan.Images[k] = v
	}

	build := batcher{
		shaders: &uniqueShaders{
			shaders: plan.Shaders,
			unique:  map[render.Shader]int{},
		},
	}

	trac.Log("add ops %d", len(doc.Ops))
	for _, op := range doc.Ops {
		build.Add(op)
	}

	trac.Log("alloc data")
	plan.FaceData = make([]uint32, 0, build.faces)
	plan.AttrData = make([]float32, 0, build.attrs)

	batches := mergeBatches(build.batches)

	for _, batch := range batches {

		offset := len(plan.AttrData)
		count := 0
		for _, mesh := range batch.meshes {
			for _, v := range mesh.Verts {
				plan.AttrData = append(plan.AttrData, v.X, v.Y)
				count += 2
			}
		}
		batch.pass.Attrs = append(batch.pass.Attrs, render.Attr{
			Name:   "a_vert",
			Offset: offset,
			Count:  count,
		})

		offset = len(plan.AttrData)
		count = 0
		for _, mesh := range batch.meshes {
			for _, v := range mesh.UV {
				plan.AttrData = append(plan.AttrData, v.X, v.Y)
				count += 2
			}
		}
		batch.pass.Attrs = append(batch.pass.Attrs, render.Attr{
			Name:   "a_uv",
			Offset: offset,
			Count:  count,
		})

		var faceOffset uint32
		for _, mesh := range batch.meshes {
			for _, face := range mesh.Faces {
				plan.FaceData = append(plan.FaceData,
					uint32(face[0])+faceOffset,
					uint32(face[1])+faceOffset,
					uint32(face[2])+faceOffset,
				)
			}
			faceOffset += uint32(len(mesh.Verts))
		}

		for key, attr := range batch.attrs {
			offset = len(plan.AttrData)
			count = 0

			for _, val := range attr.vals {
				plan.AttrData = append(plan.AttrData, val...)
				count += len(val)
			}

			batch.pass.Attrs = append(batch.pass.Attrs, render.Attr{
				Name:    key,
				Offset:  offset,
				Count:   count,
				Divisor: attr.divisor,
			})
		}
		plan.Passes = append(plan.Passes, batch.pass)
	}

	return plan
}

type uniqueShaders struct {
	shaders map[int]render.Shader
	unique  map[render.Shader]int
	nextID  int
}

func (us *uniqueShaders) Add(s render.Shader) int {
	id, ok := us.unique[s]
	if !ok {
		us.nextID++
		id = us.nextID
		us.unique[s] = id
		us.shaders[id] = s
	}
	return id
}
