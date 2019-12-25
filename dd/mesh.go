package dd

type Face [3]int

type Mesh struct {
	Verts []XY
	Faces []Face
}

func (m Mesh) Size() int {
	return len(m.Verts)
}

func (m Mesh) Mesh() Mesh {
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
