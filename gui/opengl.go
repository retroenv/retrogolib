//go:build !nesgo && !nogui && !noopengl

package gui

import (
	"fmt"
	"image"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/retroenv/retrogolib/input"
)

func init() {
	Setup = setupOpenGLGui
}

var openGLKeyMapping = map[glfw.Key]input.Key{
	glfw.KeyUnknown:      input.Unknown,
	glfw.KeySpace:        input.Space,
	glfw.KeyApostrophe:   input.Apostrophe,
	glfw.KeyComma:        input.Comma,
	glfw.KeyMinus:        input.Minus,
	glfw.KeyPeriod:       input.Period,
	glfw.KeySlash:        input.Slash,
	glfw.Key0:            input.Key0,
	glfw.Key1:            input.Key1,
	glfw.Key2:            input.Key2,
	glfw.Key3:            input.Key3,
	glfw.Key4:            input.Key4,
	glfw.Key5:            input.Key5,
	glfw.Key6:            input.Key6,
	glfw.Key7:            input.Key7,
	glfw.Key8:            input.Key8,
	glfw.Key9:            input.Key9,
	glfw.KeySemicolon:    input.Semicolon,
	glfw.KeyEqual:        input.Equal,
	glfw.KeyA:            input.A,
	glfw.KeyB:            input.B,
	glfw.KeyC:            input.C,
	glfw.KeyD:            input.D,
	glfw.KeyE:            input.E,
	glfw.KeyF:            input.F,
	glfw.KeyG:            input.G,
	glfw.KeyH:            input.H,
	glfw.KeyI:            input.I,
	glfw.KeyJ:            input.J,
	glfw.KeyK:            input.K,
	glfw.KeyL:            input.L,
	glfw.KeyM:            input.M,
	glfw.KeyN:            input.N,
	glfw.KeyO:            input.O,
	glfw.KeyP:            input.P,
	glfw.KeyQ:            input.Q,
	glfw.KeyR:            input.R,
	glfw.KeyS:            input.S,
	glfw.KeyT:            input.T,
	glfw.KeyU:            input.U,
	glfw.KeyV:            input.V,
	glfw.KeyW:            input.W,
	glfw.KeyX:            input.X,
	glfw.KeyY:            input.Y,
	glfw.KeyZ:            input.Z,
	glfw.KeyLeftBracket:  input.LeftBracket,
	glfw.KeyBackslash:    input.Backslash,
	glfw.KeyRightBracket: input.RightBracket,
	glfw.KeyEscape:       input.Escape,
	glfw.KeyEnter:        input.Enter,
	glfw.KeyTab:          input.Tab,
	glfw.KeyBackspace:    input.Backspace,
	glfw.KeyInsert:       input.Insert,
	glfw.KeyDelete:       input.Delete,
	glfw.KeyRight:        input.Right,
	glfw.KeyLeft:         input.Left,
	glfw.KeyDown:         input.Down,
	glfw.KeyUp:           input.Up,
	glfw.KeyPageUp:       input.PageUp,
	glfw.KeyPageDown:     input.PageDown,
	glfw.KeyHome:         input.Home,
	glfw.KeyEnd:          input.End,
	glfw.KeyCapsLock:     input.CapsLock,
	glfw.KeyScrollLock:   input.ScrollLock,
	glfw.KeyNumLock:      input.NumLock,
	glfw.KeyPrintScreen:  input.PrintScreen,
	glfw.KeyPause:        input.Pause,
	glfw.KeyF1:           input.F1,
	glfw.KeyF2:           input.F2,
	glfw.KeyF3:           input.F3,
	glfw.KeyF4:           input.F4,
	glfw.KeyF5:           input.F5,
	glfw.KeyF6:           input.F6,
	glfw.KeyF7:           input.F7,
	glfw.KeyF8:           input.F8,
	glfw.KeyF9:           input.F9,
	glfw.KeyF10:          input.F10,
	glfw.KeyF11:          input.F11,
	glfw.KeyF12:          input.F12,
	glfw.KeyF13:          input.F13,
	glfw.KeyF14:          input.F14,
	glfw.KeyF15:          input.F15,
	glfw.KeyF16:          input.F16,
	glfw.KeyF17:          input.F17,
	glfw.KeyF18:          input.F18,
	glfw.KeyF19:          input.F19,
	glfw.KeyF20:          input.F20,
	glfw.KeyF21:          input.F21,
	glfw.KeyF22:          input.F22,
	glfw.KeyF23:          input.F23,
	glfw.KeyF24:          input.F24,
	glfw.KeyF25:          input.F25,
	glfw.KeyKP0:          input.KP0,
	glfw.KeyKP1:          input.KP1,
	glfw.KeyKP2:          input.KP2,
	glfw.KeyKP3:          input.KP3,
	glfw.KeyKP4:          input.KP4,
	glfw.KeyKP5:          input.KP5,
	glfw.KeyKP6:          input.KP6,
	glfw.KeyKP7:          input.KP7,
	glfw.KeyKP8:          input.KP8,
	glfw.KeyKP9:          input.KP9,
	glfw.KeyKPDecimal:    input.KPDecimal,
	glfw.KeyKPDivide:     input.KPDivide,
	glfw.KeyKPMultiply:   input.KPMultiply,
	glfw.KeyKPSubtract:   input.KPSubtract,
	glfw.KeyKPAdd:        input.KPAdd,
	glfw.KeyKPEnter:      input.KPEnter,
	glfw.KeyKPEqual:      input.KPEqual,
	glfw.KeyLeftShift:    input.LeftShift,
	glfw.KeyLeftControl:  input.LeftControl,
	glfw.KeyLeftAlt:      input.LeftAlt,
	glfw.KeyLeftSuper:    input.LeftSuper,
	glfw.KeyRightShift:   input.RightShift,
	glfw.KeyRightControl: input.RightControl,
	glfw.KeyRightAlt:     input.RightAlt,
	glfw.KeyRightSuper:   input.RightSuper,
	glfw.KeyMenu:         input.Menu,
}

