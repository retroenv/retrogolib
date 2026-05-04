//go:build windows

package opengl

import (
	"fmt"
	"syscall"
)

func setupLibrary() error {
	libName, err := getOpenGLSystemLibrary()
	if err != nil {
		return fmt.Errorf("getting OpenGL library: %w", err)
	}

	lib, err := loadLibrary(libName)
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

	lib, err = loadLibrary(libName)
	if err != nil {
		return fmt.Errorf("loading GLFW library: %w", err)
	}

	if err := registerFunctions(lib, "GLFW", importsGlfw); err != nil {
		return err
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
