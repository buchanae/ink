package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/tess"
)

type Node struct {
	next  *Node
	prev  *Node
	link  *Node
	xy    XY
	alpha float32
	entry bool
}

func pathToNodes(path Path) []*Node {
	size := len(path)
	nodes := make([]*Node, size)
	last := size - 1

	for i, _ := range nodes {
		nodes[i] = new(Node)
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

func findIntersects(sNodes, cNodes []*Node) ([]*Node, []*Node) {
	sIntersects := make([]*Node, 0, 50)
	cIntersects := make([]*Node, 0, 50)

	for _, s := range sNodes {
		for _, c := range cNodes {
			sLine := Line{s.xy, s.next.xy}
			cLine := Line{c.xy, c.next.xy}

			pt, ok := IntersectLines(sLine, cLine)
			if !ok {
				continue
			}

			// TODO sorted insert when one line intersects multiple times
			sIntersect := &Node{
				prev:  s,
				next:  s.next,
				xy:    pt,
				alpha: pt.Sub(s.xy).Length() / sLine.Length(),
			}
			cIntersect := &Node{
				prev:  c,
				next:  c.next,
				link:  sIntersect,
				xy:    pt,
				alpha: pt.Sub(c.xy).Length() / cLine.Length(),
			}

			sIntersect.link = cIntersect
			s.next.prev = sIntersect
			s.next = sIntersect

			c.next.prev = cIntersect
			c.next = cIntersect

			sIntersects = append(sIntersects, sIntersect)
			cIntersects = append(cIntersects, cIntersect)
		}
	}

	return sIntersects, cIntersects
}

func markEntry(nodes []*Node, entry bool) {
	for _, n := range nodes {
		n.entry = entry
		entry = !entry
	}
}

func buildPoly(start *Node) []XY {

	points := make([]XY, 0, 50)
	points = append(points, start.xy)
	forward := start.entry

	var it *Node
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

func clipPaths(subject, clip Path, sEntry, cEntry bool) []XY {
	sNodes := pathToNodes(subject)
	cNodes := pathToNodes(clip)

	sIntersects, cIntersects := findIntersects(sNodes, cNodes)
	markEntry(sIntersects, sEntry)
	markEntry(cIntersects, cEntry)
	points := buildPoly(sIntersects[0])
	return points
}

func Difference(subject, clip Path) []XY {
	return clipPaths(subject, clip, true, true)
}

func Union(subject, clip Path) []XY {
	return clipPaths(subject, clip, true, false)
}

func Intersection(subject, clip Path) []XY {
	return clipPaths(subject, clip, false, true)
}

func Ink(doc gfx.Doc) {
	subject := RectCenter(Center, XY{.5, .5})

	clip := Circle{
		XY:       XY{.25, .35},
		Radius:   .25,
		Segments: 50,
	}

	/*
		gfx.Fill{
			Shape: subject,
			Color: Red,
		}.Draw(doc)

		gfx.Fill{
			Shape: clip,
			Color: Blue,
		}.Draw(doc)
	*/

	//points := Union(clip.Path(), subject.Path())
	//points := Intersection(clip.Path(), subject.Path())
	//points := Difference(clip.Path(), subject.Path())
	points := Difference(subject.Path(), clip.Path())

	clipped := tess.Tesselate(points)
	gfx.Fill{
		Shape: Triangles(clipped),
		Color: Yellow,
	}.Draw(doc)

	//gfx.Dots{points, Black, .005}.Draw(doc)

}

func flipColor(a bool, b, c RGBA) RGBA {
	if a {
		return b
	}
	return c
}

func PathToLines(p Path) []Line {
	lines := make([]Line, len(p))

	for i, c := range p {
		lines[i] = c.(Line)
	}
	return lines
}

func IntersectLines(a, b Line) (XY, bool) {

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
