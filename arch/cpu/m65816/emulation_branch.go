package m65816

// Branch and jump instructions.

func bcc(c *CPU, params ...any) error {
	c.branch(c.Flags.C == 0, params[0].(uint16))
	return nil
}

func bcs(c *CPU, params ...any) error {
	c.branch(c.Flags.C != 0, params[0].(uint16))
	return nil
}

func beq(c *CPU, params ...any) error {
	c.branch(c.Flags.Z != 0, params[0].(uint16))
	return nil
}

func bmi(c *CPU, params ...any) error {
	c.branch(c.Flags.N != 0, params[0].(uint16))
	return nil
}

func bne(c *CPU, params ...any) error {
	c.branch(c.Flags.Z == 0, params[0].(uint16))
	return nil
}

func bpl(c *CPU, params ...any) error {
	c.branch(c.Flags.N == 0, params[0].(uint16))
	return nil
}

func bra(c *CPU, params ...any) error {
	c.branch(true, params[0].(uint16))
	return nil
}

func brl(c *CPU, params ...any) error {
	// Branch Long: always taken, 16-bit offset, already resolved to absolute
	c.PC = params[0].(uint16)
	c.pcChanged = true
	return nil
}

func bvc(c *CPU, params ...any) error {
	c.branch(c.Flags.V == 0, params[0].(uint16))
	return nil
}

func bvs(c *CPU, params ...any) error {
	c.branch(c.Flags.V != 0, params[0].(uint16))
	return nil
}

// jmp - Jump (same bank).
func jmp(c *CPU, params ...any) error {
	switch p := params[0].(type) {
	case Absolute16:
		c.PC = uint16(p)
	case DPIndirect:
		c.PC = uint16(p)
	case DPIndirectX:
		c.PC = uint16(p)
	default:
		return nil
	}
	c.pcChanged = true
	return nil
}

// jml - Jump Long (sets PB).
func jml(c *CPU, params ...any) error {
	switch p := params[0].(type) {
	case AbsLong:
		c.PB = uint8(uint32(p) >> 16)
		c.PC = uint16(p)
		c.pcChanged = true
	}
	return nil
}

// jsr - Jump to Subroutine (saves PC-1 onto stack).
func jsr(c *CPU, params ...any) error {
	// JSR pushes PC+2 (address of last byte of JSR instruction, i.e. PC-1 from next instruction).
	// Instruction size = 3 bytes; return address = PC+2 (pointing to last byte).
	retAddr := c.PC + 2
	c.push16(retAddr)
	switch p := params[0].(type) {
	case Absolute16:
		c.PC = uint16(p)
		c.pcChanged = true
	case DPIndirectX:
		c.PC = uint16(p)
		c.pcChanged = true
	}
	return nil
}

// jsl - Jump to Subroutine Long.
// 65816-native: uses full 16-bit SP (no page-1 wrap between bytes).
func jsl(c *CPU, params ...any) error {
	// Pushes PB, then PC+3 (last byte of JSL instruction).
	c.push8raw(c.PB)
	retAddr := c.PC + 3
	c.push16raw(retAddr)
	c.fixEmuSP()
	switch p := params[0].(type) {
	case AbsLong:
		c.PB = uint8(uint32(p) >> 16)
		c.PC = uint16(p)
		c.pcChanged = true
	}
	return nil
}
