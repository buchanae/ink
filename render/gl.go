// +build !js

package render

import (
	"fmt"
	"log"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type glEnum uint32
type glTexture uint32
type glBuffer uint32
type glFramebuffer uint32
type glProgram uint32
type glUniform uint32
type glAttribute uint32
type glTextureUnit uint32
type glVAO uint32

func glEnable(flag uint32) {
	gl.Enable(flag)
	glLogErr("Enable")
}

func glDisable(flag uint32) {
	gl.Disable(flag)
	glLogErr("Disable")
}

func glBlendFunc(sfactor glEnum, dfactor glEnum) {
	gl.BlendFunc(uint32(sfactor), uint32(dfactor))
	glLogErr("BlendFunc")
}

func glViewport(x, y, width, height int32) {
	gl.Viewport(x, y, width, height)
	glLogErr("Viewport")
}

func glClearColor(red, green, blue, alpha float32) {
	gl.ClearColor(red, green, blue, alpha)
	glLogErr("ClearColor")
}

func glClear(mask uint32) {
	gl.Clear(mask)
	glLogErr("Clear")
}

/*

Create, Delete, Bind








*/

func glCreateBuffer() glBuffer {
	var id uint32
	gl.GenBuffers(1, &id)
	glLogErr("GenBuffers")
	return glBuffer(id)
}

func glCreateTexture() glTexture {
	var id uint32
	gl.GenTextures(1, &id)
	glLogErr("GenTextures")
	return glTexture(id)
}

func glCreateFramebuffer() glFramebuffer {
	var id uint32
	gl.GenFramebuffers(1, &id)
	glLogErr("GenFramebuffers")
	return glFramebuffer(id)
}

func glCreateVAO() glVAO {
	var id uint32
	gl.GenVertexArrays(1, &id)
	glLogErr("GenVertexArrays")
	return glVAO(id)
}

func glBindVertexArray(vao glVAO) {
	gl.BindVertexArray(uint32(vao))
	glLogErr("BindVertexArray")
}

func glDeleteTexture(tex glTexture) {
	id := uint32(tex)
	gl.DeleteTextures(1, &id)
	glLogErr("DeleteTextures")
}

func glDeleteFramebuffer(fbo glFramebuffer) {
	id := uint32(fbo)
	gl.DeleteFramebuffers(1, &id)
	glLogErr("DeleteFramebuffers")
}

func glDeleteVAO(vao glVAO) {
	id := uint32(vao)
	gl.DeleteVertexArrays(1, &id)
	glLogErr("DeleteVertexArrays")
}

func glDeleteBuffer(buf glBuffer) {
	id := uint32(buf)
	gl.DeleteBuffers(1, &id)
	glLogErr("DeleteBuffers")
}

func glDeleteProgram(program glProgram) {
	gl.DeleteProgram(uint32(program))
	glLogErr("DeleteProgram")
}

func glBindBuffer(target glEnum, buffer glBuffer) {
	gl.BindBuffer(uint32(target), uint32(buffer))
	glLogErr("BindBuffer")
}

func glBindTexture(target glEnum, texture glTexture) {
	gl.BindTexture(uint32(target), uint32(texture))
	glLogErr("BindTexture")
}

func glBindFramebuffer(target glEnum, framebuffer glFramebuffer) {
	gl.BindFramebuffer(uint32(target), uint32(framebuffer))
	glLogErr("BindFramebuffer")
}

func glVertexAttribPointer(
	index glAttribute,
	size int32,
	xtype glEnum,
	normalized bool,
	stride int32,
	offset int) {
	ptr := gl.PtrOffset(offset)
	gl.VertexAttribPointer(uint32(index), size, uint32(xtype), normalized, stride, ptr)
	glLogErr("VertexAttribPointer")
}

func glVertexAttribDivisor(index glAttribute, divisor uint32) {
	gl.VertexAttribDivisor(uint32(index), divisor)
	glLogErr("VertexAttribDivisor")
}

func glEnableVertexAttribArray(index glAttribute) {
	gl.EnableVertexAttribArray(uint32(index))
	glLogErr("EnableVertexAttribArray")
}

func glUseProgram(program glProgram) {
	gl.UseProgram(uint32(program))
	glLogErr("UseProgram")
}

func glTexParameteri(target, pname glEnum, param int32) {
	gl.TexParameteri(uint32(target), uint32(pname), param)
	glLogErr("TexParameteri")
}

func glBufferData(target uint32, size int, data interface{}, usage uint32) {
	ptr := gl.Ptr(data)
	gl.BufferData(target, size, ptr, usage)
	glLogErr("BufferData")
}

func glBufferSubData(target uint32, offset int, size int, data interface{}) {
	ptr := gl.Ptr(data)
	gl.BufferSubData(target, offset, size, ptr)
	glLogErr("BufferSubData")
}

func glDrawElementsInstanced(mode uint32, count int32, xtype uint32, offset int, instancecount int32) {
	ptr := gl.PtrOffset(offset)
	gl.DrawElementsInstanced(mode, count, xtype, ptr, instancecount)
	glLogErr("DrawElementsInstanced")
}

/*

Uniform binding






*/

func glUniform1f(location glUniform, v0 float32) {
	gl.Uniform1f(int32(location), v0)
	glLogErr("Uniform1f")
}

func glUniform2f(location glUniform, v0 float32, v1 float32) {
	gl.Uniform2f(int32(location), v0, v1)
	glLogErr("Uniform2f")
}

func glUniform3f(location glUniform, v0 float32, v1 float32, v2 float32) {
	gl.Uniform3f(int32(location), v0, v1, v2)
	glLogErr("Uniform3f")
}

func glUniform4f(location glUniform, v0 float32, v1 float32, v2 float32, v3 float32) {
	gl.Uniform4f(int32(location), v0, v1, v2, v3)
	glLogErr("Uniform4f")
}

func glUniform1i(location glUniform, v0 int32) {
	gl.Uniform1i(int32(location), v0)
	glLogErr("Uniform1i")
}

func glUniform1ui(location glUniform, v0 uint32) {
	gl.Uniform1ui(int32(location), v0)
	glLogErr("Uniform1ui")
}

func glActiveTexture(texture glTextureUnit) {
	gl.ActiveTexture(uint32(texture))
	glLogErr("ActiveTexture")
}

/*

Textures





*/

func glFramebufferTexture2D(
	target,
	attachment,
	textarget glEnum,
	texture glTexture,
	level int32) {
	gl.FramebufferTexture2D(uint32(target), uint32(attachment), uint32(textarget), uint32(texture), level)
	glLogErr("FramebufferTexture2D")
}

func glTexImage2DMultisample(
	target uint32,
	samples int32,
	internalformat uint32,
	width int32, height int32,
	fixedsamplelocations bool) {
	gl.TexImage2DMultisample(target, samples, internalformat, width, height, fixedsamplelocations)
	glLogErr("TexImage2DMultisample")
}

func glBlitFramebuffer(
	srcX0, srcY0, srcX1, srcY1 int32,
	dstX0, dstY0, dstX1, dstY1 int32,
	mask uint32, filter uint32) {
	gl.BlitFramebuffer(srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1, mask, filter)
	glLogErr("BlitFramebuffer")
}

func glTexImage2D(
	target uint32,
	level int32,
	internalformat int32,
	width int32, height int32,
	border int32,
	format glEnum,
	xtype glEnum,
	pixels []uint8) {

	var ptr unsafe.Pointer
	if pixels == nil {
		ptr = gl.Ptr(nil)
	} else {
		ptr = gl.Ptr(pixels)
	}
	gl.TexImage2D(target, level, internalformat, width, height, border, uint32(format), uint32(xtype), ptr)
	glLogErr("TexImage2D")
}

func glReadPixels(x int32, y int32, width int32, height int32, format uint32, xtype uint32, pixels []uint8) {
	ptr := gl.Ptr(pixels)
	gl.ReadPixels(x, y, width, height, format, xtype, ptr)
	glLogErr("ReadPixels")
}

/*

Error log






*/

func glLogErr(name string) {
	err := glCheckErr(name)
	if err != nil {
		log.Printf("%v", err)
	}
}

func glCheckErr(name string) error {
	err := gl.GetError()
	switch err {
	case gl.NO_ERROR:
		return nil
	case gl.INVALID_ENUM:
		return fmt.Errorf("%s: invalid enum", name)
	case gl.INVALID_VALUE:
		return fmt.Errorf("%s: invalid value", name)
	case gl.INVALID_OPERATION:
		return fmt.Errorf("%s: invalid operation", name)
	case gl.STACK_OVERFLOW:
		return fmt.Errorf("%s: stack overflow", name)
	case gl.STACK_UNDERFLOW:
		return fmt.Errorf("%s: stack underflow", name)
	case gl.OUT_OF_MEMORY:
		return fmt.Errorf("%s: out of memory", name)
	case gl.INVALID_FRAMEBUFFER_OPERATION:
		return fmt.Errorf("%s: invalid framebuffer operation", name)
	default:
		return fmt.Errorf("%s: unknown error code", name)
	}
}

func glGetProgramInfoLog(id uint32) string {
	var logLength int32
	gl.GetProgramiv(id, gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)

	gl.GetProgramInfoLog(id, logLength, nil, &logBuffer[0])
	glLogErr("GetProgramInfoLog")

	return gl.GoStr(&logBuffer[0])
}

func glGetShaderInfoLog(id uint32) string {

	var logLength int32
	gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &logLength)
	glLogErr("GetShaderiv")
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)
	gl.GetShaderInfoLog(id, logLength, nil, &logBuffer[0])
	glLogErr("GetShaderInfoLog")

	return gl.GoStr(&logBuffer[0])
}

