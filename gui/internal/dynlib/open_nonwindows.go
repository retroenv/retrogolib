//go:build !windows

package dynlib

import (
	"fmt"

	"github.com/ebitengine/purego"
)

func open(name string) (uintptr, error) {
	lib, err := purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return 0, fmt.Errorf("opening dynamic library %q: %w", name, err)
	}
	return lib, nil
}
