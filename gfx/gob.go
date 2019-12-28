package gfx

// TODO want to remove these dependencies from code used by sketches
import (
	"encoding/gob"
	"image"

	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
)

func init() {
	gob.Register(Shader{})
	gob.Register(color.RGBA{})
	gob.Register(dd.XY{})
	gob.Register(dd.Mesh{})
	gob.Register(dd.Rect{})
	gob.Register(dd.Quad{})
	gob.Register(dd.Triangle{})
	gob.Register(dd.Circle{})
	gob.Register(dd.Triangles{})
	gob.Register([]color.RGBA{})
	gob.Register([]dd.XY{})
	gob.Register(image.RGBA{})
}
