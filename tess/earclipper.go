package tess

import (
	"log"

	"github.com/buchanae/ink/dd"
)

type earclipper struct {
	vert, ear *list
	reflex    map[*dd.XY]struct{}
}

func (ec *earclipper) triangulate(points []dd.XY) []dd.Triangle {
	ec.vert = newList()
	ec.ear = newList()
	ec.reflex = map[*dd.XY]struct{}{}

	// Build a circular, doubly-linked list of points.
	// Allocate a contiguous array of list items upfront.
	items := make([]ecItem, len(points))

	for i, _ := range points {
		item := &items[i]
		item.idx = i
		item.XY = &points[i]
		ec.addVert(item)
	}

	// Build lists of reflex points and potential ears.
	for i, _ := range items {
		item := &items[i]
		if ec.isConcave(item) {
			ec.addReflex(item)
		} else {
			ec.addEar(item)
		}
	}

	// Eliminate potential ears which contain a reflex point.
	it := ec.ear.Iter()
	for it.Next() {
		item := it.Item().ecItem
		if ec.containsReflex(item) {
			//      log.Println("removing potential ear", item.idx)
			ec.delEar(item)
		}
	}

	var tris []dd.Triangle

	// Start clipping ears.
	it = ec.ear.Iter()
	for it.Next() {
		item := it.Item().ecItem

		if item.ear == nil {
			log.Fatalf("WTF: %+v\n", item)
		}

		prev := ec.vert.Prev(item.vert).ecItem
		next := ec.vert.Next(item.vert).ecItem
		// If prev or next point to the current point,
		// then there's only 2 points left and we're done.
		if prev == item || next == item {
			break
		}

		tri := ec.itemTriangle(item)
		tris = append(tris, tri)
		ec.delVert(item)

		// Update adjacent points, recalculating whether they are reflex.
		if prev.isReflex && !ec.isConcave(prev) {
			ec.delReflex(prev)
		}
		if next.isReflex && !ec.isConcave(next) {
			ec.delReflex(next)
		}

		// Update adjacent points, recalculating whether they are ears.
		if !prev.isReflex && !ec.containsReflex(prev) {
			ec.insertEarBefore(prev, item)
		} else {
			ec.delEar(prev)
		}
		if !next.isReflex && !ec.containsReflex(next) {
			ec.insertEarAfter(next, item)
		} else {
			ec.delEar(next)
		}

		ec.delEar(item)
	}

	return tris
}

func (ec *earclipper) containsReflex(item *ecItem) bool {
	prev := ec.vert.Prev(item.vert).XY
	next := ec.vert.Next(item.vert).XY

	for point, _ := range ec.reflex {
		if point == prev || point == next {
			continue
		}

		tri := ec.itemTriangle(item)
		if tri.Contains(*point) {
			return true
		}
	}
	return false
}

func (ec *earclipper) isConcave(item *ecItem) bool {
	prev := *ec.vert.Prev(item.vert).XY
	next := *ec.vert.Next(item.vert).XY
	concave := Concave(prev, *item.XY, next)
	return concave
}

func (ec *earclipper) addEar(item *ecItem) {
	if item.ear != nil {
		// already an ear
		return
	}
	item.ear = &listItem{ecItem: item}
	ec.ear.Append(item.ear)
}

func (ec *earclipper) addVert(item *ecItem) {
	if item.vert != nil {
		// already a vert
		return
	}
	item.vert = &listItem{ecItem: item}
	ec.vert.Append(item.vert)
}

func (ec *earclipper) insertEarBefore(item, index *ecItem) {
	if item.ear != nil {
		// already an ear
		return
	}
	if index.ear == nil {
		panic("cannot insert before nil index")
	}
	item.ear = &listItem{ecItem: item}
	ec.ear.InsertBefore(item.ear, index.ear)
}

func (ec *earclipper) insertEarAfter(item, index *ecItem) {
	if item.ear != nil {
		// already an ear
		return
	}
	if index.ear == nil {
		panic("cannot insert after nil index")
	}
	item.ear = &listItem{ecItem: item}
	ec.ear.InsertAfter(item.ear, index.ear)
}

func (ec *earclipper) delEar(item *ecItem) {
	if item.ear == nil {
		return
	}
	ec.ear.Delete(item.ear)
	item.ear = nil
}

func (ec *earclipper) delVert(item *ecItem) {
	if item.vert == nil {
		return
	}
	ec.vert.Delete(item.vert)
	item.clipped = true
	item.vert = nil
}

func (ec *earclipper) addReflex(item *ecItem) {
	if item.isReflex {
		return
	}
	ec.reflex[item.XY] = struct{}{}
	item.isReflex = true
}

func (ec *earclipper) delReflex(item *ecItem) {
	if !item.isReflex {
		return
	}
	delete(ec.reflex, item.XY)
	item.isReflex = false
}

// ecItem is an item in a linked list used by Triangulate.
type ecItem struct {
	*dd.XY
	idx       int
	vert, ear *listItem
	isReflex  bool
	clipped   bool
}

/*
func Concave(a, b, c dd.XY) bool {
  d := Line{b, a}
  e := Line{b, c}
  ang := RelativeAngle(e, d)
  log.Println("ANG", ang)
  // TODO need to be careful with ordering of points.
  return ang < 0
}
*/

func Concave(a, b, c dd.XY) bool {
	dx1 := b.X - a.X
	dy1 := b.Y - a.Y
	dx2 := c.X - b.X
	dy2 := c.Y - b.Y
	z := dx1*dy2 - dy1*dx2
	return z > 0
}

func (ec *earclipper) itemTriangle(it *ecItem) dd.Triangle {
	prev := ec.vert.Prev(it.vert).ecItem
	next := ec.vert.Next(it.vert).ecItem
	return dd.Triangle{*prev.XY, *it.XY, *next.XY}
}

func dumpItems(items []ecItem) {
	for _, item := range items {
		log.Printf("Item %d: ear %t, reflex %t, clipped %t\n",
			item.idx, item.ear != nil, item.isReflex, item.clipped)
	}
}
