package main

// vptree is a tree data structure used for efficient
// nearest neighbor queries on a list of patches.
// TODO currently a mess.
type vptree struct {
	centroid     patch
	radius       int
	inner, outer []patch
	left, right  *vptree
}

func newVPTree(ns []patch) *vptree {
	centroid := centroid(ns)
	radius, inner, outer := split(centroid, ns)

	var left, right *vptree
	if len(inner) > 10 {
		left = newVPTree(inner)
	}
	if len(outer) > 10 {
		right = newVPTree(outer)
	}
	return &vptree{centroid, radius, inner, outer, left, right}
}

// findClosest returns the patch in the tree that is
// closest to the query "q".
func (vp *vptree) findClosest(q patch) (int, patch) {
	d := distance(vp.centroid, q)
	var best patch
	var bestd int

	if d <= vp.radius {
		if vp.left != nil {
			bestd, best = vp.left.findClosest(q)
		} else {
			bestd, best = closest(q, vp.inner)
		}

		if vp.radius-d < bestd {
			var other patch
			var otherd int
			if vp.right != nil {
				otherd, other = vp.right.findClosest(q)
			} else {
				otherd, other = closest(q, vp.outer)
			}
			if otherd < bestd {
				bestd, best = otherd, other
			}
		}
	} else {
		if vp.right != nil {
			bestd, best = vp.right.findClosest(q)
		} else {
			bestd, best = closest(q, vp.outer)
		}

		if d-vp.radius < bestd {
			var other patch
			var otherd int
			if vp.left != nil {
				otherd, other = vp.left.findClosest(q)
			} else {
				otherd, other = closest(q, vp.inner)
			}
			if otherd < bestd {
				bestd, best = otherd, other
			}
		}
	}

	return bestd, best
}

func (vp *vptree) findLeaf(q patch) *vptree {
	d := distance(vp.centroid, q)
	if d <= vp.radius {
		if vp.left != nil {
			return vp.left.findLeaf(q)
		} else {
			return vp
		}
	}

	if vp.right != nil {
		return vp.right.findLeaf(q)
	} else {
		return vp
	}
}

// closest returns the patch from "ns" that is closest to "q".
func closest(q patch, ns []patch) (int, patch) {
	var best patch
	var bestd int
	for i, n := range ns {
		d := distance(q, n)
		if i == 0 || d < bestd {
			best = n
			bestd = d
		}
	}
	return bestd, best
}

// split splits the patches from "ns" into two lists
// based on the distance to "center".
func split(center patch, ns []patch) (radius int, inner, outer []patch) {

	diffs := make([]int, len(ns))
	for i, n := range ns {
		d := distance(center, n)
		diffs[i] = d
		radius += d
	}
	radius /= len(ns)

	for i, d := range diffs {
		if d <= radius {
			inner = append(inner, ns[i])
		} else {
			outer = append(outer, ns[i])
		}
	}
	return
}

func centroid(ns []patch) patch {
	centroid := make([]rgb, patchLen())

	for _, n := range ns {
		for i, c := range n {
			centroid[i].r += c.r
			centroid[i].g += c.g
			centroid[i].b += c.b
		}
	}

	for i := range centroid {
		centroid[i].r /= len(ns)
		centroid[i].g /= len(ns)
		centroid[i].b /= len(ns)
	}
	return centroid
}
