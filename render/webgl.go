// +build js

package render

import (
	"fmt"
	"log"
	"reflect"
	"syscall/js"
	"unsafe"
)

type glEnum js.Value
type glTexture js.Value
type glBuffer js.Value
type glFramebuffer js.Value
type glProgram js.Value
type glUniform js.Value
type glAttribute js.Value
type glTextureUnit js.Value
type glVAO js.Value

func (v glEnum) JSValue() js.Value        { return js.Value(v) }
func (v glTexture) JSValue() js.Value     { return js.Value(v) }
func (v glBuffer) JSValue() js.Value      { return js.Value(v) }
func (v glFramebuffer) JSValue() js.Value { return js.Value(v) }
func (v glProgram) JSValue() js.Value     { return js.Value(v) }
func (v glUniform) JSValue() js.Value     { return js.Value(v) }
func (v glAttribute) JSValue() js.Value   { return js.Value(v) }
func (v glTextureUnit) JSValue() js.Value { return js.Value(v) }
func (v glVAO) JSValue() js.Value         { return js.Value(v) }

func glEnable(flag glEnum) {
	glCall("enable", flag)
}

func glDisable(flag glEnum) {
	glCall("disable", flag)
}

func glUseProgram(program glProgram) {
	glCall("useProgram", program)
}

func glBlendFunc(sfactor, dfactor glEnum) {
	glCall("blendFunc", sfactor, dfactor)
}

func glClearColor(r, g, b, a float32) {
	glCall("clearColor", r, g, b, a)
}

func glClear(mask uint32) {
	glCall("clear", mask)
}

func glViewport(x, y, w, h int32) {
	log.Print("viewport", x, y, w, h)
	glCall("viewport", x, y, w, h)
}

/*

Gen



*/

func glCreateVAO() glVAO {
	return glVAO(vaoExt.Call("createVertexArrayOES"))
}

func glDeleteVAO(vao glVAO) {
	vaoExt.Call("deleteVertexArrayOES", vao)
}

func glBindVAO(vao glVAO) {
	vaoExt.Call("bindVertexArrayOES", vao)
}

func glCreateBuffer() glBuffer {
	return glBuffer(glCall("createBuffer"))
}

func glCreateTexture() glTexture {
	return glTexture(glCall("createTexture"))
}

func glCreateFramebuffer() glFramebuffer {
	return glFramebuffer(glCall("createFramebuffer"))
}

/*

Delete



*/

func glDeleteProgram(program glProgram) {
	glCall("deleteProgram", program)
}

func glDeleteTexture(texture glTexture) {
	glCall("deleteTexture", texture)
}

func glDeleteFramebuffer(fb glFramebuffer) {
	glCall("deleteFramebuffer", fb)
}

func glDeleteBuffer(buf glBuffer) {
	glCall("deleteBuffer", buf)
}

/*

Bind



*/

func glBindTexture(target glEnum, texture glTexture) {
	glCall("bindTexture", target, texture)
}

func glBindFramebuffer(target glEnum, framebuffer glFramebuffer) {
	glCall("bindFramebuffer", target, framebuffer)
}

func glBindBuffer(target glEnum, buffer glBuffer) {
	glCall("bindBuffer", target, buffer)
}

func glUniform1f(location glUniform, v0 float32) {
	glCall("uniform1f", location, v0)
}

func glUniform2f(location glUniform, v0, v1 float32) {
	glCall("uniform2f", location, v0, v1)
}

func glUniform3f(location glUniform, v0, v1, v2 float32) {
	glCall("uniform3f", location, v0, v1, v2)
}

func glUniform4f(location glUniform, v0, v1, v2, v3 float32) {
	glCall("uniform4f", location, v0, v1, v2, v3)
}

func glUniform1i(location glUniform, v0 int32) {
	glCall("uniform1i", location, v0)
}

func glUniform1ui(location glUniform, v0 uint32) {
	glCall("uniform1ui", location, v0)
}

func glActiveTexture(tu glTextureUnit) {
	glCall("activeTexture", tu)
}

func glEnableVertexAttribArray(index glAttribute) {
	glCall("enableVertexAttribArray", index)
}

func glVertexAttribDivisor(index glAttribute, divisor uint32) {
	instExt.Call("vertexAttribDivisorANGLE", index, divisor)
}

func glTexParameteri(target, pname, param glEnum) {
	glCall("texParameteri", target, pname, param)
}

func glBufferDataSize(target glEnum, size int, usage glEnum) {
}

