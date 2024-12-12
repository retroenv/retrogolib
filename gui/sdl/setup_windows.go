//go:build windows

package sdl

import (
	"fmt"
	"syscall"
)

func setupLibrary() error {
	libName, err := getSDLSystemLibrary()
	if err != nil {
		return fmt.Errorf("getting SDL library: %w", err)
	}

	lib, err := loadLibrary(libName)
	if err != nil {
		return fmt.Errorf("loading SDL library: %w", err)
	}

	for name, ptr := range imports {
		if err := registerFunction(lib, name, ptr); err != nil {
			return err
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
