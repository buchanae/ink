package gfx

import "github.com/buchanae/ink/dd"

type Image struct {
	ID   int
	Rect dd.Rect
}

func (img Image) Draw(out Layer) {
	out.AddShader(&Shader{
		Vert: DefaultVert,
		Frag: CopyFrag,
		Mesh: img.Rect.Fill(),
		Attrs: Attrs{
			"u_image": img,
		},
	})
}
