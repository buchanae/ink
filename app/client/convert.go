package client

import (
	"log"

	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
)

func convertAttr(val interface{}, verts int) []float32 {

	if val == nil {
		return nil
	}

	switch z := val.(type) {

	case []float32:
		return z

	case float32:
		data := make([]float32, verts)
		for i := range data {
			data[i] = z
		}
		return data

	case float64:
		data := make([]float32, verts)
		for i := range data {
			data[i] = float32(z)
		}
		return data

	case []color.RGBA:
		data := make([]float32, 0, len(z)*4)
		for _, v := range z {
			data = append(data, v.R, v.G, v.B, v.A)
		}
		return data

	case color.RGBA:
		data := make([]float32, 0, verts*4)
		for i := 0; i < verts; i++ {
			data = append(data, z.R, z.G, z.B, z.A)
		}
		return data

	case []dd.XY:
		data := make([]float32, 0, len(z)*2)
		for _, v := range z {
			data = append(data, v.X, v.Y)
		}
		return data

	case dd.XY:
		data := make([]float32, 0, verts*2)
		for i := 0; i < verts; i++ {
			data = append(data, z.X, z.Y)
		}
		return data

	default:
		log.Printf("unsupported attribute type %T", z)
		return nil
	}
}

func convertUniform(v interface{}) interface{} {
	switch z := v.(type) {
	case dd.XY:
		return [2]float32{z.X, z.Y}
	case color.RGBA:
		return [4]float32{z.R, z.G, z.B, z.A}
	case gfx.Image:
		return z.ID
	default:
		return v
	}
}
