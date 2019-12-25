package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
)

func Ink(z *Layer) {

	p := Path{}
	p.MoveTo(XY{0.2, 0.2})
	p.LineTo(XY{0.2, 0.3})
	p.LineTo(XY{0.3, 0.3})
	p.LineTo(XY{0.5, 0.7})
	p.LineTo(XY{0.7, 0.7})
	p.LineTo(XY{0.7, 0.9})
	p.LineTo(XY{0.5, 0.9})
	p.LineTo(XY{0.3, 0.7})

	m := p.Stroke(0.01)
	f := Fill{m, Red}
	z.Draw(f)
}
