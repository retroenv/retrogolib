package m68000

// Exception vector numbers.
const (
	VectorResetSSP    = 0  // Reset: initial SSP
	VectorResetPC     = 1  // Reset: initial PC
	VectorBusError    = 2  // Bus error
	VectorAddressErr  = 3  // Address error
	VectorIllegal     = 4  // Illegal instruction
	VectorDivZero     = 5  // Divide by zero
	VectorCHK         = 6  // CHK instruction
	VectorTRAPV       = 7  // TRAPV instruction
	VectorPrivilege   = 8  // Privilege violation
	VectorTrace       = 9  // Trace
	VectorLineA       = 10 // Line 1010 emulator
	VectorLineF       = 11 // Line 1111 emulator
	VectorSpurious    = 24 // Spurious interrupt
	VectorAutoVector1 = 25 // Autovector level 1
	VectorAutoVector2 = 26 // Autovector level 2
	VectorAutoVector3 = 27 // Autovector level 3
	VectorAutoVector4 = 28 // Autovector level 4
	VectorAutoVector5 = 29 // Autovector level 5
	VectorAutoVector6 = 30 // Autovector level 6
	VectorAutoVector7 = 31 // Autovector level 7
	VectorTrap0       = 32 // TRAP #0
	VectorTrap15      = 47 // TRAP #15
)

// processException processes an exception with the given vector number.
// It saves the current state and loads the new PC from the vector table.
func (c *CPU) processException(vector int) error {
	// Save current SR.
	oldSR := c.GetSR()

	// Enter supervisor mode and clear trace.
	c.sr |= MaskSupervisor
	c.sr &^= MaskTrace

	// If switching from user to supervisor, swap stack pointers.
	if oldSR&MaskSupervisor == 0 {
		c.USP = c.sp
		c.sp = c.SSP
	}

	// Push PC and SR onto the supervisor stack.
	c.push32(c.PC)
	c.push16(oldSR)

	// Load new PC from vector table.
	vectorAddr := uint32(vector) * 4
	c.PC = c.bus.ReadLong(vectorAddr)

	c.stopped = false

	return nil
}

// processInterruptException processes an interrupt exception for the given level.
func (c *CPU) processInterruptException(level uint8) {
	// Save current SR.
	oldSR := c.GetSR()

	// Enter supervisor mode, clear trace, set interrupt mask.
	c.sr |= MaskSupervisor
	c.sr &^= MaskTrace
	c.sr = (c.sr & ^uint16(MaskIPM)) | (uint16(level) << FlagIPM0)

	// If switching from user to supervisor, swap stack pointers.
	if oldSR&MaskSupervisor == 0 {
		c.USP = c.sp
		c.sp = c.SSP
	}

	// Push PC and SR.
	c.push32(c.PC)
	c.push16(oldSR)

	// Get vector from bus.
	vector := c.bus.IRQAcknowledge(level)

	// Load new PC from vector table.
	vectorAddr := vector * 4
	c.PC = c.bus.ReadLong(vectorAddr)

	c.stopped = false
	c.cycles += 44
}

// TriggerIRQ triggers a maskable interrupt at the given level (1-7).
func (c *CPU) TriggerIRQ(_ uint8) {
	// IRQ level is read from bus on each step.
}

// checkInterrupts checks for pending interrupts and processes them.
// Returns true if an interrupt was processed.
func (c *CPU) checkInterrupts() bool {
	level := c.bus.IRQLevel()
	if level == 0 {
		return false
	}

	mask := c.InterruptMask()

	// Level 7 is non-maskable. Other levels must be higher than mask.
	if level < 7 && level <= mask {
		return false
	}

	c.processInterruptException(level)
	return true
}
