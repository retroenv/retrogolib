package m68000

// Shift and rotate instructions: ASL, ASR, LSL, LSR, ROL, ROR, ROXL, ROXR.

func (c *CPU) execASL(d DecodedOpcode) error {
	if d.Extra&0x40 != 0 {
		// Memory shift (count=1, size=word).
		return c.shiftMemory(d, func(val uint32) uint32 {
			msb := val & 0x8000
			result := (val << 1) & 0xFFFF
			setFlag(&c.Flags.C, msb != 0)
			c.Flags.X = c.Flags.C
			setFlag(&c.Flags.V, (result&0x8000) != msb)
			return result
		})
	}

	count := c.shiftCount(d)
	val := c.getRegD(d.DstReg, d.Size)
	msb := msbMask(d.Size)

	result := val
	c.Flags.V = 0
	c.Flags.C = 0

	for range count {
		oldMsb := result & msb
		result = (result << 1) & maskValue(0xFFFFFFFF, d.Size)
		setFlag(&c.Flags.C, oldMsb != 0)
		c.Flags.X = c.Flags.C
		if (result & msb) != oldMsb {
			c.Flags.V = 1
		}
	}

	c.setRegD(d.DstReg, result, d.Size)
	c.setFlagN(result, d.Size)
	c.setFlagZ(result, d.Size)

	if count == 0 {
		c.Flags.C = 0
	}

	return nil
}

func (c *CPU) execASR(d DecodedOpcode) error {
	if d.Extra&0x40 != 0 {
		return c.shiftMemory(d, func(val uint32) uint32 {
			lsb := val & 1
			msb := val & 0x8000
			result := (val >> 1) | msb
			setFlag(&c.Flags.C, lsb != 0)
			c.Flags.X = c.Flags.C
			c.Flags.V = 0
			return result
		})
	}

	count := c.shiftCount(d)
	val := c.getRegD(d.DstReg, d.Size)
	msb := msbMask(d.Size)

	result := val
	c.Flags.V = 0

	for range count {
		lsb := result & 1
		signBit := result & msb
		result = (result >> 1) | signBit
		setFlag(&c.Flags.C, lsb != 0)
		c.Flags.X = c.Flags.C
	}

	result = maskValue(result, d.Size)
	c.setRegD(d.DstReg, result, d.Size)
	c.setFlagN(result, d.Size)
	c.setFlagZ(result, d.Size)

	if count == 0 {
		c.Flags.C = 0
	}

	return nil
}

func (c *CPU) execLSL(d DecodedOpcode) error {
	if d.Extra&0x40 != 0 {
		return c.shiftMemory(d, func(val uint32) uint32 {
			msb := val & 0x8000
			result := (val << 1) & 0xFFFF
			setFlag(&c.Flags.C, msb != 0)
			c.Flags.X = c.Flags.C
			c.Flags.V = 0
			return result
		})
	}

	count := c.shiftCount(d)
	val := c.getRegD(d.DstReg, d.Size)

	result := val
	c.Flags.V = 0

	for range count {
		msb := result & msbMask(d.Size)
		result = (result << 1) & maskValue(0xFFFFFFFF, d.Size)
		setFlag(&c.Flags.C, msb != 0)
		c.Flags.X = c.Flags.C
	}

	c.setRegD(d.DstReg, result, d.Size)
	c.setFlagN(result, d.Size)
	c.setFlagZ(result, d.Size)

	if count == 0 {
		c.Flags.C = 0
	}

	return nil
}

func (c *CPU) execLSR(d DecodedOpcode) error {
	if d.Extra&0x40 != 0 {
		return c.shiftMemory(d, func(val uint32) uint32 {
			lsb := val & 1
			result := val >> 1
			setFlag(&c.Flags.C, lsb != 0)
			c.Flags.X = c.Flags.C
			c.Flags.V = 0
			return result
		})
	}

	count := c.shiftCount(d)
	val := c.getRegD(d.DstReg, d.Size)

	result := val
	c.Flags.V = 0

	for range count {
		lsb := result & 1
		result >>= 1
		setFlag(&c.Flags.C, lsb != 0)
		c.Flags.X = c.Flags.C
	}

	result = maskValue(result, d.Size)
	c.setRegD(d.DstReg, result, d.Size)
	c.setFlagN(result, d.Size)
	c.setFlagZ(result, d.Size)

	if count == 0 {
		c.Flags.C = 0
	}

	return nil
}

