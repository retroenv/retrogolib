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

// VIA (MC6522) address range.
const (
	// VIAStart is the first VIA register address.
	VIAStart = 0xD000
	// VIAEnd is the last VIA register address.
	VIAEnd = 0xD00F

	// VIAMirrorStart is the start of the VIA register mirror.
	VIAMirrorStart = 0xD800
	// VIAMirrorEnd is the end of the VIA register mirror.
	VIAMirrorEnd = 0xD80F
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
)

// Standard cartridge sizes in bytes.
const (
	CartridgeSize4K  = 4 * 1024
	CartridgeSize8K  = 8 * 1024
	CartridgeSize16K = 16 * 1024
	CartridgeSize32K = 32 * 1024
)
