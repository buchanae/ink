package dd

// TODO unfinished. need to stroke half width in both directions.
//      currently stroking full width in one direciton.
func Stroke(lines []Line, width float32, closed bool) []Triangle {
	if len(lines) < 2 {
		closed = false
	}

	// TODO a single line is a special case
	if len(lines) == 1 {
		l := lines[0]
		n := l.Normal().SetLength(width)
		return []Triangle{
			{l.A, l.B, l.A.Add(n)},
			{l.A.Add(n), l.B, l.B.Add(n)},
		}
	}

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

	return tris
}

func miterPoint(a, b Line, width float32) XY {
	n := a.Normal()
	miter := n.Add(b.Normal()).Normalize()
	miterWidth := width / miter.Dot(n)
	return a.B.Add(miter.SetLength(miterWidth))
}
