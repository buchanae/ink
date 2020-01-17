package dd

type Triangle struct {
	A, B, C XY
}

// TODO http://jwilson.coe.uga.edu/EMAT6680/Dunbar/Assignment4/Assignment4_KD.htm
// https://www.eurekasparks.org/2019/07/centers-of-triangle.html<Paste>
func (t Triangle) Centroid() XY {
	return t.A.Add(t.B).Add(t.C).DivScalar(3)
}

func (t Triangle) Contains(p XY) bool {
	// https://stackoverflow.com/questions/2049582/how-to-determine-if-a-point-is-in-a-2d-triangle
	sign := func(a, b, c XY) float32 {
		return (a.X-c.X)*(b.Y-c.Y) - (b.X-c.X)*(a.Y-c.Y)
	}

	d1 := sign(p, t.A, t.B)
	d2 := sign(p, t.B, t.C)
	d3 := sign(p, t.C, t.A)

	hasNeg := (d1 < 0) || (d2 < 0) || (d3 < 0)
	hasPos := (d1 > 0) || (d2 > 0) || (d3 > 0)

	return !(hasNeg && hasPos)
}

func (t Triangle) Mesh() Mesh {
	return Mesh{
		Verts: []XY{t.A, t.B, t.C},
		Faces: []Face{{0, 1, 2}},
	}
}

func (t Triangle) Stroke(opt StrokeOpt) Mesh {
	return Stroke(Path{
		Line{t.A, t.B},
		Line{t.B, t.C},
		Line{t.C, t.A},
	}, opt)
}

/*
func (t Triangle) Path() *Path {
	p := &Path{}
	p.MoveTo(t.A)
	p.LineTo(t.B)
	p.LineTo(t.C)
	p.LineTo(t.A)
	return p
}
*/

func (t Triangle) Edges() []Line {
	return []Line{
		{t.A, t.B},
		{t.B, t.C},
		{t.C, t.A},
	}
}

type Triangles []Triangle

func (tris Triangles) Mesh() Mesh {
	verts := make([]XY, 0, len(tris)*3)
	faces := make([]Face, 0, len(tris))

	for _, t := range tris {
		// TODO optimize vertex reuse
		l := len(verts)
		verts = append(verts, t.A, t.B, t.C)
		faces = append(faces, Face{l, l + 1, l + 2})
	}
	return Mesh{Verts: verts, Faces: faces}
}
