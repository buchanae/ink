// +build js

package render

import (
	"log"
	"syscall/js"
)

var gl, vaoExt, instExt js.Value
var (
	gl_COLOR_BUFFER_BIT uint32
)

var gl_SCREEN glFramebuffer

var (
	gl_NO_ERROR,
	gl_INVALID_ENUM,
	gl_INVALID_VALUE,
	gl_INVALID_OPERATION,
	gl_STACK_OVERFLOW,
	gl_STACK_UNDERFLOW,
	gl_OUT_OF_MEMORY,
	gl_INVALID_FRAMEBUFFER_OPERATION,

	gl_INFO_LOG_LENGTH,

	gl_FRAGMENT_SHADER,
	gl_VERTEX_SHADER,
	gl_COMPILE_STATUS,
	gl_LINK_STATUS,

	gl_ACTIVE_ATTRIBUTES,
	gl_ACTIVE_ATTRIBUTE_MAX_LENGTH,
	gl_FLOAT,
	gl_FLOAT_VEC2,
	gl_FLOAT_VEC3,
	gl_FLOAT_VEC4,
	gl_FRAMEBUFFER,
	gl_TEXTURE_2D,
	gl_RGBA,
	gl_UNSIGNED_BYTE,
	gl_TEXTURE_MAG_FILTER,
	gl_TEXTURE_MIN_FILTER,
	gl_LINEAR,
	gl_REPEAT,
	gl_TEXTURE_WRAP_S,
	gl_TEXTURE_WRAP_T,
	gl_COLOR_ATTACHMENT0,
	gl_TEXTURE_2D_MULTISAMPLE,
	gl_DRAW_FRAMEBUFFER,
	gl_READ_FRAMEBUFFER,
	gl_ELEMENT_ARRAY_BUFFER,
	gl_STATIC_DRAW,
	gl_ARRAY_BUFFER,
	gl_INT,
	gl_UNSIGNED_INT,
	gl_BOOL,
	gl_SAMPLER_2D,
	gl_TEXTURE0,
	gl_ACTIVE_UNIFORM_MAX_LENGTH,
	gl_ACTIVE_UNIFORMS,
	gl_SRC_ALPHA,
	gl_ONE_MINUS_SRC_ALPHA,
	gl_BLEND,
	gl_MULTISAMPLE,
	gl_TRIANGLES glEnum
)

func glCall(name string, args ...interface{}) js.Value {
	v := gl.Call(name, args...)
	glLogErr(name)
	return v
}

