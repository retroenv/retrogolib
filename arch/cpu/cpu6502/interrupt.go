package cpu6502

// Interrupts contains the CPU interrupt info.
type Interrupts struct {
	NMITriggered bool
	NMIRunning   bool
	IrqTriggered bool
	IrqRunning   bool
}

// TriggerIrq queues one interrupt request.
// This is a no-op for the 6507 variant, which has no IRQ pin.
func (c *CPU) TriggerIrq() {
	if c.opts.variant == Variant6507 {
		return
	}
	c.triggerIrq = true
}

// SetIRQ changes the IRQ input level.
// This is a no-op for the 6507 variant, which has no IRQ pin.
func (c *CPU) SetIRQ(active bool) {
	if c.opts.variant == Variant6507 {
		return
	}
	c.irqLine = active
}

// TriggerNMI queues one non-maskable interrupt request.
// This is a no-op for the 6507 variant, which has no NMI pin.
func (c *CPU) TriggerNMI() {
	if c.opts.variant == Variant6507 {
		return
	}
	c.triggerNmi = true
}

// CheckInterrupts services a pending interrupt at an instruction boundary.
// It returns true if it serviced an interrupt.
func (c *CPU) CheckInterrupts() bool {
	if c.stallCycles != 0 {
		return false
	}
	if c.triggerNmi {
		c.nmi()
		return true
	}
	if (c.triggerIrq || c.irqLine) && c.Flags.I == 0 {
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

	c.executeInterrupt(NMIAddress)
}

func (c *CPU) irq() {
	c.mu.Lock()
	c.triggerIrq = false
	c.irqRunning = true
	c.mu.Unlock()

	c.executeInterrupt(IrqAddress)
}

func (c *CPU) executeInterrupt(vectorAddress uint16) {
	c.push16(c.PC)
	// Hardware interrupts put a zero in the stacked B-bit position.
	flags := c.GetFlags()&^0b0001_0000 | 0b0010_0000
	c.push(flags)

	c.Flags.I = 1
	// The 65C02 clears D after it pushes the status byte.
	if c.opts.variant >= Variant65C02 {
		c.Flags.D = 0
	}
	c.cycles += 7
	// Read the vector now because mapped memory can change between interrupts.
	c.PC = c.memory.ReadWord(vectorAddress)
}
