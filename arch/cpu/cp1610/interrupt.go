package cp1610

// InterruptVector is the Intellivision interrupt entry address.
const InterruptVector = 0x1004

// ResetAddress is the Intellivision reset entry address.
const ResetAddress = 0x1000

// TriggerIRQ asserts the maskable interrupt request line.
func (c *CPU) TriggerIRQ() {
	c.irqMu.Lock()
	c.irqPending = true
	c.irqMu.Unlock()
}

func (c *CPU) serviceIRQ() bool {
	c.irqMu.Lock()
	defer c.irqMu.Unlock()
	if !c.irqPending || !c.Flags.I || !c.lastInterruptible {
		return false
	}
	c.irqPending = false
	c.memory.Write(c.R[StackPointer], c.R[ProgramCounter])
	c.R[StackPointer]++
	c.R[ProgramCounter] = InterruptVector
	c.cycles += 12
	c.lastInterruptible = false
	return true
}