func glBuildShader(src string, shaderType uint32) (uint32, error) {
	id := gl.CreateShader(shaderType)
	glLogErr("CreateShader")

	src = "#version 330 core\n\n" + src
	source, free := gl.Strs(src + "\000")
	defer free()

	gl.ShaderSource(id, 1, source, nil)
	glLogErr("ShaderSource")

	gl.CompileShader(id)
	glLogErr("CompileShader")

	var status int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &status)
	glLogErr("GetShaderiv")

	if status != 1 {
		return 0, fmt.Errorf("compiling shader: %s", glGetShaderInfoLog(id))
	}
	return id, nil
}

func glBuildProgram(vert, frag, geom, out string) (program, error) {
	prog := program{}

	// VERTEX SHADER
	vs, err := glBuildShader(vert, gl.VERTEX_SHADER)
	if err != nil {
		return prog, fmt.Errorf("creating vert shader: %s", err)
	}

	// FRAGMENT SHADER
	fs, err := glBuildShader(frag, gl.FRAGMENT_SHADER)
	if err != nil {
		return prog, fmt.Errorf("creating frag shader: %s", err)
	}

	// CREATE PROGRAM
	programID := gl.CreateProgram()
	glLogErr("CreateProgram")
	gl.AttachShader(programID, vs)
	glLogErr("AttachShader")
	gl.AttachShader(programID, fs)
	glLogErr("AttachShader")

	// GEOMETRY SHADER
	if geom != "" {
		gs, err := glBuildShader(geom, gl.GEOMETRY_SHADER)
		if err != nil {
			return prog, fmt.Errorf("creating geo shader: %s", err)
		}
		gl.AttachShader(programID, gs)
		glLogErr("AttachShader")
	}

	gl.BindFragDataLocation(programID, 0, gl.Str(out+"\000"))
	glLogErr("BindFragDataLocation")
	gl.LinkProgram(programID)
	glLogErr("LinkProgram")

	var linkstatus int32
	gl.GetProgramiv(programID, gl.LINK_STATUS, &linkstatus)
	if linkstatus != 1 {
		return prog, fmt.Errorf("linking program: %s", glGetProgramInfoLog(programID))
	}

	gl.UseProgram(programID)
	glLogErr("UseProgram")

	return program{
		id:         glProgram(programID),
		attributes: inspectAttributes(programID),
		uniforms:   inspectUniforms(programID),
	}, nil
}

