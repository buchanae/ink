package render

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

func glEnable(flag uint32) {
	gl.Enable(flag)
	glLogErr("Enable")
}

func glDisable(flag uint32) {
	gl.Disable(flag)
	glLogErr("Disable")
}

func glGetProgramiv(program uint32, pname uint32, params *int32) {
	gl.GetProgramiv(program, pname, params)
	glLogErr("GetProgramiv")
}

func glCreateShader(xtype uint32) uint32 {
	return gl.CreateShader(xtype)
}

func glShaderSource(shader uint32, count int32, xstring **uint8, length *int32) {
	gl.ShaderSource(shader, count, xstring, length)
	glLogErr("ShaderSource")
}

func glGetShaderiv(shader uint32, pname uint32, params *int32) {
	gl.GetShaderiv(shader, pname, params)
	glLogErr("GetShaderiv")
}

func glCompileShader(shader uint32) {
	gl.CompileShader(shader)
	glLogErr("CompileShader")
}

func glAttachShader(program uint32, shader uint32) {
	gl.AttachShader(program, shader)
	glLogErr("AttachShader")
}

func glCreateProgram() uint32 {
	return gl.CreateProgram()
}

func glBindFragDataLocation(program uint32, color uint32, name *uint8) {
	gl.BindFragDataLocation(program, color, name)
	glLogErr("BindFragDataLocation")
}

func glLinkProgram(program uint32) {
	gl.LinkProgram(program)
	glLogErr("LinkProgram")
}

func glGenVertexArrays(n int32, arrays *uint32) {
	gl.GenVertexArrays(n, arrays)
	glLogErr("GenVertexArrays")
}

func glBindVertexArray(array uint32) {
	gl.BindVertexArray(array)
	glLogErr("BindVertexArray")
}

func glDeleteVertexArrays(n int32, arrays *uint32) {
	gl.DeleteVertexArrays(n, arrays)
	glLogErr("DeleteVertexArrays")
}

func glDeleteBuffers(n int32, buffers *uint32) {
	gl.DeleteBuffers(n, buffers)
	glLogErr("DeleteBuffers")
}

func glGenBuffers(n int32, buffers *uint32) {
	gl.GenBuffers(n, buffers)
	glLogErr("GenBuffers")
}

func glBindBuffer(target uint32, buffer uint32) {
	gl.BindBuffer(target, buffer)
	glLogErr("BindBuffer")
}

func glEnableVertexArrayAttrib(vaobj uint32, index uint32) {
	gl.EnableVertexArrayAttrib(vaobj, index)
	glLogErr("EnableVertexArrayAttrib")
}

func glVertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer) {
	gl.VertexAttribPointer(index, size, xtype, normalized, stride, pointer)
	glLogErr("VertexAttribPointer")
}

func glVertexAttribDivisor(index uint32, divisor uint32) {
	gl.VertexAttribDivisor(index, divisor)
	glLogErr("VertexAttribDivisor")
}

func glUseProgram(program uint32) {
	gl.UseProgram(program)
	glLogErr("UseProgram")
}

func glGetActiveAttrib(program uint32, index uint32, bufSize int32, length *int32, size *int32, xtype *uint32, name *uint8) {
	gl.GetActiveAttrib(program, index, bufSize, length, size, xtype, name)
	glLogErr("GetActiveAttrib")
}

func glGetAttribLocation(program uint32, name *uint8) int32 {
	return gl.GetAttribLocation(program, name)
}

func glDeleteProgram(program uint32) {
	gl.DeleteProgram(program)
	glLogErr("DeleteProgram")
}

func glGenTextures(n int32, textures *uint32) {
	gl.GenTextures(n, textures)
	glLogErr("GenTextures")
}

func glGenFramebuffers(n int32, framebuffers *uint32) {
	gl.GenFramebuffers(n, framebuffers)
	glLogErr("GenFramebuffers")
}

func glBindTexture(target uint32, texture uint32) {
	gl.BindTexture(target, texture)
	glLogErr("BindTexture")
}

func glBindFramebuffer(target uint32, framebuffer uint32) {
	gl.BindFramebuffer(target, framebuffer)
	glLogErr("BindFramebuffer")
}

