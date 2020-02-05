package main

import (
	"github.com/buchanae/ink/app"
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

const (
	lineWidth = 0.001
	gridTweak = 0.007
	// size is a little confusing. number of boxes is size - 1.
	size = 40
)

func Ink(doc *app.Doc) {
	rand.SeedNow()

	grid := Grid{Rows: size, Cols: size}
	boxes := MergedGridQuades(grid, 15, gridTweak)
	center := XY{0.5, 0.5}
	// TODO redo monochromatic color picker
	colors := rand.Palette()

	for _, b := range boxes {
		dist := b.Bounds().Center().Distance(center)

		//if rand.Bool(0.5) {
		//colors = color.Monochromatic(color.Blue, 8)
		//}

		// TODO scripting in Go: number type conversions are annoying
		i := int(dist*float32(len(colors))) + 1
		if i >= len(colors) {
			i = len(colors) - 1
		}
		c := colors[i]

		//if rand.Bool(0.7) {
		//c = color.Lighter(c, rand.Range(0.1, 1.3))
		//}

		m := gfx.NewShader(b.Fill())
		//m.Frag = "stained_glass.frag"
		m.Set("u_offset", rand.XYRange(0, 200))
		m.Set("a_color", c)
		m.Set("a_uv", []XY{
			{0, 0},
			{0, 1},
			{1, 1},

			{0, 0},
			{1, 1},
			{1, 0},
		})
		m.Draw(doc)
	}

	for _, b := range boxes {
		gfx.Stroke{
			Target: b,
			Width:  lineWidth,
			Color:  color.Black,
		}.Draw(doc)
	}
}

/*
TODO ideas:
- likeliness to merge based on some function, such as distance to point
  could create radial pattern with smaller pieces at edges and larger at center
*/
func MergedGridQuades(grid Grid, size int, tweakAmt float32) []Quad {

	type corner struct {
		row, col   int
		w, h       int
		a, b, c, d int
		offsets    struct {
			a, b, c, d XY
		}
	}

	type link struct {
		a, b, c, d *corner
	}

	var corners []*corner
	merged := make([]bool, len(grid.Cells))
	links := make([]link, len(grid.Cells))

	for i, cell := range grid.Cells {
		if merged[i] {
			continue
		}
		if cell.Row == grid.Rows-1 {
			continue
		}
		if cell.Col == grid.Cols-1 {
			continue
		}

		stopRow := cell.Row + rand.Intn(size)
		stopCol := cell.Col + rand.Intn(size)
		maxRow, maxCol := cell.Row, cell.Col

		for row := cell.Row; row < grid.Rows-1 && row < stopRow; row++ {
			for col := cell.Col; col < grid.Cols-1 && col < stopCol; col++ {

				i := grid.Index(row, col)
				if merged[i] {
					stopCol = col
					break
				}

				if row > maxRow {
					maxRow = row
				}
				if col > maxCol {
					maxCol = col
				}

				merged[i] = true
			}
		}

		c := &corner{
			row: cell.Row,
			col: cell.Col,
			w:   maxCol - cell.Col + 1,
			h:   maxRow - cell.Row + 1,
			a:   grid.Index(cell.Row, cell.Col),
			b:   grid.Index(cell.Row, maxCol+1),
			c:   grid.Index(maxRow+1, maxCol+1),
			d:   grid.Index(maxRow+1, cell.Col),
		}
		links[c.a].a = c
		links[c.b].b = c
		links[c.c].c = c
		links[c.d].d = c
		corners = append(corners, c)
	}

	for _, cor := range corners {
		// TODO
		continue
		if cor.w < 2 {
			continue
		}
		la := links[cor.a]
		if la.b == nil {
			continue
		}
		if la.c != nil || la.d != nil {
			continue
		}
		ld := links[cor.d]
		if ld.c == nil {
			continue
		}

		ta := rand.Range(-tweakAmt, tweakAmt)

		for row := 0; row < cor.h; row++ {
			i := grid.Index(cor.row+row, cor.col)
			l := links[i]
			p := 1 - (float32(row) / float32(cor.h))
			tx := ta * p

			if l.a != nil {
				l.a.offsets.a.X = tx
			}
			if l.b != nil {
				l.b.offsets.b.X = tx
			}
			if l.c != nil {
				l.c.offsets.c.X = tx
			}
			if l.d != nil {
				l.d.offsets.d.X = tx
			}
		}
	}

	for _, cor := range corners {
		ld := links[cor.d]
		if ld.a == nil {
			continue
		}
		if ld.b != nil || ld.c != nil {
			continue
		}
		if links[cor.c].b == nil {
			continue
		}

		ta := rand.Range(-tweakAmt, tweakAmt)
		w := cor.c - cor.d

		for col := 0; col < w; col++ {
			l := links[cor.d+col]
			p := 1 - (float32(col) / float32(w))
			ty := ta * p

			if l.a != nil {
				l.a.offsets.a.Y = ty
			}
			if l.b != nil {
				l.b.offsets.b.Y = ty
			}
			if l.c != nil {
				l.c.offsets.c.Y = ty
			}
			if l.d != nil {
				l.d.offsets.d.Y = ty
			}
		}
	}

	boxes := make([]Quad, 0, len(corners))
	for _, cor := range corners {
		boxes = append(boxes, Quad{
			A: grid.Cells[cor.a].XY.Add(cor.offsets.a),
			B: grid.Cells[cor.b].XY.Add(cor.offsets.b),
			C: grid.Cells[cor.c].XY.Add(cor.offsets.c),
			D: grid.Cells[cor.d].XY.Add(cor.offsets.d),
		})
	}
	return boxes
}
