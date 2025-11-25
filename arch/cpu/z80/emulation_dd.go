package z80

// DD prefix instruction implementations - IX register operations

// ddLdIXnn implements DD 21: LD IX,nn.
func ddLdIXnn(c *CPU, _ ...any) error {
	// Read 16-bit immediate value from memory at PC+2 and PC+3
	low := c.memory.Read(c.PC + 2)
	high := c.memory.Read(c.PC + 3)
	c.IX = uint16(high)<<8 | uint16(low)
	return nil
}

// IX register operations
func ddIncIX(c *CPU) error { c.IX++; return nil }
func ddDecIX(c *CPU) error { c.IX--; return nil }

func ddAddIXBc(c *CPU, _ ...any) error { c.IX += uint16(c.B)<<8 | uint16(c.C); return nil }
func ddAddIXDe(c *CPU, _ ...any) error { c.IX += uint16(c.D)<<8 | uint16(c.E); return nil }
func ddAddIXIX(c *CPU, _ ...any) error { c.IX += c.IX; return nil }
func ddAddIXSp(c *CPU, _ ...any) error { c.IX += c.SP; return nil }

// ddLdNnIX implements DD 22: LD (nn),IX.
func ddLdNnIX(c *CPU, params ...any) error {
	addr := uint16(params[1].(uint8))<<8 | uint16(params[0].(uint8))
	c.memory.Write(addr, uint8(c.IX))
	c.memory.Write(addr+1, uint8(c.IX>>8))
	return nil
}

// ddLdIXNn implements DD 2A: LD IX,(nn).
func ddLdIXNn(c *CPU, params ...any) error {
	addr := uint16(params[1].(uint8))<<8 | uint16(params[0].(uint8))
	low := c.memory.Read(addr)
	high := c.memory.Read(addr + 1)
	c.IX = uint16(high)<<8 | uint16(low)
	return nil
}

// IX indexed load operations - Load register from (IX+d)

// ddLdBIXd implements DD 46: LD B,(IX+d).
func ddLdBIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	c.B = c.memory.Read(addr)
	return nil
}

// ddLdCIXd implements DD 4E: LD C,(IX+d).
func ddLdCIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	c.C = c.memory.Read(addr)
	return nil
}

// ddLdDIXd implements DD 56: LD D,(IX+d).
func ddLdDIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	c.D = c.memory.Read(addr)
	return nil
}

// ddLdEIXd implements DD 5E: LD E,(IX+d).
func ddLdEIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	c.E = c.memory.Read(addr)
	return nil
}

// ddLdHIXd implements DD 66: LD H,(IX+d).
func ddLdHIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	c.H = c.memory.Read(addr)
	return nil
}

// ddLdLIXd implements DD 6E: LD L,(IX+d).
func ddLdLIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	c.L = c.memory.Read(addr)
	return nil
}

// ddLdAIXd implements DD 7E: LD A,(IX+d).
func ddLdAIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	c.A = c.memory.Read(addr)
	return nil
}

// IX indexed store operations - Store register to (IX+d)

// ddLdIXdB implements DD 70: LD (IX+d),B.
func ddLdIXdB(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	c.memory.Write(addr, c.B)
	return nil
}

// ddLdIXdC implements DD 71: LD (IX+d),C.
func ddLdIXdC(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	c.memory.Write(addr, c.C)
	return nil
}

// ddLdIXdD implements DD 72: LD (IX+d),D.
func ddLdIXdD(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	c.memory.Write(addr, c.D)
	return nil
}

// ddLdIXdE implements DD 73: LD (IX+d),E.
func ddLdIXdE(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	c.memory.Write(addr, c.E)
	return nil
}

// ddLdIXdH implements DD 74: LD (IX+d),H.
func ddLdIXdH(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	c.memory.Write(addr, c.H)
	return nil
}

// ddLdIXdL implements DD 75: LD (IX+d),L.
func ddLdIXdL(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	c.memory.Write(addr, c.L)
	return nil
}

// ddLdIXdA implements DD 77: LD (IX+d),A.
func ddLdIXdA(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	c.memory.Write(addr, c.A)
	return nil
}

// ddLdIXdN implements DD 36: LD (IX+d),n.
func ddLdIXdN(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	value := params[1].(uint8)
	c.memory.Write(addr, value)
	return nil
}

// ddIncIXd implements DD 34: INC (IX+d).
func ddIncIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	value := c.memory.Read(addr)
	result := value + 1
	c.memory.Write(addr, result)
	c.setSZ(result)
	c.setH((value & 0x0F) == 0x0F)
	c.setPOverflow(value == 0x7F)
	c.setN(false)
	return nil
}

// ddDecIXd implements DD 35: DEC (IX+d).
func ddDecIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	value := c.memory.Read(addr)
	result := value - 1
	c.memory.Write(addr, result)
	c.setSZ(result)
	c.setH((value & 0x0F) == 0x00)
	c.setPOverflow(value == 0x80)
	c.setN(true)
	return nil
}

// IX arithmetic operations with accumulator

