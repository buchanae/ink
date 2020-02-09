package gfx

import (
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
)

type Fill struct {
	Shape dd.Fillable
	Color color.RGBA
}

func (f Fill) Shader() *Shader {
	return &Shader{
		Name: "Fill",
		Vert: DefaultVert,
		Frag: DefaultFrag,
		Mesh: f.Shape.Fill(),
		Attrs: Attrs{
			"a_color": f.Color,
		},
	}
}

func (f Fill) Draw(l Layer) {
	f.Shader().Draw(l)
}
