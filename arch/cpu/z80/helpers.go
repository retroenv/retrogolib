package z80

import "fmt"

// Error message constants for consistent error reporting.
const (
	ErrUnimplementedEDInstruction = "unimplemented ED instruction: ED %02X"
	ErrUnimplementedDDInstruction = "unimplemented DD instruction: DD %02X"
	ErrUnimplementedFDInstruction = "unimplemented FD instruction: FD %02X"
)

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

// CreateUnimplementedError creates a formatted error for unimplemented instructions.
func CreateUnimplementedError(format string, opcodeByte uint8) error {
	return fmt.Errorf(format, opcodeByte)
}
