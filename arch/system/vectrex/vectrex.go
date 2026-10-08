package vectrex

// Address space constants.
const (
	// AddressSpaceSize is the total addressable range (64 KB).
	AddressSpaceSize = 0x10000
)

// Cartridge ROM address range.
const (
	// CartridgeStart is the first address of the cartridge ROM window.
	CartridgeStart = 0x0000
	// CartridgeEnd is the last address of the cartridge ROM window.
	CartridgeEnd = 0x7FFF
	// CartridgeMaxSize is the maximum cartridge ROM size (32 KB).
	CartridgeMaxSize = 0x8000
)

// RAM address range.
const (
	// RAMStart is the first byte of RAM.
	RAMStart = 0xC800
	// RAMEnd is the last byte of RAM.
	RAMEnd = 0xCBFF
	// RAMSize is the total RAM in bytes (1 KB).
	RAMSize = 1024

	// RAMMirrorStart is the start of the RAM mirror.
	RAMMirrorStart = 0xCC00
	// RAMMirrorEnd is the end of the RAM mirror.
	RAMMirrorEnd = 0xCFFF
)

// VIA (6522) address range. The VIA decodes only A0-A3, so the 16 registers
// repeat through the complete region.
const (
	// VIAStart is the first VIA register address.
	VIAStart = 0xD000
	// VIAEnd is the last VIA register address.
	VIAEnd = 0xD00F

	// VIARegionEnd is the last address that selects the VIA.
	VIARegionEnd = 0xD7FF
	// VIAMirrorMask selects the VIA register offset from a mirrored address.
	VIAMirrorMask = 0x000F
)

// Address range that selects RAM and VIA at the same time. Software must not use it.
const (
	// DualSelectStart is the first address that selects both RAM and VIA.
	DualSelectStart = 0xD800
	// DualSelectEnd is the last address that selects both RAM and VIA.
	DualSelectEnd = 0xDFFF
)

// System ROM address range.
const (
	// ROMStart is the first address of the system ROM (executive/BIOS).
	ROMStart = 0xE000
	// ROMEnd is the last address of the system ROM.
	ROMEnd = 0xFFFF
	// ROMSize is the size of the system ROM (8 KB).
	ROMSize = 0x2000
)

// Interrupt vector addresses (within system ROM).
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
	CartridgeSize4K  = 4 * 1024
	CartridgeSize8K  = 8 * 1024
	CartridgeSize16K = 16 * 1024
	CartridgeSize32K = 32 * 1024
)
