package gfx

import (
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
)

type Stroke struct {
	Target Strokeable
	dd.StrokeOpt
	Color color.RGBA
	// TODO don't want to put this on every gfx type
	Blend Blend
}

type Strokeable interface {
	Stroke(dd.StrokeOpt) dd.Mesh
}

func (s Stroke) Draw(out Layer) {
	out.AddShader(&Shader{
		Name:  "Stroke",
		Vert:  DefaultVert,
		Frag:  DefaultFrag,
		Blend: s.Blend,
		Mesh:  s.Target.Stroke(s.StrokeOpt),
		Attrs: Attrs{
			"a_color": s.Color,
		},
	})
}
