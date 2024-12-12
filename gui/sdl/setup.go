//go:build !windows

package sdl

import (
	"fmt"

	"github.com/ebitengine/purego"
)

func setupLibrary() error {
	libName, err := getSDLSystemLibrary()
	if err != nil {
		return fmt.Errorf("getting SDL system library: %w", err)
	}

	lib, err := purego.Dlopen(libName, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return fmt.Errorf("loading SDL system library: %w", err)
	}

	for name, ptr := range imports {
		if err := registerFunction(lib, name, ptr); err != nil {
			return err
		}
	}
	return nil
}
