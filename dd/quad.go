package dd

type Quad struct {
	A, B, C, D XY
}

func (q Quad) Path() Path {
	return Path{
		Line{q.A, q.B},
		Line{q.B, q.C},
		Line{q.C, q.D},
		Line{q.D, q.A},
	}
}

func (q Quad) Stroke(opt StrokeOpt) Mesh {
	return Stroke(q.Path(), opt)
}

func (q Quad) Fill() Mesh {
	return Triangles{
		{q.A, q.B, q.C},
		{q.A, q.C, q.D},
	}.Fill()
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

func (q Quad) Rotate(angle float32) Quad {
	return q.RotateAround(angle, q.Centroid())
}

func (q Quad) RotateAround(angle float32, pivot XY) Quad {
	return Quad{
		q.A.RotateAround(angle, pivot),
		q.B.RotateAround(angle, pivot),
		q.C.RotateAround(angle, pivot),
		q.D.RotateAround(angle, pivot),
	}
}

func (q Quad) Translate(xy XY) Quad {
	return Quad{
		q.A.Add(xy),
		q.B.Add(xy),
		q.C.Add(xy),
		q.D.Add(xy),
	}
}
