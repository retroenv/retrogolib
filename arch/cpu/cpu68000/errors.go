package cpu68000

import "errors"

// Common errors for 68000 emulation.
var (
	ErrNilBus             = errors.New("bus cannot be nil")
	ErrAddressError       = errors.New("address error: word/long access at odd address")
	ErrBusError           = errors.New("bus error")
	ErrInvalidAddressMode = errors.New("invalid addressing mode")
	ErrInvalidOperandSize = errors.New("invalid operand size")
)
