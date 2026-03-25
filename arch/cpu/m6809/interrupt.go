package m6809

// Interrupt state and external interrupt triggering.

// TriggerNMI causes a non-maskable interrupt on the next Step.
func (c *CPU) TriggerNMI() {
	c.mu.Lock()
	c.triggerNMI = true
	c.waiting = false // SYNC exits on interrupt
	c.mu.Unlock()
}

// TriggerIRQ causes an IRQ interrupt request on the next Step (if I flag permits).
func (c *CPU) TriggerIRQ() {
	c.mu.Lock()
	c.triggerIRQ = true
	c.waiting = false
	c.mu.Unlock()
}

// TriggerFIRQ causes a fast IRQ interrupt request on the next Step (if F flag permits).
func (c *CPU) TriggerFIRQ() {
	c.mu.Lock()
	c.triggerFIRQ = true
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
	if c.triggerFIRQ && c.Flags.F == 0 {
		c.handleFIRQ()
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

	// NMI saves entire state (E=1)
	c.Flags.E = 1
	c.pushEntireState()
	c.Flags.I = 1
	c.Flags.F = 1
	c.PC = c.memory.ReadVector(VectorNMI)
}

func (c *CPU) handleFIRQ() {
	c.mu.Lock()
	c.triggerFIRQ = false
	c.irqRunning = true
	c.mu.Unlock()

	// FIRQ saves only CC and PC (E=0)
	c.Flags.E = 0
	c.pushS16(c.PC)
	c.pushS8(c.GetCC())
	c.Flags.I = 1
	c.Flags.F = 1
	c.PC = c.memory.ReadVector(VectorFIRQ)
}

func (c *CPU) handleIRQ() {
	c.mu.Lock()
	c.triggerIRQ = false
	c.irqRunning = true
	c.mu.Unlock()

	// IRQ saves entire state (E=1)
	c.Flags.E = 1
	c.pushEntireState()
	c.Flags.I = 1
	c.PC = c.memory.ReadVector(VectorIRQ)
}
