package render

import (
	"log"
)

type build struct {
	tracer

	passes []*pass
	faces  []uint32
	// Total number of bytes needed to store all attributes.
	attrBytes int
}

type pass struct {
	name          string
	prog          program
	uniforms      map[string]interface{}
	layer         int
	vao           glVAO
	bindings      []binding
	vertexCount   int
	faceOffset    int
	faceCount     int
	instanceCount int
}

type binding struct {
	attr    attribute
	values  []bindingVal
	divisor int
}

type bindingVal struct {
	value interface{}
	size  int
}

func (pb *build) build(plan Plan) {
	pb.trace("start build")

	pb.passes = make([]*pass, 0, len(plan.Shaders))
	pb.faces = make([]uint32, 0, 500)

	if len(plan.Shaders) == 0 {
		return
	}

	pb.trace("add shaders")
	for _, s := range plan.Shaders {
		pb.addShader(s)
	}

	pb.batch()
	pb.trace("end build")
}

func (pb *build) addShader(shader *Shader) {

	prog, err := compile(shaderOpt{
		shader.Vert, shader.Frag, shader.Geom,
	})
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	uniforms := map[string]interface{}{}
	for _, uni := range prog.uniforms {
		val := shader.Uniforms[uni.Name]
		if val == nil {
			// TODO move to a "warnings" or "errors" list
			log.Printf("missing uniform: %s", uni.Name)
			continue
		}
		uniforms[uni.Name] = val
	}

	if shader.Vertices == 0 {
		log.Print("empty verts")
	}

	p := &pass{
		prog:          prog,
		layer:         shader.Layer,
		name:          shader.Name,
		vertexCount:   shader.Vertices,
		instanceCount: shader.Instances,
		uniforms:      uniforms,
		faceOffset:    len(pb.faces),
		faceCount:     len(shader.Faces),
	}
	pb.passes = append(pb.passes, p)
	pb.faces = append(pb.faces, shader.Faces...)

	for _, attr := range prog.attributes {
		desc, ok := shader.Attrs[attr.Name]

		if !ok || desc.Value == nil || desc.Size == 0 {
			continue
		}

		p.bindings = append(p.bindings, binding{
			attr:    attr,
			values:  []bindingVal{{desc.Value, desc.Size}},
			divisor: desc.Divisor,
		})
		pb.attrBytes += desc.Size
	}
}

func (pb *build) batch() {
	pb.trace("batch")

	var batched []*pass
	var last *pass

	for i, p := range pb.passes {
		if i == 0 {
			last = p
			continue
		}

		if pb.mergeable(last, p) {
			pb.merge(last, p)
		} else if last != nil {
			batched = append(batched, last)
			last = p
		} else {
			last = p
		}
	}
	if last != nil {
		batched = append(batched, last)
	}
	pb.trace("  merged passes %d to %d", len(pb.passes), len(batched))
	pb.passes = batched
}

func (pb *build) mergeable(a, b *pass) bool {
	if a.prog.id != b.prog.id {
		pb.trace("program IDs differ")
		return false
	}
	if len(a.bindings) != len(b.bindings) {
		pb.trace("bindings differ", len(a.bindings), len(b.bindings))
		return false
	}
	if len(a.uniforms) != len(b.uniforms) {
		pb.trace("uniforms differ")
		return false
	}
	if a.layer != b.layer {
		pb.trace("layers differ")
		return false
	}
	for i := range b.bindings {
		if a.bindings[i].attr != b.bindings[i].attr {
			return false
		}
		if a.bindings[i].divisor != b.bindings[i].divisor {
			return false
		}
	}
	for k, v := range a.uniforms {
		c, ok := b.uniforms[k]
		if !ok {
			return false
		}
		if c != v {
			return false
		}
	}
	return true
}

func (pb *build) merge(a, b *pass) {

	// merge faces
	vc := uint32(a.vertexCount)
	for i := 0; i < b.faceCount; i++ {
		pb.faces[b.faceOffset+i] += vc
	}
	a.faceCount += b.faceCount
	a.vertexCount += b.vertexCount

	for i := range a.bindings {
		for _, bv := range b.bindings[i].values {
			a.bindings[i].values = append(a.bindings[i].values, bv)
		}
	}
}
