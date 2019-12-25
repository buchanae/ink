package dd

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

func (a XY) SetLength(s float32) XY {
	l := a.Length()
	return a.Normalize().MulScalar(s / l)
}

func (a XY) Distance(b XY) float32 {
	return b.Sub(a).Length()
}

func (a XY) Dot(b XY) float32 {
	return a.X*b.X + a.Y*b.Y
}

func (a XY) Rotate(rad float32, pivot XY) XY {
	cr := cos(rad)
	sr := sin(rad)
	p := a.Sub(pivot)
	b := XY{
		X: p.X*cr - p.Y*sr,
		Y: p.X*sr + p.Y*cr,
	}
	return b.Add(pivot)
}
