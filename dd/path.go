package dd

func NewPath(xys ...XY) *Path {
	p := &Path{}
	if len(xys) == 0 {
		return p
	}
	p.MoveTo(xys[0])
	for i := 1; i < len(xys); i++ {
		p.LineTo(xys[i])
	}
	return p
}

type Path struct {
	prev     XY
	current  *subpath
	subpaths []*subpath
}

func (p *Path) MoveTo(xy XY) {
	p.prev = xy
	p.current = nil
}

func (p *Path) Move(xy XY) {
	p.MoveTo(p.prev.Add(xy))
}

func (p *Path) LineTo(xy XY) {
	s := Line{p.prev, xy}
	p.pushSegment(s)
	p.prev = xy
}

func (p *Path) Line(xy XY) {
	p.LineTo(p.prev.Add(xy))
}

func (p *Path) CubicTo(xy, a, b XY) {
	s := Cubic{p.prev, xy, a, b}
	p.pushSegment(s)
	p.prev = xy
}

func (p *Path) Cubic(xy, a, b XY) {
	p.CubicTo(
		p.prev.Add(xy),
		p.prev.Add(a),
		p.prev.Add(b),
	)
}

func (p *Path) QuadraticTo(xy, a XY) {
	s := Quadratic{p.prev, xy, a}
	p.pushSegment(s)
	p.prev = xy
}

func (p *Path) Quadratic(xy, a XY) {
	p.QuadraticTo(
		p.prev.Add(xy),
		p.prev.Add(a),
	)
}

func (p *Path) Close() {
	// TODO if first line was a curve, would calling close
	//      draw a straight line? if so, this check could change.
	if p.current == nil || len(p.current.segments) < 2 {
		return
	}
	first := p.current.segments[0].Start()
	p.LineTo(first)
}

// TODO subdivide options
func (p *Path) Lines() []Line {
	var lines []Line
	for _, sub := range p.subpaths {
		for _, seg := range sub.segments {
			switch z := seg.(type) {
			case Line:
				lines = append(lines, z)
			case Cubic:
				// TODO
				lines = append(lines, Line{z.Start(), z.End()})
			case Quadratic:
				xys := Subdivide(z, 10)
				ls := Connect(xys...)
				lines = append(lines, ls...)
			default:
				panic("unknown segment type")
			}
		}
	}
	return lines
}

// TODO closed path seems to have a glitch at the final miter joint?
func (p *Path) Stroke() Stroke {
	return Stroke{Lines: p.Lines()}
}

type subpath struct {
	segments []Curve
}

func (p *Path) pushSegment(s Curve) {
	if p.current == nil {
		p.current = &subpath{}
		p.subpaths = append(p.subpaths, p.current)
	}
	p.current.segments = append(p.current.segments, s)
}
