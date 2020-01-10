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
	commands chan func()
	window   *glfw.Window
}

// NewWindow opens a new window.
func NewWindow(conf Config) *Window {
	win := &Window{
		conf:     conf,
		commands: make(chan func(), 1000),
	}
	return win
}

// Do queues a function for execution on the main thread.
// OS windows typically require that code which accesses
// windows
func (win *Window) Do(cmd func()) {
	win.commands <- cmd
}

func (win *Window) Swap() {
	win.window.SwapBuffers()
}

func (win *Window) SetSize(w, h int) {
	win.window.SetSize(w, h)
}

func (win *Window) SetPos(x, y int) {
	win.window.SetPos(x, y)
}

func (win *Window) SetTitle(title string) {
	win.window.SetTitle(title)
}

func (win *Window) Show() {
	win.window.Show()
}

func (win *Window) GetFramebufferSize() (x, y int) {
	return win.window.GetFramebufferSize()
}

// TODO leaking glfw. leak all of it? or wrap everything?
func (win *Window) SetKeyCallback(cb glfw.KeyCallback) glfw.KeyCallback {
	return win.window.SetKeyCallback(cb)
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
