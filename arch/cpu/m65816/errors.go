package m65816

import "errors"

// Common errors for m65816 emulation.
var (
	ErrInvalidOpcode             = errors.New("invalid opcode")
	ErrInvalidParameterType      = errors.New("invalid parameter type")
	ErrMissingParameter          = errors.New("missing required parameter")
	ErrNilMemory                 = errors.New("memory is nil")
	ErrUnsupportedAddressingMode = errors.New("unsupported addressing mode")
)
