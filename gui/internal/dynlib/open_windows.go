//go:build windows

package dynlib

import (
	"fmt"
	"syscall"
)

// Open loads a dynamic library.
func Open(name string) (handle uintptr, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("loading library %q: %v", name, r)
		}
	}()

	handle = syscall.NewLazyDLL(name).Handle()
	return handle, err
}
