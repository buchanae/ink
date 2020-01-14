package dd

func NewPen() *Pen {
	return &Pen{}
}

type Pen struct {
	prev     XY
	current  *subpath
	subpaths []*subpath
}

func (p *Pen) MoveTo(xy XY) {
	p.prev = xy
	p.current = nil
}

func (p *Pen) Move(xy XY) {
	p.MoveTo(p.prev.Add(xy))
}

func (p *Pen) LineTo(xy XY) {
	s := Line{p.prev, xy}
	p.pushSegment(s)
	p.prev = xy
}

func (p *Pen) Line(xy XY) {
	p.LineTo(p.prev.Add(xy))
}

func (p *Pen) CubicTo(xy, a, b XY) {
	s := Cubic{p.prev, xy, a, b}
	p.pushSegment(s)
	p.prev = xy
}

func (p *Pen) Cubic(xy, a, b XY) {
	p.CubicTo(
		p.prev.Add(xy),
		p.prev.Add(a),
		p.prev.Add(b),
	)
}

func (p *Pen) QuadraticTo(xy, a XY) {
	s := Quadratic{p.prev, xy, a}
	p.pushSegment(s)
	p.prev = xy
}

func (p *Pen) Quadratic(xy, a XY) {
	p.QuadraticTo(
		p.prev.Add(xy),
		p.prev.Add(a),
	)
}

func (p *Pen) Close() {
	// TODO if first line was a curve, would calling close
	//      draw a straight line? if so, this check could change.
	if p.current == nil || len(p.current.segments) < 2 {
		return
	}
	first := p.current.segments[0].Start()
	p.LineTo(first)
}

func (p *Pen) Stroke(opt StrokeOpt) Mesh {
	mesh := Mesh{}
	for _, sub := range p.subpaths {
		path := Path(sub.segments)
		mesh = StrokeTo(path, mesh, opt)
	}
	return mesh
}

func (p *Pen) Paths() []Path {
	var paths []Path
	for _, sub := range p.subpaths {
		paths = append(paths, Path(sub.segments))
	}
	return paths
}

type subpath struct {
	segments []Curve
}

func (p *Pen) pushSegment(s Curve) {
	if p.current == nil {
		p.current = &subpath{}
		p.subpaths = append(p.subpaths, p.current)
	}
	p.current.segments = append(p.current.segments, s)
}
