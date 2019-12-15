package dd

type XY struct {
	X, Y float32
}

func (a XY) Add(b XY) XY {
	return XY{a.X + b.X, a.Y + b.Y}
}

func (a XY) AddScalar(s float32) XY {
	return XY{a.X + s, a.Y + s}
}

func (a XY) Sub(b XY) XY {
	return XY{a.X - b.X, a.Y - b.Y}
}

func (a XY) MulScalar(s float32) XY {
	return XY{a.X * s, a.Y * s}
}

func (a XY) SetLength(s float32) XY {
	l := a.Length()
	return a.Normalize().MulScalar(s / l)
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

func (a XY) Distance(b XY) float32 {
	return b.Sub(a).Length()
}

func (a XY) Dot(b XY) float32 {
	return a.X*b.X + a.Y*b.Y
}
