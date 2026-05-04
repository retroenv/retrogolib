// Package opengl provides an OpenGL GUI renderer.
package opengl

import (
	"errors"
	"fmt"
	"image"
	"runtime"

	"github.com/retroenv/retrogolib/gui"
	"github.com/retroenv/retrogolib/gui/internal/framebuffer"
)

// Setup initializes the OpenGL library and returns a render and cleanup function.
func Setup(backend gui.Backend) (guiRender func() (bool, error), guiCleanup func(), err error) {
	dimensions := backend.Dimensions()
	if err := framebuffer.ValidateDimensions(dimensions); err != nil {
		return nil, nil, fmt.Errorf("validating dimensions: %w", err)
	}

	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	defer func() {
		if err != nil {
			runtime.UnlockOSThread()
		}
	}()

	window, texture, err := setupOpenGL(dimensions, backend)
	if err != nil {
		return nil, nil, err
	}

	render := func() (bool, error) {
		img := backend.Image()
		if err := renderOpenGL(dimensions, img, window, texture); err != nil {
			return false, err
		}
		return glfwWindowShouldClose(window) == GLFW_FALSE, nil
	}

	cleanup := func() {
		glDeleteTextures(1, &texture)
		glfwTerminate()
	}
	return render, cleanup, nil
}

func setupOpenGL(dimensions gui.Dimensions, backend gui.Backend) (uintptr, uint32, error) {
	if err := framebuffer.ValidateDimensions(dimensions); err != nil {
		return uintptr(0), 0, fmt.Errorf("validating dimensions: %w", err)
	}

	if err := setupLibrary(); err != nil {
		return uintptr(0), 0, fmt.Errorf("setting up OpenGL library: %w", err)
	}

	setErrorCallback()

	resetLastError()
	if ret := glfwInit(); ret != GLFW_TRUE {
		return uintptr(0), 0, fmt.Errorf("initializing GLFW: %w", getLastError())
	}

	glfwWindowHint(GLFW_RESIZABLE, GLFW_FALSE)
	glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 2)
	glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 1)

	height := int32(float64(dimensions.Height) * dimensions.ScaleFactor)
	width := int32(float64(dimensions.Width) * dimensions.ScaleFactor)
	resetLastError()
	window := glfwCreateWindow(width, height, backend.WindowTitle(), uintptr(0), uintptr(0))
	if window == 0 {
		glfwTerminate()
		return uintptr(0), 0, fmt.Errorf("creating GLFW window: %w", getLastError())
	}

	setupKeyCallback(window, backend)
	glfwMakeContextCurrent(window)
	glfwSwapInterval(1)

	// setup OpenGL
	glEnable(GL_TEXTURE_2D)
	var texture uint32
	glGenTextures(1, &texture)
	if texture == 0 {
		glfwTerminate()
		return uintptr(0), 0, errors.New("generating OpenGL texture")
	}
	glBindTexture(GL_TEXTURE_2D, texture)
	// Disable filtering once to keep emulator pixels crisp.
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST)
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST)

	img := backend.Image()
	pixels, err := framebuffer.RGBAPointer(dimensions, img)
	if err != nil {
		glDeleteTextures(1, &texture)
		glfwTerminate()
		return uintptr(0), 0, fmt.Errorf("getting initial image pixels: %w", err)
	}
	glTexImage2D(GL_TEXTURE_2D, 0, GL_RGBA, int32(dimensions.Width), int32(dimensions.Height),
		0, GL_RGBA, GL_UNSIGNED_BYTE, pixels)

	return window, texture, nil
}

func renderOpenGL(dimensions gui.Dimensions, img *image.RGBA, window uintptr, texture uint32) error {
	pixels, err := framebuffer.RGBAPointer(dimensions, img)
	if err != nil {
		return fmt.Errorf("getting image pixels: %w", err)
	}

	glBindTexture(GL_TEXTURE_2D, texture)
	glTexSubImage2D(GL_TEXTURE_2D, 0, 0, 0, int32(dimensions.Width),
		int32(dimensions.Height), GL_RGBA, GL_UNSIGNED_BYTE, pixels)

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
	return nil
}
