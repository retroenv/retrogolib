package cpu6809

// System instructions: RTS, RTI, SWI, CWAI, SYNC.

// pushEntireState pushes the entire CPU state with the supplied return address.
// Order: PC, U, Y, X, DP, B, A, CC (last pushed = first on stack).
func (c *CPU) pushEntireState(returnPC uint16) {
	c.pushS16(returnPC)
	c.pushS16(c.U)
	c.pushS16(c.Y)
	c.pushS16(c.X)
	c.pushS8(c.DP)
	c.pushS8(c.B)
	c.pushS8(c.A)
	c.pushS8(c.GetCC())
}

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
		// The opcode table accounts for the six-cycle short-frame path.
		c.cycles += 9
		c.A = c.popS8()
		c.B = c.popS8()
		c.DP = c.popS8()
		c.X = c.popS16()
		c.Y = c.popS16()
		c.U = c.popS16()
	}

	c.PC = c.popS16()
	c.pcChanged = true

	return nil
}

// swiFn - Software Interrupt 1.
// Always saves entire state (sets E=1).
// Disables both IRQ and FIRQ.
func swiFn(c *CPU) error {
	c.Flags.E = 1
	c.pushEntireState(c.nextPC)
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
	c.pushEntireState(c.nextPC)
	c.PC = c.memory.ReadVector(VectorSWI2)
	c.pcChanged = true
	return nil
}

// swi3Fn - Software Interrupt 3.
// Saves entire state (sets E=1).
// Does NOT disable interrupts.
func swi3Fn(c *CPU) error {
	c.Flags.E = 1
	c.pushEntireState(c.nextPC)
	c.PC = c.memory.ReadVector(VectorSWI3)
	c.pcChanged = true
	return nil
}

// cwaiFn - AND CC then Wait for Interrupt.
// ANDs the CC with the immediate value, sets E=1, pushes entire state,
// then waits for an interrupt.
func cwaiFn(c *CPU, param any) error {
	value, err := immediate8Param(param)
	if err != nil {
		return err
	}

	mask := uint8(value)
	c.SetCC(c.GetCC() & mask)
	c.Flags.E = 1
	c.pushEntireState(c.nextPC)
	c.setWaitMode(waitCWAI)

	return nil
}

func syncFn(c *CPU) error {
	c.setWaitMode(waitSync)

	return nil
}
