package dd

type Face [3]int

type Mesh struct {
	Verts []XY
	Faces []Face
	UV    []XY
}

func NewMesh(tris []Triangle) Mesh {
	return Mesh{}.AddTriangles(tris)
}

func (m Mesh) Size() int {
	return len(m.Verts)
}

func (m Mesh) Mesh() Mesh {
	return m
}

func (m Mesh) Triangles() []Triangle {
	tris := make([]Triangle, len(m.Faces))
	for i, face := range m.Faces {
		tris[i] = Triangle{
			A: m.Verts[face[0]],
			B: m.Verts[face[1]],
			C: m.Verts[face[2]],
		}
	}
	return tris
}

func (m Mesh) Copy() Mesh {
	c := Mesh{}
	c.Verts = append(c.Verts, m.Verts...)
	c.Faces = append(c.Faces, m.Faces...)
	c.UV = append(c.UV, m.UV...)
	return c
}

func (m Mesh) AddTriangles(tris []Triangle) Mesh {
	for _, t := range tris {
		l := len(m.Verts)
		m.Verts = append(m.Verts, t.A, t.B, t.C)
		m.Faces = append(m.Faces, Face{l, l + 1, l + 2})
	}
	return m
}

func Merge(srcs ...Mesh) Mesh {
	dst := Mesh{}
	for _, src := range srcs {
		offset := len(dst.Verts)
		dst.Verts = append(dst.Verts, src.Verts...)
		for _, face := range src.Faces {
			dst.Faces = append(dst.Faces, Face{
				offset + face[0],
				offset + face[1],
				offset + face[2],
			})
		}
	}
	return dst
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
*/
