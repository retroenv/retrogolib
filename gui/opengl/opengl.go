// Package opengl provides an OpenGL GUI renderer.
package opengl

import (
	"fmt"
	"image"
	"runtime"
	"unsafe"

	"github.com/retroenv/retrogolib/gui"
)

func Setup(backend gui.Backend) (guiRender func() (bool, error), guiCleanup func(), err error) {
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
		//return !window.ShouldClose(), nil
		return false, nil
	}

	cleanup := func() {
		glDeleteTextures(1, &texture)
		glfwTerminate()
	}
	return render, cleanup, nil
}

func setupOpenGL(dimensions gui.Dimensions, backend gui.Backend) (uintptr, uint32, error) {
	if err := setupLibrary(); err != nil {
		return uintptr(0), 0, fmt.Errorf("setting up OpenGL library: %w", err)
	}

	setErrorCallback()

	if ret := glfwInit(); ret != GLFW_TRUE {
		return uintptr(0), 0, fmt.Errorf("initializing GLFW: %w", getLastError())
	}

	glfwWindowHint(GLFW_RESIZABLE, GLFW_FALSE)
	glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 2)
	glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 1)

	height := int32(float64(dimensions.Height) * dimensions.ScaleFactor)
	width := int32(float64(dimensions.Width) * dimensions.ScaleFactor)
	window := glfwCreateWindow(width, height, backend.WindowTitle(), uintptr(0), uintptr(0))
	if window == 0 {
		return uintptr(0), 0, fmt.Errorf("creating GLFW window: %w", getLastError())
	}

	keyCallback = onGLFWKey(backend)
	glfwSetKeyCallback(window, uintptr(unsafe.Pointer(&keyCallback)))
	glfwMakeContextCurrent(window)
	glfwSwapInterval(1)

	// setup OpenGL
	glEnable(TEXTURE_2D)
	var texture uint32
	glGenTextures(1, &texture)
	glBindTexture(TEXTURE_2D, texture)
	img := backend.Image()
	glTexImage2D(TEXTURE_2D, 0, RGBA, int32(dimensions.Width), int32(dimensions.Height),
		0, RGBA, UNSIGNED_BYTE, uintptr(unsafe.Pointer(&img.Pix[0])))

	return window, texture, nil
}

func renderOpenGL(dimensions gui.Dimensions, img *image.RGBA, window uintptr, texture uint32) {
	glBindTexture(TEXTURE_2D, texture)
	glTexSubImage2D(TEXTURE_2D, 0, 0, 0, int32(dimensions.Width),
		int32(dimensions.Height), RGBA, UNSIGNED_BYTE, uintptr(unsafe.Pointer(&img.Pix[0])))

	// disable any filtering to avoid blurring the texture
	glTexParameteri(TEXTURE_2D, TEXTURE_MAG_FILTER, GL_NEAREST)
	glTexParameteri(TEXTURE_2D, TEXTURE_MIN_FILTER, GL_NEAREST)

	// set an orthogonal projection (2D) with the size of the screen
	glMatrixMode(GL_PROJECTION)
	glLoadIdentity()
	glOrtho(0.0, float64(dimensions.Width), 0.0, float64(dimensions.Height), -1.0, 1.0)
	glMatrixMode(GL_MODELVIEW)

	// render a single quad with the size of the screen and with the
	// contents of the emulator frame buffer
	glBegin(GL_QUADS)
	glTexCoord2d(0.0, 1.0)
	glVertex2d(0.0, 0.0)
	glTexCoord2d(1.0, 1.0)
	glVertex2d(float64(dimensions.Width), 0.0)
	glTexCoord2d(1.0, 0.0)
	glVertex2d(float64(dimensions.Width), float64(dimensions.Height))
	glTexCoord2d(0.0, 0.0)
	glVertex2d(0.0, float64(dimensions.Height))
	glEnd()

	glfwSwapBuffers(window)
	glfwPollEvents()
}
