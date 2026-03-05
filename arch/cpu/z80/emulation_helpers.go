package z80

import "math/bits"

// ldi executes Load and Increment operation.
func (c *CPU) ldi() {
	hl := c.hl()
	de := c.de()
	bc := c.bc()

	// Copy byte from (HL) to (DE)
	value := c.bus.Read(hl)
	c.bus.Write(de, value)

	// Increment HL and DE, decrement BC
	c.setHL(hl + 1)
	c.setDE(de + 1)
	c.setBC(bc - 1)

	// Set P/V flag based on BC
	c.setPOverflow(bc != 1) // P/V set if BC-1 != 0
	c.setH(false)
	c.setN(false)
}

// ldd executes Load and Decrement operation.
func (c *CPU) ldd() {
	hl := c.hl()
	de := c.de()
	bc := c.bc()

	// Copy byte from (HL) to (DE)
	value := c.bus.Read(hl)
	c.bus.Write(de, value)

	// Decrement HL, DE, and BC
	c.setHL(hl - 1)
	c.setDE(de - 1)
	c.setBC(bc - 1)

	// Set P/V flag based on BC
	c.setPOverflow(bc != 1) // P/V set if BC-1 != 0
	c.setH(false)
	c.setN(false)
}

// boolToUint8 converts a boolean to 1 or 0 as uint8.
func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

// calculateParity calculates the parity of a byte (true if even parity).
func calculateParity(value uint8) bool {
	count := bits.OnesCount8(value)
	return count%2 == 0
}

// performShiftRotateOperation performs shift/rotate operations based on opcode.
func performShiftRotateOperation(value, opcode, oldCarry uint8) (uint8, bool) {
	switch {
	case opcode <= 0x07: // RLC
		carry := (value & 0x80) != 0
		return (value << 1) | boolToUint8(carry), carry
	case opcode <= 0x0F: // RRC
		carry := (value & 0x01) != 0
		return (value >> 1) | (boolToUint8(carry) << 7), carry
	case opcode <= 0x17: // RL
		carry := (value & 0x80) != 0
		return (value << 1) | oldCarry, carry
	case opcode <= 0x1F: // RR
		carry := (value & 0x01) != 0
		return (value >> 1) | (oldCarry << 7), carry
	case opcode <= 0x27: // SLA
		carry := (value & 0x80) != 0
		return value << 1, carry
	case opcode <= 0x2F: // SRA
		carry := (value & 0x01) != 0
		return (value >> 1) | (value & 0x80), carry
	case opcode <= 0x37: // SLL
		carry := (value & 0x80) != 0
		return (value << 1) | 0x01, carry
	default: // SRL
		carry := (value & 0x01) != 0
		return value >> 1, carry
	}
}

// setShiftRotateFlags sets flags for shift/rotate operations.
func setShiftRotateFlags(c *CPU, result uint8, carry bool) {
	c.setSZ(result)
	c.setPOverflow(calculateParity(result))
	c.setH(false)
	c.setN(false)
	c.setC(carry)
}
