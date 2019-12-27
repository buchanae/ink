// +build !sendonly

package app

import (
	"log"

	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/render"
	"github.com/buchanae/ink/trace"
	"github.com/buchanae/ink/win"
)

// App can render a gfx.Layer to a window.
type App struct {
	conf     Config
	win      *win.Window
	renderer *render.Renderer
	doc      *gfx.Layer
}

// NewApp will open a new window.
func NewApp(conf Config) (*App, error) {
	return &App{
		conf: conf,
		win:  win.NewWindow(conf.Window),
		renderer: render.NewRenderer(
			conf.Window.Width, conf.Window.Height,
		),
	}, nil
}

func (app *App) Run() {
	go func() {
		for {
			select {
			case ev := <-app.win.Events():

				switch ev {
				case win.QuitEvent:
					return

				case win.SnapshotEvent:
					if app.doc != nil {
						app.win.Do(app.snapshot)
					}

					//case events.Refresh:
					//refresh()
				}

				//case <-assets.Changed():
				//refresh()
			}
		}
	}()
	app.win.Run()
}

// Render renders the doc to the app window.
func (app *App) Render(doc *gfx.Layer) {
	app.win.Do(func() {
		app.renderer.ClearLayers()

		trace.Log("start build")
		b := builder{renderer: app.renderer}
		b.build(doc)
		trace.Log("built")

		err := app.renderer.RenderToScreen()
		if err != nil {
			log.Printf("error: rendering: %v", err)
		}
	})

	app.doc = doc
	app.win.Swap()
}
