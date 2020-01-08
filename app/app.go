package app

import (
	"encoding/gob"
	"os"
	"time"

	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/render"
	"github.com/buchanae/ink/win"
)

func Run(f func(*gfx.Doc)) {
	a, err := NewApp(DefaultConfig())
	if err != nil {
		panic(err)
	}
	doc := gfx.NewDoc()
	go func() {
		f(doc)
		a.Render(doc)
	}()
	a.Run()
}

type Frame struct {
	Time time.Time
}

func Send(f func(*gfx.Doc)) {

	doc := gfx.NewDoc()
	f(doc)

	err := gob.NewEncoder(os.Stdout).Encode(doc)
	if err != nil {
		os.Stderr.Write([]byte("sending: "))
		os.Stderr.Write([]byte(err.Error()))
		os.Stderr.Write([]byte("\n"))
	}
}

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
