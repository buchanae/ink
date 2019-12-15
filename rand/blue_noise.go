package rand

import (
	"math/rand"

	"github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/math"
)

func (r *Rand) BlueNoise(n int, w, h, d float32) []dd.XY {
	points := make([]dd.XY, 0, n)
	cellSize := d / math.Sqrt(2)
	cellsW := int(w/cellSize) + 1
	cellsH := int(h/cellSize) + 1
	cells := make([]*dd.XY, cellsW*cellsH)
	var active []*dd.XY

	cellIndex := func(pt dd.XY) int {
		return int(pt.Y/cellSize)*cellsW + int(pt.X/cellSize)
	}
	cellCoord := func(x, y float32) (int, int) {
		return int(x / cellSize), int(y / cellSize)
	}

	pt := dd.XY{w / 2, h / 2}
	ci := cellIndex(pt)
	cells[ci] = &pt
	active = append(active, &pt)
	points = append(points, pt)

	for {
		if len(active) == 0 || len(points) == n {
			break
		}
		ai := rand.Intn(len(active))
		pt := active[ai]

		var accept *dd.XY
		randRing := r.makeRandRing(30, d, 2*d)

		for _, rp := range randRing {
			// generate random point in annulus (ring) from [d, 2d]
			rp.X += pt.X
			rp.Y += pt.Y
			rx := rp.X
			ry := rp.Y

			cx, cy := cellCoord(rx, ry)
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

			if cx < 0 || cx >= cellsW || cy < 0 || cy >= cellsH {
				continue
			}

			ok := true
			for _, n := range neighbors {
				nx, ny := n[0], n[1]

				// list bounds check
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
				dist := pointDist(*nc, rp)
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

func (r *Rand) makeRandRing(n int, minR, maxR float32) []dd.XY {
	A := 2 / (maxR*maxR - minR*minR)

	points := make([]dd.XY, 0, n)
	for i := 0; i < n; i++ {
		R := math.Sqrt(2*r.Float()/A + minR*minR)
		theta := r.Float() * math.Pi * 2
		points = append(points, dd.XY{
			X: R * math.Cos(theta),
			Y: R * math.Sin(theta),
		})
	}
	return points
}

func pointDist(a, b dd.XY) float32 {
	x := b.X - a.X
	y := b.Y - a.Y
	return math.Sqrt(x*x + y*y)
}
