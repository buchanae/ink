package win

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	runtime.LockOSThread()
}

var initSDL, initGL sync.Once

type Config struct {
	Name          string
	X, Y          int
	Width, Height int
}

// Window holds a handle to an OS window.
type Window struct {
	conf     Config
	events   chan Event
	commands chan func()
	window   *sdl.Window
}

// NewWindow opens a new window.
func NewWindow(conf Config) *Window {
	win := &Window{
		conf:     conf,
		events:   make(chan Event),
		commands: make(chan func(), 1000),
	}
	return win
}

// Events returns a stream of window events,
// such as quit, keyboard, mouse, etc.
func (win *Window) Events() <-chan Event {
	return win.events
}

// Do queues a function for execution on the main thread.
// OS windows typically require that code which accesses
// windows
func (win *Window) Do(cmd func()) {
	win.commands <- cmd
}

func (win *Window) Swap() {
	win.commands <- func() {
		win.window.GLSwap()
	}
}

// run handles all sdl window operations.
func (win *Window) Run() {

	var err error
	initSDL.Do(func() {
		err = sdl.Init(sdl.INIT_EVERYTHING)
	})
	if err != nil {
		// TODO not sure what to do with errors
		log.Printf("error: initializing sdl: %v", err)
		return
	}
	// TODO this is not the right place/way to call this,
	//      because there could be multiple windows on a thread.
	defer sdl.Quit()

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

	window, err := sdl.CreateWindow(
		win.conf.Name,
		int32(win.conf.X), int32(win.conf.Y),
		int32(win.conf.Width), int32(win.conf.Height),
		sdl.WINDOW_OPENGL,
	)
	if err != nil {
		// TODO not sure what to do with errors
		log.Printf("error: creating window: %v", err)
		return
	}
	defer window.Destroy()

	context, err := window.GLCreateContext()
	if err != nil {
		// TODO not sure what to do with errors
		log.Printf("error: creating opengl context: %v", err)
		return
	}
	defer sdl.GLDeleteContext(context)

	// TODO once the Renderer methods are thread safe,
	//      move this to Renderer.render()
	initGL.Do(func() {
		err = gl.Init()
	})
	if err != nil {
		log.Printf("error: initializing OpenGL: %v", err)
		return
	}

	win.window = window
	check := time.Tick(20 * time.Millisecond)

	for {
		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
			switch z := ev.(type) {

			case *sdl.WindowEvent:

			case *sdl.QuitEvent:
				win.events <- QuitEvent
				// TODO figure out how to quit gracefully
				//      currently not closing events/commands channels
				return

			case *sdl.KeyboardEvent:
				if z.State == sdl.PRESSED && z.Keysym.Scancode == sdl.SCANCODE_X {
					win.events <- SnapshotEvent
				}
				if z.State == sdl.PRESSED && z.Keysym.Scancode == sdl.SCANCODE_R {
					win.events <- RefreshEvent
				}
				if z.State == sdl.PRESSED && z.Keysym.Scancode == sdl.SCANCODE_RETURN {
					win.events <- ReturnEvent
				}
			}
		}

	outer:
		for {
			select {
			case cmd := <-win.commands:
				cmd()
				//case <-time.After(20 * time.Millisecond):
			case <-check:
				break outer
			}
		}
	}
}
