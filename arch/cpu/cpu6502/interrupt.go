package cpu6502

// Interrupts contains the CPU interrupt info.
// Running fields describe handler lifecycle through RTI; they do not identify
// which instruction or interrupt caused an individual StackEvent.
type Interrupts struct {
	// NMITriggered reports whether an NMI is pending service.
	NMITriggered bool
	// NMIRunning reports whether an NMI handler is active until RTI executes.
	NMIRunning bool
	// IrqTriggered reports whether an IRQ request or input level is active.
	IrqTriggered bool
	// IrqRunning reports whether an IRQ or BRK handler is active until RTI executes.
	IrqRunning bool
}

// TriggerIrq queues one interrupt request.
// This is a no-op for the 6507 variant, which has no IRQ pin.
func (c *CPU) TriggerIrq() {
	if c.opts.variant == Variant6507 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.triggerIrq = true
}

// SetIRQ changes the IRQ input level.
// This is a no-op for the 6507 variant, which has no IRQ pin.
func (c *CPU) SetIRQ(active bool) {
	if c.opts.variant == Variant6507 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.irqLine = active
}

// TriggerNMI queues one non-maskable interrupt request.
// This is a no-op for the 6507 variant, which has no NMI pin.
func (c *CPU) TriggerNMI() {
	if c.opts.variant == Variant6507 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.triggerNmi = true
}

// CheckInterrupts services a pending interrupt at an instruction boundary.
// It returns true if it serviced an interrupt.
func (c *CPU) CheckInterrupts() bool {
	nmi, irq := c.pendingInterrupts()
	if nmi {
		c.nmi()
		return true
	}
	if irq {
		c.irq()
		return true
	}
	return false
}

// pendingInterrupts reads the interrupt inputs under the lock, because other
// goroutines can change them.
func (c *CPU) pendingInterrupts() (nmi, irq bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.stallCycles != 0 || c.jammed {
		return false, false
	}
	if c.opts.cycleHook != nil {
		return c.nmiPolled, c.irqPolled
	}
	return c.triggerNmi, (c.triggerIrq || c.irqLine) && c.Flags.I == 0
}

func (c *CPU) nmi() {
	c.mu.Lock()
	c.triggerNmi = false
	c.nmiRunning = true
	c.mu.Unlock()

	c.executeInterruptFrom(NMIAddress, InterruptNMI)
}

func (c *CPU) irq() {
	c.mu.Lock()
	c.triggerIrq = false
	c.irqRunning = true
	c.mu.Unlock()

	c.executeInterruptFrom(IrqAddress, InterruptIRQ)
}

func (c *CPU) executeInterrupt(vectorAddress uint16) {
	c.executeInterruptFrom(vectorAddress, InterruptNone)
}

func (c *CPU) executeInterruptFrom(vectorAddress uint16, source InterruptSource) {
	if c.opts.cycleHook != nil && c.opts.variant < Variant65C02 {
		c.interruptCycles(vectorAddress, source)
		return
	}
	if c.opts.stackEventHook != nil {
		c.executionPC = c.PC
		c.executionCycle = c.cycles
		// Interrupt entry does not decode an instruction. Reading PC here would
		// let an observation hook change mapped-memory state.
		c.opcode = 0
		c.interrupt = source
	}
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
