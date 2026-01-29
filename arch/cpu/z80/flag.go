package z80

import "math/bits"

// Flags contains the status flags of the CPU.
// Standard Z80 flag register layout:
// Bit No.   7   6   5   4   3   2   1   0
// Flag      S   Z   Y   H   X   P   N   C
//
// S (Sign): Set if result is negative (bit 7 set)
// Z (Zero): Set if result is zero
// Y (Bit 5): Copy of bit 5 of result (undocumented)
// H (Half Carry): Set if carry from bit 3 to bit 4
// X (Bit 3): Copy of bit 3 of result (undocumented)
// P (Parity/Overflow): Parity for logical ops, overflow for arithmetic
// N (Add/Subtract): Set for subtract operations, cleared for add
// C (Carry): Set if carry out of bit 7
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
	return c.Flags.C |
		c.Flags.N<<1 |
		c.Flags.P<<2 |
		c.Flags.X<<3 |
		c.Flags.H<<4 |
		c.Flags.Y<<5 |
		c.Flags.Z<<6 |
		c.Flags.S<<7
}

// setZ updates zero flag based on arithmetic/logical result.
func (c *CPU) setZ(value uint8) {
	setFlag(&c.Flags.Z, value == 0)
}

// setS updates sign flag based on MSB of result (two's complement).
func (c *CPU) setS(value uint8) {
	setFlag(&c.Flags.S, value&0x80 != 0)
}

// setP calculates and sets parity flag (even parity = 1).
func (c *CPU) setP(value uint8) {
	count := bits.OnesCount8(value)
	setFlag(&c.Flags.P, count%2 == 0) // even parity
}

// setPOverflow directly sets P/V flag for overflow conditions.
func (c *CPU) setPOverflow(set bool) {
	setFlag(&c.Flags.P, set)
}

// setH updates half carry flag for BCD and carry detection.
func (c *CPU) setH(set bool) {
	setFlag(&c.Flags.H, set)
}

// setN indicates operation type for BCD correction (1=subtract, 0=add).
func (c *CPU) setN(set bool) {
	setFlag(&c.Flags.N, set)
}

// setC updates carry flag for arithmetic overflow conditions.
func (c *CPU) setC(set bool) {
	setFlag(&c.Flags.C, set)
}

// setSZP updates S/Z/P flags and undocumented X/Y flags from result.
func (c *CPU) setSZP(value uint8) {
	c.setS(value)
	c.setZ(value)
	c.setP(value)
	c.setXY(value) // Set undocumented X and Y flags
}

// setXY copies bits 3 and 5 to undocumented flags (Z80 quirk).
func (c *CPU) setXY(value uint8) {
	c.Flags.X = (value >> 3) & 1 // bit 3
	c.Flags.Y = (value >> 5) & 1 // bit 5
}

// setSZ updates S/Z flags and undocumented X/Y flags from result.
func (c *CPU) setSZ(value uint8) {
	c.setS(value)
	c.setZ(value)
	c.setXY(value) // Set undocumented X and Y flags
}

// setFlags restores complete flag register state from byte value.
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

// setFlag helper converts boolean to Z80 flag bit value.
func setFlag(flag *uint8, condition bool) {
	if condition {
		*flag = 1
	} else {
		*flag = 0
	}
}
