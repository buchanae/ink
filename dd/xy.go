package dd

import "github.com/buchanae/ink/math"

func Unit() XY {
	return XY{1, 1}
}

type XY struct {
	X, Y float32
}

func (a XY) Add(b XY) XY {
	return XY{a.X + b.X, a.Y + b.Y}
}

func (a XY) Sub(b XY) XY {
	return XY{a.X - b.X, a.Y - b.Y}
}

func (a XY) Mul(b XY) XY {
	return XY{a.X * b.X, a.Y * b.Y}
}

func (a XY) Div(b XY) XY {
	return XY{a.X / b.X, a.Y / b.Y}
}

func (a XY) AddScalar(s float32) XY {
	return XY{a.X + s, a.Y + s}
}

func (a XY) SubScalar(s float32) XY {
	return XY{a.X - s, a.Y - s}
}

func (a XY) MulScalar(s float32) XY {
	return XY{a.X * s, a.Y * s}
}

func (a XY) DivScalar(s float32) XY {
	return XY{a.X / s, a.Y / s}
}

func (a XY) Normalize() XY {
	return a.DivScalar(a.Length())
}

func (a XY) Length() float32 {
	return sqrt(a.X*a.X + a.Y*a.Y)
}

// TODO possibly wrong
func (a XY) SetLength(s float32) XY {
	l := a.Length()
	return a.Normalize().MulScalar(s / l)
}

func (a XY) Distance(b XY) float32 {
	return b.Sub(a).Length()
}

func (a XY) DistanceToLine(l Line) float32 {
	n := a.nearestToLine(l)
	return a.Distance(n)
}

func (a XY) nearestToLine(l Line) XY {
	d := l.Length()
	x1 := l.A.X
	y1 := l.A.Y
	x2 := l.B.X
	y2 := l.B.Y
	x3 := a.X
	y3 := a.Y
	u := ((x3-x1)*(x2-x1) + (y3-y1)*(y2-y1)) / (d * d)

	// if closest point is one of the l vertices
	if u < 0 {
		return l.A
	}
	if u > 1 {
		return l.B
	}

	ix := x1 + u*(x2-x1)
	iy := y1 + u*(y2-y1)
	return XY{ix, iy}
}

func (a XY) Dot(b XY) float32 {
	return a.X*b.X + a.Y*b.Y
}

func (a XY) Rotate(angle float32) XY {
	return a.RotateAround(angle, XY{})
}

func (a XY) RotateAround(angle float32, pivot XY) XY {
	cr := cos(angle)
	sr := sin(angle)
	p := a.Sub(pivot)
	b := XY{
		X: p.X*cr - p.Y*sr,
		Y: p.X*sr + p.Y*cr,
	}
	return b.Add(pivot)
}

func (a XY) Clamp(min, max XY) XY {
	return XY{
		math.Clamp(a.X, min.X, max.X),
		math.Clamp(a.Y, min.Y, max.Y),
	}
}
