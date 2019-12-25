package voronoi

import (
	"github.com/buchanae/ink/dd"
	"github.com/pzsz/voronoi"
	"github.com/pzsz/voronoi/utils"
)

type Cell struct {
	XY    dd.XY
	Edges []dd.Line
	Tris  []dd.Triangle
}

type Voronoi struct {
	diagram *voronoi.Diagram
	bbox    voronoi.BBox
}

func (v Voronoi) Edges() []dd.Line {
	var edges []dd.Line

	for _, edge := range v.diagram.Edges {
		edges = append(edges, dd.Line{
			A: dd.XY{
				X: float32(edge.Va.X),
				Y: float32(edge.Va.Y),
			},
			B: dd.XY{
				X: float32(edge.Vb.X),
				Y: float32(edge.Vb.Y),
			},
		})
	}
	return edges
}

func (v Voronoi) Triangulate() []dd.Triangle {
	var out []dd.Triangle

	for _, cell := range v.diagram.Cells {
		center := dd.XY{float32(cell.Site.X), float32(cell.Site.Y)}

		var prev dd.XY
		for i, he := range cell.Halfedges {
			other := he.Edge.GetOtherCell(cell)
			if other == nil {
				continue
			}
			site := other.Site
			p := dd.XY{float32(site.X), float32(site.Y)}
			if i > 0 {
				tri := dd.Triangle{center, prev, p}
				// TODO dedupe
				out = append(out, tri)
			}
			prev = p
		}
	}

	return out
}

func (v Voronoi) Cells() []Cell {
	cells := make([]Cell, 0, len(v.diagram.Cells))
	for _, c := range v.diagram.Cells {

		center := dd.XY{float32(c.Site.X), float32(c.Site.Y)}
		edges := make([]dd.Line, 0, len(c.Halfedges))
		tris := make([]dd.Triangle, 0, len(c.Halfedges))

		for _, he := range c.Halfedges {
			a := dd.XY{
				X: float32(he.Edge.Va.X),
				Y: float32(he.Edge.Va.Y),
			}
			b := dd.XY{
				X: float32(he.Edge.Vb.X),
				Y: float32(he.Edge.Vb.Y),
			}
			edges = append(edges, dd.Line{a, b})
			tris = append(tris, dd.Triangle{center, a, b})
		}

		cells = append(cells, Cell{center, edges, tris})
	}
	return cells
}

func (v Voronoi) Relax() Voronoi {
	verts := utils.LloydRelaxation(v.diagram.Cells)
	diagram := voronoi.ComputeDiagram(verts, v.bbox, true)
	return Voronoi{diagram, v.bbox}
}

func New(xys []dd.XY, rect dd.Rect) Voronoi {
	bbox := voronoi.NewBBox(
		float64(rect.A.X), float64(rect.B.X),
		float64(rect.A.Y), float64(rect.B.Y),
	)

	verts := make([]voronoi.Vertex, 0, len(xys))
	for _, xy := range xys {
		verts = append(verts, voronoi.Vertex{float64(xy.X), float64(xy.Y)})
	}

	diagram := voronoi.ComputeDiagram(verts, bbox, true)
	v := Voronoi{
		diagram: diagram,
		bbox:    bbox,
	}

	return v
}
