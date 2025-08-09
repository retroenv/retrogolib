package z80

import "errors"

// Common errors for Z80 emulation
var (
	ErrUnsupportedAddressingMode = errors.New("unsupported addressing mode")
	ErrInvalidParameterType      = errors.New("invalid parameter type")
	ErrMissingParameter          = errors.New("missing required parameter")
	ErrInvalidRegisterType       = errors.New("invalid register type")
	ErrUnsupportedOpcode         = errors.New("unsupported or unimplemented opcode")
	ErrInvalidBitNumber          = errors.New("invalid bit number (must be 0-7)")
	ErrInvalidPortAddress        = errors.New("invalid port address")
	ErrInvalidConditionCode      = errors.New("invalid condition code")
	ErrInvalidInterruptMode      = errors.New("invalid interrupt mode (must be 0, 1, or 2)")
	ErrOpcodeNotImplemented      = errors.New("opcode not implemented")
	ErrInvalidInstruction        = errors.New("invalid instruction format")
)
