//go:build windows

package opengl

import (
	"fmt"
	"syscall"
)

func setupLibrary() error {
	libName, err := getOpenGLSystemLibrary()
	if err != nil {
		return fmt.Errorf("getting OpenGL system library: %w", err)
	}

	lib := syscall.NewLazyDLL(libName).Handle()

	for name, ptr := range importsGl {
		if err := registerFunction(lib, name, ptr); err != nil {
			return fmt.Errorf("registering OpenGL function '%s': %w", name, err)
		}
	}

	libName, err = getGlfwSystemLibrary()
	if err != nil {
		return fmt.Errorf("getting GLUT system library: %w", err)
	}

	lib := syscall.NewLazyDLL(libName).Handle()

	for name, ptr := range importsGlfw {
		if err := registerFunction(lib, name, ptr); err != nil {
			return fmt.Errorf("registering GLUT function '%s': %w", name, err)
		}
	}
	return nil
}
