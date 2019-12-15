package dd

type Face [3]int

type Mesh struct {
	Verts []XY
	Faces []Face
}

func (m Mesh) Size() int {
	return len(m.Verts)
}

/*
type Edge [2]uint32

// TODO
func (m *Mesh) Triangles(faces ...Face) []Triangle {
	return nil
}

// TODO
func (m *Mesh) Lines(edges ...Edge) []Line {
	return nil
}

// TODO
func (m *Mesh) Wireframe() []Edge {
	return nil
}

// TODO
func (m *Mesh) Hull() []Edge {
	return nil
}

func (m *Mesh) HullLines() []Line {
	p := Path{}
	for i, edge := range m.Hull {
		a := m.Verts[edge[0]]
		b := m.Verts[edge[1]]
		if i == 0 {
			p.MoveTo(a)
			p.LineTo(b)
			continue
		}
		p.LineTo(b)
	}
	return p.Lines()
}

func Merge(dst *Mesh, srcs ...*Mesh) {
	offset := len(dst.Verts)
	for _, src := range srcs {
		dst.Verts = append(dst.Verts, src.Verts...)
		for _, face := range src.Faces {
			dst.Faces = append(dst.Faces, Face{
				offset + face[0],
				offset + face[1],
				offset + face[2],
			})
		}
	}
}
*/

// TODO need a better naming scheme for things
//      the generate meshes?
// TODO have the triangulation return a mesh
//      with edges, hull, etc already filled in.
func Triangles(tris []Triangle) Mesh {
	verts := make([]XY, 0, len(tris)*3)
	faces := make([]Face, 0, len(tris))

	for _, t := range tris {
		// TODO optimize vertex reuse
		l := len(verts)
		verts = append(verts, t.A, t.B, t.C)
		faces = append(faces, Face{l, l + 1, l + 2})
	}
	return Mesh{verts, faces}
}

func StrokeTriangles(tris []Triangle, width float32) Mesh {

	p := &Path{}
	for _, t := range tris {
		p.MoveTo(t.A)
		p.LineTo(t.B)
		p.LineTo(t.C)
		p.LineTo(t.A)
	}

	return p.Stroke(width)
}