// ddAddAIXd implements DD 86: ADD A,(IX+d).
func ddAddAIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	value := c.memory.Read(addr)
	result := uint16(c.A) + uint16(value)

	c.setC(result > 0xFF)
	c.setH((c.A&0x0F)+(value&0x0F) > 0x0F)
	c.setPOverflow(((c.A ^ value ^ 0x80) & (value ^ uint8(result)) & 0x80) != 0)
	c.setN(false)
	c.A = uint8(result)
	c.setSZ(c.A)
	return nil
}

// ddAdcAIXd implements DD 8E: ADC A,(IX+d).
func ddAdcAIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	value := c.memory.Read(addr)
	carry := uint16(0)
	if c.Flags.C != 0 {
		carry = 1
	}
	result := uint16(c.A) + uint16(value) + carry

	c.setC(result > 0xFF)
	c.setH((c.A&0x0F)+(value&0x0F)+uint8(carry) > 0x0F)
	c.setPOverflow(((c.A ^ value ^ 0x80) & (value ^ uint8(result)) & 0x80) != 0)
	c.setN(false)
	c.A = uint8(result)
	c.setSZ(c.A)
	return nil
}

// ddSubAIXd implements DD 96: SUB (IX+d).
func ddSubAIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	value := c.memory.Read(addr)
	result := c.A - value

	c.setC(c.A < value)
	c.setH((c.A & 0x0F) < (value & 0x0F))
	c.setPOverflow(((c.A ^ value) & (c.A ^ result) & 0x80) != 0)
	c.setN(true)
	c.A = result
	c.setSZ(c.A)
	return nil
}

// ddSbcAIXd implements DD 9E: SBC A,(IX+d).
func ddSbcAIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	value := c.memory.Read(addr)
	carry := uint8(0)
	if c.Flags.C != 0 {
		carry = 1
	}
	result := c.A - value - carry

	c.setC(uint16(c.A) < uint16(value)+uint16(carry))
	c.setH((c.A & 0x0F) < (value&0x0F)+carry)
	c.setPOverflow(((c.A ^ value) & (c.A ^ result) & 0x80) != 0)
	c.setN(true)
	c.A = result
	c.setSZ(c.A)
	return nil
}

// ddAndAIXd implements DD A6: AND (IX+d).
func ddAndAIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	value := c.memory.Read(addr)
	c.A &= value
	c.setLogicalFlags(c.A, true)
	return nil
}

// ddXorAIXd implements DD AE: XOR (IX+d).
func ddXorAIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	value := c.memory.Read(addr)
	c.A ^= value
	c.setLogicalFlags(c.A, false)
	return nil
}

// ddOrAIXd implements DD B6: OR (IX+d).
func ddOrAIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	value := c.memory.Read(addr)
	c.A |= value
	c.setLogicalFlags(c.A, false)
	return nil
}

// ddCpAIXd implements DD BE: CP (IX+d).
func ddCpAIXd(c *CPU, params ...any) error {
	addr := c.calculateIndexedAddress(c.IX, params...)
	value := c.memory.Read(addr)
	result := c.A - value

	c.setSZ(result)
	c.setC(c.A < value)
	c.setH((c.A & 0x0F) < (value & 0x0F))
	c.setPOverflow(((c.A ^ value) & (c.A ^ result) & 0x80) != 0)
	c.setN(true)
	return nil
}

// IX stack and jump operations
func ddJpIX(c *CPU) error { c.PC = c.IX; return nil }

// ddExSpIX implements DD E3: EX (SP),IX.
func ddExSpIX(c *CPU) error {
	// Exchange IX with the word at the top of the stack
	low := c.memory.Read(c.SP)
	high := c.memory.Read(c.SP + 1)

	c.memory.Write(c.SP, uint8(c.IX))
	c.memory.Write(c.SP+1, uint8(c.IX>>8))

	c.IX = uint16(high)<<8 | uint16(low)
	return nil
}
func ddPushIX(c *CPU) error { c.push16(c.IX); return nil }
func ddPopIX(c *CPU) error  { c.IX = c.pop16(); return nil }

// DDCB operations - bit operations on (IX+d)

func ddcbShift(c *CPU, params ...any) error {
	displacement := int8(params[0].(uint8))
	opcode := params[1].(uint8)
	addr := uint16(int32(c.IX) + int32(displacement))
	value := c.memory.Read(addr)

	result, carry := performShiftRotateOperation(value, opcode, c.Flags.C)
	c.memory.Write(addr, result)
	setShiftRotateFlags(c, result, carry)

	return nil
}

func ddcbBit(c *CPU, params ...any) error {
	displacement := int8(params[0].(uint8))
	opcode := params[1].(uint8)
	addr := uint16(int32(c.IX) + int32(displacement))
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

func ddcbRes(c *CPU, params ...any) error {
	displacement := int8(params[0].(uint8))
	opcode := params[1].(uint8)
	addr := uint16(int32(c.IX) + int32(displacement))
	value := c.memory.Read(addr)

	bit := (opcode >> 3) & 0x07
	result := value & ^(1 << bit)
	c.memory.Write(addr, result)
	return nil
}

func ddcbSet(c *CPU, params ...any) error {
	displacement := int8(params[0].(uint8))
	opcode := params[1].(uint8)
	addr := uint16(int32(c.IX) + int32(displacement))
	value := c.memory.Read(addr)

	bit := (opcode >> 3) & 0x07
	result := value | (1 << bit)
	c.memory.Write(addr, result)
	return nil
}
