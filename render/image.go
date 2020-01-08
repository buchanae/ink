package render

import (
	"image"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Image struct {
	id int
	// OpenGL texture handle.
	tex uint32
}

func (r *Renderer) AddImage(id int, img image.Image) {

	if _, ok := r.images[id]; ok {
		return
	}

	loaded := Image{id: id}

	glGenTextures(1, &loaded.tex)
	glBindTexture(gl.TEXTURE_2D, loaded.tex)
	glTexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	glTexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	glTexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	glTexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

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
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(w),
		int32(h),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		glPtr(pixels),
	)

	r.images[id] = loaded
}
