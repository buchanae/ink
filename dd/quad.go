package dd

type Quad struct {
	A, B, C, D XY
}

/*
func (q Quad) Path() *Path {
	path := NewPath(q.A, q.B, q.C, q.D)
	path.Close()
	return path
}
*/

func (q Quad) Stroke(opt StrokeOpt) Mesh {
	return Stroke(Path{
		Line{q.A, q.B},
		Line{q.B, q.C},
		Line{q.C, q.D},
		Line{q.D, q.A},
	}, opt)
}

func (q Quad) Mesh() Mesh {
	return Triangles{
		{q.A, q.B, q.C},
		{q.A, q.C, q.D},
	}.Mesh()
}

func (q Quad) Bounds() Rect {
	return Bounds([]XY{q.A, q.B, q.C, q.D})
}

func (q Quad) Centroid() XY {
	return q.A.
		Add(q.B).
		Add(q.C).
		Add(q.D).
		DivScalar(4)
}

func (q Quad) Rotate(rad float32) Quad {
	return q.RotateAround(rad, q.Centroid())
}

func (q Quad) RotateAround(rad float32, pivot XY) Quad {
	return Quad{
		A: q.A.Rotate(rad, pivot),
		B: q.B.Rotate(rad, pivot),
		C: q.C.Rotate(rad, pivot),
		D: q.D.Rotate(rad, pivot),
	}
}
