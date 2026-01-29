package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestNewBasicMemory(t *testing.T) {
	memory := NewBasicMemory()
	assert.NotNil(t, memory)
}

func TestNewGameBoyMemory(t *testing.T) {
	memory := NewGameBoyMemory()
	assert.NotNil(t, memory)
	assert.Equal(t, uint8(1), memory.GetROMBank(), "ROM bank should be initialized to 1")
	assert.Equal(t, uint8(0), memory.GetRAMBank(), "RAM bank should be initialized to 0")
}

func TestBasicMemory_ReadWrite(t *testing.T) {
	memory := NewBasicMemory()

	// Test basic read/write
	memory.Write(0x1000, 0x42)
	value := memory.Read(0x1000)
	assert.Equal(t, uint8(0x42), value)

	// Test boundary values
	memory.Write(0x0000, 0xFF)
	assert.Equal(t, uint8(0xFF), memory.Read(0x0000))

	memory.Write(0xFFFF, 0xAA)
	assert.Equal(t, uint8(0xAA), memory.Read(0xFFFF))
}

func TestBasicMemory_ReadWriteWord(t *testing.T) {
	memory := NewBasicMemory()

	// Test 16-bit read/write (little-endian)
	memory.WriteWord(0x1000, 0x1234)
	value := memory.ReadWord(0x1000)
	assert.Equal(t, uint16(0x1234), value)

	// Verify little-endian storage
	assert.Equal(t, uint8(0x34), memory.Read(0x1000)) // Low byte
	assert.Equal(t, uint8(0x12), memory.Read(0x1001)) // High byte
}

func TestBasicMemory_LoadROM(t *testing.T) {
	memory := NewBasicMemory()

	// Test small ROM
	rom := []byte{0x01, 0x02, 0x03, 0x04}
	memory.LoadROM(rom)

	for i, expected := range rom {
		actual := memory.Read(uint16(i))
		assert.Equal(t, expected, actual)
	}

	// Test oversized ROM (should be truncated)
	const largeROMSize = 0x20000 // 128KB, larger than 64KB memory
	largeROM := make([]byte, largeROMSize)
	for i := range largeROMSize {
		largeROM[i] = uint8(i & 0xFF)
	}

	memory2 := NewBasicMemory()
	memory2.LoadROM(largeROM)

	// Should only load first 64KB
	for i := range 0x10000 {
		expected := uint8(i & 0xFF)
		actual := memory2.Read(uint16(i))
		assert.Equal(t, expected, actual)
	}
}

func TestGameBoyMemory_Banking(t *testing.T) {
	memory := NewGameBoyMemory()

	// Test ROM bank switching
	memory.Write(0x2000, 0x05) // Set ROM bank to 5
	assert.Equal(t, uint8(0x05), memory.GetROMBank())

	// Test ROM bank 0 -> 1 conversion
	memory.Write(0x2000, 0x00) // Bank 0 should become 1
	assert.Equal(t, uint8(0x01), memory.GetROMBank())

	// Test RAM bank switching in MBC1 mode
	memory.SetBankingMode(true)
	memory.Write(0x4000, 0x02) // Set RAM bank to 2
	assert.Equal(t, uint8(0x02), memory.GetRAMBank())

	// Test upper ROM bank bits when not in MBC1 mode
	memory.SetBankingMode(false)
	memory.Write(0x2000, 0x10)                        // Set lower ROM bank
	memory.Write(0x4000, 0x01)                        // Set upper ROM bank bits
	assert.Equal(t, uint8(0x30), memory.GetROMBank()) // 0x10 | (0x01 << 5)

	// Test banking mode toggle
	memory.Write(0x6000, 0x01) // Enable MBC1 mode
	assert.True(t, memory.mbc1Mode)

	memory.Write(0x6000, 0x00) // Disable MBC1 mode
	assert.False(t, memory.mbc1Mode)
}
