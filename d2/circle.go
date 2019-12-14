package d2

type Circle struct {
	XY
	Radius float32
}

func (c Circle) Bounds() Rect {
	return Rect{
		A: XY{
			c.X - c.Radius,
			c.Y - c.Radius,
		},
		B: XY{
			c.X + c.Radius,
			c.Y + c.Radius,
		},
	}
}

func (c Circle) IntersectsCircle(d Circle) bool {
	l := Line{c.XY, d.XY}
	// distance between circle centers
	// must be greater than the sum of the two radii.
	return l.Length() < c.Radius+d.Radius
}

func (c Circle) IntersectsLine(l Line) bool {
	// Given:
	//   C: center of the circle
	//   P: the point on the line closest to C
	//
	// If the distance from C to P is less than
	// the radius of the circle, then the line
	// intersects the circle.
	//
	// http://paulbourke.net/geometry/pointlineplane/

	x1 := l.A.X
	y1 := l.A.Y
	x2 := l.B.X
	y2 := l.B.Y
	x3 := c.X
	y3 := c.Y
	d := l.Length()
	u := ((x3-x1)*(x2-x1) + (y3-y1)*(y2-y1)) / (d * d)

	if u < 0 || u > 1 {
		// The closest point is one of the vertices of the line.
		l1 := Line{c.XY, l.A}
		l2 := Line{c.XY, l.B}
		return l1.Length() < c.Radius || l2.Length() < c.Radius
	}

	// Line from C to P.
	t := Line{
		A: c.XY,
		B: XY{
			X: x1 + u*(x2-x1),
			Y: y1 + u*(y2-y1),
		},
	}
	return t.Length() < c.Radius
}
