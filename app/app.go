// +build !sendonly

package app

import (
	"runtime"

	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/render"
	"github.com/buchanae/ink/render/opengl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// Code that interacts with the OS window APIs,
	// (GLFW and OpenGL calls) must run on the main thread.
	//
	// This locks the main Go goroutine to the main thread,
	// i.e. main() will run on the main thread. App.Run must
	// be called from this main goroutine. After that, all
	// interaction with OpenGL and GLFW must go through App.Do.
	runtime.LockOSThread()
}

type App struct {
	conf      Config
	win       *glfw.Window
	renderer  *opengl.Renderer
	commands  chan func()
	plan      render.Plan
	shown     bool
	keyEvents chan KeyEvent
	addKeycb  chan KeyCallback
}

func NewApp(conf Config) (*App, error) {
	return &App{
		conf:      conf,
		commands:  make(chan func()),
		keyEvents: make(chan KeyEvent, 100),
		addKeycb:  make(chan KeyCallback),
	}, nil
}

func (app *App) Run() error {
	go app.processKeyEvents()
	return app.runWindow()
}

func (app *App) Close() {
	app.Do(func() {
		app.win.SetShouldClose(true)
	})
}

func (app *App) initRenderer() {
	if app.renderer != nil {
		return
	}
	w, h := app.win.GetFramebufferSize()
	app.renderer = opengl.NewRenderer(w, h)
}

func (app *App) SetConfig(b Config) {
	app.Do(func() {
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
	})
}

func (app *App) Render(doc *gfx.Doc) {
	//app.SetConfig(doc.Config)
	plan := doc.Plan()
	app.RenderPlan(plan)
}

func (app *App) RenderPlan(plan render.Plan) {
	app.Do(func() {
		if !app.conf.Window.Hidden && !app.shown {
			app.win.Show()
			app.shown = true
		}
		app.initRenderer()

		app.renderer.Render(plan)
		app.renderer.ToScreen(plan.RootLayer)

		app.plan = plan
		app.win.SwapBuffers()
	})
}

// Do queues a function for execution on the main thread.
// OS windows typically require that code which accesses
// windows
func (app *App) Do(f func()) {
	done := make(chan struct{})
	app.commands <- func() {
		f()
		close(done)
	}
	<-done
}

func Run(f func(*gfx.Doc)) {
	conf := DefaultConfig()
	a, err := NewApp(conf)
	if err != nil {
		panic(err)
	}
	doc := gfx.NewDoc()
	//doc.Config = conf
	go func() {
		f(doc)
		a.Render(doc)
	}()
	a.Run()
}
