package sdl

import (
	"fmt"
	"runtime"

	"github.com/ebitengine/purego"
)

var (
	// glfwInit initializes GLFW.
	glfwInit func() int

	// glfwWindowHint sets hints for the window.
	glfwWindowHint func(target, hint int32)

	// glfwWindowShouldClose checks if the window should close.
	glfwWindowShouldClose func(window uintptr) int

	// glfwTerminate terminates GLFW.
	glfwTerminate func()

	// glfwCreateWindow creates a window.
	glfwCreateWindow func(width, height int32, title string, monitor uintptr, share uintptr) uintptr

	// glfwSetKeyCallback sets the key callback.
	glfwSetKeyCallback func(window uintptr, cb uintptr)

	// glfwMakeContextCurrent makes the context current.
	glfwMakeContextCurrent func(window uintptr)

	// glfwSwapInterval sets the swap interval.
	glfwSwapInterval func(interval int32)

	// glfwSwapBuffers swaps the buffers.
	glfwSwapBuffers func(window uintptr)

	// glfwPollEvents polls the events.
	glfwPollEvents func()

	// glfwSetWindowShouldClose sets the window should close.
	glfwSetWindowShouldClose func(window uintptr, value int32)

	// glDeleteTextures deletes textures.
	glDeleteTextures func(n int32, textures *uint32)

	// glEnable enables a capability.
	glEnable func(cap int32)

	// glGenTextures generates textures.
	glGenTextures func(n int32, textures *uint32)

	// glBindTexture binds a texture.
	glBindTexture func(target, texture uint32)

	// glTexImage2D sets a texture image.
	glTexImage2D func(target, level, internalFormat, width, height, border, format, xtype int32, pixels uintptr)

	// glTexSubImage2D sets a texture sub image.
	glTexSubImage2D func(target, level, xoffset, yoffset, width, height, format, xtype int32, pixels uintptr)

	// glTexParameteri sets a texture parameter.
	glTexParameteri func(target, pname, param int32)

	// glMatrixMode sets the matrix mode.
	glMatrixMode func(mode int32)

	// glLoadIdentity loads the identity matrix.
	glLoadIdentity func()

	// glOrtho sets the orthographic projection.
	glOrtho func(left, right, bottom, top, near, far float64)

	glBegin func(mode int32)

	// glTexCoord2d sets the texture coordinates.
	glTexCoord2d func(s, t float64)

	// glVertex2d sets the vertex coordinates.
	glVertex2d func(x, y float64)

	// glEnd ends the drawing.
	glEnd func()
)

var importsGl = map[string]any{
	"glDeleteTextures": &glDeleteTextures,
	"glEnable":         &glEnable,
	"glGenTextures":    &glGenTextures,
	"glBindTexture":    &glBindTexture,
	"glTexImage2D":     &glTexImage2D,
	"glTexSubImage2D":  &glTexSubImage2D,
	"glTexParameteri":  &glTexParameteri,
	"glMatrixMode":     &glMatrixMode,
	"glLoadIdentity":   &glLoadIdentity,
	"glOrtho":          &glOrtho,
	"glBegin":          &glBegin,
	"glTexCoord2d":     &glTexCoord2d,
	"glVertex2d":       &glVertex2d,
	"glEnd":            &glEnd,
}

var importsGlfw = map[string]any{
	"glfwWindowShouldClose":    &glfwWindowShouldClose,
	"glfwInit":                 &glfwInit,
	"glfwTerminate":            &glfwTerminate,
	"glfwWindowHint":           &glfwWindowHint,
	"glfwCreateWindow":         &glfwCreateWindow,
	"glfwSetKeyCallback":       &glfwSetKeyCallback,
	"glfwMakeContextCurrent":   &glfwMakeContextCurrent,
	"glfwSwapInterval":         &glfwSwapInterval,
	"glfwSwapBuffers":          &glfwSwapBuffers,
	"glfwPollEvents":           &glfwPollEvents,
	"glfwSetWindowShouldClose": &glfwSetWindowShouldClose,
}

func setupLibrary() error {
	libName, err := getOpenGLSystemLibrary()
	if err != nil {
		return fmt.Errorf("getting OpenGL system library: %w", err)
	}

	lib, err := purego.Dlopen(libName, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("loading OpenGL system library: %w", err)
	}

	for name, ptr := range importsGl {
		if err := registerFunction(lib, name, ptr); err != nil {
			return fmt.Errorf("registering OpenGL function '%s': %w", name, err)
		}
	}

	libName, err = getGlfwSystemLibrary()
	if err != nil {
		return fmt.Errorf("getting GLUT system library: %w", err)
	}

	lib, err = purego.Dlopen(libName, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("loading GLUT system library: %w", err)
	}

	for name, ptr := range importsGlfw {
		if err := registerFunction(lib, name, ptr); err != nil {
			return fmt.Errorf("registering GLUT function '%s': %w", name, err)
		}
	}
	return nil
}

func registerFunction(lib uintptr, name string, ptr any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("registering function '%s': %v", name, r)
		}
	}()

	purego.RegisterLibFunc(ptr, lib, name)
	return nil
}

func getOpenGLSystemLibrary() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return "libGL.dylib", nil
	case "freebsd":
		return "libGL.so", nil
	case "linux":
		return "libGL.so", nil
	case "windows":
		return "opengl32.dll", nil
	default:
		return "", fmt.Errorf("GOOS=%s is not supported", runtime.GOOS)
	}
}

func getGlfwSystemLibrary() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return "libglfw.dylib", nil
	case "freebsd":
		return "libglfw.so.3", nil
	case "linux":
		return "libglfw.so.3", nil
	case "windows":
		return "glew32.dll", nil
	//case "windows":
	//	return "glu32.dll", nil
	default:
		return "", fmt.Errorf("GOOS=%s is not supported", runtime.GOOS)
	}
}
