package dynlib

import (
	"fmt"
	"sort"

	"github.com/ebitengine/purego"
)

// LoadFunctions opens a dynamic library and registers functions from it.
func LoadFunctions(name string, imports map[string]any) (uintptr, error) {
	lib, err := open(name)
	if err != nil {
		return 0, err
	}

	if err := registerFunctions(lib, imports); err != nil {
		return 0, fmt.Errorf("registering functions: %w", err)
	}
	return lib, nil
}

func registerFunction(lib uintptr, name string, ptr any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("registering function %q: %v", name, r)
		}
	}()

	purego.RegisterLibFunc(ptr, lib, name)
	return nil
}

func registerFunctions(lib uintptr, imports map[string]any) error {
	names := make([]string, 0, len(imports))
	for name := range imports {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		if err := registerFunction(lib, name, imports[name]); err != nil {
			return err
		}
	}
	return nil
}
