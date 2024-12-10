package m6502

// Interrupts contains the CPU interrupt info.
type Interrupts struct {
	NMITriggered bool
	NMIRunning   bool
	IrqTriggered bool
	IrqRunning   bool
}

// TriggerIrq causes an interrupt request to occur on the next cycle.
func (c *CPU) TriggerIrq() {
	c.triggerIrq = true
}

// TriggerNMI causes a non-maskable interrupt to occur on the next cycle.
func (c *CPU) TriggerNMI() {
	c.triggerNmi = true
}

// CheckInterrupts checks if an interrupt is triggered and executes it.
// It returns true if an interrupt was executed.
func (c *CPU) CheckInterrupts() bool {
	if c.triggerNmi {
		c.nmi()
		return true
	}
	if c.triggerIrq {
		c.irq()
		return true
	}
	return false
}

func (c *CPU) nmi() {
	c.mu.Lock()
	c.triggerNmi = false
	c.nmiRunning = true
	c.mu.Unlock()

	c.executeInterrupt(c.nmiAddress)
}

func (c *CPU) irq() {
	c.mu.Lock()
	c.triggerIrq = false
	c.irqRunning = true
	c.mu.Unlock()

	c.executeInterrupt(c.irqAddress)
}

func (c *CPU) executeInterrupt(funAddress uint16) {
	c.push16(c.PC)
	php(c)

	if funAddress != 0 {
		c.Flags.I = 1
		c.cycles += 7
		c.PC = funAddress
	}
}