func (c *CPU) execROL(d DecodedOpcode) error {
	if d.Extra&0x40 != 0 {
		return c.shiftMemory(d, func(val uint32) uint32 {
			msb := val & 0x8000
			result := ((val << 1) & 0xFFFF) | (msb >> 15)
			setFlag(&c.Flags.C, msb != 0)
			c.Flags.V = 0
			return result
		})
	}

	count := c.shiftCount(d)
	val := c.getRegD(d.DstReg, d.Size)
	msb := msbMask(d.Size)
	bits := uint32(d.Size) * 8

	result := val
	c.Flags.V = 0

	for range count {
		topBit := result & msb
		result = maskValue((result<<1)|(topBit>>(bits-1)), d.Size)
		setFlag(&c.Flags.C, topBit != 0)
	}

	c.setRegD(d.DstReg, result, d.Size)
	c.setFlagN(result, d.Size)
	c.setFlagZ(result, d.Size)

	if count == 0 {
		c.Flags.C = 0
	}

	return nil
}

func (c *CPU) execROR(d DecodedOpcode) error {
	if d.Extra&0x40 != 0 {
		return c.shiftMemory(d, func(val uint32) uint32 {
			lsb := val & 1
			result := (val >> 1) | (lsb << 15)
			setFlag(&c.Flags.C, lsb != 0)
			c.Flags.V = 0
			return result
		})
	}

	count := c.shiftCount(d)
	val := c.getRegD(d.DstReg, d.Size)
	bits := uint32(d.Size) * 8

	result := val
	c.Flags.V = 0

	for range count {
		lsb := result & 1
		result = maskValue((result>>1)|(lsb<<(bits-1)), d.Size)
		setFlag(&c.Flags.C, lsb != 0)
	}

	c.setRegD(d.DstReg, result, d.Size)
	c.setFlagN(result, d.Size)
	c.setFlagZ(result, d.Size)

	if count == 0 {
		c.Flags.C = 0
	}

	return nil
}

func (c *CPU) execROXL(d DecodedOpcode) error {
	if d.Extra&0x40 != 0 {
		return c.shiftMemory(d, func(val uint32) uint32 {
			msb := val & 0x8000
			result := ((val << 1) & 0xFFFF) | uint32(c.Flags.X)
			setFlag(&c.Flags.C, msb != 0)
			c.Flags.X = c.Flags.C
			c.Flags.V = 0
			return result
		})
	}

	count := c.shiftCount(d)
	val := c.getRegD(d.DstReg, d.Size)
	msb := msbMask(d.Size)

	result := val
	c.Flags.V = 0

	for range count {
		topBit := result & msb
		result = maskValue((result<<1)|uint32(c.Flags.X), d.Size)
		setFlag(&c.Flags.C, topBit != 0)
		c.Flags.X = c.Flags.C
	}

	c.setRegD(d.DstReg, result, d.Size)
	c.setFlagN(result, d.Size)
	c.setFlagZ(result, d.Size)
	c.Flags.C = c.Flags.X

	return nil
}

func (c *CPU) execROXR(d DecodedOpcode) error {
	if d.Extra&0x40 != 0 {
		return c.shiftMemory(d, func(val uint32) uint32 {
			lsb := val & 1
			result := (val >> 1) | (uint32(c.Flags.X) << 15)
			setFlag(&c.Flags.C, lsb != 0)
			c.Flags.X = c.Flags.C
			c.Flags.V = 0
			return result
		})
	}

	count := c.shiftCount(d)
	val := c.getRegD(d.DstReg, d.Size)
	bits := uint32(d.Size) * 8

	result := val
	c.Flags.V = 0

	for range count {
		lsb := result & 1
		result = maskValue((result>>1)|(uint32(c.Flags.X)<<(bits-1)), d.Size)
		setFlag(&c.Flags.C, lsb != 0)
		c.Flags.X = c.Flags.C
	}

	c.setRegD(d.DstReg, result, d.Size)
	c.setFlagN(result, d.Size)
	c.setFlagZ(result, d.Size)
	c.Flags.C = c.Flags.X

	return nil
}

// shiftCount returns the shift/rotate count from the opcode extra field.
func (c *CPU) shiftCount(d DecodedOpcode) uint32 {
	count := uint32(d.Extra & 7)
	if d.Extra&0x20 != 0 {
		// Count from register.
		count = c.D[d.Extra&7] % 64
	} else if count == 0 {
		count = 8
	}
	return count
}

// shiftMemory performs a memory shift/rotate operation.
func (c *CPU) shiftMemory(d DecodedOpcode, op func(uint32) uint32) error {
	ea, err := c.decodeEA(d.DstMode, d.DstReg, SizeWord)
	if err != nil {
		return err
	}
	val, err := c.readEA(ea)
	if err != nil {
		return err
	}

	result := op(val)
	c.setFlagN(result, SizeWord)
	c.setFlagZ(result, SizeWord)
	return c.writeEA(ea, result)
}
