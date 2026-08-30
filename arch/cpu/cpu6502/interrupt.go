package cpu6502

// Interrupts contains the CPU interrupt info.
type Interrupts struct {
	NMITriggered bool
	NMIRunning   bool
	IrqTriggered bool
	IrqRunning   bool
}

// TriggerIrq causes an interrupt request to occur on the next cycle.
// This is a no-op for the 6507 variant (Atari 2600) which has no IRQ pin.
func (c *CPU) TriggerIrq() {
	if c.opts.variant == Variant6507 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.triggerIrq = true
}

// TriggerNMI causes a non-maskable interrupt to occur on the next cycle.
// This is a no-op for the 6507 variant (Atari 2600) which has no NMI pin.
func (c *CPU) TriggerNMI() {
	if c.opts.variant == Variant6507 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.triggerNmi = true
}

// CheckInterrupts checks if an interrupt is triggered and executes it.
// It returns true if an interrupt was executed.
func (c *CPU) CheckInterrupts() bool {
	c.mu.Lock()
	var address uint16
	switch {
	case c.triggerNmi:
		c.triggerNmi = false
		c.nmiRunning = true
		address = c.nmiAddress
	case c.triggerIrq && c.Flags.I == 0:
		c.triggerIrq = false
		c.irqRunning = true
		address = c.irqAddress
	default:
		c.mu.Unlock()
		return false
	}
	c.mu.Unlock()

	c.executeInterrupt(address)
	return true
}

func (c *CPU) executeInterrupt(funAddress uint16) {
	c.push16(c.PC)

	// Hardware interrupts push B clear and U set; PHP and BRK push B set.
	status := c.GetFlags() &^ 0b0001_0000
	status |= 0b0010_0000
	c.push(status)
	c.Flags.I = 1
	if c.opts.variant >= Variant65C02 {
		c.Flags.D = 0
	}
	c.cycles += 7
	c.PC = funAddress
}
