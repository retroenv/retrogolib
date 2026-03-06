package m65816

// Interrupt state and external interrupt triggering.

// TriggerNMI causes a non-maskable interrupt on the next Step.
func (c *CPU) TriggerNMI() {
	c.mu.Lock()
	c.triggerNMI = true
	c.waiting = false // WAI exits on interrupt
	c.mu.Unlock()
}

// TriggerIRQ causes an IRQ interrupt request on the next Step (if I flag permits).
func (c *CPU) TriggerIRQ() {
	c.mu.Lock()
	c.triggerIRQ = true
	c.waiting = false
	c.mu.Unlock()
}

// CheckInterrupts processes any pending interrupts.
// Returns true if an interrupt was handled.
func (c *CPU) CheckInterrupts() bool {
	if c.triggerNMI {
		c.handleNMI()
		return true
	}
	if c.triggerIRQ && c.Flags.I == 0 {
		c.handleIRQ()
		return true
	}
	return false
}

func (c *CPU) handleNMI() {
	c.mu.Lock()
	c.triggerNMI = false
	c.nmiRunning = true
	c.mu.Unlock()

	c.executeInterrupt(VectorNativeNMI, VectorEmuNMI)
}

func (c *CPU) handleIRQ() {
	c.mu.Lock()
	c.triggerIRQ = false
	c.irqRunning = true
	c.mu.Unlock()

	c.executeInterrupt(VectorNativeIRQ, VectorEmuIRQ)
}

// executeInterrupt pushes the CPU context and loads the interrupt vector.
func (c *CPU) executeInterrupt(nativeVec, emuVec uint32) {
	if c.E {
		// Emulation mode: 6502-style interrupt sequence
		c.push16(c.PC)
		p := c.GetP() &^ MaskBreak // B=0 for hardware interrupts
		c.push8(p)
		c.Flags.I = 1
		c.Flags.D = 0
		c.cycles += 7
		c.PC = c.memory.ReadVector(emuVec)
	} else {
		// Native mode: push PB, PC, P
		c.push8(c.PB)
		c.push16(c.PC)
		c.push8(c.GetP())
		c.Flags.I = 1
		c.Flags.D = 0
		c.cycles += 8
		c.PB = 0
		c.PC = c.memory.ReadVector(nativeVec)
	}
}
