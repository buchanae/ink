package dd

type Curve interface {
	Start() XY
	End() XY
	Length() float32
	Interpolate(p float32) XY
}

type Curveable interface {
	Curve() Curve
}

type Path []Curve

func (c Path) Curve() Curve {
	return c
}

func (c Path) Start() XY {
	if len(c) == 0 {
		return XY{}
	}
	return c[0].Start()
}

func (c Path) End() XY {
	if len(c) == 0 {
		return XY{}
	}
	return c[len(c)-1].End()
}

func (c Path) Length() float32 {
	var sum float32
	for _, x := range c {
		sum += x.Length()
	}
	return sum
}

func (c Path) Interpolate(p float32) XY {
	panic("TODO not implemented")
}