// inspectProgramAttribs queries OpenGL for information on the vertex attributes
// of a compiled GLSL program. "id" is the ID of the compiled program to inspect.
func inspectAttributes(id uint32) []attribute {

	// Get the number of active uniforms.
	var numAttribs int32
	gl.GetProgramiv(id, gl.ACTIVE_ATTRIBUTES, &numAttribs)

	// Get maximum uniform name length.
	var maxnamelen int32
	gl.GetProgramiv(id, gl.ACTIVE_ATTRIBUTE_MAX_LENGTH, &maxnamelen)

	attribs := make([]attribute, 0, numAttribs)
	for i := int32(0); i < numAttribs; i++ {

		var xtype uint32
		var size, namelen int32
		namebytes := make([]byte, maxnamelen)
		gl.GetActiveAttrib(id, uint32(i), maxnamelen,
			&namelen, &size, &xtype, &namebytes[0])

		name := string(namebytes[:namelen])
		loc := gl.GetAttribLocation(id, gl.Str(name+"\x00"))

		var bytesize int
		var components int32
		var datatype uint32

		// TODO missing int attribute types?
		switch xtype {
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
			Name:       name,
			Size:       size,
			Type:       glEnum(xtype),
			Loc:        glAttribute(loc),
			Idx:        i,
			ByteSize:   bytesize,
			Components: components,
			DataType:   glEnum(datatype),
		})
	}
	return attribs
}

// inspectUniforms queries OpenGL for information on the uniforms
// of a compiled GLSL program. "id" is the OpenGL handle of the program
// to inspect.
func inspectUniforms(id uint32) []uniform {

	// Get the number of active uniforms.
	var numUnis int32
	gl.GetProgramiv(id, gl_ACTIVE_UNIFORMS, &numUnis)

	// Get maximum uniform name length.
	var maxnamelen int32
	gl.GetProgramiv(id, gl.ACTIVE_UNIFORM_MAX_LENGTH, &maxnamelen)

	texOffset := gl_TEXTURE0

	unis := make([]uniform, 0, numUnis)
	for i := int32(0); i < numUnis; i++ {

		var xtype uint32
		var size, namelen int32
		namebytes := make([]byte, maxnamelen)
		gl.GetActiveUniform(id, uint32(i), maxnamelen,
			&namelen, &size, &xtype, &namebytes[0])

		name := string(namebytes[:namelen])
		loc := gl.GetUniformLocation(id, gl.Str(name+"\x00"))

		var texUnit glTextureUnit
		if xtype == gl_SAMPLER_2D {
			texUnit = glTextureUnit(texOffset)
			texOffset++
		}

		unis = append(unis, uniform{
			Name:        name,
			Size:        size,
			Type:        glEnum(xtype),
			Loc:         glUniform(loc),
			Idx:         i,
			TextureUnit: texUnit,
		})
	}
	return unis
}
