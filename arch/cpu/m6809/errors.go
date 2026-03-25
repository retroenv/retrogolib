package m6809

import "errors"

// Common errors for m6809 emulation.
var (
	ErrInvalidIndexPostbyte      = errors.New("invalid index postbyte")
	ErrInvalidOpcode             = errors.New("invalid opcode")
	ErrInvalidParameterType      = errors.New("invalid parameter type")
	ErrMissingParameter          = errors.New("missing required parameter")
	ErrNilMemory                 = errors.New("memory is nil")
	ErrUnsupportedAddressingMode = errors.New("unsupported addressing mode")
)
