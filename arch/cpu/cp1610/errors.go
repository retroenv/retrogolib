package cp1610

import "errors"

var (
	// ErrInvalidOpcode reports an unsupported opcode form.
	ErrInvalidOpcode = errors.New("invalid opcode")
	// ErrNilMemory reports a missing memory bus.
	ErrNilMemory = errors.New("memory is nil")
)
