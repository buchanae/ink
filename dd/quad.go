package dd

type Quad struct {
	A, B, C, D XY
}

func (q Quad) Stroke() Stroke {
	path := NewPath(q.A, q.B, q.C, q.D)
	path.Close()
	stroke := path.Stroke()
	stroke.Closed = true
	return stroke
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

/*
TODO what is the quad center?
func (q Quad) Rotate(rad float32) Quad {
	return q.RotateAround(rad, q.Center())
}
*/

func (q Quad) RotateAround(rad float32, pivot XY) Quad {
	return Quad{
		A: q.A.Rotate(rad, pivot),
		B: q.B.Rotate(rad, pivot),
		C: q.C.Rotate(rad, pivot),
		D: q.D.Rotate(rad, pivot),
	}
}
