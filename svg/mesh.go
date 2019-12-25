package svg

import (
	"github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/tess"
)

// TODO get rid of width/height
func Mesh(raw string, width, height float32) (dd.Mesh, error) {

	path, err := Parse(raw, width, height)
	if err != nil {
		return dd.Mesh{}, err
	}
	_ = path

	var points []dd.XY

	for i, p := range points {
		points[i] = dd.XY{p.X / width, p.Y / height}
	}

	triangles := tess.Tesselate(points)
	mesh := dd.Triangles(triangles).Mesh()

	// TODO not sure why I can't flip verts before triangulation
	for i, v := range mesh.Verts {
		v.Y = 1 - v.Y
		mesh.Verts[i] = v
	}

	return mesh, nil
	//sort.Sort(sortPoints{points, center})
}

type sortPoints struct {
	points []dd.XY
	center dd.XY
}

func (s sortPoints) Len() int {
	return len(s.points)
}
func (s sortPoints) Swap(i, j int) {
	s.points[i], s.points[j] = s.points[j], s.points[i]
}
func (s sortPoints) Less(i, j int) bool {
	a := s.points[i]
	b := s.points[j]

	if a.X-s.center.X >= 0 && b.X-s.center.X < 0 {
		return true
	}
	if a.X-s.center.X < 0 && b.X-s.center.X >= 0 {
		return false
	}
	if a.X-s.center.X == 0 && b.X-s.center.X == 0 {
		if a.Y-s.center.Y >= 0 || b.Y-s.center.Y >= 0 {
			return a.Y > b.Y
		}
		return b.Y > a.Y
	}

	// compute the cross product of vectors (s.center -> a) x (s.center -> b)
	det := (a.X-s.center.X)*(b.Y-s.center.Y) - (b.X-s.center.X)*(a.Y-s.center.Y)
	if det < 0 {
		return true
	}
	if det > 0 {
		return false
	}

	// points a and b are on the same line from the s.center
	// check which point is closer to the s.center
	d1 := (a.X-s.center.X)*(a.X-s.center.X) + (a.Y-s.center.Y)*(a.Y-s.center.Y)
	d2 := (b.X-s.center.X)*(b.X-s.center.X) + (b.Y-s.center.Y)*(b.Y-s.center.Y)
	return d1 > d2
}
