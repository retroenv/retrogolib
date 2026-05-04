//go:build !windows

package dynlib

import (
	"fmt"

	"github.com/ebitengine/purego"
)

// Open loads a dynamic library with process-global symbols.
func Open(name string) (uintptr, error) {
	lib, err := purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return 0, fmt.Errorf("opening dynamic library %q: %w", name, err)
	}
	return lib, nil
}
