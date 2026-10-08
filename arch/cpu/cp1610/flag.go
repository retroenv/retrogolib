package cp1610

// Flags holds the status bits of the CP1610.
type Flags struct {
	S bool // Sign
	Z bool // Zero
	O bool // Overflow
	C bool // Carry
	I bool // Interrupt enable
	D bool // Double byte data for the next instruction
}

func (c *CPU) setSZ(value uint16) {
	c.Flags.S = value&0x8000 != 0
	c.Flags.Z = value == 0
}

func (c *CPU) add(a, b uint16) uint16 {
	wide := uint32(a) + uint32(b)
	result := uint16(wide)
	c.setSZ(result)
	c.Flags.C = wide > 0xFFFF
	c.Flags.O = (^(a ^ b) & (a ^ result) & 0x8000) != 0
	return result
}

func (c *CPU) subtract(a, b uint16) uint16 {
	result := b - a
	c.setSZ(result)
	c.Flags.C = b >= a
	c.Flags.O = ((b ^ a) & (b ^ result) & 0x8000) != 0
	return result
}

func (c *CPU) statusWord() uint16 {
	var word uint16
	if c.Flags.S {
		word |= 0x80
	}
	if c.Flags.Z {
		word |= 0x40
	}
	if c.Flags.O {
		word |= 0x20
	}
	if c.Flags.C {
		word |= 0x10
	}
	return word | word<<8
}

func (c *CPU) setStatusWord(word uint16) {
	c.Flags.S = word&0x80 != 0
	c.Flags.Z = word&0x40 != 0
	c.Flags.O = word&0x20 != 0
	c.Flags.C = word&0x10 != 0
}
