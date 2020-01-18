package render

import "fmt"

// program contains information about a compiled GLSL program.
type program struct {
	// id is the OpenGL program ID.
	id glProgram
	// Uniforms contains static information about the uniforms
	// defined in the program's shader code.
	attributes []attribute
	// Attributes contains static information about the attributes
	// defined in the program's shader code.
	uniforms []uniform
}

// attribute contains information about a vertex attribute in a GLSL program.
// attributes are created by inspecting a compiled program, via inspectAttributes.
type attribute struct {
	Name string
	Type glEnum
	Size int32
	Loc  glAttribute
	Idx  int32
	// Number of bytes of the data type,
	// eg: "float" is 4 bytes (32-bit float)
	//     "vec2" is also 4, because it is a vector of floats.
	ByteSize int
	// Number of components in the attribute.
	// eg: "float" is 1, "vec2" is 2, "vec3" is 3, etc.
	Components int32
	// OpenGL enum value describing the data type
	// of a single component, e.g. a vec3 would have a DataType
	// of float, because a single component of the vec3 is a
	// float value.
	DataType glEnum
}

// uniform contains information about a uniform value in a GLSL program.
// uniforms are usually created by inspecting a compiled program with inspectUniforms.
type uniform struct {
	Name        string
	Type        glEnum
	Size        int32
	Loc         glUniform
	Idx         int32
	TextureUnit glTextureUnit
}

// Bind attempts to bind the value to the uniform.
// It's expected that the program is already active (ie: glUseProgram has been called).
func (u *uniform) Bind(val interface{}) error {

	if u.Size > 1 {
		return fmt.Errorf("uniform arrays are not supported (yet): %q\n", u.Name)
	}

	// TODO don't convert types on every frame (i.e. optimize uniform binding)
	switch u.Type {
	case gl_FLOAT:
		z, ok := val.(float32)
		if !ok {
			return fmt.Errorf("type mismatch: want float32 got %T", val)
		}
		glUniform1f(u.Loc, z)

	case gl_FLOAT_VEC2:
		z, ok := val.([2]float32)
		if !ok {
			return fmt.Errorf("type mismatch: want [2]float32 got %T", val)
		}
		glUniform2f(u.Loc, z[0], z[1])

	case gl_FLOAT_VEC3:
		z, ok := val.([3]float32)
		if !ok {
			return fmt.Errorf("type mismatch: want [3]float32 got %T", val)
		}
		glUniform3f(u.Loc, z[0], z[1], z[2])

	case gl_FLOAT_VEC4:
		z, ok := val.([4]float32)
		if !ok {
			return fmt.Errorf("type mismatch: want [4]float32 got %T", val)
		}
		glUniform4f(u.Loc, z[0], z[1], z[2], z[3])

		/*
			case gl_DOUBLE:
			case gl_DOUBLE_VEC2:
			case gl_DOUBLE_VEC3:
			case gl_DOUBLE_VEC4:
		*/

	case gl_INT:
		z, ok := val.(int32)
		if !ok {
			return fmt.Errorf("type mismatch: want int32 got %T", val)
		}
		glUniform1i(u.Loc, z)

		/*
			case gl_INT_VEC2:
			case gl_INT_VEC3:
			case gl_INT_VEC4:
		*/

	case gl_UNSIGNED_INT:
		z, ok := val.(uint32)
		if !ok {
			return fmt.Errorf("type mismatch: want uint32 got %T", val)
		}
		glUniform1ui(u.Loc, z)

		/*
			case gl_UNSIGNED_INT_VEC2:
			case gl_UNSIGNED_INT_VEC3:
			case gl_UNSIGNED_INT_VEC4:
		*/

	case gl_BOOL:
		b, ok := val.(bool)
		if !ok {
			return fmt.Errorf("type mismatch: expected bool, got %T", val)
		}
		if b {
			glUniform1i(u.Loc, 1)
		} else {
			glUniform1i(u.Loc, 0)
		}

		/*
			case gl_BOOL_VEC2:
			case gl_BOOL_VEC3:
			case gl_BOOL_VEC4:
		*/

	case gl_SAMPLER_2D:
		var tex glTexture
		switch z := val.(type) {
		case Image:
			tex = z.tex
		case msaa:
			tex = z.Read.Tex
		default:
			return fmt.Errorf("type mismatch: expected texture, got %T", val)
		}
		glUniform1i(u.Loc, u.Idx)
		glActiveTexture(u.TextureUnit)
		glBindTexture(gl_TEXTURE_2D, tex)

		/*
			case gl_FLOAT_MAT2:
			case gl_FLOAT_MAT3:
			case gl_FLOAT_MAT4:
			case gl_FLOAT_MAT2x3:
			case gl_FLOAT_MAT2x4:
			case gl_FLOAT_MAT3x2:
			case gl_FLOAT_MAT3x4:
			case gl_FLOAT_MAT4x2:
			case gl_FLOAT_MAT4x3:
			case gl_DOUBLE_MAT2:
			case gl_DOUBLE_MAT3:
			case gl_DOUBLE_MAT4:
			case gl_DOUBLE_MAT2x3:
			case gl_DOUBLE_MAT2x4:
			case gl_DOUBLE_MAT3x2:
			case gl_DOUBLE_MAT3x4:
			case gl_DOUBLE_MAT4x2:
			case gl_DOUBLE_MAT4x3:
			case gl_SAMPLER_1D:
			case gl_SAMPLER_3D:
			case gl_SAMPLER_CUBE:
			case gl_SAMPLER_1D_SHADOW:
			case gl_SAMPLER_2D_SHADOW:
			case gl_SAMPLER_1D_ARRAY:
			case gl_SAMPLER_2D_ARRAY:
			case gl_SAMPLER_1D_ARRAY_SHADOW:
			case gl_SAMPLER_2D_ARRAY_SHADOW:
			case gl_SAMPLER_2D_MULTISAMPLE:
			case gl_SAMPLER_2D_MULTISAMPLE_ARRAY:
			case gl_SAMPLER_CUBE_SHADOW:
			case gl_SAMPLER_BUFFER:
			case gl_SAMPLER_2D_RECT:
			case gl_SAMPLER_2D_RECT_SHADOW:
			case gl_INT_SAMPLER_1D:
			case gl_INT_SAMPLER_2D:
			case gl_INT_SAMPLER_3D:
			case gl_INT_SAMPLER_CUBE:
			case gl_INT_SAMPLER_1D_ARRAY:
			case gl_INT_SAMPLER_2D_ARRAY:
			case gl_INT_SAMPLER_2D_MULTISAMPLE:
			case gl_INT_SAMPLER_2D_MULTISAMPLE_ARRAY:
			case gl_INT_SAMPLER_BUFFER:
			case gl_INT_SAMPLER_2D_RECT:
			case gl_UNSIGNED_INT_SAMPLER_1D:
			case gl_UNSIGNED_INT_SAMPLER_2D:
			case gl_UNSIGNED_INT_SAMPLER_3D:
			case gl_UNSIGNED_INT_SAMPLER_CUBE:
			case gl_UNSIGNED_INT_SAMPLER_1D_ARRAY:
			case gl_UNSIGNED_INT_SAMPLER_2D_ARRAY:
			case gl_UNSIGNED_INT_SAMPLER_2D_MULTISAMPLE:
			case gl_UNSIGNED_INT_SAMPLER_2D_MULTISAMPLE_ARRAY:
			case gl_UNSIGNED_INT_SAMPLER_BUFFER:
			case gl_UNSIGNED_INT_SAMPLER_2D_RECT:
		*/
	}
	return nil
}
