// +build !sendonly

package app

import (
	"github.com/buchanae/ink/render"
	"github.com/buchanae/ink/win"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type App struct {
	conf     Config
	win      *win.Window
	renderer *render.Renderer
	doc      *Doc
	plan     render.Plan
	shown    bool
}

func NewApp(conf Config) (*App, error) {
	return &App{
		conf: conf,
		win: win.NewWindow(win.Config{
			Title:  conf.Window.Title,
			Width:  conf.Window.Width,
			Height: conf.Window.Height,
		}),
	}, nil
}

func (app *App) Run() {
	go app.Do(func() {
		app.win.SetKeyCallback(app.keyCallback)
	})
	app.win.Run()
}

func (app *App) keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

	if key == glfw.KeyX && action == glfw.Press {
		app.snapshot()
	}

	/* TODO redo key events

	case *sdl.KeyboardEvent:
		if z.State == sdl.PRESSED && z.Keysym.Scancode == sdl.SCANCODE_R {
			win.events <- RefreshEvent
		}
		if z.State == sdl.PRESSED && z.Keysym.Scancode == sdl.SCANCODE_RETURN {
			win.events <- ReturnEvent
		}
	*/
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
		if !app.shown {
			app.win.Show()
			app.shown = true
		}
		plan := buildPlan(doc)
		app.initRenderer()

		if doc.Trace {
			app.renderer.StartTrace()
		}

		app.renderer.Render(plan)
		app.renderer.ToScreen(doc.LayerID())

		if doc.Trace {
			app.renderer.EndTrace()
		}

		app.doc = doc
		app.plan = plan

		app.win.Swap()
	})
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
