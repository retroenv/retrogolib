package m65816

// Processor status flag bit positions in the P register.
const (
	FlagCarry    = 0 // C - Carry
	FlagZero     = 1 // Z - Zero
	FlagIRQ      = 2 // I - IRQ Disable
	FlagDecimal  = 3 // D - Decimal Mode
	FlagIndex    = 4 // X - Index register width (native mode: 0=16-bit, 1=8-bit)
	FlagBreak    = 4 // B - Break (emulation mode only, same bit as X)
	FlagMemory   = 5 // M - Accumulator/memory width (native mode: 0=16-bit, 1=8-bit)
	FlagOverflow = 6 // V - Overflow
	FlagNegative = 7 // N - Negative
)

// Processor status flag masks.
const (
	MaskCarry    = 1 << FlagCarry
	MaskZero     = 1 << FlagZero
	MaskIRQ      = 1 << FlagIRQ
	MaskDecimal  = 1 << FlagDecimal
	MaskIndex    = 1 << FlagIndex
	MaskBreak    = 1 << FlagBreak
	MaskMemory   = 1 << FlagMemory
	MaskOverflow = 1 << FlagOverflow
	MaskNegative = 1 << FlagNegative
)

// Flags contains the processor status register (P) broken out as individual fields.
// In emulation mode (E=1), the X/M bit positions serve as B (break) and always-1 bit.
type Flags struct {
	C uint8 // carry flag
	Z uint8 // zero flag
	I uint8 // interrupt disable flag
	D uint8 // decimal mode flag
	X uint8 // index register width flag (native) / break flag (emulation)
	M uint8 // memory/accumulator width flag (native) / always 1 (emulation)
	V uint8 // overflow flag
	N uint8 // negative flag
}

// Get returns the flags as a single byte (P register value).
func (f *Flags) Get() uint8 {
	return f.C |
		(f.Z << 1) |
		(f.I << 2) |
		(f.D << 3) |
		(f.X << 4) |
		(f.M << 5) |
		(f.V << 6) |
		(f.N << 7)
}

// Set decodes a P register byte into the individual flag fields.
func (f *Flags) Set(p uint8) {
	f.C = (p >> 0) & 1
	f.Z = (p >> 1) & 1
	f.I = (p >> 2) & 1
	f.D = (p >> 3) & 1
	f.X = (p >> 4) & 1
	f.M = (p >> 5) & 1
	f.V = (p >> 6) & 1
	f.N = (p >> 7) & 1
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

// compare8 compares two 8-bit values and sets N, Z, C flags.
func (c *CPU) compare8(a, b uint8) {
	result := int(a) - int(b)
	setFlag(&c.Flags.C, a >= b)
	c.setZN8(uint8(result))
}

// compare16 compares two 16-bit values and sets N, Z, C flags.
func (c *CPU) compare16(a, b uint16) {
	result := int32(a) - int32(b)
	setFlag(&c.Flags.C, a >= b)
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
