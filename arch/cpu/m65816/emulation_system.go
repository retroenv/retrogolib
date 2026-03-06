package m65816

// System, processor status, and miscellaneous instructions.

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
	mask := uint8(params[0].(Immediate8))
	p := c.GetP() &^ mask
	c.SetP(p)
	return nil
}

// sep - Set Processor Status Bits (set bits specified by immediate mask).
func sep(c *CPU, params ...any) error {
	mask := uint8(params[0].(Immediate8))
	p := c.GetP() | mask
	c.SetP(p)
	return nil
}

// xce - Exchange Carry and Emulation flags.
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
		// Entering emulation mode: force M=1, X=1, high byte of SP = $01
		c.Flags.M = 1
		c.Flags.X = 1
		c.SP = 0x0100 | (c.SP & 0x00FF)
	} else {
		// Entering native mode: M=1, X=1 (stays until REP changes them)
		c.Flags.M = 1
		c.Flags.X = 1
	}
	return nil
}

// xba - Exchange B and A (swap high and low bytes of accumulator C).
func xba(c *CPU) error {
	lo := uint8(c.C)
	hi := uint8(c.C >> 8)
	c.C = uint16(lo)<<8 | uint16(hi)
	// N and Z are set based on the new low byte (new A)
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
func brk(c *CPU) error {
	// BRK is 2 bytes; push PC+2 (address after signature byte)
	retAddr := c.PC + 2
	if c.E {
		c.push16(retAddr)
		p := c.GetP() | MaskBreak // B flag set when pushed in emulation mode
		c.push8(p)
		c.Flags.I = 1
		c.Flags.D = 0 // 65C02 behavior: clear D on interrupt
		vec := c.memory.ReadVector(VectorEmuIRQ)
		c.PC = vec
	} else {
		c.push8(c.PB)
		c.push16(retAddr)
		c.push8(c.GetP())
		c.Flags.I = 1
		c.Flags.D = 0
		vec := c.memory.ReadVector(VectorNativeBRK)
		c.PB = 0
		c.PC = vec
	}

	c.mu.Lock()
	c.irqRunning = true
	c.mu.Unlock()
	return nil
}

// cop - Co-Processor Enable (software interrupt via COP vector).
func cop(c *CPU) error {
	retAddr := c.PC + 2
	if c.E {
		c.push16(retAddr)
		c.push8(c.GetP())
		c.Flags.I = 1
		c.Flags.D = 0
		vec := c.memory.ReadVector(VectorEmuCOP)
		c.PC = vec
	} else {
		c.push8(c.PB)
		c.push16(retAddr)
		c.push8(c.GetP())
		c.Flags.I = 1
		c.Flags.D = 0
		vec := c.memory.ReadVector(VectorNativeCOP)
		c.PB = 0
		c.PC = vec
	}
	return nil
}

// rti - Return from Interrupt.
func rti(c *CPU) error {
	p := c.pop8()
	c.SetP(p)
	c.PC = c.pop16()
	if !c.E {
		// Native mode: also pull PB
		c.PB = c.pop8()
	}

	c.mu.Lock()
	c.irqRunning = false
	c.nmiRunning = false
	c.mu.Unlock()
	return nil
}

// rtl - Return from Subroutine Long.
func rtl(c *CPU) error {
	retAddr := c.pop16()
	c.PB = c.pop8()
	c.PC = retAddr + 1
	return nil
}

// rts - Return from Subroutine.
func rts(c *CPU) error {
	retAddr := c.pop16()
	c.PC = retAddr + 1
	return nil
}
