package render

import (
	"image"
	"log"
)

// msaa represents a multisample anti-aliased OpenGL texture.
// This is the main resource that shaders write to and read from.
//
// The renderer writes to a framebuffer containing a normal 2D texture.
// Then the renderer calls msaa.Paint() to blit that texture into an antialiased
// texture. Downstream shaders that read from this texture will read from the
// antialiased texture.
type msaa struct {
	ID   int
	Read struct {
		FBO glFramebuffer
		Tex glTexture
	}
	Write struct {
		FBO glFramebuffer
		Tex glTexture
	}
	Width, Height, Multisamples int
	DisableMultisample          bool
}

func newMsaa(id, w, h, multisamples int) msaa {

	m := msaa{
		ID:           id,
		Width:        w,
		Height:       h,
		Multisamples: multisamples,
	}

	// Create two textures:
	// 1. a multisampled texture which will be written to.
	// 2. a normal texture which will be read from.
	m.Read.Tex = glCreateTexture()
	m.Write.Tex = glCreateTexture()

	// ...and two framebuffers, one for each texture.
	//m.Read.FBO = glCreateFramebuffer()
	//m.Write.FBO = glCreateFramebuffer()
	m.Read.FBO = gl_SCREEN
	m.Write.FBO = gl_SCREEN

	// Initialize the Read texture
	glBindFramebuffer(gl_FRAMEBUFFER, m.Read.FBO)

	glBindTexture(gl_TEXTURE_2D, m.Read.Tex)
	glTexImage2D(
		gl_TEXTURE_2D,
		0,
		gl_RGBA,
		int32(m.Width),
		int32(m.Height),
		0,
		gl_RGBA,
		gl_UNSIGNED_BYTE,
		nil)

	glTexParameteri(gl_TEXTURE_2D, gl_TEXTURE_MAG_FILTER, gl_LINEAR)
	glTexParameteri(gl_TEXTURE_2D, gl_TEXTURE_MIN_FILTER, gl_LINEAR)
	glTexParameteri(gl_TEXTURE_2D, gl_TEXTURE_WRAP_S, gl_REPEAT)
	glTexParameteri(gl_TEXTURE_2D, gl_TEXTURE_WRAP_T, gl_REPEAT)

	/*
		glFramebufferTexture2D(
			gl_FRAMEBUFFER,
			gl_COLOR_ATTACHMENT0,
			gl_TEXTURE_2D,
			m.Read.Tex,
			0,
		)
	*/

	m.DisableMultisample = true

	// Initialize the Write texture
	if m.DisableMultisample {
		m.Write.FBO = m.Read.FBO
		m.Write.Tex = m.Read.Tex
	} else {
		glBindFramebuffer(gl_FRAMEBUFFER, m.Write.FBO)
		glBindTexture(gl_TEXTURE_2D_MULTISAMPLE, m.Write.Tex)

		glTexImage2DMultisample(
			gl_TEXTURE_2D_MULTISAMPLE,
			int32(m.Multisamples),
			gl_RGBA,
			int32(m.Width),
			int32(m.Height),
			false,
		)

		glFramebufferTexture2D(
			gl_FRAMEBUFFER,
			gl_COLOR_ATTACHMENT0,
			gl_TEXTURE_2D_MULTISAMPLE,
			m.Write.Tex,
			0,
		)
	}

	m.Clear()
	return m
}

func (m msaa) Clear() {
	glBindFramebuffer(gl_FRAMEBUFFER, m.Read.FBO)
	glClearColor(0, 0, 0, 1)
	glClear(gl_COLOR_BUFFER_BIT)

	glBindFramebuffer(gl_FRAMEBUFFER, m.Write.FBO)
	glClearColor(0, 0, 0, 1)
	glClear(gl_COLOR_BUFFER_BIT)
}

func (m msaa) Paint() {
	if m.DisableMultisample {
		return
	}
	// Copy the multisample texture (Write)
	// to the regular texture (Read).
	glBindFramebuffer(gl_DRAW_FRAMEBUFFER, m.Read.FBO)
	glBindFramebuffer(gl_READ_FRAMEBUFFER, m.Write.FBO)

	glBlitFramebuffer(
		0, 0, int32(m.Width), int32(m.Height),
		0, 0, int32(m.Width), int32(m.Height),
		gl_COLOR_BUFFER_BIT,
		gl_LINEAR,
	)
}

// Blit the anti-aliased "read" texture to the given framebuffer ID.
// Used during compisiting to copy textures to the screen.
func (m msaa) Blit(to glFramebuffer) {
	log.Printf("BLIT: %#v", to)
	return
	glBindFramebuffer(gl_DRAW_FRAMEBUFFER, to)
	glBindFramebuffer(gl_READ_FRAMEBUFFER, m.Read.FBO)

	glBlitFramebuffer(
		0, 0, int32(m.Width), int32(m.Height),
		0, 0, int32(m.Width), int32(m.Height),
		gl_COLOR_BUFFER_BIT,
		gl_LINEAR,
	)
}

func (m msaa) Pixels(x, y, w, h float32) []uint8 {

	glBindFramebuffer(gl_READ_FRAMEBUFFER, m.Read.FBO)
	xi := int(x * float32(m.Width))
	yi := int(y * float32(m.Height))
	wi := int(w * float32(m.Width))
	hi := int(h * float32(m.Height))

	if xi < 0 {
		xi = 0
	}
	if yi < 0 {
		yi = 0
	}
	if wi > m.Width {
		wi = m.Width
	}
	if hi > m.Height {
		hi = m.Height
	}

	pixels := make([]uint8, wi*hi*4)

	// TODO how to allow flexible querying without complex API?
	//      e.g. get only red pixels
	// 			Also, how to allow float texture?
	glReadPixels(
		int32(xi),
		int32(yi),
		int32(wi),
		int32(hi),
		gl_RGBA, gl_UNSIGNED_BYTE, pixels,
	)

	return pixels
}

func (m msaa) Image(x, y, w, h float32) image.Image {

	pixels := m.Pixels(x, y, w, h)

	r := image.Rect(0, 0, m.Width, m.Height)
	img := image.NewRGBA(r)
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			i := img.PixOffset(x, y)
			// the orientation of PNG vs OpenGL is upside-down.
			j := ((m.Height - y - 1) * m.Width * 4) + (x * 4)
			img.Pix[i+0] = pixels[j+0]
			img.Pix[i+1] = pixels[j+1]
			img.Pix[i+2] = pixels[j+2]
			// TODO difficult to retrieve pixel data where alpha
			//      hasn't been premultiplied
			img.Pix[i+3] = 255
		}
	}

	return img
}

func (m msaa) Destroy() {
	glDeleteTexture(m.Read.Tex)
	glDeleteTexture(m.Write.Tex)
	glDeleteFramebuffer(m.Read.FBO)
	glDeleteFramebuffer(m.Write.FBO)
}
