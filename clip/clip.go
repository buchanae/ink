package clip

import (
	. "github.com/buchanae/ink/dd"
)

func Difference(subject, clip Path) []XY {
	return clipPaths(subject, clip, true, true)
}

func Union(subject, clip Path) []XY {
	return clipPaths(subject, clip, true, false)
}

func Intersection(subject, clip Path) []XY {
	return clipPaths(subject, clip, false, false)
}

func clipPaths(subject, clip Path, sFlip, cFlip bool) []XY {
	sNodes := pathToNodes(subject)
	cNodes := pathToNodes(clip)

	sIntersects, cIntersects := findIntersects(sNodes, cNodes)
	// TODO what to do in this case?
	if len(sIntersects) == 0 || len(cIntersects) == 0 {
		return nil
	}

	markEntry(sNodes, sFlip)
	markEntry(cNodes, cFlip)

	points := buildPoly(sIntersects[0])
	return points
}

type node struct {
	next   *node
	prev   *node
	link   *node
	xy     XY
	alpha  float32
	entry  bool
	inside bool
}

func pathToNodes(path Path) []*node {
	size := len(path)
	nodes := make([]*node, size)
	last := size - 1

	for i, _ := range nodes {
		nodes[i] = new(node)
	}

	// TODO should be Edges()?
	for i, curve := range path {
		line := curve.(Line)

		if i == 0 {
			nodes[i].prev = nodes[last]
		} else {
			nodes[i].prev = nodes[i-1]
		}

		if i == last {
			nodes[i].next = nodes[0]
		} else {
			nodes[i].next = nodes[i+1]
		}
		nodes[i].xy = line.A
	}
	return nodes
}

func insertIntersect(head, intersect *node) {
	for {
		// Stop if next node is not an intersection.
		if head.next.link == nil {
			break
		}
		if head.next.alpha >= intersect.alpha {
			break
		}
		head = head.next
	}
	intersect.prev = head
	intersect.next = head.next
	intersect.next.prev = intersect
	intersect.prev.next = intersect
}

func findIntersects(sNodes, cNodes []*node) ([]*node, []*node) {
	sIntersects := make([]*node, 0, 50)
	cIntersects := make([]*node, 0, 50)

	sLines := make([]Line, len(sNodes))
	cLines := make([]Line, len(cNodes))
	sLengths := make([]float32, len(sNodes))
	cLengths := make([]float32, len(cNodes))

	for i, n := range sNodes {
		sLines[i] = Line{n.xy, n.next.xy}
		sLengths[i] = sLines[i].Length()
	}
	for i, n := range cNodes {
		cLines[i] = Line{n.xy, n.next.xy}
		cLengths[i] = cLines[i].Length()
	}

	for si, s := range sNodes {
		for ci, c := range cNodes {

			if intersectHorizontalRay(s.xy, cLines[ci]) {
				s.inside = !s.inside
			}
			if intersectHorizontalRay(c.xy, sLines[si]) {
				c.inside = !c.inside
			}

			pt, ok := intersectLines(sLines[si], cLines[ci])
			if !ok {
				continue
			}

			sIntersect := &node{
				xy:    pt,
				alpha: pt.Sub(s.xy).Length() / sLengths[si],
			}
			cIntersect := &node{
				link:  sIntersect,
				xy:    pt,
				alpha: pt.Sub(c.xy).Length() / cLengths[ci],
			}
			sIntersect.link = cIntersect

			insertIntersect(s, sIntersect)
			insertIntersect(c, cIntersect)

			sIntersects = append(sIntersects, sIntersect)
			cIntersects = append(cIntersects, cIntersect)
		}
	}

	return sIntersects, cIntersects
}

func markEntry(nodes []*node, flip bool) {
	start := nodes[0]
	entry := !start.inside

	if flip {
		entry = !entry
	}

	it := start.next

	for it != start {
		if it.link != nil {
			it.entry = entry
			entry = !entry
		}
		it = it.next
	}
}

func buildPoly(start *node) []XY {

	points := make([]XY, 0, 50)
	points = append(points, start.xy)
	forward := start.entry

	var it *node
	if forward {
		it = start.next
	} else {
		it = start.prev
	}

	iters := 0
	for {
		iters++
		if iters > 1000 {
			println("broken")
			break
		}

		// TODO maybe a better check that doesn't involve xy
		if it.xy == start.xy {
			break
		}

		points = append(points, it.xy)

		if it.link != nil {
			it = it.link
			forward = it.entry
		}

		if forward {
			it = it.next
		} else {
			it = it.prev
		}
	}
	return points
}

// TODO replace with some real
func intersectHorizontalRay(xy XY, l Line) bool {
	if l.A.X < xy.X && l.B.X < xy.X {
		return false
	}
	if l.A.Y < xy.Y && l.B.Y < xy.Y {
		return false
	}
	x := l.A.X
	if l.B.X > x {
		x = l.B.X
	}
	x += 100000
	xyl := Line{xy, XY{x, xy.Y}}
	_, ok := intersectLines(xyl, l)
	return ok
}

// TODO move this to a general util
func intersectLines(a, b Line) (XY, bool) {

	d1 := a.B.Sub(a.A)
	d2 := b.B.Sub(b.A)
	d3 := b.A.Sub(a.A)

	numer := (d1.X * d3.Y) - (d1.Y * d3.X)
	denom := (d1.Y * d2.X) - (d1.X * d2.Y)
	s := numer / denom

	aIsVertical := d1.X == 0

	var t float32
	if aIsVertical {
		t = ((b.A.Y + d2.Y*s) - a.A.Y) / d1.Y
	} else {
		t = ((b.A.X + d2.X*s) - a.A.X) / d1.X
	}

	if s >= 0 && s <= 1 && t >= 0 && t <= 1 {
		return b.Interpolate(s), true
	}

	return XY{}, false
}
