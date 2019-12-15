package render

import (
	"fmt"
)

func (r *Renderer) compile(s Shader) (compiled, error) {
	var prog compiled

	if s.Vert == "" {
		return prog, fmt.Errorf("empty vert shader")
	}
	if s.Frag == "" {
		return prog, fmt.Errorf("empty frag shader")
	}

	if p, ok := r.shaders[s]; ok {
		return p, nil
	}

	out := s.Out
	if out == "" {
		out = "f_color"
	}

	id, err := glBuildProgram(s.Vert, s.Frag, s.Geom, out)
	if err != nil {
		return prog, err
	}

	prog = compiled{
		id:         id,
		attributes: inspectAttributes(id),
		uniforms:   inspectUniforms(id),
	}
	r.shaders[s] = prog

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

/* TODO move to app
if vertFile == "" {
	vertFile = "!default.vert"
}
if fragFile == "" {
	fragFile = "!default.frag"
}
*/

// TODO caching by full source bytes is too slow,
//      because it requires processing all the source.
//      so cache by name, but then need a way to invalidate
//      when those files change.
//key := c.hash([]byte(vertFile), []byte(fragFile), nil)
/*
	key := cachekey{vertFile, fragFile}
	if p, ok := c.cache[key]; ok {
		// Program destruction uses ref counts, due to caching.
		// See Compiler.Cleanup.
		p.refs++
		return p, nil
	}
*/

/*
	vert, err := c.assets.Load(vertFile)
	if err != nil {
		return nil, err
	}

	frag, err := c.assets.Load(fragFile)
	if err != nil {
		return nil, err
	}

	geom, err := c.assets.Load(geomFile)
	if err != nil {
		return nil, err
	}
*/
