package app

import (
	"encoding/gob"

	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
)

func init() {
	gob.Register(gfx.Shader{})
	gob.Register(color.RGBA{})
	gob.Register(dd.XY{})
	gob.Register([]color.RGBA{})
	gob.Register([]dd.XY{})
}
