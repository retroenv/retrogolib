package m6809

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

func bccFn(c *CPU, params ...any) error {
	c.branch(c.Flags.C == 0, params[0].(uint16))
	return nil
}

func bcsFn(c *CPU, params ...any) error {
	c.branch(c.Flags.C != 0, params[0].(uint16))
	return nil
}

func beqFn(c *CPU, params ...any) error {
	c.branch(c.Flags.Z != 0, params[0].(uint16))
	return nil
}

func bgeFn(c *CPU, params ...any) error {
	c.branch(c.Flags.N == c.Flags.V, params[0].(uint16))
	return nil
}

func bgtFn(c *CPU, params ...any) error {
	c.branch(c.Flags.Z == 0 && c.Flags.N == c.Flags.V, params[0].(uint16))
	return nil
}

func bhiFn(c *CPU, params ...any) error {
	c.branch(c.Flags.C == 0 && c.Flags.Z == 0, params[0].(uint16))
	return nil
}

func bleFn(c *CPU, params ...any) error {
	c.branch(c.Flags.Z != 0 || c.Flags.N != c.Flags.V, params[0].(uint16))
	return nil
}

func blsFn(c *CPU, params ...any) error {
	c.branch(c.Flags.C != 0 || c.Flags.Z != 0, params[0].(uint16))
	return nil
}

func bltFn(c *CPU, params ...any) error {
	c.branch(c.Flags.N != c.Flags.V, params[0].(uint16))
	return nil
}

func bmiFn(c *CPU, params ...any) error {
	c.branch(c.Flags.N != 0, params[0].(uint16))
	return nil
}

func bneFn(c *CPU, params ...any) error {
	c.branch(c.Flags.Z == 0, params[0].(uint16))
	return nil
}

func bplFn(c *CPU, params ...any) error {
	c.branch(c.Flags.N == 0, params[0].(uint16))
	return nil
}

func braFn(c *CPU, params ...any) error {
	c.branch(true, params[0].(uint16))
	return nil
}

func brnFn(_ *CPU, _ ...any) error {
	return nil
}

func bsrFn(c *CPU, params ...any) error {
	// Push return address (address of next instruction, pre-computed by step.go).
	c.pushS16(c.nextPC)
	c.PC = params[0].(uint16)
	c.pcChanged = true
	return nil
}

func bvcFn(c *CPU, params ...any) error {
	c.branch(c.Flags.V == 0, params[0].(uint16))
	return nil
}

func bvsFn(c *CPU, params ...any) error {
	c.branch(c.Flags.V != 0, params[0].(uint16))
	return nil
}

// -- Long branches --

func lbccFn(c *CPU, params ...any) error {
	c.branch(c.Flags.C == 0, params[0].(uint16))
	return nil
}

func lbcsFn(c *CPU, params ...any) error {
	c.branch(c.Flags.C != 0, params[0].(uint16))
	return nil
}

func lbeqFn(c *CPU, params ...any) error {
	c.branch(c.Flags.Z != 0, params[0].(uint16))
	return nil
}

func lbgeFn(c *CPU, params ...any) error {
	c.branch(c.Flags.N == c.Flags.V, params[0].(uint16))
	return nil
}

func lbgtFn(c *CPU, params ...any) error {
	c.branch(c.Flags.Z == 0 && c.Flags.N == c.Flags.V, params[0].(uint16))
	return nil
}

func lbhiFn(c *CPU, params ...any) error {
	c.branch(c.Flags.C == 0 && c.Flags.Z == 0, params[0].(uint16))
	return nil
}

func lbleFn(c *CPU, params ...any) error {
	c.branch(c.Flags.Z != 0 || c.Flags.N != c.Flags.V, params[0].(uint16))
	return nil
}

func lblsFn(c *CPU, params ...any) error {
	c.branch(c.Flags.C != 0 || c.Flags.Z != 0, params[0].(uint16))
	return nil
}

func lbltFn(c *CPU, params ...any) error {
	c.branch(c.Flags.N != c.Flags.V, params[0].(uint16))
	return nil
}

func lbmiFn(c *CPU, params ...any) error {
	c.branch(c.Flags.N != 0, params[0].(uint16))
	return nil
}

func lbneFn(c *CPU, params ...any) error {
	c.branch(c.Flags.Z == 0, params[0].(uint16))
	return nil
}

func lbplFn(c *CPU, params ...any) error {
	c.branch(c.Flags.N == 0, params[0].(uint16))
	return nil
}

func lbraFn(c *CPU, params ...any) error {
	c.branch(true, params[0].(uint16))
	return nil
}

func lbrnFn(_ *CPU, _ ...any) error {
	return nil
}

func lbsrFn(c *CPU, params ...any) error {
	// Push return address (address of next instruction, pre-computed by step.go).
	c.pushS16(c.nextPC)
	c.PC = params[0].(uint16)
	c.pcChanged = true
	return nil
}

func lbvcFn(c *CPU, params ...any) error {
	c.branch(c.Flags.V == 0, params[0].(uint16))
	return nil
}

func lbvsFn(c *CPU, params ...any) error {
	c.branch(c.Flags.V != 0, params[0].(uint16))
	return nil
}

// -- Jump --

func jmpFn(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
	if err != nil {
		return err
	}
	c.PC = addr
	c.pcChanged = true
	return nil
}

func jsrFn(c *CPU, params ...any) error {
	addr, err := c.resolveEA(params[0])
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

func tfrFn(c *CPU, params ...any) error {
	postbyte := uint8(params[0].(RegisterPair))
	src := c.getRegisterValue(postbyte >> 4)
	c.setRegisterValue(postbyte&0x0F, src)
	return nil
}

func exgFn(c *CPU, params ...any) error {
	postbyte := uint8(params[0].(RegisterPair))
	srcReg := postbyte >> 4
	dstReg := postbyte & 0x0F
	src := c.getRegisterValue(srcReg)
	dst := c.getRegisterValue(dstReg)
	c.setRegisterValue(srcReg, dst)
	c.setRegisterValue(dstReg, src)
	return nil
}
