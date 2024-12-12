//go:build windows

package sdl

import (
	"fmt"
	"syscall"
)

func setupLibrary() error {
	libName, err := getSDLSystemLibrary()
	if err != nil {
		return fmt.Errorf("getting SDL system library: %w", err)
	}

	lib := syscall.NewLazyDLL(libName).Handle()

	for name, ptr := range imports {
		if err := registerFunction(lib, name, ptr); err != nil {
			return err
		}
	}
	return nil
}
