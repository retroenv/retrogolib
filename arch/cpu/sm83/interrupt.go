package sm83

// Interrupt vector addresses.
const (
	VectorVBlank  uint16 = 0x0040
	VectorLCDStat uint16 = 0x0048
	VectorTimer   uint16 = 0x0050
	VectorSerial  uint16 = 0x0058
	VectorJoypad  uint16 = 0x0060
)

// Interrupt bit masks.
const (
	IntVBlank  uint8 = 1 << 0
	IntLCDStat uint8 = 1 << 1
	IntTimer   uint8 = 1 << 2
	IntSerial  uint8 = 1 << 3
	IntJoypad  uint8 = 1 << 4
)

// Memory-mapped interrupt registers.
const (
	AddrIF uint16 = 0xFF0F // Interrupt Flag register
	AddrIE uint16 = 0xFFFF // Interrupt Enable register
)

// interruptVectors maps interrupt bits to their handler addresses.
var interruptVectors = [5]uint16{
	VectorVBlank,
	VectorLCDStat,
	VectorTimer,
	VectorSerial,
	VectorJoypad,
}

// HandleInterrupts processes pending interrupts.
// Called at the beginning of each Step.
// Returns true if an interrupt was serviced.
func (c *CPU) HandleInterrupts() bool {
	ie := c.memory.Read(AddrIE)
	ifReg := c.memory.Read(AddrIF)
	pending := ie & ifReg & 0x1F

	if pending == 0 {
		return false
	}

	// Any pending interrupt wakes the CPU from HALT, even if IME is disabled.
	if c.halted {
		c.halted = false
		// HALT bug: if IME is disabled and there's a pending interrupt,
		// the next instruction's first byte is read twice.
		if !c.ime {
			c.haltBug = true
			return false
		}
	}

	if !c.ime {
		return false
	}

	// Service highest priority interrupt (lowest bit).
	for i := range 5 {
		bit := uint8(1 << i)
		if pending&bit != 0 {
			c.ime = false
			c.imeDelay = false

			// Clear the interrupt flag.
			c.memory.Write(AddrIF, ifReg&^bit)

			// Push current PC and jump to vector.
			c.push16(c.PC)
			c.PC = interruptVectors[i]

			c.cycles += 5 // Interrupt dispatch takes 5 M-cycles
			return true
		}
	}

	return false
}
