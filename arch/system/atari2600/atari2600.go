package atari2600

// Address space constants. The 6507 has a 13-bit address bus,
// giving an 8 KB address space that mirrors throughout the full 64 KB range.
const (
	// AddressMask is the 13-bit address mask applied by the 6507's address bus.
	AddressMask = 0x1FFF

	// AddressSpaceSize is the total addressable range (8 KB).
	AddressSpaceSize = 0x2000
)

// TIA address ranges.
const (
	// TIAWriteStart is the first TIA write register address.
	TIAWriteStart = 0x0000
	// TIAWriteEnd is the last TIA write register address.
	TIAWriteEnd = 0x002C

	// TIAReadStart is the first TIA read register address.
	TIAReadStart = 0x0000
	// TIAReadEnd is the last TIA read register address.
	TIAReadEnd = 0x000D

	// TIAMirrorMask is the address mask for TIA register mirroring.
	// TIA only decodes the low 6 bits of the address.
	TIAMirrorMask = 0x003F
)

// RIOT (6532) address ranges.
const (
	// RAMStart is the first byte of RIOT RAM.
	RAMStart = 0x0080
	// RAMEnd is the last byte of RIOT RAM.
	RAMEnd = 0x00FF
	// RAMSize is the total RIOT RAM in bytes.
	RAMSize = 128

	// RAMMirrorStart is the start of the RIOT RAM mirror.
	RAMMirrorStart = 0x0180
	// RAMMirrorEnd is the end of the RIOT RAM mirror.
	RAMMirrorEnd = 0x01FF

	// RIOTStart is the first RIOT I/O register address.
	RIOTStart = 0x0280
	// RIOTEnd is the last RIOT I/O register address.
	RIOTEnd = 0x0297
)

// Cartridge ROM address range.
const (
	// ROMStart is the first address of the cartridge ROM window.
	ROMStart = 0x1000
	// ROMEnd is the last address of the cartridge ROM window.
	ROMEnd = 0x1FFF
	// ROMWindowSize is the size of the cartridge ROM window (4 KB).
	ROMWindowSize = 0x1000
)

// Interrupt vector addresses (within the ROM window).
// The 6507 has no IRQ or NMI pins, but the reset vector is still used.
const (
	// ResetVector is the address of the reset vector.
	ResetVector = 0x1FFC
)

// Standard cartridge sizes in bytes.
const (
	CartridgeSize2K  = 2 * 1024
	CartridgeSize4K  = 4 * 1024
	CartridgeSize8K  = 8 * 1024
	CartridgeSize12K = 12 * 1024
	CartridgeSize16K = 16 * 1024
	CartridgeSize32K = 32 * 1024
	CartridgeSize64K = 64 * 1024
)
