package sm83

// applyCBOperation reads the CB sub-opcode from PC+1, determines the target register,
// and applies the given operation. For (HL) (reg==6), reads/writes through memory.
func (c *CPU) applyCBOperation(operation func(uint8) uint8) {
	subOpcode := c.memory.Read(c.PC + 1)
	reg := subOpcode & 0x07

	if reg == 6 { // Operation on (HL)
		addr := c.hl()
		value := c.memory.Read(addr)
		result := operation(value)
		c.memory.Write(addr, result)
	} else { // Operation on register
		value := c.GetRegisterValue(reg)
		result := operation(value)
		c.SetRegisterValue(reg, result)
	}
}

// cbRlc performs RLC r (rotate left circular, CB 00-07).
// Bit 7 goes to carry and bit 0. Z, N=0, H=0, C=old bit 7.
func cbRlc(c *CPU, _ ...any) error {
	c.applyCBOperation(func(value uint8) uint8 {
		carry := (value & 0x80) >> 7
		result := (value << 1) | carry

		c.setZ(result)
		c.setN(false)
		c.setH(false)
		c.setC(carry != 0)
		return result
	})
	return nil
}

// cbRrc performs RRC r (rotate right circular, CB 08-0F).
// Bit 0 goes to carry and bit 7. Z, N=0, H=0, C=old bit 0.
func cbRrc(c *CPU, _ ...any) error {
	c.applyCBOperation(func(value uint8) uint8 {
		carry := value & 0x01
		result := (value >> 1) | (carry << 7)

		c.setZ(result)
		c.setN(false)
		c.setH(false)
		c.setC(carry != 0)
		return result
	})
	return nil
}

// cbRl performs RL r (rotate left through carry, CB 10-17).
// Old carry goes to bit 0, bit 7 goes to carry. Z, N=0, H=0, C=old bit 7.
func cbRl(c *CPU, _ ...any) error {
	c.applyCBOperation(func(value uint8) uint8 {
		newCarry := (value & 0x80) >> 7
		result := (value << 1) | c.Flags.C

		c.setZ(result)
		c.setN(false)
		c.setH(false)
		c.setC(newCarry != 0)
		return result
	})
	return nil
}

// cbRr performs RR r (rotate right through carry, CB 18-1F).
// Old carry goes to bit 7, bit 0 goes to carry. Z, N=0, H=0, C=old bit 0.
func cbRr(c *CPU, _ ...any) error {
	c.applyCBOperation(func(value uint8) uint8 {
		newCarry := value & 0x01
		result := (value >> 1) | (c.Flags.C << 7)

		c.setZ(result)
		c.setN(false)
		c.setH(false)
		c.setC(newCarry != 0)
		return result
	})
	return nil
}

// cbSla performs SLA r (shift left arithmetic, CB 20-27).
// Bit 7 goes to carry, 0 goes to bit 0. Z, N=0, H=0, C=old bit 7.
func cbSla(c *CPU, _ ...any) error {
	c.applyCBOperation(func(value uint8) uint8 {
		carry := (value & 0x80) >> 7
		result := value << 1

		c.setZ(result)
		c.setN(false)
		c.setH(false)
		c.setC(carry != 0)
		return result
	})
	return nil
}

// cbSra performs SRA r (shift right arithmetic, CB 28-2F).
// Bit 0 goes to carry, bit 7 preserved. Z, N=0, H=0, C=old bit 0.
func cbSra(c *CPU, _ ...any) error {
	c.applyCBOperation(func(value uint8) uint8 {
		carry := value & 0x01
		result := (value >> 1) | (value & 0x80)

		c.setZ(result)
		c.setN(false)
		c.setH(false)
		c.setC(carry != 0)
		return result
	})
	return nil
}

// cbSwap performs SWAP r (SM83-unique, CB 30-37).
// Swaps upper and lower nibbles. Z, N=0, H=0, C=0.
func cbSwap(c *CPU, _ ...any) error {
	c.applyCBOperation(func(value uint8) uint8 {
		result := (value>>4)&0x0F | (value&0x0F)<<4

		c.setZ(result)
		c.setN(false)
		c.setH(false)
		c.setC(false)
		return result
	})
	return nil
}

// cbSrl performs SRL r (shift right logical, CB 38-3F).
// Bit 0 goes to carry, 0 goes to bit 7. Z, N=0, H=0, C=old bit 0.
func cbSrl(c *CPU, _ ...any) error {
	c.applyCBOperation(func(value uint8) uint8 {
		carry := value & 0x01
		result := value >> 1

		c.setZ(result)
		c.setN(false)
		c.setH(false)
		c.setC(carry != 0)
		return result
	})
	return nil
}

// cbBit performs BIT b,r (test bit, CB 40-7F).
// Z = complement of tested bit. N=0, H=1, C unchanged.
func cbBit(c *CPU, _ ...any) error {
	subOpcode := c.memory.Read(c.PC + 1)
	bitNum := (subOpcode >> 3) & 0x07
	reg := subOpcode & 0x07

	var value uint8
	if reg == 6 { // BIT b,(HL)
		value = c.memory.Read(c.hl())
	} else {
		value = c.GetRegisterValue(reg)
	}

	bit := (value >> bitNum) & 1
	setFlag(&c.Flags.Z, bit == 0)
	c.setN(false)
	c.setH(true)
	return nil
}

// cbRes performs RES b,r (reset bit, CB 80-BF).
// Clears the specified bit. No flag changes.
func cbRes(c *CPU, _ ...any) error {
	subOpcode := c.memory.Read(c.PC + 1)
	bitNum := (subOpcode >> 3) & 0x07

	c.applyCBOperation(func(value uint8) uint8 {
		return value & ^(1 << bitNum)
	})
	return nil
}

// cbSet performs SET b,r (set bit, CB C0-FF).
// Sets the specified bit. No flag changes.
func cbSet(c *CPU, _ ...any) error {
	subOpcode := c.memory.Read(c.PC + 1)
	bitNum := (subOpcode >> 3) & 0x07

	c.applyCBOperation(func(value uint8) uint8 {
		return value | (1 << bitNum)
	})
	return nil
}
