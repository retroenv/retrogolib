//go:build !windows

package opengl

import (
	"fmt"

	"github.com/ebitengine/purego"
)

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
