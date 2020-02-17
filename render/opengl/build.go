package opengl

import (
	"log"

	"github.com/buchanae/ink/render"
	"github.com/buchanae/ink/trac"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type build struct {
	// opengl buffer IDs
	attrBufID uint32
	faceBufID uint32
	vaos      []uint32
}

type buildPass struct {
	render.Pass

	prog compiled
	vao  uint32
}

func (pb *build) build(plan render.Plan) []buildPass {
	trac.Log("start build")
	defer trac.Log("end build")

	if len(plan.Passes) == 0 {
		return nil
	}

	trac.Log("compile shaders")

	// compile all shader programs
	programs := map[int]compiled{}
	for id, src := range plan.Shaders {
		prog, err := compile(shaderOpt{
			src.Vert, src.Frag, src.Geom, src.Output,
		})
		if err != nil {
			log.Printf("error: compiling shader: %v", err)
			return nil
		}
		programs[id] = prog
	}

	trac.Log("upload")

	// upload faces (vertex index)
	// Element buffers are used for indexed rendering.
	glGenBuffers(1, &pb.faceBufID)
	glBindBuffer(gl.ELEMENT_ARRAY_BUFFER, pb.faceBufID)
	glBufferData(
		gl.ELEMENT_ARRAY_BUFFER,
		// Size of buffer in bytes
		// 4 bytes per index (uint32)
		len(plan.FaceData)*4,
		glPtr(plan.FaceData),
		gl.STATIC_DRAW,
	)

	// The data from all attributes is stored in one large buffer.
	glGenBuffers(1, &pb.attrBufID)
	glBindBuffer(gl.ARRAY_BUFFER, pb.attrBufID)
	glBufferData(
		gl.ARRAY_BUFFER,
		// buffer size in bytes.
		// 4 bytes per float32.
		len(plan.AttrData)*4,
		gl.Ptr(plan.AttrData),
		gl.STATIC_DRAW,
	)

	// Each pass has one VAO, which stores the configuration of all its
	// attributes: location of the data in the buffer, enabled/disable,
	// types, divisors, etc.
	pb.vaos = make([]uint32, len(plan.Passes))
	glGenVertexArrays(int32(len(pb.vaos)), &pb.vaos[0])

	passes := []buildPass{}

	for i, pass := range plan.Passes {
		vao := pb.vaos[i]
		glBindVertexArray(vao)
		glBindBuffer(gl.ELEMENT_ARRAY_BUFFER, pb.faceBufID)

		prog, ok := programs[pass.Shader]
		if !ok {
			log.Printf("missing program for shader %d", pass.Shader)
			continue
		}

		attrs := map[string]render.Attr{}
		for _, v := range pass.Attrs {
			attrs[v.Name] = v
		}

		for _, b := range prog.attributes {

			attr, ok := attrs[b.Name]
			if !ok {
				//log.Printf("missing attr %s", b.Name)
				continue
			}

			glEnableVertexAttribArray(b.Loc)
			glVertexAttribPointer(
				b.Loc,
				b.Components,
				b.Datatype,
				false, // normalized
				0,     // stride
				glPtrOffset(attr.Offset*4),
			)
			glVertexAttribDivisor(b.Loc, uint32(attr.Divisor))
		}

		passes = append(passes, buildPass{
			Pass: pass,
			vao:  vao,
			prog: prog,
		})
	}
	return passes
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
