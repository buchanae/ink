package dd

type Quad struct {
	A, B, C, D XY
}

// TODO XFromY or YToX or Y2X
func QuadFromRect(r Rect) Quad {
	return Quad{
		XY{r.A.X, r.A.Y},
		XY{r.A.X, r.B.Y},
		XY{r.B.X, r.B.Y},
		XY{r.B.X, r.A.Y},
	}
}

func (b Quad) Stroke(lineWidth float32) Mesh {
	path := NewPath(b.A, b.B, b.C, b.D, b.A)
	return path.Stroke(lineWidth)
}

func (b Quad) Mesh() Mesh {
	return Triangles([]Triangle{
		{b.A, b.B, b.C},
		{b.A, b.C, b.D},
	})
}

func (b Quad) Bounds() Rect {
	return Bounds([]XY{b.A, b.B, b.C, b.D})
}
