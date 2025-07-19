package m6502

import "errors"

// Common errors for M6502 emulation
var (
	ErrUnsupportedAddressingMode = errors.New("unsupported addressing mode")
	ErrInvalidParameterType      = errors.New("invalid parameter type")
	ErrMissingParameter          = errors.New("missing required parameter")
	ErrInvalidRegisterType       = errors.New("invalid register type")
)
