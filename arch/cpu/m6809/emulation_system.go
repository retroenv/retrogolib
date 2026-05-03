package m6809

// System instructions: RTS, RTI, SWI, CWAI, SYNC.

// pushEntireState pushes the entire CPU state onto the system stack.
// Order: PC, U, Y, X, DP, B, A, CC (last pushed = first on stack).
func (c *CPU) pushEntireState() {
	c.pushS16(c.PC)
	c.pushS16(c.U)
	c.pushS16(c.Y)
	c.pushS16(c.X)
	c.pushS8(c.DP)
	c.pushS8(c.B)
	c.pushS8(c.A)
	c.pushS8(c.GetCC())
}

// rtsFn - Return from Subroutine.
func rtsFn(c *CPU) error {
	c.PC = c.popS16()
	c.pcChanged = true
	return nil
}

// rtiFn - Return from Interrupt.
// If E flag was set (entire state saved), restore all registers.
// If E flag was clear (FIRQ), restore only CC and PC.
func rtiFn(c *CPU) error {
	cc := c.popS8()
	c.SetCC(cc)

	if c.Flags.E != 0 {
		// Entire state was saved
		c.A = c.popS8()
		c.B = c.popS8()
		c.DP = c.popS8()
		c.X = c.popS16()
		c.Y = c.popS16()
		c.U = c.popS16()
	}

	c.PC = c.popS16()
	c.pcChanged = true

	c.mu.Lock()
	c.irqRunning = false
	c.nmiRunning = false
	c.mu.Unlock()
	return nil
}

// swiFn - Software Interrupt 1.
// Always saves entire state (sets E=1).
// Disables both IRQ and FIRQ.
func swiFn(c *CPU) error {
	c.Flags.E = 1
	c.pushEntireState()
	c.Flags.I = 1
	c.Flags.F = 1
	c.PC = c.memory.ReadVector(VectorSWI)
	c.pcChanged = true
	return nil
}

// swi2Fn - Software Interrupt 2.
// Saves entire state (sets E=1).
// Does NOT disable interrupts.
func swi2Fn(c *CPU) error {
	c.Flags.E = 1
	c.pushEntireState()
	c.PC = c.memory.ReadVector(VectorSWI2)
	c.pcChanged = true
	return nil
}

// swi3Fn - Software Interrupt 3.
// Saves entire state (sets E=1).
// Does NOT disable interrupts.
func swi3Fn(c *CPU) error {
	c.Flags.E = 1
	c.pushEntireState()
	c.PC = c.memory.ReadVector(VectorSWI3)
	c.pcChanged = true
	return nil
}

// cwaiFn - AND CC then Wait for Interrupt.
// ANDs the CC with the immediate value, sets E=1, pushes entire state,
// then waits for an interrupt.
func cwaiFn(c *CPU, params ...any) error {
	mask := uint8(params[0].(Immediate8))
	c.SetCC(c.GetCC() & mask)
	c.Flags.E = 1
	c.pushEntireState()
	c.waiting = true
	return nil
}

// syncFn - Synchronize with Interrupt.
func syncFn(c *CPU) error {
	c.waiting = true
	return nil
}
