package raster

import "github.com/buchanae/ink/dd"

type Scanline struct {
	Y, A, B int
}

func Rasterize(tris []dd.Triangle, w, h int) []Scanline {
	fw := float32(w)
	fh := float32(h)

	var lines []Scanline
	for _, t := range tris {
		x1 := int(t.A.X * fw)
		y1 := int(t.A.Y * fh)
		x2 := int(t.B.X * fw)
		y2 := int(t.B.Y * fh)
		x3 := int(t.C.X * fw)
		y3 := int(t.C.Y * fh)
		lines = rasterize(x1, y1, x2, y2, x3, y3, lines)
	}
	lines = crop(lines, w, h)
	return lines
}

func crop(lines []Scanline, w, h int) []Scanline {
	var out []Scanline
	for _, line := range lines {
		if line.Y < 0 || line.Y >= h {
			continue
		}
		if line.A >= w {
			continue
		}
		if line.B < 0 {
			continue
		}
		line.A = clampInt(line.A, 0, w-1)
		line.B = clampInt(line.B, 0, w-1)
		if line.A > line.B {
			continue
		}
		out = append(out, line)
	}
	return out
}

func clampInt(x, lo, hi int) int {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func rasterize(x1, y1, x2, y2, x3, y3 int, buf []Scanline) []Scanline {
	if y1 > y3 {
		x1, x3 = x3, x1
		y1, y3 = y3, y1
	}
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	if y2 > y3 {
		x2, x3 = x3, x2
		y2, y3 = y3, y2
	}
	if y2 == y3 {
		return rasterizeBottom(x1, y1, x2, y2, x3, y3, buf)
	}
	if y1 == y2 {
		return rasterizeTop(x1, y1, x2, y2, x3, y3, buf)
	}
	x4 := x1 + int((float64(y2-y1)/float64(y3-y1))*float64(x3-x1))
	y4 := y2
	buf = rasterizeBottom(x1, y1, x2, y2, x4, y4, buf)
	buf = rasterizeTop(x2, y2, x4, y4, x3, y3, buf)
	return buf
}

func rasterizeBottom(x1, y1, x2, y2, x3, y3 int, buf []Scanline) []Scanline {
	s1 := float64(x2-x1) / float64(y2-y1)
	s2 := float64(x3-x1) / float64(y3-y1)
	ax := float64(x1)
	bx := float64(x1)
	for y := y1; y <= y2; y++ {
		a := int(ax)
		b := int(bx)
		ax += s1
		bx += s2
		if a > b {
			a, b = b, a
		}
		buf = append(buf, Scanline{y, a, b})
	}
	return buf
}

func rasterizeTop(x1, y1, x2, y2, x3, y3 int, buf []Scanline) []Scanline {
	s1 := float64(x3-x1) / float64(y3-y1)
	s2 := float64(x3-x2) / float64(y3-y2)
	ax := float64(x3)
	bx := float64(x3)
	for y := y3; y > y1; y-- {
		ax -= s1
		bx -= s2
		a := int(ax)
		b := int(bx)
		if a > b {
			a, b = b, a
		}
		buf = append(buf, Scanline{y, a, b})
	}
	return buf
}
