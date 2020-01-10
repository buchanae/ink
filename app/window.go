// +build !sendonly

package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var initGLFW, initGL sync.Once

func (app *App) runWindow() error {

	var err error
	initGLFW.Do(func() {
		err = glfw.Init()
	})
	if err != nil {
		return fmt.Errorf("initializing glfw: %v", err)
	}
	// TODO this is not the right place/way to call this,
	//      because there could be multiple windows on a thread.
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	conf := app.conf.Window
	window, err := glfw.CreateWindow(
		conf.Width, conf.Height, conf.Title,
		nil, nil,
	)
	if err != nil {
		return fmt.Errorf("creating window: %v", err)
	}
	defer window.Destroy()

	window.SetPos(conf.X, conf.Y)
	window.SetKeyCallback(app.keyCallback)
	window.MakeContextCurrent()

	initGL.Do(func() {
		err = gl.Init()
	})
	if err != nil {
		return fmt.Errorf("initializing OpenGL: %v", err)
	}

	app.win = window
	check := time.Tick(20 * time.Millisecond)

	for !window.ShouldClose() {

		glfw.PollEvents()

	outer:
		for {
			select {
			case cmd := <-app.commands:
				cmd()
			case <-check:
				break outer
			}
		}
	}
	return nil
}
