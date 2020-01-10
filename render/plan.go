package render

import (
	"image"
	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// TODO is a plan resolution independent?
//      i.e. vertex data only, no texture specs.
type Plan struct {
	Shaders []*Shader
	Images  map[int]image.Image
}

type Shader struct {
	Name             string
	Vert, Frag, Geom string
	Output           string
	Layer            int
	Vertices         int
	Instances        int
	Faces            []uint32
	Uniforms         map[string]interface{}
	Attrs            map[string]Attr
}

type Attr struct {
	Value   interface{}
	Size    int
	Divisor int
}

type build struct {
	tracer

	passes []*pass
	faces  []uint32
	// Total number of bytes needed to store all attributes.
	attrBytes int

	// opengl buffer IDs
	attrBufID uint32
	faceBufID uint32
	vaos      []uint32
}

type pass struct {
	name          string
	prog          compiled
	uniforms      map[string]interface{}
	layer         int
	vao           uint32
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
	pb.upload()

	pb.trace("end build")
}

func (pb *build) addShader(shader *Shader) {

	prog, err := compile(shaderOpt{
		shader.Vert, shader.Frag, shader.Geom, shader.Output,
	})
	if err != nil {
		log.Printf("error: compiling shader: %v", err)
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

func (pb *build) cleanup() {
	if pb.attrBufID != 0 {
		glDeleteBuffers(1, &pb.attrBufID)
	}
	if pb.faceBufID != 0 {
		glDeleteBuffers(1, &pb.faceBufID)
	}
	if len(pb.vaos) > 0 {
		glDeleteVertexArrays(int32(len(pb.vaos)), &pb.vaos[0])
	}
}

func (pb *build) uploadFaces() {
	// Element buffers are used for indexed rendering.
	glGenBuffers(1, &pb.faceBufID)
	glBindBuffer(gl.ELEMENT_ARRAY_BUFFER, pb.faceBufID)

	glBufferData(
		gl.ELEMENT_ARRAY_BUFFER,
		len(pb.faces)*4, // 4 bytes per index (uint32)
		glPtr(pb.faces),
		gl.STATIC_DRAW,
	)
}

func (pb *build) upload() {
	pb.trace("upload")

	// upload faces (vertex index)
	pb.uploadFaces()

	// The data from all attributes is stored in one large buffer.
	// A "binding" describes the slice of the buffer that holds data
	// for a single attribute.
	glGenBuffers(1, &pb.attrBufID)
	glBindBuffer(gl.ARRAY_BUFFER, pb.attrBufID)
	glBufferData(gl.ARRAY_BUFFER, pb.attrBytes, nil, gl.STATIC_DRAW)

	// Each pass has one VAO, which stores the configuration of all its
	// attributes: location of the data in the buffer, enabled/disable,
	// types, divisors, etc.
	pb.vaos = make([]uint32, len(pb.passes))
	glGenVertexArrays(int32(len(pb.passes)), &pb.vaos[0])

	offset := 0
	for i, p := range pb.passes {
		p.vao = pb.vaos[i]
		glBindVertexArray(p.vao)
		glBindBuffer(gl.ELEMENT_ARRAY_BUFFER, pb.faceBufID)

		for _, b := range p.bindings {

			glEnableVertexAttribArray(b.attr.Loc)
			glVertexAttribPointer(
				b.attr.Loc,
				b.attr.Components,
				b.attr.Datatype,
				false, // normalized
				0,     // stride
				glPtrOffset(offset),
			)
			glVertexAttribDivisor(b.attr.Loc, uint32(b.divisor))

			for _, val := range b.values {
				if val.size == 0 {
					// opengl will panic if it tries to read zero bytes
					continue
				}

				// Copy the attribute data to the GPU memory buffer.
				glBufferSubData(
					gl.ARRAY_BUFFER,
					offset,
					val.size,
					glPtr(val.value),
				)
				offset += val.size
			}
		}
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
		} else {
			batched = append(batched, last)
			last = p
		}
	}
	batched = append(batched, last)
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
