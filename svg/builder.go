package svg

import (
	"strconv"

	"github.com/buchanae/ink/dd"
)

type builder struct {
	width, height float32
	stack         []float32
	pen           dd.Pen
}

func (p *builder) MoveTo(abs bool) {
	for _, xy := range p.pop() {
		if abs {
			p.pen.MoveTo(xy)
		} else {
			p.pen.Move(xy)
		}
	}
}

func (p *builder) LineTo(abs bool) {
	for _, xy := range p.pop() {
		if abs {
			p.pen.LineTo(xy)
		} else {
			p.pen.Line(xy)
		}
	}
}

func (p *builder) CubicTo(abs bool) {
	xys := p.pop()
	for i := 0; i+2 < len(xys); i += 3 {
		a := xys[i]
		b := xys[i+1]
		c := xys[i+2]
		if abs {
			p.pen.CubicTo(c, a, b)
		} else {
			p.pen.Cubic(c, a, b)
		}
	}
}

func (p *builder) QuadraticTo(abs bool) {
	xys := p.pop()
	for i := 0; i+1 < len(xys); i += 2 {
		a := xys[i]
		b := xys[i+1]
		if abs {
			p.pen.QuadraticTo(b, a)
		} else {
			p.pen.Quadratic(b, a)
		}
	}
}

func (p *builder) ClosePath() {
	p.pen.Close()
}

func (p *builder) pop() []dd.XY {
	var points []dd.XY
	for i := 0; i < len(p.stack); i += 2 {
		a := p.stack[i]
		b := p.stack[i+1]
		xy := dd.XY{a / p.width, b / p.height}
		points = append(points, xy)
	}
	p.stack = nil
	return points
}

func (p *builder) Coord(s string) {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		// TODO
		panic(err)
	}
	p.stack = append(p.stack, float32(v))
}