func glTexParameteri(target uint32, pname uint32, param int32) {
	gl.TexParameteri(target, pname, param)
	glLogErr("TexParameteri")
}

func glFramebufferTexture2D(target uint32, attachment uint32, textarget uint32, texture uint32, level int32) {
	gl.FramebufferTexture2D(target, attachment, textarget, texture, level)
	glLogErr("FramebufferTexture2D")
}

func glTexImage2DMultisample(target uint32, samples int32, internalformat uint32, width int32, height int32, fixedsamplelocations bool) {
	gl.TexImage2DMultisample(target, samples, internalformat, width, height, fixedsamplelocations)
	glLogErr("TexImage2DMultisample")
}

func glClearColor(red float32, green float32, blue float32, alpha float32) {
	gl.ClearColor(red, green, blue, alpha)
	glLogErr("ClearColor")
}

func glClear(mask uint32) {
	gl.Clear(mask)
	glLogErr("Clear")
}

func glBlitFramebuffer(srcX0 int32, srcY0 int32, srcX1 int32, srcY1 int32, dstX0 int32, dstY0 int32, dstX1 int32, dstY1 int32, mask uint32, filter uint32) {
	gl.BlitFramebuffer(srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1, mask, filter)
	glLogErr("BlitFramebuffer")
}

func glDeleteTextures(n int32, textures *uint32) {
	gl.DeleteTextures(n, textures)
	glLogErr("DeleteTextures")
}

func glDeleteFramebuffers(n int32, framebuffers *uint32) {
	gl.DeleteFramebuffers(n, framebuffers)
	glLogErr("DeleteFramebuffers")
}

func glUniform1f(location int32, v0 float32) {
	gl.Uniform1f(location, v0)
	glLogErr("Uniform1f")
}

func glUniform2f(location int32, v0 float32, v1 float32) {
	gl.Uniform2f(location, v0, v1)
	glLogErr("Uniform2f")
}

func glUniform3f(location int32, v0 float32, v1 float32, v2 float32) {
	gl.Uniform3f(location, v0, v1, v2)
	glLogErr("Uniform3f")
}

func glUniform4f(location int32, v0 float32, v1 float32, v2 float32, v3 float32) {
	gl.Uniform4f(location, v0, v1, v2, v3)
	glLogErr("Uniform4f")
}

func glUniform1i(location int32, v0 int32) {
	gl.Uniform1i(location, v0)
	glLogErr("Uniform1i")
}

func glUniform1ui(location int32, v0 uint32) {
	gl.Uniform1ui(location, v0)
	glLogErr("Uniform1ui")
}

func glActiveTexture(texture uint32) {
	gl.ActiveTexture(texture)
	glLogErr("ActiveTexture")
}

func glGetActiveUniform(program uint32, index uint32, bufSize int32, length *int32, size *int32, xtype *uint32, name *uint8) {
	gl.GetActiveUniform(program, index, bufSize, length, size, xtype, name)
	glLogErr("GetActiveUniform")
}

func glGetUniformLocation(program uint32, name *uint8) int32 {
	return gl.GetUniformLocation(program, name)
}

func glBufferData(target uint32, size int, data unsafe.Pointer, usage uint32) {
	gl.BufferData(target, size, data, usage)
	glLogErr("BufferData")
}

func glBufferSubData(target uint32, offset int, size int, data unsafe.Pointer) {
	gl.BufferSubData(target, offset, size, data)
	glLogErr("BufferSubData")
}

func glEnableVertexAttribArray(index uint32) {
	gl.EnableVertexAttribArray(index)
	glLogErr("EnableVertexAttribArray")
}

func glPtr(data interface{}) unsafe.Pointer {
	return gl.Ptr(data)
}

func glPtrOffset(offset int) unsafe.Pointer {
	return gl.PtrOffset(offset)
}

func glFinish() {
	gl.Finish()
}

func glStr(str string) *uint8 {
	return gl.Str(str)
}

func glTexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
	gl.TexImage2D(target, level, internalformat, width, height, border, format, xtype, pixels)
	glLogErr("TexImage2D")
}

