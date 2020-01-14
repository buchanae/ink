package dd

import (
	"math"
)

type Circle struct {
	XY
	Radius   float32
	Segments int
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

// TODO Contains or Intersects or both?
func (c Circle) ContainsXY(xy XY) bool {
	l := Line{c.XY, xy}
	return l.Length() <= c.Radius
}

func (c Circle) DistanceToXY(xy XY) float32 {
	l := Line{c.XY, xy}
	return l.Length() - c.Radius
}

func (c Circle) DistanceToCircle(d Circle) float32 {
	l := Line{c.XY, d.XY}
	return l.Length() - c.Radius - d.Radius
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

func (c Circle) Interpolate(p float32) XY {
	ang := math.Pi * 2 * p
	return XY{
		X: cos(ang)*c.Radius + c.X,
		Y: sin(ang)*c.Radius + c.Y,
	}
}

func (c Circle) Ellipse() Ellipse {
	return Ellipse{
		c.XY,
		XY{c.Radius, c.Radius},
		c.Segments,
	}
}

func (c Circle) Mesh() Mesh {
	return c.Ellipse().Mesh()
}

func (c Circle) Normals() []XY {
	// TODO could every mesh have per-face normals?
	mesh := c.Mesh()
	verts := mesh.Verts
	normals := make([]XY, mesh.Size())

	// the 0 index is the center vertex.
	// perimeter vertices start at index 1.
	for i := 1; i < len(verts); i++ {
		prev := i - 1
		if prev == 0 {
			prev = len(verts) - 1
		}
		next := i + 1
		if next == len(verts) {
			next = 1
		}

		a := Line{verts[prev], verts[i]}
		b := Line{verts[i], verts[next]}
		n := a.Normal().Add(b.Normal()).Normalize()
		normals[i] = n
	}
	return normals
}

func (c Circle) Stroke(opt StrokeOpt) Mesh {
	return c.Ellipse().Stroke(opt)
}
