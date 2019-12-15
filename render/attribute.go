package render

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

// attribute contains information about a vertex attribute in a GLSL program.
// attributes are created by inspecting a compiled program, via inspectAttributes.
type attribute struct {
	Name string
	Type uint32
	Size int32
	Loc  uint32
	Idx  int32
	// Number of bytes of the data type,
	// eg: "float" is 4 bytes (32-bit float)
	//     "vec2" is also 4, because it is a vector of floats.
	Bytesize int
	// Number of components in the attribute.
	// eg: "float" is 1, "vec2" is 2, "vec3" is 3, etc.
	Components int32
	// OpenGL enum value describing the data type, eg. gl.FLOAT
	Datatype uint32
}

// inspectProgramAttribs queries OpenGL for information on the vertex attributes
// of a compiled GLSL program. "id" is the ID of the compiled program to inspect.
func inspectAttributes(id uint32) []attribute {
	glUseProgram(id)

	// Get the number of active uniforms.
	var numAttribs int32
	glGetProgramiv(id, gl.ACTIVE_ATTRIBUTES, &numAttribs)

	// Get maximum uniform name length.
	var maxnamelen int32
	glGetProgramiv(id, gl.ACTIVE_ATTRIBUTE_MAX_LENGTH, &maxnamelen)

	attribs := make([]attribute, 0, numAttribs)
	for i := int32(0); i < numAttribs; i++ {

		var xtype uint32
		var size, namelen int32
		namebytes := make([]byte, maxnamelen)
		glGetActiveAttrib(id, uint32(i), maxnamelen,
			&namelen, &size, &xtype, &namebytes[0])

		name := string(namebytes[:namelen])
		loc := glGetAttribLocation(id, glStr(name+"\x00"))

		var bytesize int
		var components int32
		var datatype uint32

		// TODO missing int attribute types?
		switch xtype {
		case gl.FLOAT:
			bytesize = 4
			components = 1
			datatype = gl.FLOAT
		case gl.FLOAT_VEC2:
			bytesize = 4
			components = 2
			datatype = gl.FLOAT
		case gl.FLOAT_VEC3:
			bytesize = 4
			components = 3
			datatype = gl.FLOAT
		case gl.FLOAT_VEC4:
			bytesize = 4
			components = 4
			datatype = gl.FLOAT
		}

		attribs = append(attribs, attribute{
			Name:       name,
			Size:       size,
			Type:       xtype,
			Loc:        uint32(loc),
			Idx:        i,
			Bytesize:   bytesize,
			Components: components,
			Datatype:   datatype,
		})
	}
	return attribs
}
