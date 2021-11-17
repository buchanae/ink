package dd

func RectWH(w, h float32) Rect {
	return Rect{
		B: XY{w, h},
	}
}

func RectAWH(a, wh XY) Rect {
	return Rect{
		A: a,
		B: a.Add(wh),
	}
}

func RectCenter(center, wh XY) Rect {
	half := wh.MulScalar(0.5)
	return Rect{
		A: center.Sub(half),
		B: center.Add(half),
	}
}

func Square(xy XY, size float32) Rect {
	return Rect{
		A: xy,
		B: xy.Add(XY{size, size}),
	}
}
func SquareCenter(center XY, size float32) Rect {
	half := size / 2
	return Rect{
		A: center.SubScalar(half),
		B: center.AddScalar(half),
	}
}

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
		if xy.Y > r.B.Y {
			r.B.Y = xy.Y
		}
	}
	return r
}

type Rect struct {
	A, B XY
}

var zeroRect Rect

func (r Rect) IsZero() bool {
	return r == zeroRect
}

func (r Rect) Center() XY {
	return r.B.Sub(r.A).MulScalar(0.5).Add(r.A)
}

func (r Rect) Area() float32 {
	size := r.Size()
	return size.X * size.Y
}

func (r Rect) Size() XY {
	return r.B.Sub(r.A)
}

func (r Rect) Contains(xy XY) bool {
	if xy.X < r.A.X {
		return false
	}
	if xy.Y < r.A.Y {
		return false
	}
	if xy.X > r.B.X {
		return false
	}
	if xy.Y > r.B.Y {
		return false
	}
	return true
}

func (r Rect) Translate(xy XY) Rect {
	return Rect{
		A: r.A.Add(xy),
		B: r.B.Add(xy),
	}
}

// TODO note, "xy" is a percentage
func (r Rect) Interpolate(xy XY) XY {
	return r.B.Sub(r.A).Mul(xy).Add(r.A)
}

func (r Rect) Rotate(angle float32) Quad {
	return r.Quad().RotateAround(angle, r.Center())
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
	if xy.Y > r.B.Y {
		r.B.Y = xy.Y
	}
	return r
}

// TODO should "amount" be a percentage?
func (r Rect) Shrink(amount float32) Rect {
	return r.ShrinkXY(XY{amount, amount})
}

// TODO should "amount" be a percentage?
func (r Rect) ShrinkXY(amount XY) Rect {
	return Rect{
		A: XY{r.A.X + amount.X, r.A.Y + amount.Y},
		B: XY{r.B.X - amount.X, r.B.Y - amount.Y},
	}
}

func (r Rect) Grow(amount float32) Rect {
	return r.Shrink(-amount)
}

// TODO should "amount" be a percentage?
func (r Rect) GrowXY(amount XY) Rect {
	return r.ShrinkXY(XY{-amount.X, -amount.Y})
}

func (r Rect) Quad() Quad {
	return Quad{
		XY{r.A.X, r.A.Y},
		XY{r.A.X, r.B.Y},
		XY{r.B.X, r.B.Y},
		XY{r.B.X, r.A.Y},
	}
}

// "sub" is percentage based
// TODO make percentage vs offset consistent
func (r Rect) SubRect(sub Rect) Rect {
	return Rect{
		A: r.Interpolate(sub.A),
		B: r.Interpolate(sub.B),
	}
}

func (r Rect) Fill() Mesh {
	return Mesh{
		Verts: []XY{
			r.A, {r.A.X, r.B.Y},
			r.B, {r.B.X, r.A.Y},
		},
		Faces: []Face{
			{0, 2, 1},
			{0, 3, 2},
		},
		UV: []XY{
			{0, 0},
			{0, 1},
			{1, 1},
			{1, 0},
		},
	}
}

func (r Rect) Stroke(opt StrokeOpt) Mesh {
	return r.Quad().Stroke(opt)
}

func (r Rect) Path() Path {
	return r.Quad().Path()
}
