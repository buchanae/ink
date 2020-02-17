package opengl

import "errors"

// TODO would be better if the cache was per GL context
var shaderCache map[shaderOpt]compiled

type shaderOpt struct {
	vert, frag, geom, out string
}

func init() {
	shaderCache = map[shaderOpt]compiled{}
}

func compile(s shaderOpt) (compiled, error) {
	var prog compiled

	if s.vert == "" {
		return prog, errors.New("empty vert shader")
	}
	if s.frag == "" {
		return prog, errors.New("empty frag shader")
	}

	if p, ok := shaderCache[s]; ok {
		return p, nil
	}

	out := s.out
	if out == "" {
		out = "f_color"
	}

	id, err := glBuildProgram(s.vert, s.frag, s.geom, out)
	if err != nil {
		return prog, err
	}

	prog = compiled{
		id:         id,
		attributes: inspectAttributes(id),
		uniforms:   inspectUniforms(id),
	}
	shaderCache[s] = prog

	return prog, nil
}

// compiled contains information about a compiled GLSL program.
type compiled struct {
	// ID is the OpenGL program ID.
	id uint32
	// Uniforms contains static information about the uniforms
	// defined in the program's shader code.
	attributes []attribute
	// Attributes contains static information about the attributes
	// defined in the program's shader code.
	uniforms []uniform
}
