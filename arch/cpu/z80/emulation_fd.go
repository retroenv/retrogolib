package z80

// FD prefix instruction implementations - IY register operations

func fdLdIYnn(c *CPU, _ ...any) error {
	// Read 16-bit immediate value from memory at PC+2 and PC+3
	low := c.memory.Read(c.PC + 2)
	high := c.memory.Read(c.PC + 3)
	c.IY = uint16(high)<<8 | uint16(low)
	return nil
}

// IY register operations
func fdIncIY(c *CPU) error { c.IY++; return nil }
func fdDecIY(c *CPU) error { c.IY--; return nil }

func fdAddIYBc(c *CPU, _ ...any) error { c.IY += uint16(c.B)<<8 | uint16(c.C); return nil }
func fdAddIYDe(c *CPU, _ ...any) error { c.IY += uint16(c.D)<<8 | uint16(c.E); return nil }
func fdAddIYIY(c *CPU, _ ...any) error { c.IY += c.IY; return nil }
func fdAddIYSp(c *CPU, _ ...any) error { c.IY += c.SP; return nil }

// fdLdNnIY implements FD 22: LD (nn),IY.
func fdLdNnIY(c *CPU, params ...any) error {
	addr := uint16(params[1].(uint8))<<8 | uint16(params[0].(uint8))
	c.memory.Write(addr, uint8(c.IY))
	c.memory.Write(addr+1, uint8(c.IY>>8))
	return nil
}

// fdLdIYNn implements FD 2A: LD IY,(nn).
func fdLdIYNn(c *CPU, params ...any) error {
	addr := uint16(params[1].(uint8))<<8 | uint16(params[0].(uint8))
	low := c.memory.Read(addr)
	high := c.memory.Read(addr + 1)
	c.IY = uint16(high)<<8 | uint16(low)
	return nil
}

// IY indexed load operations - Load register from (IY+d)

// fdLdBIYd implements FD 46: LD B,(IY+d).
func fdLdBIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	c.B = c.memory.Read(addr)
	return nil
}

// fdLdCIYd implements FD 4E: LD C,(IY+d).
func fdLdCIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	c.C = c.memory.Read(addr)
	return nil
}

// fdLdDIYd implements FD 56: LD D,(IY+d).
func fdLdDIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	c.D = c.memory.Read(addr)
	return nil
}

// fdLdEIYd implements FD 5E: LD E,(IY+d).
func fdLdEIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	c.E = c.memory.Read(addr)
	return nil
}

// fdLdHIYd implements FD 66: LD H,(IY+d).
func fdLdHIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	c.H = c.memory.Read(addr)
	return nil
}

// fdLdLIYd implements FD 6E: LD L,(IY+d).
func fdLdLIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	c.L = c.memory.Read(addr)
	return nil
}

// fdLdAIYd implements FD 7E: LD A,(IY+d).
func fdLdAIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	c.A = c.memory.Read(addr)
	return nil
}

// IY indexed store operations - Store register to (IY+d)

// fdLdIYdB implements FD 70: LD (IY+d),B.
func fdLdIYdB(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	c.memory.Write(addr, c.B)
	return nil
}

// fdLdIYdC implements FD 71: LD (IY+d),C.
func fdLdIYdC(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	c.memory.Write(addr, c.C)
	return nil
}

// fdLdIYdD implements FD 72: LD (IY+d),D.
func fdLdIYdD(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	c.memory.Write(addr, c.D)
	return nil
}

// fdLdIYdE implements FD 73: LD (IY+d),E.
func fdLdIYdE(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	c.memory.Write(addr, c.E)
	return nil
}

// fdLdIYdH implements FD 74: LD (IY+d),H.
func fdLdIYdH(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	c.memory.Write(addr, c.H)
	return nil
}

// fdLdIYdL implements FD 75: LD (IY+d),L.
func fdLdIYdL(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	c.memory.Write(addr, c.L)
	return nil
}

// fdLdIYdA implements FD 77: LD (IY+d),A.
func fdLdIYdA(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	c.memory.Write(addr, c.A)
	return nil
}

// fdLdIYdN implements FD 36: LD (IY+d),n.
func fdLdIYdN(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	value := params[1].(uint8)
	c.memory.Write(addr, value)
	return nil
}

// fdIncIYd implements FD 34: INC (IY+d).
func fdIncIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	value := c.memory.Read(addr)
	result := value + 1
	c.memory.Write(addr, result)
	c.setSZ(result)
	c.setH((value & 0x0F) == 0x0F)
	c.setPOverflow(value == 0x7F)
	c.setN(false)
	return nil
}

// fdDecIYd implements FD 35: DEC (IY+d).
func fdDecIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	value := c.memory.Read(addr)
	result := value - 1
	c.memory.Write(addr, result)
	c.setSZ(result)
	c.setH((value & 0x0F) == 0x00)
	c.setPOverflow(value == 0x80)
	c.setN(true)
	return nil
}

// IY arithmetic operations

// fdAddAIYd implements FD 86: ADD A,(IY+d).
func fdAddAIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	value := c.memory.Read(addr)
	result := c.A + value

	c.setSZ(result)
	c.setC(uint16(c.A)+uint16(value) > 0xFF)
	c.setH((c.A&0x0F)+(value&0x0F) > 0x0F)
	c.setPOverflow(((c.A ^ value ^ 0x80) & (result ^ c.A) & 0x80) != 0)
	c.setN(false)
	c.A = result
	return nil
}

