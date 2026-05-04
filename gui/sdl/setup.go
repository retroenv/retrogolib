//go:build !windows

package sdl

import (
	"fmt"

	"github.com/retroenv/retrogolib/gui/internal/dynlib"
)

func setupLibrary() error {
	libName, err := getSDLSystemLibrary()
	if err != nil {
		return fmt.Errorf("getting SDL library: %w", err)
	}

	lib, err := dynlib.Open(libName)
	if err != nil {
		return fmt.Errorf("loading SDL library: %w", err)
	}

	if err := dynlib.RegisterFunctions(lib, "SDL", imports); err != nil {
		return fmt.Errorf("registering SDL functions: %w", err)
	}
	return nil
}
