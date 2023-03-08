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

// TODO: complete
var openGLKeyMapping = map[glfw.Key]input.Key{
	glfw.KeyUp:        input.Up,
	glfw.KeyDown:      input.Down,
	glfw.KeyLeft:      input.Left,
	glfw.KeyRight:     input.Right,
	glfw.KeyZ:         input.A,
	glfw.KeyX:         input.B,
	glfw.KeyEnter:     input.Enter,
	glfw.KeyBackspace: input.Backspace,
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
