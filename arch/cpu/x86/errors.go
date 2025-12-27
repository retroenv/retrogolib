package x86

import "errors"

// Common x86 CPU errors.
var (
	ErrNilMemory = errors.New("memory is nil")
)
