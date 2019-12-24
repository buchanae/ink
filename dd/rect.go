package dd

func RectAWH(a, wh XY) Rect {
	return Rect{
		A: a,
		B: a.Add(wh),
	}
}

func Square(xy XY, size float32) Rect {
	return Rect{
		A: xy.Add(XY{size / 2, size / 2}),
		B: xy.Add(XY{size, size}),
	}
}

type Rect struct {
	A, B XY
}

func (r Rect) Center() XY {
	return r.B.Sub(r.A).MulScalar(0.5).Add(r.A)
	//TODO wish ((r.B - r.A) * 0.5) + r.A
}

// Include will change the min/max bounds of the box
// if x/y is greater/less than the current min/max points.
// This is useful for updating bounding boxes.
func (r Rect) Include(xy XY) Rect {
	if xy.X < r.A.X {
		r.A.X = xy.X
	}
	if xy.Y < r.A.Y {
		r.A.Y = xy.Y
	}
	if xy.X > r.B.X {
		r.B.X = xy.X
	}
	if xy.Y < r.B.Y {
		r.B.Y = xy.Y
	}
	return r
}

// TODO should "amount" be a percentage?
func (r Rect) Shrink(amount float32) Rect {
	return Rect{
		A: XY{r.A.X + amount, r.A.Y + amount},
		B: XY{r.B.X - amount, r.B.Y - amount},
	}
}

func (r Rect) Quad() Quad {
	return Quad{
		XY{r.A.X, r.A.Y},
		XY{r.A.X, r.B.Y},
		XY{r.B.X, r.B.Y},
		XY{r.B.X, r.A.Y},
	}
}

func (r Rect) Grow(amount float32) Rect {
	return Rect{
		A: XY{r.A.X - amount, r.A.Y - amount},
		B: XY{r.B.X + amount, r.B.Y + amount},
	}
}

func (r Rect) Mesh() Mesh {
	return Mesh{
		Verts: []XY{
			r.A, {r.A.X, r.B.Y},
			r.B, {r.B.X, r.A.Y},
		},
		Faces: []Face{
			{0, 2, 1},
			{0, 3, 2},
		},
	}
}

func (r Rect) Stroke(width float32) Mesh {
	path := NewPath(
		r.A, XY{r.A.X, r.B.Y},
		r.B, XY{r.B.X, r.A.Y},
		r.A,
	)
	return path.Stroke(width)
}

// Bounding box
func Bounds(xys []XY) Rect {
	r := Rect{}

	if len(xys) == 0 {
		return r
	}

	r.A = xys[0]
	r.B = xys[0]

	for i := 1; i < len(xys); i++ {
		xy := xys[i]
		if xy.X < r.A.X {
			r.A.X = xy.X
		}
		if xy.Y < r.A.Y {
			r.A.Y = xy.Y
		}
		if xy.X > r.B.X {
			r.B.X = xy.X
		}
		if xy.Y < r.B.Y {
			r.B.Y = xy.Y
		}
	}
	return r
}
