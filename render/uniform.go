package render

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// uniform contains information about a uniform value in a GLSL program.
// uniforms are usually created by inspecting a compiled program with inspectUniforms.
type uniform struct {
	Name string
	Type uint32
	Size int32
	Loc  int32
	Tex  int32
}

// Bind attempts to bind the value to the uniform.
// It's expected that the program is already active (ie: glUseProgram has been called).
func (u *uniform) Bind(val interface{}) error {

	if u.Size > 1 {
		return fmt.Errorf("uniform arrays are not supported (yet): %q\n", u.Name)
	}

	// TODO don't convert types on every frame (i.e. optimize uniform binding)
	switch u.Type {
	case gl.FLOAT:
		z, ok := val.(float32)
		if !ok {
			return fmt.Errorf("type mismatch: want float32 got %T", val)
		}
		glUniform1f(u.Loc, z)

	case gl.FLOAT_VEC2:
		z, ok := val.([2]float32)
		if !ok {
			return fmt.Errorf("type mismatch: want [2]float32 got %T", val)
		}
		glUniform2f(u.Loc, z[0], z[1])

	case gl.FLOAT_VEC3:
		z, ok := val.([3]float32)
		if !ok {
			return fmt.Errorf("type mismatch: want [3]float32 got %T", val)
		}
		glUniform3f(u.Loc, z[0], z[1], z[2])

	case gl.FLOAT_VEC4:
		z, ok := val.([4]float32)
		if !ok {
			return fmt.Errorf("type mismatch: want [4]float32 got %T", val)
		}
		glUniform4f(u.Loc, z[0], z[1], z[2], z[3])

		/*
			case gl.DOUBLE:
			case gl.DOUBLE_VEC2:
			case gl.DOUBLE_VEC3:
			case gl.DOUBLE_VEC4:
		*/

	case gl.INT:
		z, ok := val.(int32)
		if !ok {
			return fmt.Errorf("type mismatch: want int32 got %T", val)
		}
		glUniform1i(u.Loc, z)

		/*
			case gl.INT_VEC2:
			case gl.INT_VEC3:
			case gl.INT_VEC4:
		*/

	case gl.UNSIGNED_INT:
		z, ok := val.(uint32)
		if !ok {
			return fmt.Errorf("type mismatch: want uint32 got %T", val)
		}
		glUniform1ui(u.Loc, z)

		/*
			case gl.UNSIGNED_INT_VEC2:
			case gl.UNSIGNED_INT_VEC3:
			case gl.UNSIGNED_INT_VEC4:
		*/

	case gl.BOOL:
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
			case gl.BOOL_VEC2:
			case gl.BOOL_VEC3:
			case gl.BOOL_VEC4:
		*/

	case gl.SAMPLER_2D:
		var tex uint32
		switch z := val.(type) {
		case Image:
			tex = z.tex
		case *Layer:
			tex = z.tex.Read.Tex
		case msaa:
			tex = z.Read.Tex
		default:
			return fmt.Errorf("type mismatch: expected texture, got %T", val)
		}
		glUniform1i(u.Loc, u.Tex)
		glActiveTexture(uint32(gl.TEXTURE0 + u.Tex))
		glBindTexture(gl.TEXTURE_2D, tex)

		/*
			case gl.FLOAT_MAT2:
			case gl.FLOAT_MAT3:
			case gl.FLOAT_MAT4:
			case gl.FLOAT_MAT2x3:
			case gl.FLOAT_MAT2x4:
			case gl.FLOAT_MAT3x2:
			case gl.FLOAT_MAT3x4:
			case gl.FLOAT_MAT4x2:
			case gl.FLOAT_MAT4x3:
			case gl.DOUBLE_MAT2:
			case gl.DOUBLE_MAT3:
			case gl.DOUBLE_MAT4:
			case gl.DOUBLE_MAT2x3:
			case gl.DOUBLE_MAT2x4:
			case gl.DOUBLE_MAT3x2:
			case gl.DOUBLE_MAT3x4:
			case gl.DOUBLE_MAT4x2:
			case gl.DOUBLE_MAT4x3:
			case gl.SAMPLER_1D:
			case gl.SAMPLER_3D:
			case gl.SAMPLER_CUBE:
			case gl.SAMPLER_1D_SHADOW:
			case gl.SAMPLER_2D_SHADOW:
			case gl.SAMPLER_1D_ARRAY:
			case gl.SAMPLER_2D_ARRAY:
			case gl.SAMPLER_1D_ARRAY_SHADOW:
			case gl.SAMPLER_2D_ARRAY_SHADOW:
			case gl.SAMPLER_2D_MULTISAMPLE:
			case gl.SAMPLER_2D_MULTISAMPLE_ARRAY:
			case gl.SAMPLER_CUBE_SHADOW:
			case gl.SAMPLER_BUFFER:
			case gl.SAMPLER_2D_RECT:
			case gl.SAMPLER_2D_RECT_SHADOW:
			case gl.INT_SAMPLER_1D:
			case gl.INT_SAMPLER_2D:
			case gl.INT_SAMPLER_3D:
			case gl.INT_SAMPLER_CUBE:
			case gl.INT_SAMPLER_1D_ARRAY:
			case gl.INT_SAMPLER_2D_ARRAY:
			case gl.INT_SAMPLER_2D_MULTISAMPLE:
			case gl.INT_SAMPLER_2D_MULTISAMPLE_ARRAY:
			case gl.INT_SAMPLER_BUFFER:
			case gl.INT_SAMPLER_2D_RECT:
			case gl.UNSIGNED_INT_SAMPLER_1D:
			case gl.UNSIGNED_INT_SAMPLER_2D:
			case gl.UNSIGNED_INT_SAMPLER_3D:
			case gl.UNSIGNED_INT_SAMPLER_CUBE:
			case gl.UNSIGNED_INT_SAMPLER_1D_ARRAY:
			case gl.UNSIGNED_INT_SAMPLER_2D_ARRAY:
			case gl.UNSIGNED_INT_SAMPLER_2D_MULTISAMPLE:
			case gl.UNSIGNED_INT_SAMPLER_2D_MULTISAMPLE_ARRAY:
			case gl.UNSIGNED_INT_SAMPLER_BUFFER:
			case gl.UNSIGNED_INT_SAMPLER_2D_RECT:
		*/
	}
	return nil
}

// inspectUniforms queries OpenGL for information on the uniforms
// of a compiled GLSL program. "id" is the OpenGL handle of the program
// to inspect.
func inspectUniforms(id uint32) []uniform {
	glUseProgram(id)

	// Get the number of active uniforms.
	var numUnis int32
	glGetProgramiv(id, gl.ACTIVE_UNIFORMS, &numUnis)

	// Get maximum uniform name length.
	var maxnamelen int32
	glGetProgramiv(id, gl.ACTIVE_UNIFORM_MAX_LENGTH, &maxnamelen)

	unis := make([]uniform, 0, numUnis)
	for i := int32(0); i < numUnis; i++ {

		var xtype uint32
		var size, namelen int32
		namebytes := make([]byte, maxnamelen)
		glGetActiveUniform(id, uint32(i), maxnamelen,
			&namelen, &size, &xtype, &namebytes[0])

		name := string(namebytes[:namelen])
		loc := glGetUniformLocation(id, glStr(name+"\x00"))

		unis = append(unis, uniform{
			Name: name,
			Size: size,
			Type: xtype,
			Loc:  loc,
			Tex:  i,
		})
	}
	return unis
}
