package m6809

// PSH/PUL stack operations.
//
// Register bitmask in postbyte:
//   bit 0: CC    bit 4: S (PSHU/PULU) or U (PSHS/PULS)
//   bit 1: A     bit 5: Y
//   bit 2: B     bit 6: X (note: X and Y are swapped from natural order)
//   bit 3: DP    bit 7: PC
//
// Actually the correct order per Motorola documentation:
//   bit 0: CC    bit 4: S (for PSHU/PULU) or U (for PSHS/PULS)
//   bit 1: A     bit 5: Y
//   bit 2: B     bit 6: X
//   bit 3: DP    bit 7: PC
//
// Push order: PC first (highest bit), then U/S, Y, X, DP, B, A, CC last.
// Pull order: CC first, A, B, DP, X, Y, U/S, PC last.

// pshsFn - Push registers onto the System stack.
func pshsFn(c *CPU, params ...any) error {
	mask := uint8(params[0].(StackMask))
	if mask&0x80 != 0 {
		c.pushS16(c.PC)
	}
	if mask&0x40 != 0 {
		c.pushS16(c.U)
	}
	if mask&0x20 != 0 {
		c.pushS16(c.Y)
	}
	if mask&0x10 != 0 {
		c.pushS16(c.X)
	}
	if mask&0x08 != 0 {
		c.pushS8(c.DP)
	}
	if mask&0x04 != 0 {
		c.pushS8(c.B)
	}
	if mask&0x02 != 0 {
		c.pushS8(c.A)
	}
	if mask&0x01 != 0 {
		c.pushS8(c.GetCC())
	}
	return nil
}

// pulsFn - Pull registers from the System stack.
func pulsFn(c *CPU, params ...any) error {
	mask := uint8(params[0].(StackMask))
	if mask&0x01 != 0 {
		c.SetCC(c.popS8())
	}
	if mask&0x02 != 0 {
		c.A = c.popS8()
	}
	if mask&0x04 != 0 {
		c.B = c.popS8()
	}
	if mask&0x08 != 0 {
		c.DP = c.popS8()
	}
	if mask&0x10 != 0 {
		c.X = c.popS16()
	}
	if mask&0x20 != 0 {
		c.Y = c.popS16()
	}
	if mask&0x40 != 0 {
		c.U = c.popS16()
	}
	if mask&0x80 != 0 {
		c.PC = c.popS16()
		c.pcChanged = true
	}
	return nil
}

// pshuFn - Push registers onto the User stack.
func pshuFn(c *CPU, params ...any) error {
	mask := uint8(params[0].(StackMask))
	if mask&0x80 != 0 {
		c.pushU16(c.PC)
	}
	if mask&0x40 != 0 {
		c.pushU16(c.S) // PSHU pushes S, not U
	}
	if mask&0x20 != 0 {
		c.pushU16(c.Y)
	}
	if mask&0x10 != 0 {
		c.pushU16(c.X)
	}
	if mask&0x08 != 0 {
		c.pushU8(c.DP)
	}
	if mask&0x04 != 0 {
		c.pushU8(c.B)
	}
	if mask&0x02 != 0 {
		c.pushU8(c.A)
	}
	if mask&0x01 != 0 {
		c.pushU8(c.GetCC())
	}
	return nil
}

// puluFn - Pull registers from the User stack.
func puluFn(c *CPU, params ...any) error {
	mask := uint8(params[0].(StackMask))
	if mask&0x01 != 0 {
		c.SetCC(c.popU8())
	}
	if mask&0x02 != 0 {
		c.A = c.popU8()
	}
	if mask&0x04 != 0 {
		c.B = c.popU8()
	}
	if mask&0x08 != 0 {
		c.DP = c.popU8()
	}
	if mask&0x10 != 0 {
		c.X = c.popU16()
	}
	if mask&0x20 != 0 {
		c.Y = c.popU16()
	}
	if mask&0x40 != 0 {
		c.S = c.popU16() // PULU pulls S, not U
	}
	if mask&0x80 != 0 {
		c.PC = c.popU16()
		c.pcChanged = true
	}
	return nil
}
