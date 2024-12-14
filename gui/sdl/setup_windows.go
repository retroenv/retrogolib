//go:build windows

package sdl

import (
	"fmt"
	"sort"
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

	names := make([]string, 0, len(imports))
	for name := range imports {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		ptr := imports[name]
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
