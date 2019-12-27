package render

import (
	"unsafe"

	"github.com/buchanae/ink/trace"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type pass struct {
	name        string
	prog        compiled
	uniforms    map[string]interface{}
	output      msaa
	vao         uint32
	bindings    []binding
	vertexCount int
	faceOffset  int
	faceCount   int
}

type bindingVal struct {
	value unsafe.Pointer
	size  int
}

type binding struct {
	attr   attribute
	values []bindingVal
	// byte offset into the attribute data buffer
	//offset  int
	divisor   uint32
	totalSize int
}

type passBuilder struct {
	passes []*pass
	faces  []uint32
	// Total number of bytes needed to store all attributes.
	attrBytes int

	// opengl buffer IDs
	attrBufID uint32
	faceBufID uint32
	vaos      []uint32
}

func newPassBuilder() *passBuilder {
	return &passBuilder{
		passes: make([]*pass, 0, 500),
		faces:  make([]uint32, 0, 5000),
	}
}

func (pb *passBuilder) Cleanup() {
	glDeleteBuffers(1, &pb.attrBufID)
	glDeleteBuffers(1, &pb.faceBufID)
	glDeleteVertexArrays(int32(len(pb.vaos)), &pb.vaos[0])
}

func (pb *passBuilder) AddLayer(layer *Layer, output msaa) {
	p := &pass{
		prog:        layer.prog,
		output:      output,
		name:        layer.name,
		vertexCount: layer.vertexCount,
		uniforms:    layer.uniforms,
		faceOffset:  len(pb.faces),
		faceCount:   len(layer.faces),
	}
	pb.faces = append(pb.faces, layer.faces...)
	pb.passes = append(pb.passes, p)

	for _, attr := range layer.prog.attributes {
		val, ok := layer.attrs[attr.Name]
		if !ok {
			continue
		}

		p.bindings = append(p.bindings, binding{
			attr:   attr,
			values: []bindingVal{val},
		})
		pb.attrBytes += val.size
	}
}

func (pb *passBuilder) Passes() []*pass {
	if len(pb.passes) == 0 {
		return nil
	}

	pb.batch()
	pb.upload()

	return pb.passes
}

func (pb *passBuilder) uploadFaces() {
	// Element buffers are used for indexed rendering.
	glGenBuffers(1, &pb.faceBufID)
	glBindBuffer(gl.ELEMENT_ARRAY_BUFFER, pb.faceBufID)

	glBufferData(
		gl.ELEMENT_ARRAY_BUFFER,
		// 4 bytes per index (uint32)
		len(pb.faces)*4,
		glPtr(pb.faces),
		gl.STATIC_DRAW,
	)
}

func (pb *passBuilder) upload() {
	trace.Log("upload")

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

			for _, val := range b.values {
				if val.size == 0 {
					// protection against weird things panicing.
					continue
				}

				// Copy the attribute data to the GPU memory buffer.
				glBufferSubData(
					gl.ARRAY_BUFFER,
					offset,
					val.size,
					val.value,
				)
				offset += val.size
				//glVertexAttribDivisor(b.attr.Loc, b.divisor)
			}
		}
	}
}

func (pb *passBuilder) batch() {
	trace.Log("batch")

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
	trace.Log("  merged passes %d to %d", len(pb.passes), len(batched))
	pb.passes = batched
}

func (pb *passBuilder) mergeable(a, b *pass) bool {
	if a.prog.id != b.prog.id {
		//trace.Log("not mergeable: prog.ID")
		return false
	}
	if len(a.bindings) != len(b.bindings) {
		//trace.Log("not mergeable: len(bindings)")
		return false
	}
	if len(a.uniforms) != len(b.uniforms) {
		//trace.Log("not mergeable: len(uniforms)")
		return false
	}
	if a.output != b.output {
		//trace.Log("not mergeable: output")
		return false
	}
	for i := range b.bindings {
		if a.bindings[i].attr != b.bindings[i].attr {
			//trace.Log("not mergeable: binding.attr %v %v", a.bindings[i].attr, b.bindings[i].attr)
			return false
		}
		if a.bindings[i].divisor != b.bindings[i].divisor {
			//trace.Log("not mergeable: binding.divisor")
			return false
		}
	}
	for k, v := range a.uniforms {
		c, ok := b.uniforms[k]
		if !ok {
			//trace.Log("not mergeable: uniform !ok")
			return false
		}
		if c != v {
			//trace.Log("not mergeable: uniform value")
			return false
		}
	}
	return true
}

func (pb *passBuilder) merge(a, b *pass) {

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
			//a.bindings[i].size += bv.size
		}
	}
}
