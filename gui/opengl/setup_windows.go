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

	lib, err := loadLibrary(libName)
	if err != nil {
		return fmt.Errorf("loading OpenGL library: %w", err)
	}

	for name, ptr := range importsGl {
		if err := registerFunction(lib, name, ptr); err != nil {
			return fmt.Errorf("registering OpenGL function '%s': %w", name, err)
		}
	}

	libName, err = getGlfwSystemLibrary()
	if err != nil {
		return fmt.Errorf("getting GLFW system library: %w", err)
	}

	lib, err = loadLibrary(libName)
	if err != nil {
		return fmt.Errorf("loading GLFW library: %w", err)
	}

	for name, ptr := range importsGlfw {
		if err := registerFunction(lib, name, ptr); err != nil {
			return fmt.Errorf("registering GLFW function '%s': %w", name, err)
		}
	}
	return nil
}

func loadLibrary(libName string) (handle uintptr, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("loading library '%s': %v", libName, r)
		}
	}()

	handle = syscall.NewLazyDLL(libName).Handle()
	return handle, err
}
