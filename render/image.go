package render

import (
	"image"
)

type Image struct {
	id int
	// OpenGL texture handle.
	tex glTexture
}

func (r *Renderer) AddImage(id int, img image.Image) {

	if _, ok := r.images[id]; ok {
		return
	}

	loaded := Image{
		id:  id,
		tex: glCreateTexture(),
	}

	glBindTexture(gl_TEXTURE_2D, loaded.tex)
	glTexParameteri(gl_TEXTURE_2D, gl_TEXTURE_MAG_FILTER, gl_LINEAR)
	glTexParameteri(gl_TEXTURE_2D, gl_TEXTURE_MIN_FILTER, gl_LINEAR)
	glTexParameteri(gl_TEXTURE_2D, gl_TEXTURE_WRAP_S, gl_REPEAT)
	glTexParameteri(gl_TEXTURE_2D, gl_TEXTURE_WRAP_T, gl_REPEAT)

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	pixels := make([]uint8, w*h*4)

	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			k := ((h - j - 1) * w * 4) + (i * 4)
			c := img.At(i, j)
			r, g, b, a := c.RGBA()
			pixels[k+0] = uint8(r)
			pixels[k+1] = uint8(g)
			pixels[k+2] = uint8(b)
			pixels[k+3] = uint8(a)
		}
	}

	glTexImage2D(
		gl_TEXTURE_2D,
		0,
		gl_RGBA,
		int32(w),
		int32(h),
		0,
		gl_RGBA,
		gl_UNSIGNED_BYTE,
		pixels,
	)

	r.images[id] = loaded
}
