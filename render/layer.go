package render

func (r *Renderer) NewLayer(s Shader) *Layer {
	r.layerID++
	l := &Layer{
		id:       r.layerID,
		shader:   s,
		attrs:    map[string]bindingVal{},
		uniforms: map[string]interface{}{},
	}
	r.layers = append(r.layers, l)
	return l
}

func (r *Renderer) ClearLayers() {
	r.layers = nil
}

type Shader struct {
	Vert, Frag, Geom, Out string
}

type Layer struct {
	id          int
	name        string
	shader      Shader
	vertexCount int
	faces       []uint32
	attrs       map[string]bindingVal
	uniforms    map[string]interface{}
	hide        bool
}

func (l *Layer) Name(n string) {
	l.name = n
}

func (l *Layer) VertexCount(i int) {
	l.vertexCount = i
}

func (l *Layer) Faces(f []uint32) {
	l.faces = f
}

func (l *Layer) Uniform(key string, val interface{}) {
	l.uniforms[key] = val
}

func (l *Layer) FloatAttr(key string, val []float32) {
	// 4 bytes in a float32
	l.UnsafeAttr(key, val, len(val)*4)
}

func (l *Layer) UnsafeAttr(key string, val interface{}, bytes int) {
	// skip empty attributes to avoid panics
	if bytes == 0 {
		return
	}
	l.attrs[key] = bindingVal{
		value: glPtr(val),
		size:  bytes,
	}
}
