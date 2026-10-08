package cp1610

func (c *CPU) executeShift(op Opcode) {
	register := op.Destination
	value := c.R[register]
	count := uint(op.Variant) + 1
	var result uint16

	switch op.Instruction.Name {
	case SwapName:
		if count == 1 {
			result = value<<8 | value>>8
		} else {
			result = value&0xFF | value<<8
		}
		c.Flags.S = result&0x80 != 0

	case SllName:
		result = value << count
		c.Flags.S = result&0x8000 != 0

	case SllcName:
		result = value << count
		c.setShiftLeftCarry(value, count)
		c.Flags.S = result&0x8000 != 0

	case RlcName:
		if count == 1 {
			result = value<<1 | boolWord(c.Flags.C)
		} else {
			result = value<<2 | boolWord(c.Flags.C)<<1 | boolWord(c.Flags.O)
		}
		c.setShiftLeftCarry(value, count)
		c.Flags.S = result&0x8000 != 0

	case SlrName:
		result = value >> count
		c.Flags.S = result&0x80 != 0

	case SarName:
		result = uint16(int16(value) >> count)
		c.Flags.S = result&0x80 != 0

	case RrcName:
		if count == 1 {
			result = value>>1 | boolWord(c.Flags.C)<<15
		} else {
			result = value>>2 | boolWord(c.Flags.C)<<14 | boolWord(c.Flags.O)<<15
		}
		c.setShiftRightCarry(value, count)
		c.Flags.S = result&0x80 != 0

	case SarcName:
		result = uint16(int16(value) >> count)
		c.setShiftRightCarry(value, count)
		c.Flags.S = result&0x80 != 0
	}
	c.Flags.Z = result == 0
	c.R[register] = result
}

func (c *CPU) setShiftLeftCarry(value uint16, count uint) {
	c.Flags.C = value&0x8000 != 0
	if count == 2 {
		c.Flags.O = value&0x4000 != 0
	}
}

func (c *CPU) setShiftRightCarry(value uint16, count uint) {
	c.Flags.C = value&1 != 0
	if count == 2 {
		c.Flags.O = value&2 != 0
	}
}

func boolWord(value bool) uint16 {
	if value {
		return 1
	}
	return 0
}
