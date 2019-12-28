package main

func main() {

	// point
	a := XY{0.4, 0.3}
	doc.Point(a)

	b := XY{0.1, 0.2}
	c := a.Add(b)
	doc.Point(c)

	d := c.Sub(XY{0.3, 0.3}).Multiply(XY{2, 2})
	doc.Point(d)

	// Add, Sub, Multiply,
	// Distance
	// Dot

	// line
	// Interpolate

	// intersect lines
}
