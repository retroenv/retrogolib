package cpu6809

import "math/bits"

// PSH/PUL stack operations.
//
// Register bitmask in postbyte:
//   bit 0: CC    bit 4: S (PSHU/PULU) or U (PSHS/PULS)
//   bit 1: A     bit 5: Y
//   bit 2: B     bit 6: X (note: X and Y are swapped from natural order)
//   bit 3: DP    bit 7: PC
// Push order: PC first (highest bit), then U/S, Y, X, DP, B, A, CC last.
// Pull order: CC first, A, B, DP, X, Y, U/S, PC last.

func pshsFn(c *CPU, param any) error {
	mask, err := stackMask(param)
	if err != nil {
		return err
	}

	if mask&0x80 != 0 {
		c.pushS16(c.nextPC)
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
	c.cycles += stackByteCount(mask)

	return nil
}

func pulsFn(c *CPU, param any) error {
	mask, err := stackMask(param)
	if err != nil {
		return err
	}

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
	c.cycles += stackByteCount(mask)

	return nil
}

func pshuFn(c *CPU, param any) error {
	mask, err := stackMask(param)
	if err != nil {
		return err
	}

	if mask&0x80 != 0 {
		c.pushU16(c.nextPC)
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
	c.cycles += stackByteCount(mask)

	return nil
}

func puluFn(c *CPU, param any) error {
	mask, err := stackMask(param)
	if err != nil {
		return err
	}

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
	c.cycles += stackByteCount(mask)

	return nil
}

func stackMask(param any) (uint8, error) {
	mask, err := stackMaskParam(param)

	return uint8(mask), err
}

func stackByteCount(mask uint8) uint64 {
	return uint64(bits.OnesCount8(mask&0x0F) + 2*bits.OnesCount8(mask&0xF0))
}
