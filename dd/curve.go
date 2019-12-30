package dd

type Cubic struct {
	A, B         XY
	CtrlA, CtrlB XY
}

func (c Cubic) Start() XY {
	return c.A
}

func (c Cubic) End() XY {
	return c.B
}

func (c Cubic) Interpolate(p float32) XY {
	panic("not implemented")
}

type Quadratic struct {
	A, B, Ctrl XY
}

func (q Quadratic) Start() XY {
	return q.A
}

func (q Quadratic) End() XY {
	return q.B
}

func (q Quadratic) Interpolate(p float32) XY {
	al := Line{q.A, q.Ctrl}
	bl := Line{q.Ctrl, q.B}
	a := al.Interpolate(p)
	b := bl.Interpolate(p)
	cl := Line{a, b}
	return cl.Interpolate(p)
}

type Curve interface {
	//Segments() []Curve
	//Length() float32
	Start() XY
	End() XY
	Interpolate(p float32) XY
}

func Midpoint(c Curve) XY {
	return c.Interpolate(.5)
}

func Subdivide(c Curve, n int) []XY {
	var xys []XY
	for i := 0; i <= n; i++ {
		p := float32(i) / float32(n)
		xys = append(xys, c.Interpolate(p))
	}
	return xys
}
