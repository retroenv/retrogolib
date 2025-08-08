package z80

import "fmt"

// Error message constants for consistent error reporting.
const (
	ErrUnimplementedCBInstruction = "unimplemented CB instruction: CB %02X"
	ErrUnimplementedEDInstruction = "unimplemented ED instruction: ED %02X"
	ErrUnimplementedDDInstruction = "unimplemented DD instruction: DD %02X"
	ErrUnimplementedFDInstruction = "unimplemented FD instruction: FD %02X"
)

// Opcode creation helpers to reduce boilerplate and improve performance.

// NewImpliedOpcode creates an opcode with implied addressing.
func NewImpliedOpcode(instruction *Instruction, timing, size byte) Opcode {
	return Opcode{
		Instruction: instruction,
		Addressing:  ImpliedAddressing,
		Timing:      timing,
		Size:        size,
	}
}

// NewRegisterOpcode creates an opcode with register addressing.
func NewRegisterOpcode(instruction *Instruction, timing, size byte) Opcode {
	return Opcode{
		Instruction: instruction,
		Addressing:  RegisterAddressing,
		Timing:      timing,
		Size:        size,
	}
}

// Timing calculation helpers for complex instruction patterns.

// GetCBTiming calculates timing for CB-prefixed instructions based on operation and register.
func GetCBTiming(opcodeByte, reg uint8) byte {
	switch {
	case opcodeByte <= 0x3F && reg == 6: // Rotate/shift (HL)
		return 15
	case opcodeByte <= 0x7F && reg == 6: // BIT n,(HL)
		return 12
	case opcodeByte >= 0x80 && reg == 6: // RES/SET n,(HL)
		return 15
	default:
		return 8
	}
}

// GetCBInstruction returns the appropriate CB instruction based on opcode byte.
func GetCBInstruction(opcodeByte uint8) *Instruction {
	switch {
	case opcodeByte <= 0x07:
		return CBRlc
	case opcodeByte <= 0x0F:
		return CBRrc
	case opcodeByte <= 0x17:
		return CBRl
	case opcodeByte <= 0x1F:
		return CBRr
	case opcodeByte <= 0x27:
		return CBSla
	case opcodeByte <= 0x2F:
		return CBSra
	case opcodeByte <= 0x37:
		return CBSll
	case opcodeByte <= 0x3F:
		return CBSrl
	case opcodeByte <= 0x7F:
		return CBBit
	case opcodeByte <= 0xBF:
		return CBRes
	default:
		return CBSet
	}
}

// CreateUnimplementedError creates a formatted error for unimplemented instructions.
func CreateUnimplementedError(format string, opcodeByte uint8) error {
	return fmt.Errorf(format, opcodeByte)
}
