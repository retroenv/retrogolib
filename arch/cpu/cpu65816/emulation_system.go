package cpu65816

// System, processor status, and miscellaneous instructions.

// enterInterrupt sets the interrupt flags and loads the vector into PB:PC.
func (c *CPU) enterInterrupt(vector uint32) {
	c.Flags.I = 1
	c.Flags.D = 0 // 65C02 behavior: clear D on interrupt
	c.PB = 0
	c.PC = c.memory.ReadWord(vector)
	c.pcChanged = true
}

// markIRQRunning records that a software interrupt handler is active until RTI.
func (c *CPU) markIRQRunning() {
	c.mu.Lock()
	c.irqRunning = true
	c.mu.Unlock()
}

// clc - Clear Carry.
func clc(c *CPU) error { c.Flags.C = 0; return nil }

// cld - Clear Decimal.
func cld(c *CPU) error { c.Flags.D = 0; return nil }

// cli - Clear Interrupt Disable.
func cli(c *CPU) error { c.Flags.I = 0; return nil }

// clv - Clear Overflow.
func clv(c *CPU) error { c.Flags.V = 0; return nil }

// sec - Set Carry.
func sec(c *CPU) error { c.Flags.C = 1; return nil }

// sed - Set Decimal.
func sed(c *CPU) error { c.Flags.D = 1; return nil }

// sei - Set Interrupt Disable.
func sei(c *CPU) error { c.Flags.I = 1; return nil }

// rep - Reset Processor Status Bits (clear bits specified by immediate mask).
func rep(c *CPU, params ...any) error {
	mask, err := operand[Immediate8](params)
	if err != nil {
		return err
	}
	c.SetP(c.GetP() &^ uint8(mask))
	return nil
}

// sep - Set Processor Status Bits (set bits specified by immediate mask).
func sep(c *CPU, params ...any) error {
	mask, err := operand[Immediate8](params)
	if err != nil {
		return err
	}
	c.SetP(c.GetP() | uint8(mask))
	return nil
}

// xce - Exchange Carry and Emulation flags.
// Entry into emulation mode forces M=1 and X=1, clears the high bytes of X
// and Y, and moves SP to page 1. Entry into native mode keeps M and X.
func xce(c *CPU) error {
	oldE := c.E
	oldC := c.Flags.C

	if oldE {
		c.Flags.C = 1
	} else {
		c.Flags.C = 0
	}

	c.E = oldC != 0
	if c.E {
		c.Flags.M = 1
		c.Flags.X = 1
		c.X &= 0x00FF
		c.Y &= 0x00FF
		c.SP = 0x0100 | (c.SP & 0x00FF)
	}
	return nil
}

// xba - Exchange B and A (swap high and low bytes of accumulator C).
func xba(c *CPU) error {
	lo := uint8(c.C)
	hi := uint8(c.C >> 8)
	c.C = uint16(lo)<<8 | uint16(hi)
	// N and Z follow the new low byte (the new A).
	c.setZN8(hi)
	return nil
}

// stp - Stop the Processor (halts until RESET).
func stp(c *CPU) error {
	c.stopped = true
	return nil
}

// wai - Wait for Interrupt.
func wai(c *CPU) error {
	c.waiting = true
	return nil
}

// brk - Software Interrupt.
// BRK is 2 bytes; the pushed return address is PC+2 (after the signature byte).
func brk(c *CPU) error {
	retAddr := c.PC + 2
	if c.E {
		c.push16(retAddr)
		p := c.GetP() | MaskBreak // B flag set when pushed in emulation mode
		c.push8(p)
		c.enterInterrupt(VectorEmuIRQ)
	} else {
		c.push8(c.PB)
		c.push16(retAddr)
		c.push8(c.GetP())
		c.enterInterrupt(VectorNativeBRK)
	}
	c.markIRQRunning()
	return nil
}

// cop - Co-Processor Enable (software interrupt via COP vector).
func cop(c *CPU) error {
	retAddr := c.PC + 2
	if c.E {
		c.push16(retAddr)
		c.push8(c.GetP())
		c.enterInterrupt(VectorEmuCOP)
	} else {
		c.push8(c.PB)
		c.push16(retAddr)
		c.push8(c.GetP())
		c.enterInterrupt(VectorNativeCOP)
	}
	c.markIRQRunning()
	return nil
}

// rti - Return from Interrupt.
func rti(c *CPU) error {
	p := c.pop8()
	c.SetP(p)
	c.PC = c.pop16()
	if !c.E {
		// Native mode also pulls PB.
		c.PB = c.pop8()
	}
	c.pcChanged = true

	c.mu.Lock()
	c.irqRunning = false
	c.nmiRunning = false
	c.mu.Unlock()
	return nil
}

// rtl - Return from Subroutine Long.
// 65816-native: uses full 16-bit SP (no page-1 wrap between bytes).
func rtl(c *CPU) error {
	retAddr := c.pop16raw()
	c.PB = c.pop8raw()
	c.fixEmuSP()
	c.PC = retAddr + 1
	c.pcChanged = true
	return nil
}

// rts - Return from Subroutine.
func rts(c *CPU) error {
	retAddr := c.pop16()
	c.PC = retAddr + 1
	c.pcChanged = true
	return nil
}
