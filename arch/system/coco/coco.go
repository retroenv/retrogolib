package coco

// Address space constants.
const (
	// AddressSpaceSize is the total addressable range (64 KB).
	AddressSpaceSize = 0x10000
)

// RAM address range.
const (
	// RAMStart is the first byte of RAM.
	RAMStart = 0x0000
	// RAMEnd is the last byte of standard 32 KB RAM.
	RAMEnd = 0x7FFF
	// RAMSize32K is the standard RAM size.
	RAMSize32K = 32 * 1024
	// RAMSize64K is the extended RAM size.
	RAMSize64K = 64 * 1024
)

// ROM address ranges.
const (
	// ExtendedBASICStart is the start of Extended Color BASIC ROM.
	ExtendedBASICStart = 0x8000
	// ExtendedBASICEnd is the end of Extended Color BASIC ROM.
	ExtendedBASICEnd = 0x9FFF
	// ExtendedBASICSize is the size of Extended BASIC ROM (8 KB).
	ExtendedBASICSize = 0x2000

	// ColorBASICStart is the start of Color BASIC ROM.
	ColorBASICStart = 0xA000
	// ColorBASICEnd is the end of Color BASIC ROM.
	ColorBASICEnd = 0xBFFF
	// ColorBASICSize is the size of Color BASIC ROM (8 KB).
	ColorBASICSize = 0x2000

	// CartridgeStart is the start of the cartridge ROM space.
	CartridgeStart = 0xC000
	// CartridgeEnd is the end of the cartridge ROM space.
	CartridgeEnd = 0xFEFF
)

// I/O address ranges.
const (
	// PIA0Start is the first PIA 0 register address.
	PIA0Start = 0xFF00
	// PIA0End is the last PIA 0 register address.
	PIA0End = 0xFF03

	// PIA1Start is the first PIA 1 register address.
	PIA1Start = 0xFF20
	// PIA1End is the last PIA 1 register address.
	PIA1End = 0xFF23

	// FloppyStart is the start of floppy disk controller registers.
	FloppyStart = 0xFF40
	// FloppyEnd is the end of floppy disk controller registers.
	FloppyEnd = 0xFF5F

	// SAMStart is the start of SAM register addresses.
	SAMStart = 0xFFC0
	// SAMEnd is the end of SAM register addresses.
	SAMEnd = 0xFFDF
)

// Interrupt vector addresses.
const (
	// ResetVector is the address of the reset vector.
	ResetVector = 0xFFFE
	// NMIVector is the address of the NMI vector.
	NMIVector = 0xFFFC
	// SWIVector is the address of the SWI vector.
	SWIVector = 0xFFFA
	// IRQVector is the address of the IRQ vector.
	IRQVector = 0xFFF8
	// FIRQVector is the address of the FIRQ vector.
	FIRQVector = 0xFFF6
	// SWI2Vector is the address of the SWI2 vector.
	SWI2Vector = 0xFFF4
	// SWI3Vector is the address of the SWI3 vector.
	SWI3Vector = 0xFFF2
)

// Standard cartridge sizes in bytes.
const (
	CartridgeSize8K  = 8 * 1024
	CartridgeSize16K = 16 * 1024
	CartridgeSize32K = 32 * 1024
)
