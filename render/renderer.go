package render

import (
	"image"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Renderer struct {
	width, height int
	multisamples  int
	layerID       int
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
	main.Blit(0)
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
1. don't bother writing to a separate texture if the output is never used.
   just draw to the main framebuffer.
4. reuse textures. kinda like a compiler would reuse registers.
   tricky though, because the app could be holding on to a texture ID across frames,
	 so would require knowing that the app doesn't hold a reference to the texture.
*/
func (r *Renderer) render(dst msaa) error {
	defer traceTime("render", time.Now())

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

	pb := passBuilder{}
	links := findLinks(r.layers)

	start := time.Now()
	for _, layer := range r.layers {

		output := dst
		if _, ok := links[layer.id]; ok {
			output = r.texture(layer.id)
		}

		prog, err := r.compile(layer.shader)
		if err != nil {
			return err
		}

		pb.AddLayer(prog, layer, output)

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

		if !layer.hide {
			//pb.Composite(dst, output)
		}
	}
	passes := pb.Passes()

	trace("prep %s", time.Since(start))
	trace("passes %d", len(passes))

	for _, p := range passes {
		r.renderPasses(p)
	}

	// TODO garbage collect resources
	return nil
}

func (r *Renderer) renderPasses(p *pass) {

	renderPassStart := time.Now()
	trace("render pass %s", p.name)

	// TODO clear existing program entirely
	// TODO need to cleanup cached buffers
	// TODO determine number of clones (instances).
	//count := 1

	start := time.Now()

	glBindFramebuffer(gl.FRAMEBUFFER, p.output.Write.FBO)
	glUseProgram(p.prog.id)
	r.bindUniforms(p)
	glBindVertexArray(p.vao)

	trace("  render config %s", time.Since(start))

	start = time.Now()
	glDrawElements(
		gl.TRIANGLES,
		int32(p.faceCount),
		gl.UNSIGNED_INT,
		// 4 bytes in each uint32 face index
		glPtrOffset(p.faceOffset*4),
	)
	trace("  draw elements %s", time.Since(start))

	start = time.Now()
	p.output.Paint()

	trace("  output paint %s", time.Since(start))
	trace("  render pass took %s", time.Since(renderPassStart))
}

func (r *Renderer) bindUniforms(p *pass) {
	for _, uni := range p.prog.uniforms {
		val, ok := p.uniforms[uni.Name]
		if !ok {
			log("  missing uniform: %s", uni.Name)
			continue
		}
		trace("  bind uniform %s to %#v", uni.Name, val)

		err := uni.Bind(val)
		if err != nil {
			log("  ERR: binding uniform %q: %s", uni.Name, err)
		}
	}
}

func (r *Renderer) texture(id int) msaa {
	t, ok := r.textures[id]
	if !ok {
		t = newMsaa(r.width, r.height, r.multisamples)
		r.textures[id] = t
	}
	return t
}

func findLinks(layers []*Layer) map[int]struct{} {
	links := map[int]struct{}{}
	for _, layer := range layers {
		for _, uni := range layer.uniforms {
			if l, ok := uni.(*Layer); ok {
				links[l.id] = struct{}{}
			}
		}
	}
	return links
}
