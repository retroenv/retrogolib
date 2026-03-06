package m68000

// Bit manipulation instructions: BTST, BSET, BCLR, BCHG.

func (c *CPU) execBTST(d DecodedOpcode) error {
	bitNum := c.getBitNumber(d)

	// For register operands, bit number modulo 32. For memory, modulo 8.
	if d.DstMode == 0 {
		bitNum %= 32
		val := c.D[d.DstReg]
		setFlag(&c.Flags.Z, val&(1<<bitNum) == 0)
		return nil
	}

	bitNum %= 8
	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, SizeByte)
	if err != nil {
		return err
	}
	val, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	setFlag(&c.Flags.Z, val&(1<<bitNum) == 0)
	return nil
}

func (c *CPU) execBSET(d DecodedOpcode) error {
	bitNum := c.getBitNumber(d)

	if d.DstMode == 0 {
		bitNum %= 32
		val := c.D[d.DstReg]
		setFlag(&c.Flags.Z, val&(1<<bitNum) == 0)
		c.D[d.DstReg] = val | (1 << bitNum)
		return nil
	}

	bitNum %= 8
	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, SizeByte)
	if err != nil {
		return err
	}
	val, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	setFlag(&c.Flags.Z, val&(1<<bitNum) == 0)
	return c.writeEA(dstEA, val|(1<<bitNum))
}

func (c *CPU) execBCLR(d DecodedOpcode) error {
	bitNum := c.getBitNumber(d)

	if d.DstMode == 0 {
		bitNum %= 32
		val := c.D[d.DstReg]
		setFlag(&c.Flags.Z, val&(1<<bitNum) == 0)
		c.D[d.DstReg] = val &^ (1 << bitNum)
		return nil
	}

	bitNum %= 8
	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, SizeByte)
	if err != nil {
		return err
	}
	val, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	setFlag(&c.Flags.Z, val&(1<<bitNum) == 0)
	return c.writeEA(dstEA, val&^(1<<bitNum))
}

func (c *CPU) execBCHG(d DecodedOpcode) error {
	bitNum := c.getBitNumber(d)

	if d.DstMode == 0 {
		bitNum %= 32
		val := c.D[d.DstReg]
		setFlag(&c.Flags.Z, val&(1<<bitNum) == 0)
		c.D[d.DstReg] = val ^ (1 << bitNum)
		return nil
	}

	bitNum %= 8
	dstEA, err := c.decodeEA(d.DstMode, d.DstReg, SizeByte)
	if err != nil {
		return err
	}
	val, err := c.readEA(dstEA)
	if err != nil {
		return err
	}

	setFlag(&c.Flags.Z, val&(1<<bitNum) == 0)
	return c.writeEA(dstEA, val^(1<<bitNum))
}

// getBitNumber returns the bit number from the source operand.
func (c *CPU) getBitNumber(d DecodedOpcode) uint32 {
	if d.SrcMode == 7 && d.SrcReg == 4 {
		// Immediate bit number.
		return c.readImmediate(SizeByte)
	}
	// Register bit number.
	return c.D[d.SrcReg]
}
