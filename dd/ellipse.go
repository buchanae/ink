package dd

import (
	"github.com/buchanae/ink/math"
)

type Ellipse struct {
	XY
	Size     XY
	Segments int
}

// TODO interpolate should probably be a percentage?
func (e Ellipse) Interpolate(t float32) XY {
	xy := XY{
		e.Size.X * math.Cos(t),
		e.Size.Y * math.Sin(t),
	}
	return xy.Add(e.XY)
}

func (e Ellipse) Mesh() Mesh {
	segments := e.Segments
	if segments <= 0 {
		segments = 40
	}

	faces := make([]Face, 0, segments)
	verts := make([]XY, 0, segments+1)
	verts = append(verts, e.XY)
	inc := (math.Pi * 2) / float32(segments)
	t := float32(0)

	for i := 0; i < segments; i++ {
		// the 0 index is the center vertex.
		// perimeter vertices start at index 1.
		current := i + 1
		previous := current - 1
		if previous == 0 {
			previous = segments
		}

		verts = append(verts, XY{
			X: e.Size.X*math.Cos(t) + e.X,
			Y: e.Size.Y*math.Sin(t) + e.Y,
		})
		faces = append(faces, Face{
			0,
			current,
			previous,
		})

		t += inc
	}

	return Mesh{Verts: verts, Faces: faces}
}
