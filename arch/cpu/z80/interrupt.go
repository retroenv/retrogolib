package z80

// InterruptMode defines the Z80 interrupt modes.
type InterruptMode uint8

const (
	InterruptMode0 InterruptMode = 0 // Execute instruction on data bus (usually RST)
	InterruptMode1 InterruptMode = 1 // Jump to 0x0038
	InterruptMode2 InterruptMode = 2 // Vector table lookup using I register
)

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
	if mode > InterruptMode2 {
		return ErrInvalidInterruptMode
	}
	cpu.im = mode
	return nil
}

// GetInterruptMode returns the current interrupt mode.
func (cpu *CPU) GetInterruptMode() InterruptMode {
	return cpu.im
}

// InterruptsEnabled returns whether maskable interrupts are enabled.
func (cpu *CPU) InterruptsEnabled() bool {
	return cpu.iff1
}

// CheckInterrupts accepts a pending NMI or enabled IRQ without executing an ISR instruction.
// It shares Step's interrupt behavior, including the one-instruction delay after EI.
func (cpu *CPU) CheckInterrupts() bool {
	cpu.mu.Lock()
	defer cpu.mu.Unlock()
	return cpu.handleInterrupts()
}

func (cpu *CPU) handleInterrupts() bool {
	if cpu.triggerNmi {
		cpu.triggerNmi = false
		cpu.halted = false
		cpu.iff1 = false
		// IFF2 retains the state saved before the first NMI, including nested NMIs.
		cpu.lastWasLdAIR = false
		cpu.incrementRefresh(false)
		cpu.push16(cpu.PC)
		cpu.PC, cpu.MEMPTR = 0x66, 0x66
		cpu.cycles += 11
		return true
	}
	if !cpu.triggerIrq || !cpu.iff1 || cpu.eiPending {
		return false
	}

	cpu.triggerIrq = false
	cpu.halted = false
	if cpu.lastWasLdAIR {
		cpu.Flags.P = 0
	}
	cpu.lastWasLdAIR = false
	cpu.iff1, cpu.iff2 = false, false
	cpu.incrementRefresh(false)

	// Sample the vector before stack writes, which may overlap a mapped device.
	vector := uint16(0x38)
	switch cpu.im {
	case InterruptMode0:
		data := cpu.bus.IRQData()
		if data&0xC7 == 0xC7 {
			vector = uint16(data & 0x38)
		}
		// Non-RST IM0 opcodes retain the legacy RST 38h fallback.
	case InterruptMode2:
		vector = uint16(cpu.I)<<8 | uint16(cpu.bus.IRQData())
	}
	cpu.push16(cpu.PC)
	cpu.cycles += 13
	if cpu.im == InterruptMode2 {
		vector = cpu.read16(vector)
		cpu.cycles += 6
	}
	cpu.PC, cpu.MEMPTR = vector, vector
	return true
}
