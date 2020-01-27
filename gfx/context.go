package gfx

import (
	"github.com/buchanae/ink/color"
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

func (ctx Context) Fill(m Meshable) {
	Fill{m, ctx.FillColor}.Draw(ctx.Output)
}

func (ctx Context) Clear(c color.RGBA) {
	Clear(ctx.Output, c)
}

func (ctx Context) Stroke(s Strokeable) {
	Stroke{
		Target: s,
		Width:  ctx.StrokeWidth,
		Color:  ctx.StrokeColor,
	}.Draw(ctx.Output)
}
