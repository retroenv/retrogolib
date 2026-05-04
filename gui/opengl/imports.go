package opengl

import (
	"fmt"
	"runtime"
	"sort"

	"github.com/ebitengine/purego"
)

var (
	// glfwSetErrorCallback sets the error callback.
	glfwSetErrorCallback func(cb uintptr)

	// glfwInit initializes GLFW.
	glfwInit func() int32

	// glfwWindowHint sets hints for the window.
	glfwWindowHint func(target, hint int32)

	// glfwWindowShouldClose checks if the window should close.
	glfwWindowShouldClose func(window uintptr) int32

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
	glEnable func(capability int32)

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
	"glBegin":          &glBegin,
	"glBindTexture":    &glBindTexture,
	"glDeleteTextures": &glDeleteTextures,
	"glEnable":         &glEnable,
	"glEnd":            &glEnd,
	"glGenTextures":    &glGenTextures,
	"glLoadIdentity":   &glLoadIdentity,
	"glMatrixMode":     &glMatrixMode,
	"glOrtho":          &glOrtho,
	"glTexCoord2d":     &glTexCoord2d,
	"glTexImage2D":     &glTexImage2D,
	"glTexParameteri":  &glTexParameteri,
	"glTexSubImage2D":  &glTexSubImage2D,
	"glVertex2d":       &glVertex2d,
}

var importsGlfw = map[string]any{
	"glfwCreateWindow":         &glfwCreateWindow,
	"glfwInit":                 &glfwInit,
	"glfwMakeContextCurrent":   &glfwMakeContextCurrent,
	"glfwPollEvents":           &glfwPollEvents,
	"glfwSetErrorCallback":     &glfwSetErrorCallback,
	"glfwSetKeyCallback":       &glfwSetKeyCallback,
	"glfwSetWindowShouldClose": &glfwSetWindowShouldClose,
	"glfwSwapBuffers":          &glfwSwapBuffers,
	"glfwSwapInterval":         &glfwSwapInterval,
	"glfwTerminate":            &glfwTerminate,
	"glfwWindowHint":           &glfwWindowHint,
	"glfwWindowShouldClose":    &glfwWindowShouldClose,
}

func registerFunction(lib uintptr, name string, ptr any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("registering function %q: %v", name, r)
		}
	}()

	purego.RegisterLibFunc(ptr, lib, name)
	return nil
}

func registerFunctions(lib uintptr, group string, imports map[string]any) error {
	names := make([]string, 0, len(imports))
	for name := range imports {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		if err := registerFunction(lib, name, imports[name]); err != nil {
			return fmt.Errorf("registering %s function %q: %w", group, name, err)
		}
	}
	return nil
}

func getOpenGLSystemLibrary() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return "/System/Library/Frameworks/OpenGL.framework/OpenGL", nil
	case "freebsd":
		return "libGL.so.1", nil
	case "linux":
		return "libGL.so.1", nil
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
		return "glfw3.dll", nil
	default:
		return "", fmt.Errorf("GOOS=%s is not supported", runtime.GOOS)
	}
}
