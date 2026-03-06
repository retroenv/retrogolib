package m65816

// Data movement instructions: LDA, LDX, LDY, STA, STX, STY, STZ,
// TAX, TAY, TCD, TCS, TDC, TSC, TSX, TXA, TXS, TXY, TYA, TYX.

func lda(c *CPU, params ...any) error {
	val, err := c.readOperandAcc(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		c.C = uint16(c.B())<<8 | uint16(uint8(val))
		c.setZN8(uint8(c.C))
	} else {
		c.C = val
		c.setZN16(c.C)
	}
	return nil
}

func ldx(c *CPU, params ...any) error {
	val, err := c.readOperandIdx(params[0])
	if err != nil {
		return err
	}
	if c.IdxWidth() == 1 {
		c.X = uint16(val & 0xFF)
		c.setZN8(uint8(c.X))
	} else {
		c.X = val
		c.setZN16(c.X)
	}
	return nil
}

func ldy(c *CPU, params ...any) error {
	val, err := c.readOperandIdx(params[0])
	if err != nil {
		return err
	}
	if c.IdxWidth() == 1 {
		c.Y = uint16(val & 0xFF)
		c.setZN8(uint8(c.Y))
	} else {
		c.Y = val
		c.setZN16(c.Y)
	}
	return nil
}

func sta(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		c.writeMem8(addr, uint8(c.C))
	} else {
		c.writeMem16(addr, c.C)
	}
	return nil
}

func stx(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	if c.IdxWidth() == 1 {
		c.writeMem8(addr, uint8(c.X))
	} else {
		c.writeMem16(addr, c.X)
	}
	return nil
}

func sty(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	if c.IdxWidth() == 1 {
		c.writeMem8(addr, uint8(c.Y))
	} else {
		c.writeMem16(addr, c.Y)
	}
	return nil
}

func stz(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	if c.AccWidth() == 1 {
		c.writeMem8(addr, 0)
	} else {
		c.writeMem16(addr, 0)
	}
	return nil
}

// -- Register transfers --

func tax(c *CPU) error {
	if c.IdxWidth() == 1 {
		c.X = uint16(uint8(c.C))
		c.setZN8(uint8(c.X))
	} else {
		c.X = c.C
		c.setZN16(c.X)
	}
	return nil
}

func tay(c *CPU) error {
	if c.IdxWidth() == 1 {
		c.Y = uint16(uint8(c.C))
		c.setZN8(uint8(c.Y))
	} else {
		c.Y = c.C
		c.setZN16(c.Y)
	}
	return nil
}

func tcd(c *CPU) error {
	c.DP = c.C
	c.setZN16(c.DP)
	return nil
}

func tcs(c *CPU) error {
	if c.E {
		// Emulation: only low byte transferred, SP high byte stays $01
		c.SP = 0x0100 | uint16(uint8(c.C))
	} else {
		c.SP = c.C
	}
	return nil
}

func tdc(c *CPU) error {
	c.C = c.DP
	c.setZN16(c.C)
	return nil
}

func tsc(c *CPU) error {
	c.C = c.SP
	c.setZN16(c.C)
	return nil
}

func tsx(c *CPU) error {
	if c.IdxWidth() == 1 {
		c.X = uint16(uint8(c.SP))
		c.setZN8(uint8(c.X))
	} else {
		c.X = c.SP
		c.setZN16(c.X)
	}
	return nil
}

func txa(c *CPU) error {
	if c.AccWidth() == 1 {
		v := uint8(c.X)
		c.C = uint16(c.B())<<8 | uint16(v)
		c.setZN8(v)
	} else {
		c.C = c.X
		c.setZN16(c.C)
	}
	return nil
}

func txs(c *CPU) error {
	if c.E {
		c.SP = 0x0100 | uint16(uint8(c.X))
	} else {
		c.SP = c.X
	}
	return nil
}

func txy(c *CPU) error {
	if c.IdxWidth() == 1 {
		c.Y = uint16(uint8(c.X))
		c.setZN8(uint8(c.Y))
	} else {
		c.Y = c.X
		c.setZN16(c.Y)
	}
	return nil
}

func tya(c *CPU) error {
	if c.AccWidth() == 1 {
		v := uint8(c.Y)
		c.C = uint16(c.B())<<8 | uint16(v)
		c.setZN8(v)
	} else {
		c.C = c.Y
		c.setZN16(c.C)
	}
	return nil
}

func tyx(c *CPU) error {
	if c.IdxWidth() == 1 {
		c.X = uint16(uint8(c.Y))
		c.setZN8(uint8(c.X))
	} else {
		c.X = c.Y
		c.setZN16(c.X)
	}
	return nil
}
