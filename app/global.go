package app

import (
	"github.com/buchanae/ink/gfx"
)

// Render renders the doc with default config and blocks until the window is closed.
// If there is no window open, Render will open one.
// If an error occurs while opening the window, Render panics.
// If an error occurs while rendering the doc, Render logs the error.
func Render(doc *gfx.Doc) {

	app, err := NewApp(DefaultConfig())
	if err != nil {
		panic(err)
	}
	go app.Render(doc)

	// Most access to the window must be done on a single OS thread,
	// so this code locks itself to the OS thread and handles all communication
	// via SDL queues and Go channels.
	app.Run()
}
