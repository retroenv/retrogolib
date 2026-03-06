package m65816

// Stack and related push/pull instructions.

// pha - Push Accumulator.
func pha(c *CPU) error {
	if c.AccWidth() == 1 {
		c.push8(uint8(c.C))
	} else {
		c.push16(c.C)
	}
	return nil
}

// phb - Push Data Bank Register.
func phb(c *CPU) error {
	c.push8(c.DB)
	return nil
}

// phd - Push Direct Page Register (always 16-bit).
// 65816-native: uses full 16-bit SP (no page-1 wrap between bytes).
func phd(c *CPU) error {
	c.push16raw(c.DP)
	c.fixEmuSP()
	return nil
}

// phk - Push Program Bank Register.
func phk(c *CPU) error {
	c.push8(c.PB)
	return nil
}

// php - Push Processor Status.
func php(c *CPU) error {
	c.push8(c.GetP())
	return nil
}

// phx - Push X Register.
func phx(c *CPU) error {
	if c.IdxWidth() == 1 {
		c.push8(uint8(c.X))
	} else {
		c.push16(c.X)
	}
	return nil
}

// phy - Push Y Register.
func phy(c *CPU) error {
	if c.IdxWidth() == 1 {
		c.push8(uint8(c.Y))
	} else {
		c.push16(c.Y)
	}
	return nil
}

// pla - Pull Accumulator.
func pla(c *CPU) error {
	if c.AccWidth() == 1 {
		v := c.pop8()
		c.C = uint16(c.B())<<8 | uint16(v)
		c.setZN8(v)
	} else {
		c.C = c.pop16()
		c.setZN16(c.C)
	}
	return nil
}

// plb - Pull Data Bank Register.
// 65816-native: uses full 16-bit SP (no page-1 wrap).
func plb(c *CPU) error {
	c.DB = c.pop8raw()
	c.fixEmuSP()
	c.setZN8(c.DB)
	return nil
}

// pld - Pull Direct Page Register.
// 65816-native: uses full 16-bit SP (no page-1 wrap between bytes).
func pld(c *CPU) error {
	c.DP = c.pop16raw()
	c.fixEmuSP()
	c.setZN16(c.DP)
	return nil
}

// plp - Pull Processor Status.
func plp(c *CPU) error {
	p := c.pop8()
	c.SetP(p)
	return nil
}

// plx - Pull X Register.
func plx(c *CPU) error {
	if c.IdxWidth() == 1 {
		v := c.pop8()
		c.X = uint16(v)
		c.setZN8(v)
	} else {
		c.X = c.pop16()
		c.setZN16(c.X)
	}
	return nil
}

// ply - Pull Y Register.
func ply(c *CPU) error {
	if c.IdxWidth() == 1 {
		v := c.pop8()
		c.Y = uint16(v)
		c.setZN8(v)
	} else {
		c.Y = c.pop16()
		c.setZN16(c.Y)
	}
	return nil
}

// pea - Push Effective Absolute Address.
// Pushes the 16-bit absolute address from the instruction stream (not the contents).
// 65816-native: uses full 16-bit SP (no page-1 wrap between bytes).
func pea(c *CPU) error {
	b1 := c.fetchByte(1)
	b2 := c.fetchByte(2)
	addr := uint16(b2)<<8 | uint16(b1)
	c.push16raw(addr)
	c.fixEmuSP()
	return nil
}

// pei - Push Effective Indirect Address.
// Reads 16-bit pointer from (DP+dp) and pushes it.
// 65816-native: uses full 16-bit SP (no page-1 wrap between bytes).
func pei(c *CPU) error {
	dp := c.fetchByte(1)
	ptr := c.dpAddr(dp)
	val := c.readMem16(ptr)
	c.push16raw(val)
	c.fixEmuSP()
	return nil
}

// per - Push Effective Relative Address.
// Pushes PC+3+signed16offset (the effective absolute address, not the contents).
// 65816-native: uses full 16-bit SP (no page-1 wrap between bytes).
func per(c *CPU) error {
	b1 := c.fetchByte(1)
	b2 := c.fetchByte(2)
	offset := int16(uint16(b2)<<8 | uint16(b1))
	eff := uint16(int32(c.PC) + 3 + int32(offset))
	c.push16raw(eff)
	c.fixEmuSP()
	return nil
}
