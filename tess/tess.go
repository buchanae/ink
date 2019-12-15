package tess

import "github.com/buchanae/ink/dd"

func Tesselate(xys []dd.XY) []dd.Triangle {
	if len(xys) < 3 {
		return nil
	}
	if len(xys) == 3 {
		// TODO does order matter?
		return []dd.Triangle{{xys[0], xys[1], xys[2]}}
	}
	// TODO is 4 points a special case? 5?

	ec := earclipper{}
	return ec.triangulate(xys)
}
