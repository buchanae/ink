package gfx

import (
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
)

type Stroke struct {
	Shape dd.Strokeable
	Width float32
	Color color.RGBA
}

func (s Stroke) Draw(out Layer) {
	out.AddShader(&Shader{
		Name: "Stroke",
		Vert: DefaultVert,
		Frag: DefaultFrag,
		Mesh: s.Shape.Stroke(dd.StrokeOpt{
			Width: s.Width,
		}),
		Attrs: Attrs{
			"a_color": s.Color,
		},
	})
}
