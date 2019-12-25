// +build !sendonly

package app

import (
	"github.com/buchanae/ink/gfx"
)

// Render renders the doc with default config and blocks until the window is closed.
// If there is no window open, Render will open one.
// If an error occurs while opening the window, Render panics.
// If an error occurs while rendering the doc, Render logs the error.
func Render(layer *gfx.Layer) {
	app, err := NewApp(DefaultConfig())
	if err != nil {
		panic(err)
	}
	go app.Render(layer)
	app.Run()
}
