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
func (cpu *CPU) EnableInterrupts() {
	cpu.iff1 = true
	cpu.iff2 = true
}

// DisableInterrupts disables maskable interrupts (clears IFF1 and IFF2).
func (cpu *CPU) DisableInterrupts() {
	cpu.iff1 = false
	cpu.iff2 = false
}

// SetInterruptMode sets the interrupt mode (0, 1, or 2).
func (cpu *CPU) SetInterruptMode(mode InterruptMode) error {
	if mode > 2 {
		return ErrInvalidInterruptMode
	}
	cpu.im = uint8(mode)
	return nil
}

// GetInterruptMode returns the current interrupt mode.
func (cpu *CPU) GetInterruptMode() InterruptMode {
	return InterruptMode(cpu.im)
}

// InterruptsEnabled returns whether maskable interrupts are enabled.
func (cpu *CPU) InterruptsEnabled() bool {
	return cpu.iff1
}

// CheckInterrupts checks if an interrupt is triggered and executes it.
// It returns true if an interrupt was executed.
func (cpu *CPU) CheckInterrupts() bool {
	// Non-maskable interrupt has highest priority
	if cpu.triggerNmi {
		cpu.executeNMI()
		return true
	}

	// Maskable interrupt (only if enabled)
	if cpu.triggerIrq && cpu.iff1 {
		cpu.executeIRQ()
		return true
	}

	return false
}

// executeNMI handles non-maskable interrupt execution.
func (cpu *CPU) executeNMI() {
	cpu.triggerNmi = false

	// Save IFF1 to IFF2 and disable interrupts
	cpu.iff2 = cpu.iff1
	cpu.iff1 = false

	// Push PC to stack
	cpu.SP -= 2
	cpu.memory.WriteWord(cpu.SP, cpu.PC)

	// Jump to NMI vector
	cpu.PC = 0x0066
	cpu.cycles += 11
}

// executeIRQ handles maskable interrupt execution based on interrupt mode.
func (cpu *CPU) executeIRQ() {
	cpu.triggerIrq = false

	// Disable interrupts
	cpu.iff1 = false
	cpu.iff2 = false

	// Push PC to stack
	cpu.SP -= 2
	cpu.memory.WriteWord(cpu.SP, cpu.PC)

	switch InterruptMode(cpu.im) {
	case InterruptMode0:
		// Execute instruction on data bus (usually RST)
		// For simplicity, we'll execute RST 38H
		cpu.PC = 0x0038
		cpu.cycles += 13

	case InterruptMode1:
		// Jump to fixed address 0x0038
		cpu.PC = 0x0038
		cpu.cycles += 13

	case InterruptMode2:
		// Vector table lookup using I register
		vector := uint16(cpu.I)<<8 | uint16(cpu.memory.Read(0xFFFF))
		cpu.PC = cpu.memory.ReadWord(vector)
		cpu.cycles += 19
	}
}
