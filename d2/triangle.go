package d2

func NewTriangle(a, b, c XY) Triangle {
	return Triangle{a, b, c}
}

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
