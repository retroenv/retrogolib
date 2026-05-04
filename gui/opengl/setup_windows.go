//go:build windows

package opengl

import (
	"fmt"

	"github.com/retroenv/retrogolib/gui/internal/dynlib"
)

func setupLibrary() error {
	libName, err := getOpenGLSystemLibrary()
	if err != nil {
		return fmt.Errorf("getting OpenGL library: %w", err)
	}

	lib, err := dynlib.Open(libName)
	if err != nil {
		return fmt.Errorf("loading OpenGL library: %w", err)
	}

	if err := dynlib.RegisterFunctions(lib, "OpenGL", importsGl); err != nil {
		return fmt.Errorf("registering OpenGL functions: %w", err)
	}

	libName, err = getGlfwSystemLibrary()
	if err != nil {
		return fmt.Errorf("getting GLFW library: %w", err)
	}

	lib, err = dynlib.Open(libName)
	if err != nil {
		return fmt.Errorf("loading GLFW library: %w", err)
	}

	if err := dynlib.RegisterFunctions(lib, "GLFW", importsGlfw); err != nil {
		return fmt.Errorf("registering GLFW functions: %w", err)
	}
	return nil
}
