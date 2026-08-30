package cpu6502

func dcp(c *CPU, params ...any) error {
	if err := dec(c, params...); err != nil {
		return err
	}
	val, err := c.memory.ReadAddressModes(false, params...)
	if err != nil {
		return err
	}
	c.compare(c.A, val)
	return nil
}

func isc(c *CPU, params ...any) error {
	if err := inc(c, params...); err != nil {
		return err
	}
	return sbc(c, params...)
}

// las implements LAS/LAR (0xBB): result = mem & SP; A = X = SP = result.
func las(c *CPU, params ...any) error {
	val, err := c.memory.ReadAddressModes(false, params...)
	if err != nil {
		return err
	}
	result := val & c.SP
	c.A = result
	c.X = result
	c.SP = result
	c.setZN(result)
	return nil
}

func lax(c *CPU, params ...any) error {
	val, err := c.memory.ReadAddressModes(false, params...)
	if err != nil {
		return err
	}
	c.A = val
	c.X = c.A
	c.setZN(c.A)
	return nil
}

// lxa implements LXA (0xAB): A = X = (A | magic) & imm.
// This is a highly unstable instruction; the magic constant varies by chip.
func lxa(c *CPU, params ...any) error {
	val, err := c.memory.ReadAddressModes(true, params...)
	if err != nil {
		return err
	}
	// 0xEE matches SingleStepTests/65x02 hardware captures.
	const magicConstant = 0xEE
	c.A = (c.A | magicConstant) & val
	c.X = c.A
	c.setZN(c.A)
	return nil
}

func nopUnofficial(c *CPU, params ...any) error {
	if len(params) > 0 {
		_, err := c.memory.ReadAddressModes(false, params...)
		return err
	}
	return nil
}

func rla(c *CPU, params ...any) error {
	if err := rol(c, params...); err != nil {
		return err
	}
	return and(c, params...)
}

func rra(c *CPU, params ...any) error {
	if err := ror(c, params...); err != nil {
		return err
	}
	return adc(c, params...)
}

func sax(c *CPU, params ...any) error {
	val := c.A & c.X
	return c.memory.WriteAddressModes(val, params...)
}

// sha implements SHA/AHX for absolute,Y and indirect,Y addressing.
// A page crossing replaces the write address high byte with the stored value.
func sha(c *CPU, params ...any) error {
	baseAddr, indexReg := shBaseAddr(c, params)
	shWrite(c, c.A&c.X, baseAddr, indexReg)
	return nil
}

// shBaseAddr extracts the address before indexing and its index register.
// Indirect,Y parameters contain only the resolved address, so reread the base.
func shBaseAddr(c *CPU, params []any) (uint16, uint8) {
	if _, ok := params[0].(Absolute); ok {
		return uint16(params[0].(Absolute)), *params[1].(*uint8)
	}
	zp := c.memory.Read(c.PC + 1)
	baseAddr := c.memory.ReadWordBug(uint16(zp))
	return baseAddr, c.Y
}

// shWrite emulates the SH-family page-crossing bus conflict.
func shWrite(c *CPU, value uint8, baseAddr uint16, indexReg uint8) {
	andValue := value & (byte(baseAddr>>8) + 1)
	effectiveAddr := baseAddr + uint16(indexReg)
	pageCrossed := effectiveAddr&0xFF00 != baseAddr&0xFF00

	var writeAddr uint16
	if pageCrossed {
		writeAddr = (uint16(andValue) << 8) | (effectiveAddr & 0xFF)
	} else {
		writeAddr = effectiveAddr
	}
	c.memory.Write(writeAddr, andValue)
}

// shx implements SHX/SXA (0x9E): stores X & (base_addr_hi + 1).
func shx(c *CPU, params ...any) error {
	baseAddr := uint16(params[0].(Absolute))
	indexReg := *params[1].(*uint8) // Y
	shWrite(c, c.X, baseAddr, indexReg)
	return nil
}

// shy implements SHY/SYA (0x9C): stores Y & (base_addr_hi + 1).
func shy(c *CPU, params ...any) error {
	baseAddr := uint16(params[0].(Absolute))
	indexReg := *params[1].(*uint8) // X
	shWrite(c, c.Y, baseAddr, indexReg)
	return nil
}

func slo(c *CPU, params ...any) error {
	if err := asl(c, params...); err != nil {
		return err
	}
	return ora(c, params...)
}

func sre(c *CPU, params ...any) error {
	if err := lsr(c, params...); err != nil {
		return err
	}
	return eor(c, params...)
}

// tas implements TAS/XAS (0x9B): SP = A & X, then stores SP & (base_addr_hi + 1).
func tas(c *CPU, params ...any) error {
	c.SP = c.A & c.X
	baseAddr := uint16(params[0].(Absolute))
	indexReg := *params[1].(*uint8) // Y
	shWrite(c, c.SP, baseAddr, indexReg)
	return nil
}

// alr - AND with accumulator, then LSR.
func alr(c *CPU, params ...any) error {
	if err := and(c, params...); err != nil {
		return err
	}
	c.Flags.C = c.A & 1
	c.A >>= 1
	c.setZN(c.A)
	return nil
}

// anc - AND with accumulator, copy N flag to C flag.
func anc(c *CPU, params ...any) error {
	if err := and(c, params...); err != nil {
		return err
	}
	c.Flags.C = c.Flags.N
	return nil
}

// ane implements ANE/XAA (0x8B): A = (A | magic) & X & imm.
// This is a highly unstable instruction; the magic constant varies by chip.
func ane(c *CPU, params ...any) error {
	val, err := c.memory.ReadAddressModes(true, params...)
	if err != nil {
		return err
	}
	// 0xEE matches SingleStepTests/65x02 hardware captures.
	const magicConstant = 0xEE
	c.A = (c.A | magicConstant) & c.X & val
	c.setZN(c.A)
	return nil
}

// arr applies AND followed by ROR with ARR-specific flag behavior.
func arr(c *CPU, params ...any) error {
	if err := and(c, params...); err != nil {
		return err
	}
	andResult := c.A

	oldCarry := c.Flags.C
	c.A = (andResult >> 1) | (oldCarry << 7)
	c.setZN(c.A)
	c.Flags.V = (c.A>>6)&1 ^ (c.A>>5)&1

	if c.Flags.D != 0 && c.opts.variant != VariantNES6502 {
		r := c.A
		if (andResult & 0x0F) >= 5 {
			r = (r & 0xF0) | ((r&0x0F + 6) & 0x0F)
		}
		if (andResult & 0xF0) >= 0x50 {
			r += 0x60
			c.Flags.C = 1
		} else {
			c.Flags.C = 0
		}
		c.A = r
	} else {
		c.Flags.C = (c.A >> 6) & 1
	}
	return nil
}

// axs - (A AND X) minus immediate, store in X.
func axs(c *CPU, params ...any) error {
	value, err := c.memory.ReadAddressModes(true, params...)
	if err != nil {
		return err
	}

	val := c.A & c.X
	result := int(val) - int(value)
	setFlag(&c.Flags.C, result >= 0)
	c.X = uint8(result)
	c.setZN(c.X)
	return nil
}
