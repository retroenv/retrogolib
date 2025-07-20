package gameboy

// Memory map constants for Game Boy DMG
const (
	// ROM areas
	ROMBank0Start = 0x0000 // ROM Bank 0 (16KB)
	ROMBank0End   = 0x3FFF
	ROMBank1Start = 0x4000 // ROM Bank 1-N (16KB, switchable)
	ROMBank1End   = 0x7FFF

	// Video RAM
	VRAMStart = 0x8000 // Video RAM (8KB)
	VRAMEnd   = 0x9FFF

	// External RAM
	ExternalRAMStart = 0xA000 // External RAM (8KB, switchable)
	ExternalRAMEnd   = 0xBFFF

	// Work RAM
	WorkRAMStart = 0xC000 // Work RAM (8KB)
	WorkRAMEnd   = 0xDFFF

	// Echo RAM (mirror of 0xC000-0xDDFF)
	EchoRAMStart = 0xE000
	EchoRAMEnd   = 0xFDFF

	// Object Attribute Memory
	OAMStart = 0xFE00 // OAM (160 bytes)
	OAMEnd   = 0xFE9F

	// Not usable area
	NotUsableStart = 0xFEA0
	NotUsableEnd   = 0xFEFF

	// I/O Registers
	IORegistersStart = 0xFF00
	IORegistersEnd   = 0xFF7F

	// High RAM
	HighRAMStart = 0xFF80 // High RAM (127 bytes)
	HighRAMEnd   = 0xFFFE

	// Interrupt Enable Register
	InterruptEnableRegister = 0xFFFF
)

// Cartridge header constants
const (
	CartridgeHeaderStart = 0x0100
	CartridgeHeaderEnd   = 0x014F

	// Cartridge header offsets
	EntryPoint       = 0x0100 // Entry point (4 bytes)
	NintendoLogo     = 0x0104 // Nintendo logo (48 bytes)
	Title            = 0x0134 // Title (15 bytes)
	ManufacturerCode = 0x013F // Manufacturer code (4 bytes)
	CGBFlag          = 0x0143 // CGB flag
	NewLicenseeCode  = 0x0144 // New licensee code (2 bytes)
	SGBFlag          = 0x0146 // SGB flag
	CartridgeType    = 0x0147 // Cartridge type
	ROMSize          = 0x0148 // ROM size
	RAMSize          = 0x0149 // RAM size
	DestinationCode  = 0x014A // Destination code
	OldLicenseeCode  = 0x014B // Old licensee code
	MaskROMVersion   = 0x014C // Mask ROM version
	HeaderChecksum   = 0x014D // Header checksum
	GlobalChecksum   = 0x014E // Global checksum (2 bytes)
)

// Cartridge types (MBC - Memory Bank Controller)
const (
	CartridgeTypeROMOnly                    = 0x00
	CartridgeTypeMBC1                       = 0x01
	CartridgeTypeMBC1RAM                    = 0x02
	CartridgeTypeMBC1RAMBattery             = 0x03
	CartridgeTypeMBC2                       = 0x05
	CartridgeTypeMBC2Battery                = 0x06
	CartridgeTypeROMRAM                     = 0x08
	CartridgeTypeROMRAMBattery              = 0x09
	CartridgeTypeMMM01                      = 0x0B
	CartridgeTypeMMM01RAM                   = 0x0C
	CartridgeTypeMMM01RAMBattery            = 0x0D
	CartridgeTypeMBC3TimerBattery           = 0x0F
	CartridgeTypeMBC3TimerRAMBattery        = 0x10
	CartridgeTypeMBC3                       = 0x11
	CartridgeTypeMBC3RAM                    = 0x12
	CartridgeTypeMBC3RAMBattery             = 0x13
	CartridgeTypeMBC5                       = 0x19
	CartridgeTypeMBC5RAM                    = 0x1A
	CartridgeTypeMBC5RAMBattery             = 0x1B
	CartridgeTypeMBC5Rumble                 = 0x1C
	CartridgeTypeMBC5RumbleRAM              = 0x1D
	CartridgeTypeMBC5RumbleRAMBattery       = 0x1E
	CartridgeTypeMBC6                       = 0x20
	CartridgeTypeMBC7SensorRumbleRAMBattery = 0x22
	CartridgeTypePocketCamera               = 0xFC
	CartridgeTypeBandaiTAMA5                = 0xFD
	CartridgeTypeHuC3                       = 0xFE
	CartridgeTypeHuC1RAMBattery             = 0xFF
)

// ROM sizes
const (
	ROMSize32KB  = 0x00 // 2 banks
	ROMSize64KB  = 0x01 // 4 banks
	ROMSize128KB = 0x02 // 8 banks
	ROMSize256KB = 0x03 // 16 banks
	ROMSize512KB = 0x04 // 32 banks
	ROMSize1MB   = 0x05 // 64 banks
	ROMSize2MB   = 0x06 // 128 banks
	ROMSize4MB   = 0x07 // 256 banks
	ROMSize8MB   = 0x08 // 512 banks
)

// RAM sizes
const (
	RAMSizeNone  = 0x00
	RAMSize2KB   = 0x01
	RAMSize8KB   = 0x02
	RAMSize32KB  = 0x03 // 4 banks of 8KB each
	RAMSize128KB = 0x04 // 16 banks of 8KB each
	RAMSize64KB  = 0x05 // 8 banks of 8KB each
)

// MBC1 control registers
const (
	MBC1RAMEnable   = 0x0000 // 0x0000-0x1FFF: RAM Enable
	MBC1ROMBankLow  = 0x2000 // 0x2000-0x3FFF: ROM Bank Number (low 5 bits)
	MBC1ROMBankHigh = 0x4000 // 0x4000-0x5FFF: RAM Bank Number or ROM Bank high bits
	MBC1BankingMode = 0x6000 // 0x6000-0x7FFF: Banking Mode Select
)

// MBC3 control registers
const (
	MBC3RAMTimerEnable = 0x0000 // 0x0000-0x1FFF: RAM and Timer Enable
	MBC3ROMBankNumber  = 0x2000 // 0x2000-0x3FFF: ROM Bank Number
	MBC3RAMBankRTC     = 0x4000 // 0x4000-0x5FFF: RAM Bank Number or RTC Register Select
	MBC3LatchClock     = 0x6000 // 0x6000-0x7FFF: Latch Clock Data
)

// MBC5 control registers
const (
	MBC5RAMEnable     = 0x0000 // 0x0000-0x1FFF: RAM Enable
	MBC5ROMBankLow    = 0x2000 // 0x2000-0x2FFF: ROM Bank Number (low 8 bits)
	MBC5ROMBankHigh   = 0x3000 // 0x3000-0x3FFF: ROM Bank Number (high bit)
	MBC5RAMBankNumber = 0x4000 // 0x4000-0x5FFF: RAM Bank Number
)

// Special memory areas
const (
	BootROMStart = 0x0000 // Boot ROM area
	BootROMEnd   = 0x00FF // Boot ROM is disabled after initialization

	CartridgeRAMStart = 0xA000 // Cartridge RAM area
	CartridgeRAMEnd   = 0xBFFF

	VideoRAMSize = 0x2000 // 8KB
	WorkRAMSize  = 0x2000 // 8KB
	HighRAMSize  = 0x007F // 127 bytes
	OAMSize      = 0x00A0 // 160 bytes
)
