package gfx

// TODO want to remove these dependencies from code used by sketches
import (
	"encoding/gob"
	"image"

	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
)

func init() {
	// core shader + mesh
	gob.Register(Shader{})
	gob.Register(dd.Mesh{})
	gob.Register([2]float32{})
	gob.Register([3]float32{})
	gob.Register([4]float32{})

	// supported in shader attributes
	gob.Register(color.RGBA{})
	gob.Register(dd.XY{})
	gob.Register([]color.RGBA{})
	gob.Register([]dd.XY{})

	// images
	gob.Register(Image{})
	gob.Register(&image.RGBA{})
	gob.Register(&image.NRGBA{})
	gob.Register(&image.Gray{})
}
