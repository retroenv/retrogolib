package sm83

// Flags contains the status flags of the SM83 CPU.
// SM83 flag register layout (upper nibble only, lower nibble always 0):
// Bit No.   7   6   5   4   3   2   1   0
// Flag      Z   N   H   C   0   0   0   0
//
// Z (Zero): Set if result is zero
// N (Subtract): Set for subtract operations (used for BCD)
// H (Half Carry): Set if carry from bit 3 to bit 4
// C (Carry): Set if carry out of bit 7
type Flags struct {
	C uint8 // carry flag (bit 4)
	H uint8 // half carry flag (bit 5)
	N uint8 // subtract flag (bit 6)
	Z uint8 // zero flag (bit 7)
}

// GetFlags returns the current state of flags as a byte.
// Lower nibble is always 0 on SM83.
func (c *CPU) GetFlags() uint8 {
	return c.Flags.C<<4 |
		c.Flags.H<<5 |
		c.Flags.N<<6 |
		c.Flags.Z<<7
}

// setZ updates zero flag based on result.
func (c *CPU) setZ(value uint8) {
	setFlag(&c.Flags.Z, value == 0)
}

// setH updates half carry flag.
func (c *CPU) setH(set bool) {
	setFlag(&c.Flags.H, set)
}

// setN indicates operation type for BCD correction (1=subtract, 0=add).
func (c *CPU) setN(set bool) {
	setFlag(&c.Flags.N, set)
}

// setC updates carry flag.
func (c *CPU) setC(set bool) {
	setFlag(&c.Flags.C, set)
}

// setFlags restores complete flag register state from byte value.
func (c *CPU) setFlags(flags uint8) {
	c.Flags.C = (flags >> 4) & 1
	c.Flags.H = (flags >> 5) & 1
	c.Flags.N = (flags >> 6) & 1
	c.Flags.Z = (flags >> 7) & 1
}

// setFlag helper converts boolean to flag bit value.
func setFlag(flag *uint8, condition bool) {
	if condition {
		*flag = 1
	} else {
		*flag = 0
	}
}
