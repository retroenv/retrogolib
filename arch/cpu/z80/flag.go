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
func (cpu *CPU) GetFlags() uint8 {
	return cpu.Flags.C |
		cpu.Flags.N<<1 |
		cpu.Flags.P<<2 |
		cpu.Flags.X<<3 |
		cpu.Flags.H<<4 |
		cpu.Flags.Y<<5 |
		cpu.Flags.Z<<6 |
		cpu.Flags.S<<7
}

// setZ - set the zero flag if the argument is zero.
func (cpu *CPU) setZ(value uint8) {
	setFlag(&cpu.Flags.Z, value == 0)
}

// setS - set the sign flag if the argument is negative (high bit is set).
func (cpu *CPU) setS(value uint8) {
	setFlag(&cpu.Flags.S, value&0x80 != 0)
}

// setP - set the parity flag based on the parity of the argument.
func (cpu *CPU) setP(value uint8) {
	count := bits.OnesCount8(value)
	setFlag(&cpu.Flags.P, count%2 == 0) // even parity
}

// setPOverflow - set the parity/overflow flag with a boolean value.
func (cpu *CPU) setPOverflow(set bool) {
	setFlag(&cpu.Flags.P, set)
}

// setH - set the half carry flag.
func (cpu *CPU) setH(set bool) {
	setFlag(&cpu.Flags.H, set)
}

// setN - set the add/subtract flag.
func (cpu *CPU) setN(set bool) {
	setFlag(&cpu.Flags.N, set)
}

// setC - set the carry flag.
func (cpu *CPU) setC(set bool) {
	setFlag(&cpu.Flags.C, set)
}

// setSZP - set the sign, zero, and parity flags.
func (cpu *CPU) setSZP(value uint8) {
	cpu.setS(value)
	cpu.setZ(value)
	cpu.setP(value)
	cpu.setXY(value) // Set undocumented X and Y flags
}

// setXY - set the undocumented X and Y flags from bits 3 and 5.
func (cpu *CPU) setXY(value uint8) {
	cpu.Flags.X = (value >> 3) & 1 // bit 3
	cpu.Flags.Y = (value >> 5) & 1 // bit 5
}

// setSZ - set the sign and zero flags.
func (cpu *CPU) setSZ(value uint8) {
	cpu.setS(value)
	cpu.setZ(value)
	cpu.setXY(value) // Set undocumented X and Y flags
}

// setFlags sets the flags from the given byte.
func (cpu *CPU) setFlags(flags uint8) {
	cpu.Flags.C = (flags >> 0) & 1
	cpu.Flags.N = (flags >> 1) & 1
	cpu.Flags.P = (flags >> 2) & 1
	cpu.Flags.X = (flags >> 3) & 1
	cpu.Flags.H = (flags >> 4) & 1
	cpu.Flags.Y = (flags >> 5) & 1
	cpu.Flags.Z = (flags >> 6) & 1
	cpu.Flags.S = (flags >> 7) & 1
}

// setFlag sets a flag to 1 if condition is true, 0 otherwise.
func setFlag(flag *uint8, condition bool) {
	if condition {
		*flag = 1
	} else {
		*flag = 0
	}
}
