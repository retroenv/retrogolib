package z80

// Shared indexed register operations for DD (IX) and FD (IY) prefix instructions.
// Each function takes a pointer to the index register (*uint16) to support both IX and IY.

// indexedLdRegNn loads a 16-bit immediate value into an index register.
func indexedLdRegNn(c *CPU, reg *uint16, _ ...any) error {
	low := c.bus.Read(c.PC + 2)
	high := c.bus.Read(c.PC + 3)
	*reg = uint16(high)<<8 | uint16(low)
	return nil
}

// indexedAddRegPair adds a 16-bit value to an index register.
func indexedAddRegPair(c *CPU, reg *uint16, value uint16, _ ...any) error {
	c.MEMPTR = *reg + 1
	*reg = c.add16(*reg, value)
	c.setXY(uint8(*reg >> 8))
	return nil
}

// indexedLdNnReg stores an index register to a memory address.
func indexedLdNnReg(c *CPU, reg uint16, _ ...any) error {
	addr := c.read16(c.PC + 2)
	c.bus.Write(addr, uint8(reg))
	c.bus.Write(addr+1, uint8(reg>>8))
	c.MEMPTR = addr + 1
	return nil
}

// indexedLdRegFromNn loads an index register from a memory address.
func indexedLdRegFromNn(c *CPU, reg *uint16, _ ...any) error {
	addr := c.read16(c.PC + 2)
	low := c.bus.Read(addr)
	high := c.bus.Read(addr + 1)
	*reg = uint16(high)<<8 | uint16(low)
	c.MEMPTR = addr + 1
	return nil
}

// indexedLdRegFromMem loads a register from indexed memory (IX/IY+d).
func indexedLdRegFromMem(c *CPU, dst *uint8, reg uint16, params ...any) error {
	addr := c.calculateIndexedAddress(reg, params...)
	*dst = c.bus.Read(addr)
	return nil
}

// indexedLdMemFromReg stores a register to indexed memory (IX/IY+d).
func indexedLdMemFromReg(c *CPU, src uint8, reg uint16, params ...any) error {
	addr := c.calculateIndexedAddress(reg, params...)
	c.bus.Write(addr, src)
	return nil
}

// indexedLdMemN loads an immediate value to indexed memory (IX/IY+d).
func indexedLdMemN(c *CPU, reg uint16, _ ...any) error {
	addr := c.calculateIndexedAddress(reg)
	value := c.bus.Read(c.PC + 3)
	c.bus.Write(addr, value)
	return nil
}

// indexedIncMem increments the value at indexed memory (IX/IY+d).
func indexedIncMem(c *CPU, reg uint16, params ...any) error {
	addr := c.calculateIndexedAddress(reg, params...)
	value := c.bus.Read(addr)
	result := value + 1
	c.bus.Write(addr, result)
	c.setSZ(result)
	c.setH((value & 0x0F) == 0x0F)
	c.setPOverflow(value == 0x7F)
	c.setN(false)
	return nil
}

// indexedDecMem decrements the value at indexed memory (IX/IY+d).
func indexedDecMem(c *CPU, reg uint16, params ...any) error {
	addr := c.calculateIndexedAddress(reg, params...)
	value := c.bus.Read(addr)
	result := value - 1
	c.bus.Write(addr, result)
	c.setSZ(result)
	c.setH((value & 0x0F) == 0x00)
	c.setPOverflow(value == 0x80)
	c.setN(true)
	return nil
}

// indexedAddA adds indexed memory value to accumulator.
func indexedAddA(c *CPU, reg uint16, params ...any) error {
	addr := c.calculateIndexedAddress(reg, params...)
	value := c.bus.Read(addr)
	result := uint16(c.A) + uint16(value)

	c.setC(result > 0xFF)
	c.setH((c.A&0x0F)+(value&0x0F) > 0x0F)
	c.setPOverflow(((c.A ^ value ^ 0x80) & (value ^ uint8(result)) & 0x80) != 0)
	c.setN(false)
	c.A = uint8(result)
	c.setSZ(c.A)
	return nil
}