func glBufferData(target glEnum, size int, data interface{}, usage glEnum) {
	if data == nil {
		glCall("bufferData", target, size, usage)
	} else {
		bytes := toByteSlice(data, size)
		log.Printf("DATA: %v", data)
		log.Printf("DATA: %v", bytes)
		jsbuf := js.Global().Get("Uint8Array").New(len(bytes))
		js.CopyBytesToJS(jsbuf, bytes)
		glCall("bufferData", target, jsbuf, usage)
	}
}

func glBufferSubData(target glEnum, offset int, size int, data interface{}) {
	bytes := toByteSlice(data, size)
	jsbuf := js.Global().Get("Uint8Array").New(len(bytes))
	js.CopyBytesToJS(jsbuf, bytes)
	glCall("bufferSubData", target, offset, jsbuf)
}

// TODO ANGLE extension detection
//      webgl2 detection
//      non instanced fallback or error
func glDrawElementsInstanced(
	mode glEnum,
	count int32,
	xtype glEnum,
	offset int,
	instancecount int32) {

	log.Print("draw els", mode, count, xtype, offset, instancecount)

	instExt.Call("drawElementsInstancedANGLE", mode, count, xtype, offset, instancecount)
}

func glTexImage2D(target glEnum,
	level int32,
	internalformat glEnum,
	width, height, border int32,
	format, xtype glEnum,
	pixels interface{}) {

	glCall("texImage2D", target, level, internalformat, width, height, border, format, xtype, pixels)
}

func glReadPixels(
	x, y, width, height int32,
	format, xtype glEnum,
	pixels []uint8) {
	glCall("readPixels", x, y, width, height, format, xtype, pixels)
}

func glBlitFramebuffer(
	srcX0, srcY0, srcX1, srcY1 int32,
	dstX0, dstY0, dstX1, dstY1 int32,
	mask uint32, filter glEnum) {
	glCall("blitFramebuffer", srcX0, srcY0, srcX1, srcY1, dstX0, dstY0, dstX1, dstY1, mask, filter)
}

func glVertexAttribPointer(
	index glAttribute,
	size int32,
	xtype glEnum,
	normalized bool,
	stride int32,
	byteOffset int) {
	glCall("vertexAttribPointer", index, size, xtype, normalized, stride, byteOffset)
}

func glFramebufferTexture2D(
	target glEnum,
	attachment glEnum,
	textarget glEnum,
	texture glTexture,
	level int32) {
	glCall("framebufferTexture2D", target, attachment, textarget, texture, level)
}

func glTexImage2DMultisample(
	target glEnum,
	samples int32,
	internalformat glEnum,
	width, height int32,
	fixedsamplelocations bool) {
	glCall("texImage2DMultisample", target, samples, internalformat, width, height, fixedsamplelocations)
}

/*

Errors




*/

func glLogErr(name string) {
	err := glCheckErr(name)
	if err != nil {
		log.Printf("%v", err)
	}
}

func glCheckErr(name string) error {
	err := gl.Call("getError")
	switch glEnum(err) {
	case gl_NO_ERROR:
		return nil
	case gl_INVALID_ENUM:
		return fmt.Errorf("%s: invalid enum", name)
	case gl_INVALID_VALUE:
		return fmt.Errorf("%s: invalid value", name)
	case gl_INVALID_OPERATION:
		return fmt.Errorf("%s: invalid operation", name)
	case gl_STACK_OVERFLOW:
		return fmt.Errorf("%s: stack overflow", name)
	case gl_STACK_UNDERFLOW:
		return fmt.Errorf("%s: stack underflow", name)
	case gl_OUT_OF_MEMORY:
		return fmt.Errorf("%s: out of memory", name)
	case gl_INVALID_FRAMEBUFFER_OPERATION:
		return fmt.Errorf("%s: invalid framebuffer operation", name)
	default:
		return fmt.Errorf("%s: unknown error code: %v", name, err)
	}
}

// TODO this is pretty risky and complex.
//      maybe have users of render pass in byte slices.
//      can still use unsafe for efficiency, but push unsafe usage
//      up the stack and simplify render package.
//
// https://github.com/golang/go/issues/19367 would help
func toByteSlice(s interface{}, size int) []byte {
	log.Printf("toByteslice %#v %d", s, size)

	type iface struct {
		Type, Data unsafe.Pointer
	}

	ii := *(*iface)(unsafe.Pointer(&s))
	h := (*reflect.SliceHeader)(ii.Data)
	h.Len = size
	h.Cap = size
	bs := *(*[]byte)(unsafe.Pointer(h))
	return bs
}
