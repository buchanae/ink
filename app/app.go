package app

import (
	"log"

	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/render"
	"github.com/buchanae/ink/win"
)

// App can render a gfx.Doc to a window.
type App struct {
	conf     Config
	win      *win.Window
	renderer *render.Renderer
	doc      *gfx.Doc
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
func (app *App) Render(doc *gfx.Doc) {
	// TODO it's unclear which renderer methods are thread safe
	app.renderer.ClearLayers()

	app.win.Do(func() {
		err := build(doc, app.renderer)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		err = app.renderer.RenderToScreen()
		if err != nil {
			log.Printf("error: rendering: %v", err)
		}
	})

	app.doc = doc
	app.win.Swap()
}
