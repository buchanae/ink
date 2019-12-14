package d2

type Quadratic struct {
	A, B, Ctrl XY
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
	Segments() []Curve
	Length() float32
	Start() XY
	End() XY
	Interpolate(p float32) XY
}

func Midpoint(c Curve) XY {
	return c.Interpolate(.5)
}

func Subdivide(c Curve, n int) []XY {
	var xys []XY
	for i := 1; i <= n; i++ {
		p := float32(i) / float32(n)
		xys = append(xys, c.Interpolate(p))
	}
	return xys
}
