package sm83

import "errors"

// Common errors for SM83 emulation.
var (
	ErrIllegalOpcode             = errors.New("illegal opcode")
	ErrNilMemory                 = errors.New("memory cannot be nil")
	ErrUnsupportedAddressingMode = errors.New("unsupported addressing mode")
	ErrUnsupportedOpcode         = errors.New("unsupported or unimplemented opcode")
)
