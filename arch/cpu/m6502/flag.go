package m6502

// Flags contains the status flags of the CPU.
// Bit No.   7   6   5   4   3   2   1   0
// Flag      S   V       B   D   I   Z   C
type Flags struct {
	C uint8 // carry flag
	Z uint8 // zero flag
	I uint8 // interrupt disable flag
	D uint8 // decimal mode flag
	B uint8 // break command flag
	U uint8 // unused flag
	V uint8 // overflow flag
	N uint8 // negative flag
}

// GetFlags returns the current state of flags as byte.
func (c *CPU) GetFlags() uint8 {
	return c.Flags.C |
		(c.Flags.Z << 1) |
		(c.Flags.I << 2) |
		(c.Flags.D << 3) |
		(c.Flags.B << 4) |
		(c.Flags.U << 5) |
		(c.Flags.V << 6) |
		(c.Flags.N << 7)
}

// setFlags sets the flags from the given byte.
func (c *CPU) setFlags(flags uint8) {
	c.Flags.C = (flags >> 0) & 1
	c.Flags.Z = (flags >> 1) & 1
	c.Flags.I = (flags >> 2) & 1
	c.Flags.D = (flags >> 3) & 1
	c.Flags.B = (flags >> 4) & 1
	c.Flags.U = (flags >> 5) & 1
	c.Flags.V = (flags >> 6) & 1
	c.Flags.N = (flags >> 7) & 1
}

// setZ - set the zero flag if the argument is zero.
func (c *CPU) setZ(value uint8) {
	setFlag(&c.Flags.Z, value == 0)
}

// setN - set the negative flag if the argument is negative (high bit is set).
func (c *CPU) setN(value uint8) {
	setFlag(&c.Flags.N, value&0x80 != 0)
}

// setV - set the overflow flag.
func (c *CPU) setV(set bool) {
	setFlag(&c.Flags.V, set)
}

// setZN - set the zero and negative flags.
func (c *CPU) setZN(value uint8) {
	c.setZ(value)
	c.setN(value)
}

// compare - compare two values and set the zero and negative flags.
func (c *CPU) compare(a, b byte) {
	c.setZN(a - b)
	setFlag(&c.Flags.C, a >= b)
}

// setFlag sets a flag to 1 if condition is true, 0 otherwise.
func setFlag(flag *uint8, condition bool) {
	if condition {
		*flag = 1
	} else {
		*flag = 0
	}
}