func glBlendFunc(sfactor uint32, dfactor uint32) {
	gl.BlendFunc(sfactor, dfactor)
	glLogErr("BlendFunc")
}

func glViewport(x int32, y int32, width int32, height int32) {
	gl.Viewport(x, y, width, height)
	glLogErr("Viewport")
}

func glDrawArrays(mode uint32, first int32, count int32) {
	gl.DrawArrays(mode, first, count)
	glLogErr("DrawArrays")
}

func glDrawElements(mode uint32, count int32, xtype uint32, indices unsafe.Pointer) {
	gl.DrawElements(mode, count, xtype, indices)
	glLogErr("DrawElements")
}

func glDrawElementsInstanced(mode uint32, count int32, xtype uint32, indices unsafe.Pointer, instancecount int32) {
	gl.DrawElementsInstanced(mode, count, xtype, indices, instancecount)
	glLogErr("DrawElementsInstanced")
}

func glLineWidth(width float32) {
	gl.LineWidth(width)
	glLogErr("LineWidth")
}

func glDrawArraysInstanced(mode uint32, first int32, count int32, instancecount int32) {
	gl.DrawArraysInstanced(mode, first, count, instancecount)
	glLogErr("DrawArraysInstanced")
}

func glReadPixels(x int32, y int32, width int32, height int32, format uint32, xtype uint32, pixels unsafe.Pointer) {
	gl.ReadPixels(x, y, width, height, format, xtype, pixels)
	glLogErr("ReadPixels")
}

func glGoStr(cstr *uint8) string {
	return gl.GoStr(cstr)
}

func glStrs(strs ...string) (cstrs **uint8, free func()) {
	return gl.Strs(strs...)
}

func glLogErr(name string) {
	err := glCheckErr(name)
	if err != nil {
		log("%v", err)
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
	glGetProgramiv(id, gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)

	gl.GetProgramInfoLog(id, logLength, nil, &logBuffer[0])
	glLogErr("GetProgramInfoLog")

	return glGoStr(&logBuffer[0])
}

func glGetShaderInfoLog(id uint32) string {
	var logLength int32

	glGetShaderiv(id, gl.INFO_LOG_LENGTH, &logLength)
	if logLength == 0 {
		return ""
	}

	logBuffer := make([]uint8, logLength)
	gl.GetShaderInfoLog(id, logLength, nil, &logBuffer[0])
	glLogErr("GetShaderInfoLog")

	return glGoStr(&logBuffer[0])
}

func glBuildShader(src string, shaderType uint32) (uint32, error) {
	id := glCreateShader(shaderType)
	source, free := gl.Strs(src + "\000")
	defer free()

	glShaderSource(id, 1, source, nil)
	glCompileShader(id)

	var status int32
	glGetShaderiv(id, gl.COMPILE_STATUS, &status)

	if status != 1 {
		return 0, fmt.Errorf("compiling shader: %s", glGetShaderInfoLog(id))
	}
	return id, nil
}

func glBuildProgram(vert, frag, geom, out string) (uint32, error) {
	// VERTEX SHADER
	vs, err := glBuildShader(vert, gl.VERTEX_SHADER)
	if err != nil {
		return 0, fmt.Errorf("creating vert shader: %s", err)
	}

	// FRAGMENT SHADER
	fs, err := glBuildShader(frag, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, fmt.Errorf("creating frag shader: %s", err)
	}

	// CREATE PROGRAM
	programID := glCreateProgram()
	glAttachShader(programID, vs)
	glAttachShader(programID, fs)

	// GEOMETRY SHADER
	if geom != "" {
		gs, err := glBuildShader(geom, gl.GEOMETRY_SHADER)
		if err != nil {
			return 0, fmt.Errorf("creating geo shader: %s", err)
		}
		glAttachShader(programID, gs)
	}

	glBindFragDataLocation(programID, 0, glStr(out+"\000"))
	glLinkProgram(programID)

	var linkstatus int32
	glGetProgramiv(programID, gl.LINK_STATUS, &linkstatus)
	if linkstatus != 1 {
		return 0, fmt.Errorf("linking program: %s", glGetProgramInfoLog(programID))
	}

	return programID, nil
}
