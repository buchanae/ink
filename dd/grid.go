package dd

type Cell struct {
	XY
	Row, Col int
}

type Grid struct {
	Rows, Cols int
	Cells      []Cell
}

// TODO grid size is confusing
//      because the grid is actually a point
//      grid, not a rectangle grid. Maybe
//      don't create Cells automatically,
//      and add Points() function instead.
//
// ALso row/cols is y/x, while everything else is x/y (width/height)
//
// TODO want other ways to create grids, especially non-square grids
//      with cells with similar width/height
func NewGrid(rows, cols int) Grid {
	grid := Grid{
		Rows:  rows,
		Cols:  cols,
		Cells: make([]Cell, 0, rows*cols),
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			grid.Cells = append(grid.Cells, Cell{
				Col: c,
				Row: r,
				XY: XY{
					X: float32(c) / float32(cols-1),
					Y: float32(r) / float32(rows-1),
				},
			})
		}
	}

	return grid
}

func (g *Grid) IsEdge(row, col int) bool {
	return row == 0 || col == 0 || row == g.Rows-1 || col == g.Cols-1
}

func (g *Grid) Rects() []Rect {
	// TODO precise capacity
	rects := make([]Rect, 0, len(g.Cells))

	for row := 0; row < g.Rows-1; row++ {
		for col := 0; col < g.Cols-1; col++ {
			a := g.Cell(row, col)
			b := g.Cell(row+1, col+1)
			rects = append(rects, Rect{a.XY, b.XY})
		}
	}

	return rects
}

func (g *Grid) Quads() []Quad {
	// TODO precise capacity
	quads := make([]Quad, 0, len(g.Cells))

	for row := 0; row < g.Rows-1; row++ {
		for col := 0; col < g.Cols-1; col++ {
			a := g.Cell(row, col)
			b := g.Cell(row, col+1)
			c := g.Cell(row+1, col+1)
			d := g.Cell(row+1, col)
			q := Quad{a.XY, b.XY, c.XY, d.XY}
			quads = append(quads, q)
		}
	}
	return quads
}

func (g *Grid) Index(row, col int) int {
	return (row * g.Cols) + col
}

func (g *Grid) Cell(row, col int) Cell {
	i := g.Index(row, col)
	return g.Cells[i]
}
