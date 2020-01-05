package rand

import "github.com/buchanae/ink/dd"

func TweakMesh(mesh dd.Mesh, amount float32) dd.Mesh {
	mesh = mesh.Copy()
	for i := range mesh.Verts {
		mesh.Verts[i] = TweakXY(mesh.Verts[i], amount)
	}
	return mesh
}

func TweakXY(xy dd.XY, amount float32) dd.XY {
	rx := src.Float()*amount*2 - amount
	ry := src.Float()*amount*2 - amount
	return dd.XY{xy.X + rx, xy.Y + ry}
}

func TweakQuad(b dd.Quad, amount float32) dd.Quad {
	return dd.Quad{
		A: TweakXY(b.A, amount),
		B: TweakXY(b.B, amount),
		C: TweakXY(b.C, amount),
		D: TweakXY(b.D, amount),
	}
}
