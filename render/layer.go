package render

func (r *Renderer) NewLayer(s Shader) (*Layer, error) {

	prog, err := r.compile(s)
	if err != nil {
		return nil, err
	}

	l := &Layer{
		id:          s.ID,
		name:        s.Name,
		vertexCount: s.VertexCount,
		prog:        prog,
		attrs:       map[string]bindingVal{},
		uniforms:    map[string]interface{}{},
	}
	r.layers = append(r.layers, l)
	return l, nil
}

func (r *Renderer) ClearLayers() {
	r.layers = nil
}

type Shader struct {
	ID                    int
	Name                  string
	Vert, Frag, Geom, Out string
	VertexCount           int
}

type Layer struct {
	id          int
	name        string
	prog        compiled
	vertexCount int
	faces       []uint32
	attrs       map[string]bindingVal
	uniforms    map[string]interface{}
}

func (l *Layer) ID() int {
	return l.id
}

func (l *Layer) UniformNames() []string {
	var names []string
	for _, uni := range l.prog.uniforms {
		names = append(names, uni.Name)
	}
	return names
}

func (l *Layer) AttrNames() []string {
	var names []string
	for _, attr := range l.prog.attributes {
		names = append(names, attr.Name)
	}
	return names
}

func (l *Layer) SetFaces(f []uint32) {
	l.faces = f
}

func (l *Layer) SetUniform(key string, val interface{}) {
	l.uniforms[key] = val
}

func (l *Layer) SetAttr(key string, val interface{}, bytes int) {
	// skip empty attributes to avoid panics
	if bytes == 0 {
		return
	}
	l.attrs[key] = bindingVal{
		value: glPtr(val),
		size:  bytes,
	}
}
