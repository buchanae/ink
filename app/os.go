// +build !js

package app

import (
	"image"
	"os"

	"github.com/buchanae/ink/gfx"
)

func (d *Doc) LoadImage(path string) gfx.Image {
	fh, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	img, _, err := image.Decode(fh)
	if err != nil {
		panic(err)
	}

	return d.NewImage(img)
}
