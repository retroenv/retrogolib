package gameboy

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestParseCartridgeHeader(t *testing.T) {
	// Create a minimal ROM with header
	rom := make([]byte, 0x8000)

	// Fill in some header fields
	copy(rom[0x0134:], []byte("TEST GAME\x00\x00\x00\x00\x00\x00")) // Title
	rom[0x0147] = CartridgeTypeMBC1                                 // Cartridge type
	rom[0x0148] = ROMSize64KB                                       // ROM size
	rom[0x0149] = RAMSize8KB                                        // RAM size
	rom[0x014A] = 0x01                                              // Destination code
	rom[0x014B] = 0x33                                              // Old licensee code
	rom[0x014C] = 0x00                                              // Mask ROM version
	rom[0x014D] = 0x90                                              // Header checksum (example)

	header, err := ParseCartridgeHeader(rom)
	assert.NoError(t, err, "ParseCartridgeHeader should not return error")
	assert.NotNil(t, header, "Header should not be nil")

	assert.Equal(t, uint8(CartridgeTypeMBC1), header.CartridgeType, "Cartridge type should match")
	assert.Equal(t, uint8(ROMSize64KB), header.ROMSize, "ROM size should match")
	assert.Equal(t, uint8(RAMSize8KB), header.RAMSize, "RAM size should match")
	assert.Equal(t, uint8(0x01), header.DestinationCode, "Destination code should match")
	assert.Equal(t, uint8(0x33), header.OldLicenseeCode, "Old licensee code should match")
	assert.Equal(t, uint8(0x00), header.MaskROMVersion, "Mask ROM version should match")
	assert.Equal(t, uint8(0x90), header.HeaderChecksum, "Header checksum should match")
}

func TestParseCartridgeHeaderTooSmall(t *testing.T) {
	// ROM too small to contain header
	rom := make([]byte, 0x100)

	header, err := ParseCartridgeHeader(rom)
	assert.Error(t, err, "ROM too small to contain header", "ParseCartridgeHeader should return error for small ROM")
	assert.Nil(t, header, "Header should be nil for error case")
	assert.Contains(t, err.Error(), "ROM too small", "Error should mention ROM size")
}

