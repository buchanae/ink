package app

import (
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/render"
	"github.com/buchanae/ink/win"
)

// App can render a gfx.Layer to a window.
type App struct {
	conf     Config
	win      *win.Window
	renderer *render.Renderer
	doc      *gfx.Doc
	events   chan win.Event
}

// NewApp will open a new window.
func NewApp(conf Config) (*App, error) {
	return &App{
		conf:   conf,
		win:    win.NewWindow(conf.Window),
		events: make(chan win.Event, 1000),
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

				select {
				case app.events <- ev:
				default:
				}

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

func (app *App) Events() <-chan win.Event {
	return app.events
}

func (app *App) Render(doc *gfx.Doc) {
	app.win.Do(func() {
		plan := buildPlan(doc)
		app.renderer.Render(plan)
		app.renderer.ToScreen(doc.LayerID())
	})

	app.doc = doc
	app.win.Swap()
}

func (app *App) Swap() {
	app.win.Swap()
}

func (app *App) Do(f func(*render.Renderer)) {

	done := make(chan struct{})
	app.win.Do(func() {
		f(app.renderer)
		close(done)
	})
	<-done
}
