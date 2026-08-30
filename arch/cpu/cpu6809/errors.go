package cpu6809

import "errors"

// Common 6809 emulation errors.
var (
	ErrInvalidIndexPostbyte      = errors.New("invalid index postbyte")
	ErrInvalidOpcode             = errors.New("invalid opcode")
	ErrInvalidParameterType      = errors.New("invalid parameter type")
	ErrNilMemory                 = errors.New("memory is nil")
	ErrUnsupportedAddressingMode = errors.New("unsupported addressing mode")
)
