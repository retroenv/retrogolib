package cpu65816

import "errors"

// Common 65816 emulation errors.
var (
	ErrInvalidOpcode             = errors.New("invalid opcode")
	ErrInvalidParameterType      = errors.New("invalid parameter type")
	ErrMissingParameter          = errors.New("missing required parameter")
	ErrNilMemory                 = errors.New("memory is nil")
	ErrUnsupportedAddressingMode = errors.New("unsupported addressing mode")
)
