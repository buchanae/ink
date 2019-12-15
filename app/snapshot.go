package app

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/buchanae/ink/render"
)

func (app *App) snapshot() {
	width := app.conf.Snapshot.Width
	height := app.conf.Snapshot.Height

	if width == 0 {
		width = app.conf.Window.Width
		height = app.conf.Window.Height
	}

	log.Printf("snapshot: %v x %v", width, height)
	renderer := render.NewRenderer(width, height)
	build(app.doc, renderer)

	img, err := renderer.RenderToImage()
	if err != nil {
		log.Printf("error: rendering snapshot: %v", err)
		return
	}

	dir := app.conf.Snapshot.Dir
	err = writeSnapshot(dir, img)
	if err != nil {
		log.Printf("error: writing snapshot: %v", err)
	}
}

func writeSnapshot(dir string, img image.Image) error {
	err := ensureDir(dir)
	if err != nil {
		return err
	}

	stamp := time.Now().Format("01-02-2006-15-04-05")
	name := stamp + ".png"
	name = filepath.Join(dir, name)

	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func ensureDir(path string) error {
	// Check that the data directory exists.
	s, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0700)
		if err != nil {
			return fmt.Errorf("creating data directory: %v", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("checking for data directory: %v", err)
	}

	if !s.IsDir() {
		return fmt.Errorf("%q is a file, but mailer needs to put a directory here", path)
	}
	return nil
}
