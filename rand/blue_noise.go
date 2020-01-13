package rand

import (
	"math"

	"github.com/buchanae/ink/dd"
)

type BlueNoise struct {
	Rect    dd.Rect
	Spacing float32
	Limit   int
	Initial []dd.XY
}

func (bn BlueNoise) Generate() []dd.XY {
	return bn.GenerateWith(Default)
}

func (bn BlueNoise) GenerateWith(r *Rand) []dd.XY {
	n := bn.Limit
	if n == 0 {
		n = 100000
	}

	d := bn.Spacing
	if d == 0 {
		d = 0.01
	}

	bounds := bn.Rect
	if bounds.IsZero() {
		bounds = dd.RectWH(1, 1)
	}

	initial := bn.Initial
	if initial == nil {
		initial = []dd.XY{bounds.Center()}
	}

	// final points that will be returned
	output := make([]dd.XY, 0, n)

	size := bounds.Size()
	w := size.X
	h := size.Y
	shiftedBounds := dd.Rect{
		B: bounds.B.Sub(bounds.A),
	}

	// grid of cells, to track neighbors
	cellSize := d / sqrt(2)
	cellsW := int(w/cellSize) + 1
	cellsH := int(h/cellSize) + 1
	cells := make([]*dd.XY, cellsW*cellsH)
	cellIndex := func(pt dd.XY) int {
		return int(pt.Y/cellSize)*cellsW + int(pt.X/cellSize)
	}
	cellCoord := func(x, y float32) (int, int) {
		return int(x / cellSize), int(y / cellSize)
	}

	// active points are those that may generate
	// new neighboring points. when an active point
	// fails to generate a valid neighbor, it is removed
	// from the active list.
	var active []*dd.XY

	// add initial points to the output and active lists
	for _, pt := range initial {
		ptv := pt.Sub(bounds.A)
		ci := cellIndex(ptv)
		cells[ci] = &ptv
		active = append(active, &ptv)
		output = append(output, ptv)
	}

	iter := 0

	for {
		if len(active) == 0 || len(output) >= n+len(initial) {
			break
		}
		ai := Intn(len(active))
		pt := active[ai]

		var accept *dd.XY

		// generate 30 points in a ring around "pt"
		// check each point and see if it's valid
		// to add to the output list. if so, break.
		for ri := 0; ri < 30; ri++ {
			iter++
			rp := r.makeRandRing(*pt, d, 2*d)

			cx, cy := cellCoord(rp.X, rp.Y)

			// skip any points that fall out of bounds
			if !shiftedBounds.Contains(rp) {
				continue
			}
			if cx < 0 || cx >= cellsW || cy < 0 || cy >= cellsH {
				continue
			}

			neighbors := [9][2]int{
				{cx - 1, cy - 1},
				{cx, cy - 1},
				{cx + 1, cy - 1},
				{cx - 1, cy},
				{cx, cy},
				{cx + 1, cy},
				{cx - 1, cy + 1},
				{cx, cy + 1},
				{cx + 1, cy + 1},
			}

			ok := true
			for _, n := range neighbors {
				nx, ny := n[0], n[1]

				// skip any neighbors that fall out of bounds
				if nx < 0 || nx >= cellsW || ny < 0 || ny >= cellsH {
					continue
				}

				// does a neighboring cell have a point in it?
				ni := ny*cellsW + nx
				nc := cells[ni]
				if nc == nil {
					continue
				}

				// is a neighboring point within the minimum distance?
				dist := nc.Distance(rp)
				if dist <= d {
					ok = false
					break
				}
			}

			if ok {
				accept = &rp
				break
			}
		}

		if accept != nil {
			out := *accept
			out = out.Add(bounds.A)
			output = append(output, out)

			active = append(active, accept)
			ci := cellIndex(*accept)
			cells[ci] = accept
		} else {
			// TODO better way to remove an active index?
			//      copy seems inefficient
			active = append(active[:ai], active[ai+1:]...)
		}
	}

	//log.Printf("blue: iter %d", iter)
	return output
}

// makeRandRing returns a list of random points in the ring
// with "minR" inner radius and "maxR" outer radius.
func (r *Rand) makeRandRing(src dd.XY, minR, maxR float32) dd.XY {
	A := 2 / (maxR*maxR - minR*minR)
	R := sqrt(2*r.Float()/A + minR*minR)
	theta := r.Float() * math.Pi * 2
	return dd.XY{
		X: R * cos(theta),
		Y: R * sin(theta),
	}.Add(src)
}
