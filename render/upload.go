package render

type uploaded struct {
	// opengl buffer IDs
	attrBuf glBuffer
	faceBuf glBuffer
	vaos    []glVAO
}

func uploadBuild(pb *build) uploaded {
	pb.trace("upload")

	up := uploaded{}

	pb.trace("upload faces")
	// upload faces (vertex index)
	// Element buffers are used for indexed rendering.
	up.faceBuf = glCreateBuffer()
	glBindBuffer(gl_ELEMENT_ARRAY_BUFFER, up.faceBuf)

	if len(pb.faces) > 0 {
		glBufferData(
			gl_ELEMENT_ARRAY_BUFFER,
			len(pb.faces)*4, // 4 bytes per index (uint32)
			pb.faces,
			gl_STATIC_DRAW,
		)
	}

	pb.trace("upload attributes")

	// The data from all attributes is stored in one large buffer.
	// A "binding" describes the slice of the buffer that holds data
	// for a single attribute.
	up.attrBuf = glCreateBuffer()
	glBindBuffer(gl_ARRAY_BUFFER, up.attrBuf)
	glBufferData(gl_ARRAY_BUFFER, pb.attrBytes, nil, gl_STATIC_DRAW)

	// Each pass has one VAO, which stores the configuration of all its
	// attributes: location of the data in the buffer, enabled/disable,
	// types, divisors, etc.
	up.vaos = make([]glVAO, len(pb.passes))

	offset := 0
	for i, p := range pb.passes {
		vao := glCreateVAO()
		up.vaos[i] = vao
		p.vao = vao

		glBindVAO(p.vao)
		glBindBuffer(gl_ELEMENT_ARRAY_BUFFER, up.faceBuf)

		for _, b := range p.bindings {

			glEnableVertexAttribArray(b.attr.Loc)
			glVertexAttribPointer(
				b.attr.Loc,
				b.attr.Components,
				b.attr.DataType,
				false, // normalized
				0,     // stride
				offset,
			)
			glVertexAttribDivisor(b.attr.Loc, uint32(b.divisor))

			for _, val := range b.values {
				if val.size == 0 {
					// opengl will panic if it tries to read zero bytes
					continue
				}

				// Copy the attribute data to the GPU memory buffer.
				glBufferSubData(
					gl_ARRAY_BUFFER,
					offset,
					val.size,
					val.value,
				)
				offset += val.size
			}
		}
	}

	return up
}

func (up uploaded) cleanup() {
	glDeleteBuffer(up.attrBuf)
	glDeleteBuffer(up.faceBuf)
	for _, vao := range up.vaos {
		glDeleteVAO(vao)
	}
}
