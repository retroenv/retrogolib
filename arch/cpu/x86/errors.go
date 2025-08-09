package x86

import "errors"

// Common x86 CPU errors.
var (
	ErrNilMemory              = errors.New("memory is nil")
	ErrInvalidInstruction     = errors.New("invalid instruction")
	ErrInvalidAddressingMode  = errors.New("invalid addressing mode")
	ErrInvalidRegister        = errors.New("invalid register")
	ErrInvalidOpcode          = errors.New("invalid opcode")
	ErrInvalidOperand         = errors.New("invalid operand")
	ErrDivisionByZero         = errors.New("division by zero")
	ErrInvalidSegment         = errors.New("invalid segment")
	ErrStackOverflow          = errors.New("stack overflow")
	ErrStackUnderflow         = errors.New("stack underflow")
	ErrGeneralProtectionFault = errors.New("general protection fault")
	ErrInvalidInterruptVector = errors.New("invalid interrupt vector")
)