func TestGetTitle(t *testing.T) {
	header := &CartridgeHeader{}
	copy(header.Title[:], []byte("TETRIS\x00\x00\x00\x00\x00\x00\x00\x00\x00"))

	title := header.GetTitle()
	assert.Equal(t, "TETRIS", title, "Title should be extracted correctly")

	// Test with full title (no null terminator)
	copy(header.Title[:], []byte("SUPER MARIO LAN"))
	title = header.GetTitle()
	assert.Equal(t, "SUPER MARIO LAN", title, "Full title should be extracted")

	// Test with empty title
	copy(header.Title[:], []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
	title = header.GetTitle()
	assert.Equal(t, "", title, "Empty title should return empty string")
}

func TestGetCartridgeTypeName(t *testing.T) {
	header := &CartridgeHeader{}

	testCases := []struct {
		cartridgeType uint8
		expectedName  string
	}{
		{CartridgeTypeROMOnly, "ROM ONLY"},
		{CartridgeTypeMBC1, "MBC1"},
		{CartridgeTypeMBC1RAM, "MBC1+RAM"},
		{CartridgeTypeMBC1RAMBattery, "MBC1+RAM+BATTERY"},
		{CartridgeTypeMBC2, "MBC2"},
		{CartridgeTypeMBC3, "MBC3"},
		{CartridgeTypeMBC5, "MBC5"},
		{0x99, "UNKNOWN (0x99)"},
	}

	for _, tc := range testCases {
		header.CartridgeType = tc.cartridgeType
		name := header.GetCartridgeTypeName()
		assert.Equal(t, tc.expectedName, name, "Cartridge type name should match")
	}
}

func TestGetROMSizeKB(t *testing.T) {
	header := &CartridgeHeader{}

	testCases := []struct {
		romSize      uint8
		expectedSize int
	}{
		{ROMSize32KB, 32},
		{ROMSize64KB, 64},
		{ROMSize128KB, 128},
		{ROMSize256KB, 256},
		{ROMSize512KB, 512},
		{ROMSize1MB, 1024},
		{ROMSize2MB, 2048},
		{ROMSize4MB, 4096},
		{ROMSize8MB, 8192},
		{0xFF, 0}, // Unknown size
	}

	for _, tc := range testCases {
		header.ROMSize = tc.romSize
		size := header.GetROMSizeKB()
		assert.Equal(t, tc.expectedSize, size, "ROM size should match")
	}
}

func TestGetRAMSizeKB(t *testing.T) {
	header := &CartridgeHeader{}

	testCases := []struct {
		ramSize      uint8
		expectedSize int
	}{
		{RAMSizeNone, 0},
		{RAMSize2KB, 2},
		{RAMSize8KB, 8},
		{RAMSize32KB, 32},
		{RAMSize128KB, 128},
		{RAMSize64KB, 64},
		{0xFF, 0}, // Unknown size
	}

	for _, tc := range testCases {
		header.RAMSize = tc.ramSize
		size := header.GetRAMSizeKB()
		assert.Equal(t, tc.expectedSize, size, "RAM size should match")
	}
}

func TestHasBattery(t *testing.T) {
	header := &CartridgeHeader{}

	// Test cartridge types with battery
	batteryTypes := []uint8{
		CartridgeTypeMBC1RAMBattery,
		CartridgeTypeMBC2Battery,
		CartridgeTypeROMRAMBattery,
		CartridgeTypeMMM01RAMBattery,
		CartridgeTypeMBC3TimerBattery,
		CartridgeTypeMBC3TimerRAMBattery,
		CartridgeTypeMBC3RAMBattery,
		CartridgeTypeMBC5RAMBattery,
		CartridgeTypeMBC5RumbleRAMBattery,
		CartridgeTypeMBC7SensorRumbleRAMBattery,
		CartridgeTypeHuC1RAMBattery,
	}

	for _, cartType := range batteryTypes {
		header.CartridgeType = cartType
		assert.True(t, header.HasBattery(), "Should have battery for type 0x%02X", cartType)
	}

	// Test cartridge types without battery
	noBatteryTypes := []uint8{
		CartridgeTypeROMOnly,
		CartridgeTypeMBC1,
		CartridgeTypeMBC1RAM,
		CartridgeTypeMBC2,
		CartridgeTypeMBC3,
		CartridgeTypeMBC5,
	}

	for _, cartType := range noBatteryTypes {
		header.CartridgeType = cartType
		assert.False(t, header.HasBattery(), "Should not have battery for type 0x%02X", cartType)
	}
}

func TestHasRAM(t *testing.T) {
	header := &CartridgeHeader{}

	// Test cartridge types with RAM
	ramTypes := []uint8{
		CartridgeTypeMBC1RAM,
		CartridgeTypeMBC1RAMBattery,
		CartridgeTypeMBC2, // MBC2 has built-in RAM
		CartridgeTypeMBC2Battery,
		CartridgeTypeROMRAM,
		CartridgeTypeROMRAMBattery,
		CartridgeTypeMBC3RAM,
		CartridgeTypeMBC3RAMBattery,
		CartridgeTypeMBC5RAM,
		CartridgeTypeMBC5RAMBattery,
	}

	for _, cartType := range ramTypes {
		header.CartridgeType = cartType
		assert.True(t, header.HasRAM(), "Should have RAM for type 0x%02X", cartType)
	}

	// Test cartridge types without RAM
	noRAMTypes := []uint8{
		CartridgeTypeROMOnly,
		CartridgeTypeMBC1,
		CartridgeTypeMBC3,
		CartridgeTypeMBC5,
	}

	for _, cartType := range noRAMTypes {
		header.CartridgeType = cartType
		assert.False(t, header.HasRAM(), "Should not have RAM for type 0x%02X", cartType)
	}
}

func TestHasTimer(t *testing.T) {
	header := &CartridgeHeader{}

	// Test cartridge types with timer
	timerTypes := []uint8{
		CartridgeTypeMBC3TimerBattery,
		CartridgeTypeMBC3TimerRAMBattery,
	}

	for _, cartType := range timerTypes {
		header.CartridgeType = cartType
		assert.True(t, header.HasTimer(), "Should have timer for type 0x%02X", cartType)
	}

	// Test cartridge types without timer
	noTimerTypes := []uint8{
		CartridgeTypeROMOnly,
		CartridgeTypeMBC1,
		CartridgeTypeMBC2,
		CartridgeTypeMBC3,
		CartridgeTypeMBC5,
	}

	for _, cartType := range noTimerTypes {
		header.CartridgeType = cartType
		assert.False(t, header.HasTimer(), "Should not have timer for type 0x%02X", cartType)
	}
}

func TestGetMBCType(t *testing.T) {
	header := &CartridgeHeader{}

	testCases := []struct {
		cartridgeType uint8
		expectedMBC   int
	}{
		{CartridgeTypeROMOnly, 0},
		{CartridgeTypeMBC1, 1},
		{CartridgeTypeMBC1RAM, 1},
		{CartridgeTypeMBC1RAMBattery, 1},
		{CartridgeTypeMBC2, 2},
		{CartridgeTypeMBC2Battery, 2},
		{CartridgeTypeMBC3, 3},
		{CartridgeTypeMBC3RAM, 3},
		{CartridgeTypeMBC3TimerBattery, 3},
		{CartridgeTypeMBC5, 5},
		{CartridgeTypeMBC5RAM, 5},
		{CartridgeTypeMBC6, 6},
		{CartridgeTypeMBC7SensorRumbleRAMBattery, 7},
		{CartridgeTypePocketCamera, -1}, // Unknown
	}

	for _, tc := range testCases {
		header.CartridgeType = tc.cartridgeType
		mbcType := header.GetMBCType()
		assert.Equal(t, tc.expectedMBC, mbcType, "MBC type should match for cartridge type 0x%02X", tc.cartridgeType)
	}
}

func TestValidateHeaderChecksum(t *testing.T) {
	// Create a ROM with a correct checksum
	rom := make([]byte, 0x8000)

	// Fill in header data
	copy(rom[0x0134:], []byte("TEST\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"))
	rom[0x0147] = CartridgeTypeMBC1
	rom[0x0148] = ROMSize64KB
	rom[0x0149] = RAMSize8KB

	// Calculate correct checksum
	var checksum uint8
	for i := 0x0134; i <= 0x014C; i++ {
		checksum = checksum - rom[i] - 1
	}
	rom[0x014D] = checksum

	header, err := ParseCartridgeHeader(rom)
	assert.NoError(t, err, "Should parse header successfully")

	// Test with correct checksum
	assert.True(t, header.ValidateHeaderChecksum(rom), "Checksum should be valid")

	// Test with incorrect checksum
	rom[0x014D] = checksum + 1
	// Re-parse header with incorrect checksum
	headerBad, err := ParseCartridgeHeader(rom)
	assert.NoError(t, err, "Should parse header successfully")
	assert.False(t, headerBad.ValidateHeaderChecksum(rom), "Checksum should be invalid")

	// Test with ROM too small
	smallROM := make([]byte, 0x100)
	assert.False(t, header.ValidateHeaderChecksum(smallROM), "Should fail for small ROM")
}