func setupOpenGLGui(backend Backend) (guiRender func() (bool, error), guiCleanup func(), err error) {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()

	dimensions := backend.Dimensions()

	window, texture, err := setupOpenGL(dimensions, backend)
	if err != nil {
		return nil, nil, err
	}

	render := func() (bool, error) {
		img := backend.Image()
		renderOpenGL(dimensions, img, window, texture)
		return !window.ShouldClose(), nil
	}

	cleanup := func() {
		gl.DeleteTextures(1, &texture)
		glfw.Terminate()
	}
	return render, cleanup, nil
}

func setupOpenGL(dimensions Dimensions, backend Backend) (*glfw.Window, uint32, error) {
	// setup GLFW
	if err := glfw.Init(); err != nil {
		return nil, 0, fmt.Errorf("initializing GLFW: %w", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	height := int(float64(dimensions.Height) * dimensions.ScaleFactor)
	width := int(float64(dimensions.Width) * dimensions.ScaleFactor)
	window, err := glfw.CreateWindow(width, height, backend.WindowTitle(), nil, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("creating GLFW window: %w", err)
	}

	window.SetKeyCallback(onGLFWKey(backend))
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	// setup OpenGL
	if err = gl.Init(); err != nil {
		return nil, 0, fmt.Errorf("initializing OpenGL: %w", err)
	}
	gl.Enable(gl.TEXTURE_2D)
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	img := backend.Image()
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(dimensions.Width), int32(dimensions.Height),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&img.Pix[0]))

	return window, texture, nil
}

func renderOpenGL(dimensions Dimensions, img *image.RGBA, window *glfw.Window, texture uint32) {
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0,
		int32(dimensions.Width), int32(dimensions.Height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&img.Pix[0]))

	// disable any filtering to avoid blurring the texture
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	// set an orthogonal projection (2D) with the size of the screen
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0.0, float64(dimensions.Width), 0.0, float64(dimensions.Height), -1.0, 1.0)
	gl.MatrixMode(gl.MODELVIEW)

	// render a single quad with the size of the screen and with the
	// contents of the emulator frame buffer
	gl.Begin(gl.QUADS)
	gl.TexCoord2d(0.0, 1.0)
	gl.Vertex2d(0.0, 0.0)
	gl.TexCoord2d(1.0, 1.0)
	gl.Vertex2d(float64(dimensions.Width), 0.0)
	gl.TexCoord2d(1.0, 0.0)
	gl.Vertex2d(float64(dimensions.Width), float64(dimensions.Height))
	gl.TexCoord2d(0.0, 0.0)
	gl.Vertex2d(0.0, float64(dimensions.Height))
	gl.End()

	window.SwapBuffers()
	glfw.PollEvents()
}

func onGLFWKey(backend Backend) func(window *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
	return func(window *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
		if action == glfw.Press && key == glfw.KeyEscape {
			window.SetShouldClose(true)
		}

		controllerKey, ok := openGLKeyMapping[key]
		if !ok {
			return
		}

		switch action {
		case glfw.Press:
			backend.KeyDown(controllerKey)

		case glfw.Release:
			backend.KeyUp(controllerKey)
		}
	}
}
