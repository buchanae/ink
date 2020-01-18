// +build js

package render

import (
	"fmt"
	"log"
	"regexp"
	"syscall/js"
)

var inrx = regexp.MustCompile(`(?m)^\s*(in )`)
var outrx = regexp.MustCompile(`(?m)^\s*(out )`)

func glBuildProgram(vert, frag, geom, out string) (program, error) {
	var prog program

	vert = inrx.ReplaceAllString(vert, "attribute ")
	vert = outrx.ReplaceAllString(vert, "varying ")
	frag = inrx.ReplaceAllString(frag, "varying ")

	log.Print("VERT", vert)

	// VERTEX SHADER
	vs, err := glBuildShader(vert, gl_VERTEX_SHADER)
	if err != nil {
		return prog, fmt.Errorf("creating vert shader: %s", err)
	}

	// FRAGMENT SHADER
	fs, err := glBuildShader(frag, gl_FRAGMENT_SHADER)
	if err != nil {
		return prog, fmt.Errorf("creating frag shader: %s", err)
	}

	// CREATE PROGRAM
	programID := gl.Call("createProgram")
	gl.Call("attachShader", programID, vs)
	gl.Call("attachShader", programID, fs)

	// GEOMETRY SHADER
	if geom != "" {
		return prog, fmt.Errorf("geometry shader not implemented")
	}

	//gl.Call("bindFragDataLocation", programID, 0, out)
	gl.Call("linkProgram", programID)

	linkstatus := gl.Call("getProgramParameter", programID, gl_LINK_STATUS)
	if !linkstatus.Bool() {
		log := gl.Call("getProgramInfoLog", programID)
		return prog, fmt.Errorf("linking program: %s", log)
	}

	gl.Call("useProgram", programID)

	glpid := glProgram(programID)

	return program{
		id:         glpid,
		attributes: inspectAttributes(glpid),
		uniforms:   inspectUniforms(glpid),
	}, nil
}

func glBuildShader(src string, shaderType glEnum) (js.Value, error) {
	id := gl.Call("createShader", shaderType)

	//src = "#version 300 es\n\n" + src
	gl.Call("shaderSource", id, src)
	gl.Call("compileShader", id)

	status := gl.Call("getShaderParameter", id, gl_COMPILE_STATUS)
	if !status.Bool() {
		log := gl.Call("getShaderInfoLog", id)
		return id, fmt.Errorf("compiling shader: %s", log)
	}
	return id, nil
}

// inspectProgramAttribs queries OpenGL for information on the vertex attributes
// of a compiled GLSL program. "id" is the ID of the compiled program to inspect.
func inspectAttributes(id glProgram) []attribute {

	// Get the number of active uniforms.
	numAttribs := gl.Call("getProgramParameter", id, gl_ACTIVE_ATTRIBUTES).Int()

	attribs := make([]attribute, 0, numAttribs)
	for i := 0; i < numAttribs; i++ {

		info := glGetActiveAttrib(id, uint32(i))
		loc := gl.Call("getAttribLocation", id, info.Name)

		var bytesize int
		var components int32
		var datatype glEnum

		// TODO missing int attribute types?
		switch info.Type {
		case gl_FLOAT:
			bytesize = 4
			components = 1
			datatype = gl_FLOAT
		case gl_FLOAT_VEC2:
			bytesize = 4
			components = 2
			datatype = gl_FLOAT
		case gl_FLOAT_VEC3:
			bytesize = 4
			components = 3
			datatype = gl_FLOAT
		case gl_FLOAT_VEC4:
			bytesize = 4
			components = 4
			datatype = gl_FLOAT
		}

		attribs = append(attribs, attribute{
			Name:       info.Name,
			Size:       int32(info.Size),
			Type:       info.Type,
			Loc:        glAttribute(loc),
			Idx:        int32(i),
			ByteSize:   bytesize,
			Components: components,
			DataType:   datatype,
		})
	}
	return attribs
}

// inspectUniforms queries OpenGL for information on the uniforms
// of a compiled GLSL program. "id" is the OpenGL handle of the program
// to inspect.
func inspectUniforms(id glProgram) []uniform {

	// Get the number of active uniforms.
	numUnis := gl.Call("getProgramParameter", id, gl_ACTIVE_UNIFORMS).Int()

	texOffset := gl.Get("TEXTURE0").Int()

	unis := make([]uniform, 0, numUnis)
	for i := 0; i < numUnis; i++ {

		info := glGetActiveUniform(id, uint32(i))
		loc := gl.Call("getUniformLocation", id, info.Name)

		var texUnit glTextureUnit
		if info.Type == gl_SAMPLER_2D {
			texUnit = glTextureUnit(js.ValueOf(texOffset))
			texOffset++
		}

		unis = append(unis, uniform{
			Name:        info.Name,
			Size:        int32(info.Size),
			Type:        info.Type,
			Loc:         glUniform(loc),
			Idx:         int32(i),
			TextureUnit: texUnit,
		})
	}
	return unis
}

type glActiveInfo struct {
	Name string
	Size int
	Type glEnum
}

func glGetActiveUniform(program glProgram, index uint32) glActiveInfo {
	val := gl.Call("getActiveUniform", program, index)
	return glActiveInfo{
		Name: val.Get("name").String(),
		Size: val.Get("size").Int(),
		Type: glEnum(val.Get("type")),
	}
}

func glGetActiveAttrib(program glProgram, index uint32) glActiveInfo {
	val := gl.Call("getActiveAttrib", program, index)
	return glActiveInfo{
		Name: val.Get("name").String(),
		Size: val.Get("size").Int(),
		Type: glEnum(val.Get("type")),
	}
}
