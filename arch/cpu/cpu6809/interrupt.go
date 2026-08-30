package cpu6809

// TriggerNMI requests a non-maskable interrupt on the next Step.
func (c *CPU) TriggerNMI() {
	c.requestInterrupt(&c.triggerNMI)
}

// TriggerIRQ requests an IRQ on the next Step when the I flag permits it.
func (c *CPU) TriggerIRQ() {
	c.requestInterrupt(&c.triggerIRQ)
}

// TriggerFIRQ requests a fast IRQ on the next Step when the F flag permits it.
func (c *CPU) TriggerFIRQ() {
	c.requestInterrupt(&c.triggerFIRQ)
}

// CheckInterrupts services the highest-priority pending interrupt.
func (c *CPU) CheckInterrupts() bool {
	kind, stackState := c.takePendingInterrupt()
	switch kind {
	case interruptNMI:
		c.handleNMI(stackState)
	case interruptFIRQ:
		c.handleFIRQ(stackState)
	case interruptIRQ:
		c.handleIRQ(stackState)
	default:
		return false
	}

	return true
}

func (c *CPU) requestInterrupt(pending *bool) {
	c.interruptMu.Lock()
	*pending = true
	// Any interrupt releases SYNC, even when its mask prevents servicing.
	if c.waitMode == waitSync {
		c.waitMode = waitNone
	}
	c.interruptMu.Unlock()
}

func (c *CPU) takePendingInterrupt() (interruptKind, bool) {
	c.interruptMu.Lock()
	defer c.interruptMu.Unlock()

	if c.waitMode == waitSync && (c.triggerNMI || c.triggerFIRQ || c.triggerIRQ) {
		c.waitMode = waitNone
	}

	kind := interruptNone
	switch {
	case c.triggerNMI:
		c.triggerNMI = false
		kind = interruptNMI
	case c.triggerFIRQ && c.Flags.F == 0:
		c.triggerFIRQ = false
		kind = interruptFIRQ
	case c.triggerIRQ && c.Flags.I == 0:
		c.triggerIRQ = false
		kind = interruptIRQ
	}

	if kind == interruptNone {
		return interruptNone, false
	}

	// CWAI stacked the entire state before sleeping, so interrupt acceptance
	// only needs to load the vector and update masks.
	stackState := c.waitMode != waitCWAI
	c.waitMode = waitNone

	return kind, stackState
}

func (c *CPU) handleNMI(stackState bool) {
	if stackState {
		c.Flags.E = 1
		c.pushEntireState(c.PC)
	}
	c.Flags.I = 1
	c.Flags.F = 1
	c.PC = c.memory.ReadVector(VectorNMI)
}

func (c *CPU) handleFIRQ(stackState bool) {
	if stackState {
		c.Flags.E = 0
		c.pushS16(c.PC)
		c.pushS8(c.GetCC())
	}
	c.Flags.I = 1
	c.Flags.F = 1
	c.PC = c.memory.ReadVector(VectorFIRQ)
}

func (c *CPU) handleIRQ(stackState bool) {
	if stackState {
		c.Flags.E = 1
		c.pushEntireState(c.PC)
	}
	c.Flags.I = 1
	c.PC = c.memory.ReadVector(VectorIRQ)
}
