package dynlib

import (
	"fmt"
	"sort"

	"github.com/ebitengine/purego"
)

// RegisterFunction registers a single dynamic library function.
func RegisterFunction(lib uintptr, name string, ptr any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("registering function %q: %v", name, r)
		}
	}()

	purego.RegisterLibFunc(ptr, lib, name)
	return nil
}

// RegisterFunctions registers dynamic library functions in deterministic order.
func RegisterFunctions(lib uintptr, group string, imports map[string]any) error {
	names := make([]string, 0, len(imports))
	for name := range imports {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		if err := RegisterFunction(lib, name, imports[name]); err != nil {
			return fmt.Errorf("registering %s function %q: %w", group, name, err)
		}
	}
	return nil
}
