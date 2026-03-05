package m68000

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestBasicMemory_ReadWrite(t *testing.T) {
	mem := NewBasicMemory()

	mem.Write(0x1000, 0x42)
	assert.Equal(t, uint8(0x42), mem.Read(0x1000))

	mem.Write(0x0000, 0xFF)
	assert.Equal(t, uint8(0xFF), mem.Read(0x0000))

	mem.Write(0xFFFFFF, 0xAA)
	assert.Equal(t, uint8(0xAA), mem.Read(0xFFFFFF))
}

func TestBasicMemory_ReadWriteWord(t *testing.T) {
	mem := NewBasicMemory()

	mem.WriteWord(0x1000, 0x1234)
	value := mem.ReadWord(0x1000)
	assert.Equal(t, uint16(0x1234), value)

	// Verify big-endian storage.
	assert.Equal(t, uint8(0x12), mem.Read(0x1000)) // High byte first
	assert.Equal(t, uint8(0x34), mem.Read(0x1001)) // Low byte second
}

func TestBasicMemory_ReadWriteLong(t *testing.T) {
	mem := NewBasicMemory()

	mem.WriteLong(0x1000, 0x12345678)
	value := mem.ReadLong(0x1000)
	assert.Equal(t, uint32(0x12345678), value)

	// Verify big-endian storage.
	assert.Equal(t, uint8(0x12), mem.Read(0x1000))
	assert.Equal(t, uint8(0x34), mem.Read(0x1001))
	assert.Equal(t, uint8(0x56), mem.Read(0x1002))
	assert.Equal(t, uint8(0x78), mem.Read(0x1003))
}

func TestBasicMemory_LoadROM(t *testing.T) {
	mem := NewBasicMemory()

	rom := []byte{0x01, 0x02, 0x03, 0x04}
	mem.LoadROM(rom)

	for i, expected := range rom {
		actual := mem.Read(uint32(i))
		assert.Equal(t, expected, actual)
	}

	// Test nil ROM.
	mem2 := NewBasicMemory()
	mem2.LoadROM(nil)
	assert.Equal(t, uint8(0), mem2.Read(0))
}

func TestBasicMemory_AddressMask(t *testing.T) {
	mem := NewBasicMemory()

	// Address 0x01000000 should wrap to 0x000000 (24-bit mask).
	mem.Write(0x01000000, 0x42)
	assert.Equal(t, uint8(0x42), mem.Read(0x000000))
}
