package gfx

// TODO want to remove these dependencies from code used by sketches
import (
	"encoding/gob"
	"image"
)

func init() {
	// core shader + mesh
	gob.Register([2]float32{})
	gob.Register([3]float32{})
	gob.Register([4]float32{})

	// images
	gob.Register(Image{})
	gob.Register(&image.RGBA{})
	gob.Register(&image.NRGBA{})
	gob.Register(&image.Gray{})
}
