package color

import (
	"image/color"

	"github.com/muesli/gamut"
)

func Monochromatic(c RGBA, count int) []RGBA {
	res := gamut.Monochromatic(toGoColor(c), count)
	return convert(res)
}

func Shades(c RGBA, count int) []RGBA {
	res := gamut.Shades(toGoColor(c), count)
	return convert(res)
}

func Tints(c RGBA, count int) []RGBA {
	res := gamut.Tints(toGoColor(c), count)
	return convert(res)
}

func Tones(c RGBA, count int) []RGBA {
	res := gamut.Tones(toGoColor(c), count)
	return convert(res)
}

func Blends(a, b RGBA, count int) []RGBA {
	res := gamut.Blends(toGoColor(a), toGoColor(b), count)
	return convert(res)
}

func Lighter(c RGBA, percent float32) RGBA {
	res := gamut.Lighter(toGoColor(c), float64(percent))
	return fromGoColor(res)
}

// TODO Darker doesn't work as expected
func Darker(c RGBA, percent float32) RGBA {
	res := gamut.Darker(toGoColor(c), float64(percent))
	return fromGoColor(res)
}

func convert(in []color.Color) []RGBA {
	out := make([]RGBA, 0, len(in))
	for _, x := range in {
		out = append(out, fromGoColor(x))
	}
	return out
}

func toGoColor(c RGBA) color.Color {
	return color.RGBA{
		R: uint8(c.R * 255),
		G: uint8(c.G * 255),
		B: uint8(c.B * 255),
		A: uint8(c.A * 255),
	}
}

func fromGoColor(c color.Color) RGBA {
	r, g, b, a := c.RGBA()
	return RGBA{
		R: float32(r) / 65536,
		G: float32(g) / 65536,
		B: float32(b) / 65536,
		A: float32(a) / 65536,
	}
}
