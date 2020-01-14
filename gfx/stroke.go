package gfx

import (
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
)

type Stroke struct {
	Target Strokeable
	dd.StrokeOpt
	Color color.RGBA
}

type Strokeable interface {
	Stroke(dd.StrokeOpt) dd.Mesh
}

func (s Stroke) Draw(out Layer) {
	out.AddShader(&Shader{
		Name: "Stroke",
		Vert: DefaultVert,
		Frag: DefaultFrag,
		Mesh: s.Target.Stroke(s.StrokeOpt),
		Attrs: Attrs{
			"a_color": s.Color,
		},
	})
}