// fdAdcAIYd implements FD 8E: ADC A,(IY+d).
func fdAdcAIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	value := c.memory.Read(addr)
	carry := uint8(0)
	if c.Flags.C != 0 {
		carry = 1
	}
	result := c.A + value + carry

	c.setSZ(result)
	c.setC(uint16(c.A)+uint16(value)+uint16(carry) > 0xFF)
	c.setH((c.A&0x0F)+(value&0x0F)+carry > 0x0F)
	c.setPOverflow(((c.A ^ value ^ 0x80) & (result ^ c.A) & 0x80) != 0)
	c.setN(false)
	c.A = result
	return nil
}

// fdSubAIYd implements FD 96: SUB (IY+d).
func fdSubAIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	value := c.memory.Read(addr)
	result := c.A - value

	c.setSZ(result)
	c.setC(c.A < value)
	c.setH((c.A & 0x0F) < (value & 0x0F))
	c.setPOverflow(((c.A ^ value) & (c.A ^ result) & 0x80) != 0)
	c.setN(true)
	c.A = result
	return nil
}

// fdSbcAIYd implements FD 9E: SBC A,(IY+d).
func fdSbcAIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	value := c.memory.Read(addr)
	carry := uint8(0)
	if c.Flags.C != 0 {
		carry = 1
	}
	result := c.A - value - carry

	c.setSZ(result)
	c.setC(uint16(c.A) < uint16(value)+uint16(carry))
	c.setH((c.A & 0x0F) < (value&0x0F)+carry)
	c.setPOverflow(((c.A ^ value) & (c.A ^ result) & 0x80) != 0)
	c.setN(true)
	c.A = result
	return nil
}

// fdAndAIYd implements FD A6: AND (IY+d).
func fdAndAIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	value := c.memory.Read(addr)
	c.A &= value
	c.setLogicalFlags(c.A, true)
	return nil
}

// fdXorAIYd implements FD AE: XOR (IY+d).
func fdXorAIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	value := c.memory.Read(addr)
	c.A ^= value
	c.setLogicalFlags(c.A, false)
	return nil
}

// fdOrAIYd implements FD B6: OR (IY+d).
func fdOrAIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	value := c.memory.Read(addr)
	c.A |= value
	c.setLogicalFlags(c.A, false)
	return nil
}

// fdCpAIYd implements FD BE: CP (IY+d).
func fdCpAIYd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IY, params...)
	value := c.memory.Read(addr)
	result := c.A - value

	c.setSZ(result)
	c.setC(c.A < value)
	c.setH((c.A & 0x0F) < (value & 0x0F))
	c.setPOverflow(((c.A ^ value) & (c.A ^ result) & 0x80) != 0)
	c.setN(true)
	return nil
}

// IY stack and jump operations
func fdJpIY(c *CPU) error { c.PC = c.IY; return nil }

// fdExSpIY implements FD E3: EX (SP),IY.
func fdExSpIY(c *CPU) error {
	// Exchange IY with the word at the top of the stack
	low := c.memory.Read(c.SP)
	high := c.memory.Read(c.SP + 1)

	c.memory.Write(c.SP, uint8(c.IY))
	c.memory.Write(c.SP+1, uint8(c.IY>>8))

	c.IY = uint16(high)<<8 | uint16(low)
	return nil
}
func fdPushIY(c *CPU) error { c.push16(c.IY); return nil }
func fdPopIY(c *CPU) error  { c.IY = c.pop16(); return nil }

// FDCB operations - bit operations on (IY+d)

func fdcbShift(c *CPU, params ...any) error {
	displacement := int8(params[0].(uint8))
	opcode := params[1].(uint8)
	addr := uint16(int32(c.IY) + int32(displacement))
	value := c.memory.Read(addr)

	result, carry := performShiftRotateOperation(value, opcode, c.Flags.C)
	c.memory.Write(addr, result)
	setShiftRotateFlags(c, result, carry)

	return nil
}

func fdcbBit(c *CPU, params ...any) error {
	displacement := int8(params[0].(uint8))
	opcode := params[1].(uint8)
	addr := uint16(int32(c.IY) + int32(displacement))
	value := c.memory.Read(addr)

	bit := (opcode >> 3) & 0x07
	result := value & (1 << bit)

	c.setZ(result)
	c.setH(true)
	c.setN(false)
	if bit == 7 {
		c.setS(result)
	} else {
		c.setS(0)
	}
	return nil
}

func fdcbRes(c *CPU, params ...any) error {
	displacement := int8(params[0].(uint8))
	opcode := params[1].(uint8)
	addr := uint16(int32(c.IY) + int32(displacement))
	value := c.memory.Read(addr)

	bit := (opcode >> 3) & 0x07
	result := value & ^(1 << bit)
	c.memory.Write(addr, result)
	return nil
}

func fdcbSet(c *CPU, params ...any) error {
	displacement := int8(params[0].(uint8))
	opcode := params[1].(uint8)
	addr := uint16(int32(c.IY) + int32(displacement))
	value := c.memory.Read(addr)

	bit := (opcode >> 3) & 0x07
	result := value | (1 << bit)
	c.memory.Write(addr, result)
	return nil
}