func init() {

	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "gocanvas")

	//width = doc.Get("body").Get("clientWidth").Int()
	//height = doc.Get("body").Get("clientHeight").Int()
	//canvasEl.Set("width", width)
	//canvasEl.Set("height", height)

	gl = canvasEl.Call("getContext", "webgl")
	if gl == js.Undefined() {
		gl = canvasEl.Call("getContext", "experimental-webgl")
	}
	// once again
	if gl == js.Undefined() {
		js.Global().Call("alert", "browser might not support webgl")
		return
	}

	gl_COLOR_BUFFER_BIT = uint32(gl.Get("COLOR_BUFFER_BIT").Int())
	gl_SCREEN = glFramebuffer(js.ValueOf(nil))

	getEnum := func(name string) glEnum {
		return glEnum(gl.Get(name))
	}

	gl_NO_ERROR = getEnum("NO_ERROR")
	gl_INVALID_ENUM = getEnum("INVALID_ENUM")
	gl_INVALID_VALUE = getEnum("INVALID_VALUE")
	gl_INVALID_OPERATION = getEnum("INVALID_OPERATION")
	gl_STACK_OVERFLOW = getEnum("STACK_OVERFLOW")
	gl_STACK_UNDERFLOW = getEnum("STACK_UNDERFLOW")
	gl_OUT_OF_MEMORY = getEnum("OUT_OF_MEMORY")
	gl_INVALID_FRAMEBUFFER_OPERATION = getEnum("INVALID_FRAMEBUFFER_OPERATION")
	gl_INFO_LOG_LENGTH = getEnum("INFO_LOG_LENGTH")

	gl_FRAGMENT_SHADER = getEnum("FRAGMENT_SHADER")
	gl_VERTEX_SHADER = getEnum("VERTEX_SHADER")
	gl_COMPILE_STATUS = getEnum("COMPILE_STATUS")
	gl_LINK_STATUS = getEnum("LINK_STATUS")

	gl_ACTIVE_ATTRIBUTES = getEnum("ACTIVE_ATTRIBUTES")
	gl_ACTIVE_ATTRIBUTE_MAX_LENGTH = getEnum("ACTIVE_ATTRIBUTE_MAX_LENGTH")
	gl_FLOAT = getEnum("FLOAT")
	gl_FLOAT_VEC2 = getEnum("FLOAT_VEC2")
	gl_FLOAT_VEC3 = getEnum("FLOAT_VEC3")
	gl_FLOAT_VEC4 = getEnum("FLOAT_VEC4")
	gl_FRAMEBUFFER = getEnum("FRAMEBUFFER")
	gl_TEXTURE_2D = getEnum("TEXTURE_2D")
	gl_RGBA = getEnum("RGBA")
	gl_UNSIGNED_BYTE = getEnum("UNSIGNED_BYTE")
	gl_TEXTURE_MAG_FILTER = getEnum("TEXTURE_MAG_FILTER")
	gl_TEXTURE_MIN_FILTER = getEnum("TEXTURE_MIN_FILTER")
	gl_LINEAR = getEnum("LINEAR")
	gl_REPEAT = getEnum("REPEAT")
	gl_TEXTURE_WRAP_S = getEnum("TEXTURE_WRAP_S")
	gl_TEXTURE_WRAP_T = getEnum("TEXTURE_WRAP_T")
	gl_COLOR_ATTACHMENT0 = getEnum("COLOR_ATTACHMENT0")
	gl_TEXTURE_2D_MULTISAMPLE = getEnum("TEXTURE_2D_MULTISAMPLE")
	gl_DRAW_FRAMEBUFFER = getEnum("DRAW_FRAMEBUFFER")
	gl_READ_FRAMEBUFFER = getEnum("READ_FRAMEBUFFER")
	gl_ELEMENT_ARRAY_BUFFER = getEnum("ELEMENT_ARRAY_BUFFER")
	gl_STATIC_DRAW = getEnum("STATIC_DRAW")
	gl_ARRAY_BUFFER = getEnum("ARRAY_BUFFER")
	gl_INT = getEnum("INT")
	gl_UNSIGNED_INT = getEnum("UNSIGNED_INT")
	gl_BOOL = getEnum("BOOL")
	gl_SAMPLER_2D = getEnum("SAMPLER_2D")
	gl_TEXTURE0 = getEnum("TEXTURE0")
	gl_ACTIVE_UNIFORM_MAX_LENGTH = getEnum("ACTIVE_UNIFORM_MAX_LENGTH")
	gl_ACTIVE_UNIFORMS = getEnum("ACTIVE_UNIFORMS")
	gl_SRC_ALPHA = getEnum("SRC_ALPHA")
	gl_ONE_MINUS_SRC_ALPHA = getEnum("ONE_MINUS_SRC_ALPHA")
	gl_BLEND = getEnum("BLEND")
	gl_MULTISAMPLE = getEnum("MULTISAMPLE")
	gl_TRIANGLES = getEnum("TRIANGLES")

	vaoExt = glCall("getExtension", "OES_vertex_array_object")
	instExt = glCall("getExtension", "ANGLE_instanced_arrays")
	uintExt := glCall("getExtension", "OES_element_index_uint")

	if !vaoExt.Truthy() {
		log.Print("error: getting vao extension")
	}
	if !instExt.Truthy() {
		log.Print("error: getting inst extension")
	}
	if !uintExt.Truthy() {
		log.Print("error: getting uint extension")
	}
}
