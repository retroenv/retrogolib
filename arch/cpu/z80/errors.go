package z80

import "errors"

// Common errors for Z80 emulation
var (
	// CPU creation errors
	ErrNilMemory = errors.New("memory cannot be nil")

	// Parameter validation errors
	ErrUnsupportedAddressingMode = errors.New("unsupported addressing mode")
	ErrInvalidParameterType      = errors.New("invalid parameter type")
	ErrMissingParameter          = errors.New("missing required parameter")
	ErrInvalidInterruptMode      = errors.New("invalid interrupt mode (must be 0, 1, or 2)")

	// Opcode execution errors
	ErrUnsupportedOpcode   = errors.New("unsupported or unimplemented opcode")
	ErrUnsupportedEDOpcode = errors.New("unsupported ED-prefixed opcode")
)
