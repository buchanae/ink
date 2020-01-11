// +build !sendonly

package app

import "github.com/go-gl/glfw/v3.3/glfw"

type KeyCallback func(KeyEvent)

// TODO leaking gflw
type KeyEvent struct {
	Win      *glfw.Window
	Key      glfw.Key
	Scancode int
	Action   glfw.Action
	Mod      glfw.ModifierKey
}

func (app *App) defaultKeyCallback(ev KeyEvent) {
	if ev.Key == glfw.KeyX && ev.Action == glfw.Press {
		app.snapshotAndWrite()
	}
}

func (app *App) processKeyEvents() {
	keycb := []KeyCallback{app.defaultKeyCallback}

	for {
		select {
		case cb := <-app.addKeycb:
			keycb = append(keycb, cb)
		case ev := <-app.keyEvents:
			for _, cb := range keycb {
				cb(ev)
			}
		}
	}
}

func (app *App) AddKeyCallback(cb KeyCallback) {
	app.addKeycb <- cb
}

// keyCallback pipes events from glfw main thread
// to ink.App non-main thread.
func (app *App) keyCallback(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// keyCallback gets called by GLFW, so pretty sure
	// this function always runs on the main thread.
	//
	// Be careful not to call out to user code from this function.
	app.keyEvents <- KeyEvent{win, key, scancode, action, mods}
}
