package win

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

var initGLFW, initGL sync.Once

type Config struct {
	Title         string
	X, Y          int
	Width, Height int
	Visible       bool
}

// Window holds a handle to an OS window.
type Window struct {
	conf     Config
	events   chan Event
	commands chan func()
	window   *glfw.Window
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
		win.window.SwapBuffers()
	}
}

func (win *Window) SetSize(w, h int) {
	win.commands <- func() {
		win.window.SetSize(w, h)
	}
}

func (win *Window) SetPos(x, y int) {
	win.commands <- func() {
		win.window.SetPos(x, y)
	}
}

func (win *Window) SetTitle(title string) {
	win.commands <- func() {
		win.window.SetTitle(title)
	}
}

func (win *Window) Show() {
	win.commands <- func() {
		win.window.Show()
	}
}

func (win *Window) Run() {

	var err error
	initGLFW.Do(func() {
		err = glfw.Init()
	})
	if err != nil {
		// TODO not sure what to do with errors
		log.Printf("error: initializing glfw: %v", err)
		return
	}
	// TODO this is not the right place/way to call this,
	//      because there could be multiple windows on a thread.
	defer glfw.Terminate()

	if !win.conf.Visible {
		glfw.WindowHint(glfw.Visible, glfw.False)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(
		win.conf.Width, win.conf.Height,
		win.conf.Title,
		nil, nil,
	)
	if err != nil {
		// TODO not sure what to do with errors
		log.Printf("error: creating window: %v", err)
		return
	}
	defer window.Destroy()

	window.SetPos(win.conf.X, win.conf.Y)

	window.MakeContextCurrent()

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

	for !window.ShouldClose() {

		glfw.PollEvents()
		/* TODO redo key events

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
		*/

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
