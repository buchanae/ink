package dd

// TODO stroke options: miter, cap, etc.
type Stroke struct {
	Lines  []Line
	Width  float32
	Closed bool
}

func (s Stroke) Mesh() Mesh {
	return NewMesh(s.Triangles())
}

// TODO unfinished. need to stroke half width in both directions.
//      currently stroking full width in one direciton.
func (s Stroke) Triangles() []Triangle {
	if len(s.Lines) == 0 {
		return nil
	}

	width := s.Width
	if width == 0 {
		width = 0.001
	}

	if len(s.Lines) == 1 {
		l := s.Lines[0]
		n := l.Normal().SetLength(width)
		return []Triangle{
			{l.A, l.B, l.A.Add(n)},
			{l.A.Add(n), l.B, l.B.Add(n)},
		}
	}

	var tris []Triangle
	var prev [2]XY

	for i, line := range s.Lines {
		n := line.Normal()
		var next Line

		// if on the last line, cap the end points
		if i == len(s.Lines)-1 {
			if s.Closed {
				next = s.Lines[0]
			} else {
				tris = append(tris,
					Triangle{prev[0], prev[1], line.B},
					Triangle{prev[1], line.B.Add(n.SetLength(width)), line.B},
				)
				break
			}
		} else {
			next = s.Lines[i+1]
		}

		// if on the first line, initialize the start points
		if i == 0 {
			if s.Closed {
				prev[0] = s.Lines[len(s.Lines)-1].B
				prev[1] = miterPoint(s.Lines[len(s.Lines)-1], line, width)
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
