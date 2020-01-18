// +build !js

package render

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

const gl_SCREEN glFramebuffer = 0
const (
	gl_ACTIVE_ATTRIBUTES           glEnum = gl.ACTIVE_ATTRIBUTES
	gl_ACTIVE_ATTRIBUTE_MAX_LENGTH        = gl.ACTIVE_ATTRIBUTE_MAX_LENGTH
	gl_FLOAT                              = gl.FLOAT
	gl_FLOAT_VEC2                         = gl.FLOAT_VEC2
	gl_FLOAT_VEC3                         = gl.FLOAT_VEC3
	gl_FLOAT_VEC4                         = gl.FLOAT_VEC4
	gl_FRAMEBUFFER                        = gl.FRAMEBUFFER
	gl_TEXTURE_2D                         = gl.TEXTURE_2D
	gl_RGBA                               = gl.RGBA
	gl_UNSIGNED_BYTE                      = gl.UNSIGNED_BYTE
	gl_TEXTURE_MAG_FILTER                 = gl.TEXTURE_MAG_FILTER
	gl_TEXTURE_MIN_FILTER                 = gl.TEXTURE_MIN_FILTER
	gl_LINEAR                             = gl.LINEAR
	gl_REPEAT                             = gl.REPEAT
	gl_TEXTURE_WRAP_S                     = gl.TEXTURE_WRAP_S
	gl_TEXTURE_WRAP_T                     = gl.TEXTURE_WRAP_T
	gl_COLOR_ATTACHMENT0                  = gl.COLOR_ATTACHMENT0
	gl_TEXTURE_2D_MULTISAMPLE             = gl.TEXTURE_2D_MULTISAMPLE
	gl_COLOR_BUFFER_BIT                   = gl.COLOR_BUFFER_BIT
	gl_DRAW_FRAMEBUFFER                   = gl.DRAW_FRAMEBUFFER
	gl_READ_FRAMEBUFFER                   = gl.READ_FRAMEBUFFER
	gl_ELEMENT_ARRAY_BUFFER               = gl.ELEMENT_ARRAY_BUFFER
	gl_STATIC_DRAW                        = gl.STATIC_DRAW
	gl_ARRAY_BUFFER                       = gl.ARRAY_BUFFER
	gl_INT                                = gl.INT
	gl_UNSIGNED_INT                       = gl.UNSIGNED_INT
	gl_BOOL                               = gl.BOOL
	gl_SAMPLER_2D                         = gl.SAMPLER_2D
	gl_TEXTURE0                           = gl.TEXTURE0
	gl_ACTIVE_UNIFORM_MAX_LENGTH          = gl.ACTIVE_UNIFORM_MAX_LENGTH
	gl_ACTIVE_UNIFORMS                    = gl.ACTIVE_UNIFORMS
	gl_SRC_ALPHA                          = gl.SRC_ALPHA
	gl_ONE_MINUS_SRC_ALPHA                = gl.ONE_MINUS_SRC_ALPHA
	gl_BLEND                              = gl.BLEND
	gl_MULTISAMPLE                        = gl.MULTISAMPLE
	gl_TRIANGLES                          = gl.TRIANGLES
)
