package dd

type Triangle struct {
	A, B, C XY
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

func StrokeTriangles(tris []Triangle, width float32) Mesh {
	p := &Path{}
	for _, t := range tris {
		p.MoveTo(t.A)
		p.LineTo(t.B)
		p.LineTo(t.C)
		p.LineTo(t.A)
	}
	return p.Stroke(width)
}
