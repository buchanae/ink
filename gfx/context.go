package gfx

import (
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
)

type Context struct {
	Output      Layer
	FillColor   color.RGBA
	StrokeColor color.RGBA
	StrokeWidth float32
}

func NewContext(out Layer) Context {
	return Context{
		Output:      out,
		FillColor:   color.Black,
		StrokeColor: color.Black,
		StrokeWidth: 0.001,
	}
}

func (ctx Context) Fill(m dd.Fillable) {
	Fill{m, ctx.FillColor}.Draw(ctx.Output)
}

func (ctx Context) Clear(c color.RGBA) {
	Clear(ctx.Output, c)
}

func (ctx Context) Stroke(s dd.Strokeable) {
	Stroke{
		Shape: s,
		Width: ctx.StrokeWidth,
		Color: ctx.StrokeColor,
	}.Draw(ctx.Output)
}
