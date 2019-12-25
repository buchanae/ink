package rand

import (
	"math"

	"github.com/buchanae/ink/dd"
)

func (r *Rand) BlueNoise(n int, w, h, d float32) []dd.XY {
	initial := []dd.XY{{w / 2, h / 2}}
	return r.BlueNoiseInitial(n, w, h, d, initial)
}

func (r *Rand) BlueNoiseInitial(n int, w, h, d float32, initial []dd.XY) []dd.XY {
	// final point that will be returned
	points := make([]dd.XY, 0, n)

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

	var active []*dd.XY

	for _, pt := range initial {
		ptv := pt
		ci := cellIndex(ptv)
		cells[ci] = &ptv
		active = append(active, &ptv)
		points = append(points, ptv)
	}

	for {
		if len(active) == 0 || len(points) >= n+len(initial) {
			break
		}
		ai := Intn(len(active))
		pt := active[ai]

		var accept *dd.XY
		randRing := r.makeRandRing(30, d, 2*d)

		for _, rp := range randRing {
			rp.X += pt.X
			rp.Y += pt.Y

			cx, cy := cellCoord(rp.X, rp.Y)

			// skip any points that fall out of bounds
			if cx < 0 || cx >= cellsW || cy < 0 || cy >= cellsH {
				continue
			}

			neighbors := [][2]int{
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
			points = append(points, *accept)
			active = append(active, accept)
			ci := cellIndex(*accept)
			cells[ci] = accept
		} else {
			active = append(active[:ai], active[ai+1:]...)
		}
	}
	return points
}

// makeRandRing returns a list of random points in the ring
// with "minR" inner radius and "maxR" outer radius.
func (r *Rand) makeRandRing(n int, minR, maxR float32) []dd.XY {
	A := 2 / (maxR*maxR - minR*minR)

	points := make([]dd.XY, 0, n)
	for i := 0; i < n; i++ {
		R := sqrt(2*r.Float()/A + minR*minR)
		theta := r.Float() * math.Pi * 2
		points = append(points, dd.XY{
			X: R * cos(theta),
			Y: R * sin(theta),
		})
	}
	return points
}
