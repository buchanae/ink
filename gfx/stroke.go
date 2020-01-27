package gfx

import (
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
)

type Stroke struct {
	Target Strokeable
	Width  float32
	Color  color.RGBA
}

type Strokeable interface {
	Stroke(dd.StrokeOpt) dd.Mesh
}

func (s Stroke) Draw(out Layer) {
	out.AddShader(&Shader{
		Name: "Stroke",
		Vert: DefaultVert,
		Frag: DefaultFrag,
		Mesh: s.Target.Stroke(dd.StrokeOpt{
			Width: s.Width,
		}),
		Attrs: Attrs{
			"a_color": s.Color,
		},
	})
}