// indexedAdcA adds indexed memory value with carry to accumulator.
func indexedAdcA(c *CPU, reg uint16, params ...any) error {
	addr := c.calculateIndexedAddress(reg, params...)
	value := c.bus.Read(addr)
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

// indexedSubA subtracts indexed memory value from accumulator.
func indexedSubA(c *CPU, reg uint16, params ...any) error {
	addr := c.calculateIndexedAddress(reg, params...)
	value := c.bus.Read(addr)
	result := c.A - value

	c.setC(c.A < value)
	c.setH((c.A & 0x0F) < (value & 0x0F))
	c.setPOverflow(((c.A ^ value) & (c.A ^ result) & 0x80) != 0)
	c.setN(true)
	c.A = result
	c.setSZ(c.A)
	return nil
}

// indexedSbcA subtracts indexed memory value with carry from accumulator.
func indexedSbcA(c *CPU, reg uint16, params ...any) error {
	addr := c.calculateIndexedAddress(reg, params...)
	value := c.bus.Read(addr)
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

// indexedAndA performs AND of indexed memory value with accumulator.
func indexedAndA(c *CPU, reg uint16, params ...any) error {
	addr := c.calculateIndexedAddress(reg, params...)
	value := c.bus.Read(addr)
	c.A &= value
	c.setLogicalFlags(c.A, true)
	return nil
}

// indexedXorA performs XOR of indexed memory value with accumulator.
func indexedXorA(c *CPU, reg uint16, params ...any) error {
	addr := c.calculateIndexedAddress(reg, params...)
	value := c.bus.Read(addr)
	c.A ^= value
	c.setLogicalFlags(c.A, false)
	return nil
}

// indexedOrA performs OR of indexed memory value with accumulator.
func indexedOrA(c *CPU, reg uint16, params ...any) error {
	addr := c.calculateIndexedAddress(reg, params...)
	value := c.bus.Read(addr)
	c.A |= value
	c.setLogicalFlags(c.A, false)
	return nil
}

// indexedCpA compares indexed memory value with accumulator.
func indexedCpA(c *CPU, reg uint16, params ...any) error {
	addr := c.calculateIndexedAddress(reg, params...)
	value := c.bus.Read(addr)
	c.cp(c.A, value)
	return nil
}

// indexedExSp exchanges top of stack with an index register.
func indexedExSp(c *CPU, reg *uint16) error {
	low := c.bus.Read(c.SP)
	high := c.bus.Read(c.SP + 1)

	c.bus.Write(c.SP, uint8(*reg))
	c.bus.Write(c.SP+1, uint8(*reg>>8))

	*reg = uint16(high)<<8 | uint16(low)
	c.MEMPTR = *reg
	return nil
}

// indexedCBShift performs a CB-prefixed shift/rotate on indexed memory.
func indexedCBShift(c *CPU, reg uint16, _ ...any) error {
	displacement := int8(c.bus.Read(c.PC + 2))
	opcode := c.bus.Read(c.PC + 3)
	addr := uint16(int32(reg) + int32(displacement))
	c.MEMPTR = addr
	value := c.bus.Read(addr)

	result, carry := performShiftRotateOperation(value, opcode, c.Flags.C)
	c.bus.Write(addr, result)
	setShiftRotateFlags(c, result, carry)

	// Undocumented: copy result to register if low 3 bits != 6
	if r := opcode & 0x07; r != 6 {
		c.SetRegisterValue(r, result)
	}

	return nil
}

// indexedCBBit performs a CB-prefixed BIT test on indexed memory.
func indexedCBBit(c *CPU, reg uint16, _ ...any) error {
	displacement := int8(c.bus.Read(c.PC + 2))
	opcode := c.bus.Read(c.PC + 3)
	addr := uint16(int32(reg) + int32(displacement))
	c.MEMPTR = addr
	value := c.bus.Read(addr)

	bitNum := (opcode >> 3) & 0x07
	c.bitMemptr(bitNum, value, uint8(addr>>8))
	return nil
}

// indexedCBRes performs a CB-prefixed RES on indexed memory.
func indexedCBRes(c *CPU, reg uint16, _ ...any) error {
	displacement := int8(c.bus.Read(c.PC + 2))
	opcode := c.bus.Read(c.PC + 3)
	addr := uint16(int32(reg) + int32(displacement))
	c.MEMPTR = addr
	value := c.bus.Read(addr)

	bit := (opcode >> 3) & 0x07
	result := value & ^(1 << bit)
	c.bus.Write(addr, result)

	// Undocumented: copy result to register if low 3 bits != 6
	if r := opcode & 0x07; r != 6 {
		c.SetRegisterValue(r, result)
	}

	return nil
}

// indexedCBSet performs a CB-prefixed SET on indexed memory.
func indexedCBSet(c *CPU, reg uint16, _ ...any) error {
	displacement := int8(c.bus.Read(c.PC + 2))
	opcode := c.bus.Read(c.PC + 3)
	addr := uint16(int32(reg) + int32(displacement))
	c.MEMPTR = addr
	value := c.bus.Read(addr)

	bit := (opcode >> 3) & 0x07
	result := value | (1 << bit)
	c.bus.Write(addr, result)

	// Undocumented: copy result to register if low 3 bits != 6
	if r := opcode & 0x07; r != 6 {
		c.SetRegisterValue(r, result)
	}

	return nil
}
