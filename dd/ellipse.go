package dd

import (
	"github.com/buchanae/ink/math"
)

type Ellipse struct {
	XY
	Size     XY
	Segments int
}

func (e Ellipse) Interpolate(p float32) XY {
	t := math.Pi * 2 * p
	xy := XY{
		e.Size.X * math.Cos(t),
		e.Size.Y * math.Sin(t),
	}
	return xy.Add(e.XY)
}

func (e Ellipse) Stroke() Stroke {
	segments := e.Segments
	if segments <= 0 {
		segments = 50
	}

	// TODO reimplement with path.ArcTo
	//      and implement proper curve interpolation
	path := Path{}

	for i := 0; i < segments; i++ {
		p := float32(i) / float32(segments)
		xy := e.Interpolate(p)
		if i == 0 {
			path.MoveTo(xy)
		} else {
			path.LineTo(xy)
		}
	}

	s := path.Stroke()
	s.Closed = true
	return s
}

func (e Ellipse) Mesh() Mesh {
	segments := e.Segments
	if segments <= 0 {
		segments = 50
	}

	faces := make([]Face, 0, segments)
	verts := make([]XY, 0, segments+1)
	verts = append(verts, e.XY)

	for i := 0; i < segments; i++ {
		// the 0 index is the center vertex.
		// perimeter vertices start at index 1.
		current := i + 1
		previous := current - 1
		if previous == 0 {
			previous = segments
		}

		p := float32(i) / float32(segments)
		xy := e.Interpolate(p)
		verts = append(verts, xy)
		faces = append(faces, Face{
			0, current, previous,
		})
	}

	return Mesh{Verts: verts, Faces: faces}
}
