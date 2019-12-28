package render

import (
	"image"
	"log"

	"github.com/buchanae/ink/trace"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Renderer struct {
	width, height int
	multisamples  int
	images        map[image.Image]Image
	shaders       map[Shader]compiled
	textures      map[int]msaa
	layers        []*Layer
}

func NewRenderer(width, height int) *Renderer {
	return &Renderer{
		width:        width,
		height:       height,
		multisamples: 4,
		images:       map[image.Image]Image{},
		shaders:      map[Shader]compiled{},
		textures:     map[int]msaa{},
	}
}

func (r *Renderer) RenderToScreen() error {
	main := r.texture(0)
	main.Clear()

	err := r.render(main)
	if err != nil {
		return err
	}

	trace.Log("blit")
	main.Blit(0)

	trace.Log("done")
	return nil
}

func (r *Renderer) RenderToImage() (image.Image, error) {
	main := r.texture(0)
	main.Clear()

	err := r.render(main)
	if err != nil {
		return nil, err
	}
	return main.Image(), nil
}

/*
TODO
optimize:
4. reuse textures. kinda like a compiler would reuse registers.
   tricky though, because the app could be holding on to a texture ID across frames,
	 so would require knowing that the app doesn't hold a reference to the texture.
*/
func (r *Renderer) render(dst msaa) error {
	trace.Log("render")

	// TODO maybe move these to renderPasses
	glViewport(0, 0, int32(r.width), int32(r.height))
	glEnable(gl.MULTISAMPLE)
	glEnable(gl.BLEND)
	//glEnable(gl.STENCIL_TEST)

	/*
		glBlendFunc(src.factor, dst.factor)
		src comes from the texture being rendered
		dst is the data currently in the color buffer
				(rendered by previous layer/frame)
		result = src.color * src.factor + dst.color * dst.factor
	*/
	glBlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	trace.Log("build")
	pb := newPassBuilder()
	defer pb.Cleanup()

	links := findLinks(r.layers)

	for _, layer := range r.layers {

		output := dst
		if _, ok := links[layer.id]; ok {
			output = r.texture(layer.id)
		}

		pb.AddLayer(layer, output)

		/*
			TODO render image
			rw := float32(loaded.w) / float32(r.width)
			rh := float32(loaded.h) / float32(r.height)
				shader := &gfx.Shader{
					Vert: "/default.vert",
					Frag: "/copy.frag",
					Mesh: gfx.RectXYWH(0.5, 0.5, rw, rh),
					Uniforms: gfx.Uniforms{
						"u_image": loaded.tex,
					},
				}
		*/
	}

	passes := pb.Passes()

	trace.Log("passes %d", len(passes))

	for _, p := range passes {
		r.renderPass(p)
	}

	// TODO garbage collect resources
	return nil
}

func (r *Renderer) renderPass(p *pass) {

	trace.Log("render pass %s", p.name)
	trace.Log("  output to %d", p.output.ID)

	// TODO clear existing program entirely
	// TODO need to cleanup cached buffers
	// TODO redo instancing

	trace.Log("  gl config")
	glBindFramebuffer(gl.FRAMEBUFFER, p.output.Write.FBO)
	glUseProgram(p.prog.id)
	r.bindUniforms(p)
	glBindVertexArray(p.vao)

	trace.Log("  draw elements")
	glDrawElements(
		gl.TRIANGLES,
		int32(p.faceCount),
		gl.UNSIGNED_INT,
		// 4 bytes in each uint32 face index
		glPtrOffset(p.faceOffset*4),
	)

	trace.Log("  output.Paint")
	p.output.Paint()

	trace.Log("  pass done")
}

func (r *Renderer) bindUniforms(p *pass) {
	for _, uni := range p.prog.uniforms {
		val, ok := p.uniforms[uni.Name]
		if !ok {
			log.Printf("  missing uniform: %s", uni.Name)
			continue
		}
		if uni.Type == gl.SAMPLER_2D {
			id, ok := val.(int)
			if !ok {
				log.Printf("  invalid type for texture ID: %T", val)
				continue
			}

			val = r.texture(id)
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

func findLinks(layers []*Layer) map[int]struct{} {
	links := map[int]struct{}{}
	for _, layer := range layers {
		for _, uni := range layer.prog.uniforms {
			val, ok := layer.uniforms[uni.Name]
			if !ok {
				continue
			}
			if uni.Type != gl.SAMPLER_2D {
				continue
			}
			id, ok := val.(int)
			if !ok {
				continue
			}
			links[id] = struct{}{}
		}
	}
	return links
}
