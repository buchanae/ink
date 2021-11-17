package dd

type Curve interface {
	Start() XY
	End() XY
	Length() float32
	Interpolate(p float32) XY
}

func XYsToPath(xys []XY) Path {
	lines := XYsToLines(xys...)
	path := Path{}

	for _, line := range lines {
		path = append(path, line)
	}
	path = append(path, Line{xys[len(xys)-1], xys[0]})
	return path
}

type Path []Curve

func (p Path) Stroke(opt StrokeOpt) Mesh {
	return Stroke(p, opt)
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
