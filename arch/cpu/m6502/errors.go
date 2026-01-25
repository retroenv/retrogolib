package m6502

import "errors"

// Common errors for M6502 emulation
var (
	ErrInvalidParameterType      = errors.New("invalid parameter type")
	ErrInvalidRegisterType       = errors.New("invalid register type")
	ErrMissingParameter          = errors.New("missing required parameter")
	ErrUnknownOpcode             = errors.New("unknown opcode")
	ErrUnsupportedAddressingMode = errors.New("unsupported addressing mode")
)
