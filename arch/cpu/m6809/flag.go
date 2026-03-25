package m6809

// Condition code register (CC) flag bit positions.
const (
	FlagCarry     = 0 // C - Carry
	FlagOverflow  = 1 // V - Overflow
	FlagZero      = 2 // Z - Zero
	FlagNegative  = 3 // N - Negative
	FlagIRQ       = 4 // I - IRQ Disable
	FlagHalfCarry = 5 // H - Half Carry
	FlagFIRQ      = 6 // F - FIRQ Disable
	FlagEntire    = 7 // E - Entire flag (set if full state saved on stack)
)

// Condition code register flag masks.
const (
	MaskCarry     = 1 << FlagCarry
	MaskOverflow  = 1 << FlagOverflow
	MaskZero      = 1 << FlagZero
	MaskNegative  = 1 << FlagNegative
	MaskIRQ       = 1 << FlagIRQ
	MaskHalfCarry = 1 << FlagHalfCarry
	MaskFIRQ      = 1 << FlagFIRQ
	MaskEntire    = 1 << FlagEntire
)

// Flags contains the condition code register (CC) broken out as individual fields.
type Flags struct {
	C uint8 // carry flag
	V uint8 // overflow flag
	Z uint8 // zero flag
	N uint8 // negative flag
	I uint8 // IRQ disable flag
	H uint8 // half carry flag
	F uint8 // FIRQ disable flag
	E uint8 // entire flag
}

// Get returns the flags as a single byte (CC register value).
func (f *Flags) Get() uint8 {
	return f.C |
		(f.V << 1) |
		(f.Z << 2) |
		(f.N << 3) |
		(f.I << 4) |
		(f.H << 5) |
		(f.F << 6) |
		(f.E << 7)
}

// Set decodes a CC register byte into the individual flag fields.
func (f *Flags) Set(cc uint8) {
	f.C = (cc >> 0) & 1
	f.V = (cc >> 1) & 1
	f.Z = (cc >> 2) & 1
	f.N = (cc >> 3) & 1
	f.I = (cc >> 4) & 1
	f.H = (cc >> 5) & 1
	f.F = (cc >> 6) & 1
	f.E = (cc >> 7) & 1
}

// setZN8 sets the zero and negative flags based on an 8-bit value.
func (c *CPU) setZN8(value uint8) {
	setFlag(&c.Flags.Z, value == 0)
	setFlag(&c.Flags.N, value&0x80 != 0)
}

// setZN16 sets the zero and negative flags based on a 16-bit value.
func (c *CPU) setZN16(value uint16) {
	setFlag(&c.Flags.Z, value == 0)
	setFlag(&c.Flags.N, value&0x8000 != 0)
}

// compare8 compares two 8-bit values and sets N, Z, V, C flags.
func (c *CPU) compare8(a, b uint8) {
	result := int16(a) - int16(b)
	setFlag(&c.Flags.C, a >= b)
	setFlag(&c.Flags.V, (a^b)&0x80 != 0 && (a^uint8(result))&0x80 != 0)
	c.setZN8(uint8(result))
}

// compare16 compares two 16-bit values and sets N, Z, V, C flags.
func (c *CPU) compare16(a, b uint16) {
	result := int32(a) - int32(b)
	setFlag(&c.Flags.C, a >= b)
	setFlag(&c.Flags.V, (a^b)&0x8000 != 0 && (a^uint16(result))&0x8000 != 0)
	c.setZN16(uint16(result))
}

// setFlag sets or clears a processor status flag.
func setFlag(flag *uint8, condition bool) {
	if condition {
		*flag = 1
	} else {
		*flag = 0
	}
}
