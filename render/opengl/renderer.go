package opengl

import (
	"image"
	"log"

	"github.com/buchanae/ink/render"
	"github.com/buchanae/ink/trac"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Renderer struct {
	width, height int
	multisamples  int
	textures      map[int]msaa
	images        map[int]Image
}

func NewRenderer(width, height int) *Renderer {
	return &Renderer{
		width:        width,
		height:       height,
		multisamples: 4,
		textures:     map[int]msaa{},
		images:       map[int]Image{},
	}
}

func (r *Renderer) Cleanup() {
	for _, tex := range r.textures {
		tex.Destroy()
	}

	for _, img := range r.images {
		img.Destroy()
	}

	r.textures = map[int]msaa{}
	r.images = map[int]Image{}
}

func (r *Renderer) Render(plan render.Plan) {
	r.render(plan)
}

func (r *Renderer) ToScreen(layerID int) {
	r.texture(layerID).Blit(0)
}

// TODO want a better capture API that allows flexible capturing
//      possibly via a single Capture() function
func (r *Renderer) CapturePixels(layerID int, x, y, w, h float32) []uint8 {
	return r.texture(layerID).Pixels(x, y, w, h)
}

func (r *Renderer) CaptureImage(layerID int, x, y, w, h float32) image.Image {
	return r.texture(layerID).Image(x, y, w, h)
}

func (r *Renderer) render(plan render.Plan) {
	trac.Log("start render")

	for _, tex := range r.textures {
		tex.Clear()
	}

	for id, img := range plan.Images {
		r.AddImage(id, img)
	}

	pb := &build{}
	passes := pb.build(plan)
	defer pb.cleanup()

	trac.Log("passes %d", len(passes))

	for _, p := range passes {
		r.renderPass(p)
	}
	trac.Log("passes done")
}

func (r *Renderer) renderPass(p buildPass) {
	// TODO clear existing program entirely

	glViewport(0, 0, int32(r.width), int32(r.height))
	glEnable(gl.MULTISAMPLE)
	glEnable(gl.BLEND)

	/*
		glBlendFunc(src.factor, dst.factor)
		src is the new color being added
		dst is the existing color
		result = src.color * src.factor + dst.color * dst.factor
	*/
	glBlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	trac.Log("render pass %s", p.Name)
	trac.Log("  output to %d", p.Layer)

	count := 1
	if p.Instances > 1 {
		count = p.Instances
	}

	trac.Log("  gl config")
	output := r.texture(p.Layer)

	glBindFramebuffer(gl.FRAMEBUFFER, output.Write.FBO)
	glUseProgram(p.prog.id)
	r.bindUniforms(p)
	glBindVertexArray(p.vao)

	// TODO indexed elements are not always a win.
	//      if most verts are unique, then indicies are just
	//      overhead.
	elcount := p.Faces.Count * 3
	// 4 bytes in each uint32 face index
	offset := p.Faces.Offset * 4

	trac.Log("  draw elements: %d %d", elcount, offset)

	glDrawElementsInstanced(
		gl.TRIANGLES,
		int32(elcount),
		gl.UNSIGNED_INT,
		glPtrOffset(offset),
		int32(count),
	)

	output.Paint()

	trac.Log("  pass done")
}

// TODO bind uniforms should be a dead simple loop
//      without any preprocessing (logic, type checking, etc).
//      move all preprocessing to something like pass builder.
//
//      consider exposing a public resource holding a handle
//      to the preprocessed passes and resources. might
//      make a clear separation between preprocessing and
//      execution. might help abstract rendering backends later.
func (r *Renderer) bindUniforms(p buildPass) {
	for _, uni := range p.prog.uniforms {

		val, ok := p.Uniforms[uni.Name]
		if !ok {
			// TODO return error list from render
			log.Printf("  missing uniform: %s", uni.Name)
			continue
		}

		if uni.Type == gl.SAMPLER_2D {
			id, ok := val.(int)
			if !ok {
				log.Printf("  invalid type for texture ID: %T", val)
				continue
			}

			tex, ok := r.textures[id]
			if !ok {
				img, ok := r.images[id]
				if !ok {
					log.Printf("  missing texture ID: %d", id)
					continue
				} else {
					val = img
				}
			} else {
				val = tex
			}
		}

		err := uni.Bind(val)
		if err != nil {
			log.Printf("error: binding uniform %q: %s", uni.Name, err)
		}
	}
}

func (r *Renderer) texture(id int) msaa {
	t, ok := r.textures[id]
	if !ok {
		t = newMsaa(id, r.width, r.height, r.multisamples)
		r.textures[id] = t
	}
	return t
}
