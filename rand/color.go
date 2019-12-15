package rand

import "github.com/buchanae/ink/color"

func (r *Rand) Color(c []color.RGBA) color.RGBA {
	if len(c) == 0 {
		return color.Black
	}
	i := r.src.Intn(len(c))
	return c[i]
}

func Color(c []color.RGBA) color.RGBA {
	return src.Color(c)
}

func (r *Rand) Gray() color.RGBA {
	f := r.Float()
	return color.RGBA{f, f, f, 1}
}

func (r *Rand) Palette() []color.RGBA {
	i := r.src.Intn(len(color.Palettes))
	return color.Palettes[i]
}
