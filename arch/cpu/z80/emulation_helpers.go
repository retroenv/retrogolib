package z80

import "math/bits"

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

// inc16 increments a 16-bit value.
func (c *CPU) inc16(value uint16) uint16 {
	return value + 1
}

// dec16 decrements a 16-bit value.
func (c *CPU) dec16(value uint16) uint16 {
	return value - 1
}

// addHL adds a 16-bit value to HL register pair.
func (c *CPU) addHL(hl, value uint16) uint16 {
	result32 := uint32(hl) + uint32(value)
	result := uint16(result32)

	// Set flags for 16-bit addition
	c.setC(result32 > 0xFFFF)                   // Carry if result > 65535
	c.setH((hl&0x0FFF)+(value&0x0FFF) > 0x0FFF) // Half carry on bit 11
	c.setN(false)                               // Clear N flag for addition
	// Note: Z and S flags are not affected by ADD HL

	return result
}

// ldi executes Load and Increment operation.
func (c *CPU) ldi() {
	hl := c.hl()
	de := c.de()
	bc := c.bc()

	// Copy byte from (HL) to (DE)
	value := c.memory.Read(hl)
	c.memory.Write(de, value)

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
	value := c.memory.Read(hl)
	c.memory.Write(de, value)

	// Decrement HL, DE, and BC
	c.setHL(hl - 1)
	c.setDE(de - 1)
	c.setBC(bc - 1)

	// Set P/V flag based on BC
	c.setPOverflow(bc != 1) // P/V set if BC-1 != 0
	c.setH(false)
	c.setN(false)
}

// ldReg8 loads a value to a register.
func (c *CPU) ldReg8(dst *uint8, src uint8) {
	*dst = src
}

// ldImm8 loads immediate value to register.
func (c *CPU) ldImm8(dst *uint8, value uint8) {
	*dst = value
}

// ldMemToReg8 loads memory value to register.
func (c *CPU) ldMemToReg8(dst *uint8, addr uint16) {
	*dst = c.memory.Read(addr)
}

// ldRegToMem8 loads register value to memory.
func (c *CPU) ldRegToMem8(addr uint16, src uint8) {
	c.memory.Write(addr, src)
}

// jp executes unconditional jump.
func (c *CPU) jp(addr uint16) {
	c.PC = addr
}

// jr executes relative jump.
func (c *CPU) jr(offset int8) {
	c.PC = uint16(int32(c.PC) + int32(offset))
}

// jpZ executes conditional jump if Z flag is set.
func (c *CPU) jpZ(addr uint16) {
	if c.Flags.Z == 1 {
		c.PC = addr
	}
}

// djnz executes Decrement and Jump if Not Zero.
func (c *CPU) djnz(offset int8) {
	c.B--
	if c.B != 0 {
		c.PC = uint16(int32(c.PC) + int32(offset))
	}
}

// exx exchanges BC, DE, HL with shadow registers.
func (c *CPU) exx() {
	c.B, c.AltB = c.AltB, c.B
	c.C, c.AltC = c.AltC, c.C
	c.D, c.AltD = c.AltD, c.D
	c.E, c.AltE = c.AltE, c.E
	c.H, c.AltH = c.AltH, c.H
	c.L, c.AltL = c.AltL, c.L
}

// exAF exchanges AF with shadow AF.
func (c *CPU) exAF() {
	c.A, c.AltA = c.AltA, c.A
	c.Flags, c.AltFlags = c.AltFlags, c.Flags
}

// exDEHL exchanges DE and HL register pairs.
func (c *CPU) exDEHL() {
	d, e := c.D, c.E
	c.D, c.E = c.H, c.L
	c.H, c.L = d, e
}

// getFlags returns current flags as uint8.
func (c *CPU) getFlags() uint8 {
	return c.GetFlags()
}

// setFlagsFromUint8 sets flags struct from uint8 value.
func (c *CPU) setFlagsFromUint8(flags *Flags, value uint8) {
	flags.C = value & 0x01
	flags.N = (value >> 1) & 0x01
	flags.P = (value >> 2) & 0x01
	flags.X = (value >> 3) & 0x01
	flags.H = (value >> 4) & 0x01
	flags.Y = (value >> 5) & 0x01
	flags.Z = (value >> 6) & 0x01
	flags.S = (value >> 7) & 0x01
}

// getFlagsAsUint8 gets flags struct as uint8 value.
func (c *CPU) getFlagsAsUint8(flags Flags) uint8 {
	return flags.C | (flags.N << 1) | (flags.P << 2) | (flags.X << 3) |
		(flags.H << 4) | (flags.Y << 5) | (flags.Z << 6) | (flags.S << 7)
}
