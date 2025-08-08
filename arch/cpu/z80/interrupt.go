package z80

// InterruptMode defines the Z80 interrupt modes.
type InterruptMode uint8

const (
	InterruptMode0 InterruptMode = 0 // Execute instruction on data bus (usually RST)
	InterruptMode1 InterruptMode = 1 // Jump to 0x0038
	InterruptMode2 InterruptMode = 2 // Vector table lookup using I register
)

// Note: Interrupts struct and TriggerIRQ/TriggerNMI methods are defined in cpu.go

// EnableInterrupts enables maskable interrupts (sets IFF1 and IFF2).
func (c *CPU) EnableInterrupts() {
	c.iff1 = true
	c.iff2 = true
}

// DisableInterrupts disables maskable interrupts (clears IFF1 and IFF2).
func (c *CPU) DisableInterrupts() {
	c.iff1 = false
	c.iff2 = false
}

// SetInterruptMode sets the interrupt mode (0, 1, or 2).
func (c *CPU) SetInterruptMode(mode InterruptMode) error {
	if mode > 2 {
		return ErrInvalidInterruptMode
	}
	c.im = uint8(mode)
	return nil
}

// GetInterruptMode returns the current interrupt mode.
func (c *CPU) GetInterruptMode() InterruptMode {
	return InterruptMode(c.im)
}

// InterruptsEnabled returns whether maskable interrupts are enabled.
func (c *CPU) InterruptsEnabled() bool {
	return c.iff1
}

// CheckInterrupts checks if an interrupt is triggered and executes it.
// It returns true if an interrupt was executed.
func (c *CPU) CheckInterrupts() bool {
	// Non-maskable interrupt has highest priority
	if c.triggerNmi {
		c.executeNMI()
		return true
	}

	// Maskable interrupt (only if enabled)
	if c.triggerIrq && c.iff1 {
		c.executeIRQ()
		return true
	}

	return false
}

// executeNMI handles non-maskable interrupt execution.
func (c *CPU) executeNMI() {
	c.triggerNmi = false

	// Save IFF1 to IFF2 and disable interrupts
	c.iff2 = c.iff1
	c.iff1 = false

	// Push PC to stack
	c.SP -= 2
	c.memory.WriteWord(c.SP, c.PC)

	// Jump to NMI vector
	c.PC = 0x0066
	c.cycles += 11
}

// executeIRQ handles maskable interrupt execution based on interrupt mode.
func (c *CPU) executeIRQ() {
	c.triggerIrq = false

	// Disable interrupts
	c.iff1 = false
	c.iff2 = false

	// Push PC to stack
	c.SP -= 2
	c.memory.WriteWord(c.SP, c.PC)

	switch InterruptMode(c.im) {
	case InterruptMode0:
		// Execute instruction on data bus (usually RST)
		// For simplicity, we'll execute RST 38H
		c.PC = 0x0038
		c.cycles += 13

	case InterruptMode1:
		// Jump to fixed address 0x0038
		c.PC = 0x0038
		c.cycles += 13

	case InterruptMode2:
		// Vector table lookup using I register
		vector := uint16(c.I)<<8 | uint16(c.memory.Read(0xFFFF))
		c.PC = c.memory.ReadWord(vector)
		c.cycles += 19
	}
}
