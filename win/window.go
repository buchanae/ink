package win

import (
	"log"
	"runtime"
	"sync"

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
}

// NewWindow opens a new window.
func NewWindow(conf Config) *Window {
	win := &Window{
		conf:     conf,
		events:   make(chan Event),
		commands: make(chan func(), 1),
	}
	return win
}

// Events returns a stream of window events,
// such as quit, keyboard, mouse, etc.
func (win *Window) Events() <-chan Event {
	return win.events
}

const (
	codePopCommand int32 = iota
	codeSwap
)

// Do queues a function for execution on the main thread.
// OS windows typically require that code which accesses
// windows
func (win *Window) Do(cmd func()) {
	/*
		SDL has its own event queue mechansim, and it's tricky
		to coordinate that system with Go's scheduler, while
		also keeping all window/opengl code on a single OS thread.

		This is the only way I've found to work:
		1. Push "cmd" onto a work queue (buffered channel)
		2. Push an event onto the SDL queue
		3. In Window.Run, which runs only on the OpenGL thread,
		   pop the UserEvent from the SDL queue
		4. In Window.Run, pop "cmd" from the work queue and run it.
	*/
	win.commands <- cmd
	_, err := sdl.PushEvent(&sdl.UserEvent{
		Type:      sdl.USEREVENT,
		Timestamp: sdl.GetTicks(),
		Code:      codePopCommand,
	})
	if err != nil {
		// TODO I don't know what to do with these errors
		log.Print(err)
	}
}

func (win *Window) Swap() {
	_, err := sdl.PushEvent(&sdl.UserEvent{
		Type:      sdl.USEREVENT,
		Timestamp: sdl.GetTicks(),
		Code:      codeSwap,
	})
	if err != nil {
		// TODO I don't know what to do with these errors
		log.Print(err)
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

	for {
		ev := sdl.WaitEvent()
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

		case *sdl.UserEvent:
			switch z.Code {
			case codePopCommand:
				// sdl.UserEvent is used to make goroutines
				// and SDL2 event queues play nice together.
				// See Window.Do for details.
				cmd := <-win.commands
				cmd()
			case codeSwap:
				window.GLSwap()
			}
		}
	}
}
