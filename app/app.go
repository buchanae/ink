// +build !sendonly

package app

import (
	"github.com/buchanae/ink/render"
	"github.com/buchanae/ink/win"
)

type App struct {
	conf     Config
	win      *win.Window
	renderer *render.Renderer
	doc      *Doc
	events   chan win.Event
}

func NewApp(conf Config) (*App, error) {
	return &App{
		conf: conf,
		win: win.NewWindow(win.Config{
			Title:  conf.Window.Title,
			Width:  conf.Window.Width,
			Height: conf.Window.Height,
		}),
		events: make(chan win.Event, 1000),
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
				case win.SnapshotEvent:
					if app.doc != nil {
						app.win.Do(app.snapshot)
					}
				}
			}
		}
	}()
	app.win.Run()
}

func (app *App) Events() <-chan win.Event {
	return app.events
}

func (app *App) initRenderer() {
	if app.renderer != nil {
		return
	}
	w, h := app.win.GetFramebufferSize()
	app.renderer = render.NewRenderer(w, h)
}

func (app *App) updateConfig(b Config) {
	app.conf.Snapshot = b.Snapshot

	aw := app.conf.Window
	bw := b.Window
	if aw.Width != bw.Width || aw.Height != bw.Height {
		// reset renderer
		// TODO this is bound to cause some issue
		//      when rendering multiple docs.
		//      need a better way to resize a renderer
		app.renderer = nil
		app.win.SetSize(bw.Width, bw.Height)
	}
	if aw.Title != bw.Title {
		app.win.SetTitle(bw.Title)
	}

	app.conf.Window = b.Window
}

func (app *App) Render(doc *Doc) {
	app.Do(func() {
		app.updateConfig(doc.Config)
		app.win.Show()
		plan := buildPlan(doc)
		app.initRenderer()
		app.renderer.Render(plan)
		app.renderer.ToScreen(doc.LayerID())
	})

	app.doc = doc
	app.win.Swap()
}

func (app *App) Do(f func()) {
	done := make(chan struct{})
	app.win.Do(func() {
		f()
		close(done)
	})
	<-done
}

func Run(f func(*Doc)) {
	a, err := NewApp(DefaultConfig())
	if err != nil {
		panic(err)
	}
	doc := NewDoc()
	go func() {
		f(doc)
		a.Render(doc)
	}()
	a.Run()
}
