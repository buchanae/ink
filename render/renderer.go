package render

import (
	"image"
	"log"
)

type Renderer struct {
	tracer

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

func (r *Renderer) Render(plan Plan) {
	r.render(plan)
}

func (r *Renderer) ToScreen(layerID int) {
	r.texture(layerID).Blit(gl_SCREEN)
}

// TODO want a better capture API that allows flexible capturing
//      possibly via a single Capture() function
func (r *Renderer) CapturePixels(layerID int, x, y, w, h float32) []uint8 {
	return r.texture(layerID).Pixels(x, y, w, h)
}

func (r *Renderer) CaptureImage(layerID int, x, y, w, h float32) image.Image {
	return r.texture(layerID).Image(x, y, w, h)
}

func (r *Renderer) render(plan Plan) {
	r.trace("start render")

	for id, img := range plan.Images {
		r.AddImage(id, img)
	}

	pb := &build{
		tracer: r.tracer,
	}
	pb.build(plan)

	r.trace("passes %d", len(pb.passes))

	up := uploadBuild(pb)
	defer up.cleanup()

	glViewport(0, 0, int32(r.width), int32(r.height))
	//glEnable(gl_MULTISAMPLE)
	glEnable(gl_BLEND)

	for _, p := range pb.passes {
		r.renderPass(p)
		glLogErr("pass")
	}
}

func (r *Renderer) renderPass(p *pass) {

	/*
		glBlendFunc(src.factor, dst.factor)
		src is the new color being added
		dst is the existing color
		result = src.color * src.factor + dst.color * dst.factor
	*/
	glBlendFunc(gl_SRC_ALPHA, gl_ONE_MINUS_SRC_ALPHA)

	r.trace("render pass %s", p.name)
	r.trace("  output to %d", p.layer)

	// TODO clear existing program entirely
	count := 1
	if p.instanceCount > 1 {
		count = p.instanceCount
	}

	r.trace("  gl config")
	output := r.texture(p.layer)

	glBindFramebuffer(gl_FRAMEBUFFER, output.Write.FBO)
	glUseProgram(p.prog.id)
	r.bindUniforms(p)
	glBindVAO(p.vao)

	r.trace("  draw elements")
	glDrawElementsInstanced(
		gl_TRIANGLES,
		int32(p.faceCount),
		gl_UNSIGNED_INT,
		// 4 bytes in each uint32 face index
		p.faceOffset*4,
		int32(count),
	)

	output.Paint()

	r.trace("  pass done")
}

// TODO bind uniforms should be a dead simple loop
//      without any preprocessing (logic, type checking, etc).
//      move all preprocessing to something like pass builder.
//
//      consider exposing a public resource holding a handle
//      to the preprocessed passes and resources. might
//      make a clear separation between preprocessing and
//      execution. might help abstract rendering backends later.
func (r *Renderer) bindUniforms(p *pass) {
	for _, uni := range p.prog.uniforms {

		val, ok := p.uniforms[uni.Name]
		if !ok {
			// TODO return error list from render
			log.Printf("  missing uniform: %s", uni.Name)
			continue
		}

		if uni.Type == gl_SAMPLER_2D {
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
