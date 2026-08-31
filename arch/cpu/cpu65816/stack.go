package cpu65816

// push8 pushes a byte onto the stack and decrements SP.
// In emulation mode, SP wraps within page 1 ($0100-$01FF) after each byte.
func (c *CPU) push8(value uint8) {
	c.memory.Write(bank24(0, c.SP), value)
	c.SP--
	if c.E {
		c.SP = 0x0100 | (c.SP & 0x00FF)
	}
}

// push8raw pushes a byte without page-1 wrap, for 65816-native stack instructions
// that use the full 16-bit SP even in emulation mode.
func (c *CPU) push8raw(value uint8) {
	c.memory.Write(bank24(0, c.SP), value)
	c.SP--
}

// push16 pushes a 16-bit word onto the stack (high byte first).
func (c *CPU) push16(value uint16) {
	c.push8(uint8(value >> 8))
	c.push8(uint8(value))
}

// push16raw pushes a 16-bit word without per-byte page-1 wrap.
func (c *CPU) push16raw(value uint16) {
	c.push8raw(uint8(value >> 8))
	c.push8raw(uint8(value))
}

// pop16raw pops a 16-bit word without per-byte page-1 wrap.
func (c *CPU) pop16raw() uint16 {
	lo := uint16(c.pop8raw())
	hi := uint16(c.pop8raw())
	return hi<<8 | lo
}

// fixEmuSP normalizes SP to page 1 after a 65816-native stack instruction.
func (c *CPU) fixEmuSP() {
	if c.E {
		c.SP = 0x0100 | (c.SP & 0x00FF)
	}
}

// pop8 pops a byte from the stack and increments SP.
// In emulation mode, SP wraps within page 1 ($0100-$01FF) after each byte.
func (c *CPU) pop8() uint8 {
	c.SP++
	if c.E {
		c.SP = 0x0100 | (c.SP & 0x00FF)
	}
	return c.memory.Read(bank24(0, c.SP))
}

// pop8raw pops a byte without page-1 wrap, for 65816-native stack instructions
// that use the full 16-bit SP even in emulation mode.
func (c *CPU) pop8raw() uint8 {
	c.SP++
	return c.memory.Read(bank24(0, c.SP))
}

// pop16 pops a 16-bit word from the stack (low byte first).
func (c *CPU) pop16() uint16 {
	lo := uint16(c.pop8())
	hi := uint16(c.pop8())
	return hi<<8 | lo
}
