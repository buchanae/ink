package dd

type subpath struct {
	lines []Line
}

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
	l := Line{p.prev, xy}

	if p.current == nil {
		p.current = &subpath{lines: []Line{l}}
		p.subpaths = append(p.subpaths, p.current)
	} else {
		p.current.lines = append(p.current.lines, l)
	}
	p.prev = xy
}

func (p *Path) Line(xy XY) {
	p.LineTo(p.prev.Add(xy))
}

func (p *Path) CubicTo(xy, a, b XY) {
}

func (p *Path) Cubic(xy, a, b XY) {
}

func (p *Path) QuadraticTo(xy, a XY) {
}

//func (p *Path) ArcTo(center, radii XY, rot float32

func (p *Path) Close() {
	// TODO if first line was a curve, would calling close
	//      draw a straight line? if so, this check could change.
	if p.current == nil || len(p.current.lines) < 2 {
		return
	}
	first := p.current.lines[0].A
	p.LineTo(first)
}

/*
func (p *Path) Lines() []Line {
	out := make([]Line, len(p.lines))
	copy(out, p.lines)
	return out
}
*/

// TODO stroke options: closed, width, miter, cap, etc.
// TODO closed path seems to have a glitch at the final miter joint?
func (p *Path) Stroke(width float32) Mesh {
	var verts []XY
	var faces []Face

	for _, sub := range p.subpaths {
		tris := Stroke(sub.lines, width, false)
		for _, t := range tris {
			l := len(verts)
			verts = append(verts, t.A, t.B, t.C)
			faces = append(faces, Face{l, l + 1, l + 2})
		}
	}
	return Mesh{verts, faces}
}
