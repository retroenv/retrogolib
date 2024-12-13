//go:build !windows

package opengl

import (
	"fmt"
	"sort"

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

	names := make([]string, 0, len(importsGl))
	for name := range importsGl {
		names = append(names, name)
	}
	sort.Strings(names)

	for name := range names {
		ptr := importsGl[name]
		if err := registerFunction(lib, name, ptr); err != nil {
			return fmt.Errorf("registering OpenGL function '%s': %w", name, err)
		}
	}

	libName, err = getGlfwSystemLibrary()
	if err != nil {
		return fmt.Errorf("getting GLFW library: %w", err)
	}

	lib, err = purego.Dlopen(libName, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("loading GLFW library: %w", err)
	}

	names = make([]string, 0, len(importsGlfw))
	for name := range importsGlfw {
		names = append(names, name)
	}
	sort.Strings(names)

	for name := range names {
		ptr := importsGlfw[name]
		if err := registerFunction(lib, name, ptr); err != nil {
			return fmt.Errorf("registering GLFW function '%s': %w", name, err)
		}
	}
	return nil
}
