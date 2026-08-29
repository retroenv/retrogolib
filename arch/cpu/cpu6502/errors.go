package cpu6502

import "errors"

// Common 6502 emulation errors.
var (
	ErrInvalidParameterType      = errors.New("invalid parameter type")
	ErrInvalidRegisterType       = errors.New("invalid register type")
	ErrMissingParameter          = errors.New("missing required parameter")
	ErrUnknownOpcode             = errors.New("unknown opcode")
	ErrUnsupportedAddressingMode = errors.New("unsupported addressing mode")
)
