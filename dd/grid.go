package dd

type Cell struct {
	Rect
	Row, Col int
}

type Grid struct {
	Rect
	Rows, Cols int
}

func (g Grid) Cells() []Cell {
	rect := g.Rect
	if rect.IsZero() {
		rect = RectWH(1, 1)
	}
	cells := make([]Cell, 0, g.Rows*g.Cols)

	rows := float32(g.Rows)
	cols := float32(g.Cols)

	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols; c++ {

			cells = append(cells, Cell{
				Row: r,
				Col: c,
				Rect: Rect{
					A: rect.Interpolate(XY{
						X: float32(c) / cols,
						Y: float32(r) / rows,
					}),
					B: rect.Interpolate(XY{
						X: float32(c+1) / cols,
						Y: float32(r+1) / rows,
					}),
				},
			})
		}
	}
	return cells
}

func (g Grid) Cell(row, col int) Cell {
	rect := g.Rect
	if rect.IsZero() {
		rect = RectWH(1, 1)
	}

	rows := float32(g.Rows)
	cols := float32(g.Cols)

	return Cell{
		Row: row,
		Col: col,
		Rect: Rect{
			A: rect.Interpolate(XY{
				X: float32(col) / cols,
				Y: float32(row) / rows,
			}),
			B: rect.Interpolate(XY{
				X: float32(col+1) / cols,
				Y: float32(row+1) / rows,
			}),
		},
	}
}
