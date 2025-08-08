package z80

// Flags contains the status flags of the CPU.
// Bit No.   7   6   5   4   3   2   1   0
// Flag      S   Z   Y   H   X   P   N   C
type Flags struct {
	C uint8 // carry flag
	N uint8 // add/subtract flag (used for BCD operations)
	P uint8 // parity/overflow flag
	X uint8 // bit 3 of last result (undocumented flag)
	H uint8 // half carry flag
	Y uint8 // bit 5 of last result (undocumented flag)
	Z uint8 // zero flag
	S uint8 // sign flag
}

// GetFlags returns the current state of flags as byte.
func (c *CPU) GetFlags() uint8 {
	var f byte
	f |= c.Flags.C << 0
	f |= c.Flags.N << 1
	f |= c.Flags.P << 2
	f |= c.Flags.X << 3
	f |= c.Flags.H << 4
	f |= c.Flags.Y << 5
	f |= c.Flags.Z << 6
	f |= c.Flags.S << 7
	return f
}

// setFlags sets the flags from the given byte.
func (c *CPU) setFlags(flags uint8) {
	c.Flags.C = (flags >> 0) & 1
	c.Flags.N = (flags >> 1) & 1
	c.Flags.P = (flags >> 2) & 1
	c.Flags.X = (flags >> 3) & 1
	c.Flags.H = (flags >> 4) & 1
	c.Flags.Y = (flags >> 5) & 1
	c.Flags.Z = (flags >> 6) & 1
	c.Flags.S = (flags >> 7) & 1
}

// setZ - set the zero flag if the argument is zero.
func (c *CPU) setZ(value uint8) {
	setFlag(&c.Flags.Z, value == 0)
}

// setS - set the sign flag if the argument is negative (high bit is set).
func (c *CPU) setS(value uint8) {
	setFlag(&c.Flags.S, value&0x80 != 0)
}

// setP - set the parity flag based on the parity of the argument.
func (c *CPU) setP(value uint8) {
	// Count number of 1 bits
	count := 0
	for i := range 8 {
		if value&(1<<i) != 0 {
			count++
		}
	}
	setFlag(&c.Flags.P, count%2 == 0) // even parity
}

// setPOverflow - set the parity/overflow flag with a boolean value.
func (c *CPU) setPOverflow(set bool) {
	setFlag(&c.Flags.P, set)
}

// setH - set the half carry flag.
func (c *CPU) setH(set bool) {
	setFlag(&c.Flags.H, set)
}

// setN - set the add/subtract flag.
func (c *CPU) setN(set bool) {
	setFlag(&c.Flags.N, set)
}

// setC - set the carry flag.
func (c *CPU) setC(set bool) {
	setFlag(&c.Flags.C, set)
}

// setSZP - set the sign, zero, and parity flags.
func (c *CPU) setSZP(value uint8) {
	c.setS(value)
	c.setZ(value)
	c.setP(value)
}

// setSZ - set the sign and zero flags.
func (c *CPU) setSZ(value uint8) {
	c.setS(value)
	c.setZ(value)
}

// setFlag sets a flag to 1 if condition is true, 0 otherwise.
func setFlag(flag *uint8, condition bool) {
	if condition {
		*flag = 1
	} else {
		*flag = 0
	}
}
