package gfx

type Shader struct {
	Name       string
	Vert, Frag string
	Output     string
	Mesh       Meshable
	Attrs      Attrs
	Divisors   map[string]int
	Instances  int
	Blend      Blend
}

type Attrs map[string]interface{}

func NewShader(m Meshable) *Shader {
	return &Shader{
		Vert:  DefaultVert,
		Frag:  DefaultFrag,
		Mesh:  m,
		Attrs: Attrs{},
	}
}

func (s *Shader) Draw(l Layer) {
	l.AddShader(s)
}

func (s *Shader) Set(name string, val interface{}) {
	if s.Attrs == nil {
		s.Attrs = Attrs{}
	}
	s.Attrs[name] = val
}
