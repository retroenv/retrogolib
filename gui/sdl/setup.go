//go:build !windows

package sdl

import (
	"fmt"
	"sort"

	"github.com/ebitengine/purego"
)

func setupLibrary() error {
	libName, err := getSDLSystemLibrary()
	if err != nil {
		return fmt.Errorf("getting SDL library: %w", err)
	}

	lib, err := purego.Dlopen(libName, purego.RTLD_NOW|purego.RTLD_GLOBAL)
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
