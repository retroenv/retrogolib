package sm83

import "errors"

// Common errors for SM83 emulation
var (
	// CPU creation errors
	ErrNilMemory = errors.New("memory cannot be nil")

	// Parameter validation errors
	ErrUnsupportedAddressingMode = errors.New("unsupported addressing mode")
	ErrInvalidParameterType      = errors.New("invalid parameter type")
	ErrMissingParameter          = errors.New("missing required parameter")

	// Opcode execution errors
	ErrUnsupportedOpcode = errors.New("unsupported or unimplemented opcode")
	ErrIllegalOpcode     = errors.New("illegal opcode")
)
