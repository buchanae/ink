package dd

// TODO https://developer.mozilla.org/en-US/docs/Web/SVG/Attribute/stroke-miterlimit
// https://www.w3.org/TR/SVG11/painting.html#StrokeLinecapProperty
type StrokeCap int

const (
	ButtCap = iota
	RoundCap
	SquareCap
)

// TODO https://developer.mozilla.org/en-US/docs/Web/SVG/Attribute/stroke-linejoin
// https://www.w3.org/TR/SVG11/painting.html#StrokeLinejoinProperty
type StrokeJoin int

const (
	MiterJoin = iota
	RoundJoin
	BevelJoin
)

// TODO stroke options: miter, cap, etc.
// https://www.w3.org/TR/SVG11/painting.html#StrokeMiterlimitProperty
// https://www.w3.org/TR/SVG11/painting.html#StrokeLinejoinProperty
// https://www.w3.org/TR/SVG11/painting.html#StrokeDasharrayProperty
// https://www.w3.org/TR/SVG11/painting.html#FillRuleProperty
// TODO a zero length curve with a square/round cap should be drawn
type StrokeOpt struct {
	Width float32
	Cap   StrokeCap
	Join  StrokeJoin
}

// TODO unfinished. need to stroke half width in both directions.
//      currently stroking full width in one direciton.
// http://tavmjong.free.fr/SVG/LINEJOIN_STUDY/
// http://tavmjong.free.fr/SVG/LINEJOIN/index.html

func Stroke(path Path, opt StrokeOpt) Mesh {
	mesh := Mesh{}
	return StrokeTo(path, mesh, opt)
}

func StrokeTo(path Path, mesh Mesh, opt StrokeOpt) Mesh {
	lines := pathToLines(path)
	if len(lines) == 0 {
		return mesh
	}

	width := opt.Width
	if width == 0 {
		width = 0.001
	}

	if len(lines) == 1 {
		l := lines[0]
		n := l.Normal().SetLength(width)
		return mesh.AddTriangles([]Triangle{
			{l.A, l.B, l.A.Add(n)},
			{l.A.Add(n), l.B, l.B.Add(n)},
		})
	}

	closed := path.Start() == path.End()
	var tris []Triangle
	var prev [2]XY

	for i, line := range lines {
		n := line.Normal()
		var next Line

		// if on the last line, cap the end points
		if i == len(lines)-1 {
			if closed {
				next = lines[0]
			} else {
				tris = append(tris,
					Triangle{prev[0], prev[1], line.B},
					Triangle{prev[1], line.B.Add(n.SetLength(width)), line.B},
				)
				break
			}
		} else {
			next = lines[i+1]
		}

		// if on the first line, initialize the start points
		if i == 0 {
			if closed {
				prev[0] = lines[len(lines)-1].B
				prev[1] = miterPoint(lines[len(lines)-1], line, width)
			} else {
				prev[0] = line.A
				prev[1] = line.A.Add(n.SetLength(width))
			}
		}

		mp := miterPoint(line, next, width)
		tris = append(tris,
			Triangle{prev[0], line.B, prev[1]},
			Triangle{prev[1], line.B, mp},
		)

		prev[0] = line.B
		prev[1] = mp
	}

	return mesh.AddTriangles(tris)
}

func pathToLines(path Path) []Line {

	lines := make([]Line, 0, len(path))
	for _, curve := range path {
		switch z := curve.(type) {
		case Line:
			lines = append(lines, z)
		default:
			// TODO dynamic segment count using error margin
			segments := 50
			xys := Subdivide(curve, segments)
			lines = append(lines, XYsToLines(xys...)...)
		}
	}
	return lines
}

func miterPoint(a, b Line, width float32) XY {
	n := a.Normal()
	miter := n.Add(b.Normal()).Normalize()
	miterWidth := width / miter.Dot(n)
	return a.B.Add(miter.SetLength(miterWidth))
}
