package cpu6809

// Branch and jump instructions.

// getRegisterValue returns the value of a register by its TFR/EXG encoding.
// 0=D, 1=X, 2=Y, 3=U, 4=S, 5=PC, 8=A, 9=B, 10=CC, 11=DP
func (c *CPU) getRegisterValue(reg uint8) uint16 {
	switch reg {
	case 0x00:
		return c.D()
	case 0x01:
		return c.X
	case 0x02:
		return c.Y
	case 0x03:
		return c.U
	case 0x04:
		return c.S
	case 0x05:
		return c.PC
	case 0x08:
		return uint16(c.A)
	case 0x09:
		return uint16(c.B)
	case 0x0A:
		return uint16(c.GetCC())
	case 0x0B:
		return uint16(c.DP)
	default:
		return 0
	}
}

// setRegisterValue sets a register by its TFR/EXG encoding.
func (c *CPU) setRegisterValue(reg uint8, value uint16) {
	switch reg {
	case 0x00:
		c.SetD(value)
	case 0x01:
		c.X = value
	case 0x02:
		c.Y = value
	case 0x03:
		c.U = value
	case 0x04:
		c.S = value
	case 0x05:
		c.PC = value
		c.pcChanged = true
	case 0x08:
		c.A = uint8(value)
	case 0x09:
		c.B = uint8(value)
	case 0x0A:
		c.SetCC(uint8(value))
	case 0x0B:
		c.DP = uint8(value)
	}
}

func bccFn(c *CPU, param any) error {
	return branchOn(c, param, c.Flags.C == 0)
}

func bcsFn(c *CPU, param any) error {
	return branchOn(c, param, c.Flags.C != 0)
}

func beqFn(c *CPU, param any) error {
	return branchOn(c, param, c.Flags.Z != 0)
}

func bgeFn(c *CPU, param any) error {
	return branchOn(c, param, c.Flags.N == c.Flags.V)
}

func bgtFn(c *CPU, param any) error {
	return branchOn(c, param, c.Flags.Z == 0 && c.Flags.N == c.Flags.V)
}

func bhiFn(c *CPU, param any) error {
	return branchOn(c, param, c.Flags.C == 0 && c.Flags.Z == 0)
}

func bleFn(c *CPU, param any) error {
	return branchOn(c, param, c.Flags.Z != 0 || c.Flags.N != c.Flags.V)
}

func blsFn(c *CPU, param any) error {
	return branchOn(c, param, c.Flags.C != 0 || c.Flags.Z != 0)
}

func bltFn(c *CPU, param any) error {
	return branchOn(c, param, c.Flags.N != c.Flags.V)
}

func bmiFn(c *CPU, param any) error {
	return branchOn(c, param, c.Flags.N != 0)
}

func bneFn(c *CPU, param any) error {
	return branchOn(c, param, c.Flags.Z == 0)
}

func bplFn(c *CPU, param any) error {
	return branchOn(c, param, c.Flags.N == 0)
}

func braFn(c *CPU, param any) error {
	return branchOn(c, param, true)
}

func brnFn(_ *CPU, _ any) error {
	return nil
}

func bsrFn(c *CPU, param any) error {
	target, err := branchTargetParam(param)
	if err != nil {
		return err
	}

	c.pushS16(c.nextPC)
	c.PC = target
	c.pcChanged = true

	return nil
}

func bvcFn(c *CPU, param any) error {
	return branchOn(c, param, c.Flags.V == 0)
}

func bvsFn(c *CPU, param any) error {
	return branchOn(c, param, c.Flags.V != 0)
}

func branchOn(c *CPU, param any, taken bool) error {
	target, err := branchTargetParam(param)
	if err != nil {
		return err
	}

	c.branch(taken, target)

	return nil
}

// -- Jump --

func jmpFn(c *CPU, param any) error {
	addr, err := c.resolveEA(param)
	if err != nil {
		return err
	}
	c.PC = addr
	c.pcChanged = true
	return nil
}

func jsrFn(c *CPU, param any) error {
	addr, err := c.resolveEA(param)
	if err != nil {
		return err
	}
	// Push return address (address of next instruction, pre-computed by step.go).
	c.pushS16(c.nextPC)
	c.PC = addr
	c.pcChanged = true
	return nil
}

// -- TFR/EXG --

func tfrFn(c *CPU, param any) error {
	pair, err := registerPairParam(param)
	if err != nil {
		return err
	}

	postbyte := uint8(pair)
	src := c.getRegisterValue(postbyte >> 4)
	c.setRegisterValue(postbyte&0x0F, src)

	return nil
}

func exgFn(c *CPU, param any) error {
	pair, err := registerPairParam(param)
	if err != nil {
		return err
	}

	postbyte := uint8(pair)
	srcReg := postbyte >> 4
	dstReg := postbyte & 0x0F
	src := c.getRegisterValue(srcReg)
	dst := c.getRegisterValue(dstReg)
	c.setRegisterValue(srcReg, dst)
	c.setRegisterValue(dstReg, src)
	return nil
}
