package m68000

// System and privileged instructions: TRAP, TRAPV, CHK, RTE, RTS, RTR, STOP,
// RESET, ILLEGAL, TAS.

func (c *CPU) execTRAP(d DecodedOpcode) error {
	vector := int(VectorTrap0) + int(d.Extra)
	return c.processException(vector)
}

func (c *CPU) execTRAPV(_ DecodedOpcode) error {
	if c.Flags.V != 0 {
		return c.processException(VectorTRAPV)
	}
	return nil
}

func (c *CPU) execCHK(d DecodedOpcode) error {
	srcEA, err := c.decodeEA(d.SrcMode, d.SrcReg, SizeWord)
	if err != nil {
		return err
	}
	src, err := c.readEA(srcEA)
	if err != nil {
		return err
	}

	dn := int16(c.D[d.DstReg])
	upper := int16(src)

	if dn < 0 {
		c.Flags.N = 1
		return c.processException(VectorCHK)
	}

	if dn > upper {
		c.Flags.N = 0
		return c.processException(VectorCHK)
	}

	return nil
}

func (c *CPU) execRTE(_ DecodedOpcode) error {
	if !c.IsSupervisor() {
		return c.processException(VectorPrivilege)
	}

	sr := c.pop16()
	pc := c.pop32()
	c.SetSR(sr)
	c.PC = pc
	return nil
}

func (c *CPU) execRTS(_ DecodedOpcode) error {
	c.PC = c.pop32()
	return nil
}

func (c *CPU) execRTR(_ DecodedOpcode) error {
	ccr := c.pop16()
	c.SetCCR(uint8(ccr))
	c.PC = c.pop32()
	return nil
}

func (c *CPU) execSTOP(_ DecodedOpcode) error {
	if !c.IsSupervisor() {
		return c.processException(VectorPrivilege)
	}

	sr := c.readWord()
	c.SetSR(sr)
	c.stopped = true
	return nil
}

func (c *CPU) execRESET(_ DecodedOpcode) error {
	if !c.IsSupervisor() {
		return c.processException(VectorPrivilege)
	}

	c.bus.OnReset()
	return nil
}

func (c *CPU) execILLEGAL(d DecodedOpcode) error {
	// Check for Line A / Line F traps.
	opcodeWord := d.Extra
	if opcodeWord != 0 {
		lineNibble := (opcodeWord >> 8) & 0xF0
		switch lineNibble {
		case 0xA0:
			return c.processException(VectorLineA)
		case 0xF0:
			return c.processException(VectorLineF)
		}
	}
	return c.processException(VectorIllegal)
}

func (c *CPU) execTAS(d DecodedOpcode) error {
	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, SizeByte)
	if err != nil {
		return err
	}
	val, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	// Test and set: test the byte, then set bit 7.
	c.setFlagN(val, SizeByte)
	c.setFlagZ(val, SizeByte)
	c.Flags.V = 0
	c.Flags.C = 0

	return c.writeEA(dstEA, val|0x80)
}
