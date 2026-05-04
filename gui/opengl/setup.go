//go:build !windows

package opengl

import (
	"fmt"

	"github.com/ebitengine/purego"
)

func setupLibrary() error {
	libName, err := getOpenGLSystemLibrary()
	if err != nil {
		return fmt.Errorf("getting OpenGL library: %w", err)
	}

	lib, err := purego.Dlopen(libName, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("loading OpenGL library: %w", err)
	}

	if err := registerFunctions(lib, "OpenGL", importsGl); err != nil {
		return err
	}

	libName, err = getGlfwSystemLibrary()
	if err != nil {
		return fmt.Errorf("getting GLFW library: %w", err)
	}

	lib, err = purego.Dlopen(libName, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("loading GLFW library: %w", err)
	}

	if err := registerFunctions(lib, "GLFW", importsGlfw); err != nil {
		return err
	}
	return nil
}
