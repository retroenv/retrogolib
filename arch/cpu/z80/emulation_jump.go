package z80

import "fmt"

// checkCondition returns true if the condition is met based on the opcode.
func (c *CPU) checkCondition(opcode uint8) bool {
	switch opcode {
	// NZ (Not Zero) - Z flag clear
	case 0xC0, 0xC2, 0xC4:
		return c.Flags.Z == 0
	// Z (Zero) - Z flag set
	case 0xC8, 0xCA, 0xCC:
		return c.Flags.Z != 0
	// NC (Not Carry) - C flag clear
	case 0xD0, 0xD2, 0xD4:
		return c.Flags.C == 0
	// C (Carry) - C flag set
	case 0xD8, 0xDA, 0xDC:
		return c.Flags.C != 0
	// PO (Parity Odd) - P flag clear
	case 0xE0, 0xE2, 0xE4:
		return c.Flags.P == 0
	// PE (Parity Even) - P flag set
	case 0xE8, 0xEA, 0xEC:
		return c.Flags.P != 0
	// P (Plus/Positive) - S flag clear
	case 0xF0, 0xF2, 0xF4:
		return c.Flags.S == 0
	// M (Minus/Negative) - S flag set
	case 0xF8, 0xFA, 0xFC:
		return c.Flags.S != 0
	default:
		return false
	}
}

// jpAbs performs absolute jump.
func jpAbs(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	addr, ok := params[0].(Extended)
	if !ok {
		return ErrInvalidParameterType
	}

	c.PC = uint16(addr)
	c.MEMPTR = uint16(addr)
	return nil
}

// jrRel performs relative jump.
func jrRel(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	offset, ok := params[0].(Relative)
	if !ok {
		return ErrInvalidParameterType
	}

	// Calculate target address: PC after this 2-byte instruction + offset
	c.PC = uint16(int32(c.PC) + 2 + int32(offset))
	c.MEMPTR = c.PC
	return nil
}

// djnz decrements B and jumps if not zero.
func djnz(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}
	offset, ok := params[0].(Relative)
	if !ok {
		return ErrInvalidParameterType
	}

	c.B--
	if c.B != 0 {
		// Calculate target address: PC after this 2-byte instruction + offset
		c.PC = uint16(int32(c.PC) + 2 + int32(offset))
		c.MEMPTR = c.PC
		c.cycles += 5 // Extra cycles for taken branch (13 total vs 8 not taken)
	} else {
		c.PC += 2 // Advance past the 2-byte instruction
	}
	return nil
}

// jrCond performs conditional relative jump.
func jrCond(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}
	offset, ok := params[0].(Relative)
	if !ok {
		return ErrInvalidParameterType
	}

	var shouldJump bool
	opcode := c.currentOpcode
	switch opcode {
	case 0x20: // JR NZ,e - jump if zero flag is clear
		shouldJump = c.Flags.Z == 0
	case 0x28: // JR Z,e - jump if zero flag is set
		shouldJump = c.Flags.Z != 0
	case 0x30: // JR NC,e - jump if carry flag is clear
		shouldJump = c.Flags.C == 0
	case 0x38: // JR C,e - jump if carry flag is set
		shouldJump = c.Flags.C != 0
	default:
		return fmt.Errorf("unsupported jrCond opcode: 0x%02X", opcode)
	}

	if shouldJump {
		// Calculate target address: PC after this 2-byte instruction + offset
		c.PC = uint16(int32(c.PC) + 2 + int32(offset))
		c.MEMPTR = c.PC
		// Add extra cycles for taken branch (JR taken = 12 cycles, not taken = 7 cycles)
		c.cycles += 5
	} else {
		c.PC += 2 // Advance past the 2-byte instruction
	}
	return nil
}

// jpCond performs conditional jump.
func jpCond(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}
	address, ok := params[0].(Extended)
	if !ok {
		return ErrInvalidParameterType
	}

	c.MEMPTR = uint16(address)
	if c.checkCondition(c.currentOpcode) {
		c.PC = uint16(address)
	} else {
		c.PC += 3 // Advance past the 3-byte instruction
	}
	return nil
}

// callCond performs conditional call.
func callCond(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}
	address, ok := params[0].(Extended)
	if !ok {
		return ErrInvalidParameterType
	}

	c.MEMPTR = uint16(address)
	if c.checkCondition(c.currentOpcode) {
		c.push16(c.PC + 3) // Push return address (next instruction after 3-byte CALL)
		c.PC = uint16(address)
		// Add extra cycles for taken call (CALL taken = 17 cycles, not taken = 10 cycles)
		c.cycles += 7
	} else {
		c.PC += 3 // Advance past the 3-byte instruction
	}
	return nil
}

// ret returns from subroutine.
func ret(c *CPU) error {
	c.PC = c.pop16()
	c.MEMPTR = c.PC
	return nil
}

// retCond performs conditional return.
func retCond(c *CPU) error {
	if c.checkCondition(c.currentOpcode) {
		c.PC = c.pop16()
		c.MEMPTR = c.PC
		// Add extra cycles for taken return (RET taken = 11 cycles, not taken = 5 cycles)
		c.cycles += 6
	} else {
		c.PC++ // Advance past the 1-byte instruction
	}
	return nil
}

// call calls a subroutine.
func call(c *CPU, params ...any) error {
	if len(params) < 1 {
		return ErrMissingParameter
	}

	addr, ok := params[0].(Extended)
	if !ok {
		return ErrInvalidParameterType
	}

	// Push return address (PC + 3 for the 3-byte CALL instruction) to stack
	// The PC still points to the CALL instruction when this function runs
	returnAddr := c.PC + 3
	c.push16(returnAddr)

	// Jump to the called address
	c.PC = uint16(addr)
	c.MEMPTR = uint16(addr)
	return nil
}

// jpIndirect performs indirect jump.
func jpIndirect(c *CPU, _ ...any) error {
	// JP (HL) - Jump to address in HL register
	c.PC = c.hl()
	return nil
}

// rst performs restart (call to fixed address).
func rst(c *CPU, _ ...any) error {
	// RST pushes return address (next instruction) to stack and jumps to fixed address
	c.push16(c.PC + 1)

	// RST vector is encoded in bits 3-5 of the opcode: 0xC7/CF/D7/DF/E7/EF/F7/FF
	vector := uint16(c.currentOpcode & 0x38)
	c.PC = vector
	c.MEMPTR = vector
	return nil
}
