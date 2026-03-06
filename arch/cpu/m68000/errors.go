package m68000

import "errors"

// Common errors for 68000 emulation.
var (
	ErrNilBus             = errors.New("bus cannot be nil")
	ErrUnsupportedOpcode  = errors.New("unsupported or unimplemented opcode")
	ErrAddressError       = errors.New("address error: word/long access at odd address")
	ErrBusError           = errors.New("bus error")
	ErrPrivilegeViolation = errors.New("privilege violation")
	ErrIllegalInstruction = errors.New("illegal instruction")
	ErrDivideByZero       = errors.New("divide by zero")
	ErrCHK                = errors.New("CHK exception")
	ErrTRAPV              = errors.New("TRAPV exception")
	ErrUnimplemented      = errors.New("unimplemented instruction")
	ErrInvalidAddressMode = errors.New("invalid addressing mode")
	ErrInvalidOperandSize = errors.New("invalid operand size")
)
