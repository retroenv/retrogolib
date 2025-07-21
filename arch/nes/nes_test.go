package nes

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestConstants(t *testing.T) {
	t.Parallel()

	// Test that constants have expected values
	assert.Equal(t, 0x8000, CodeBaseAddress)
	assert.Equal(t, 0x4000, IORegisterStartAddress)
	assert.Equal(t, 0x401F, IORegisterEndAddress)
	assert.Equal(t, 4, NameTableCount)
	assert.Equal(t, 0x400, NameTableSize)
	assert.Equal(t, 32, PaletteSize)
	assert.Equal(t, 0x0FFF, RAMEndAddress)
}

func TestMemoryLayout(t *testing.T) {
	t.Parallel()

	// Test that memory regions don't overlap incorrectly
	assert.True(t, RAMEndAddress < IORegisterStartAddress)
	assert.True(t, IORegisterEndAddress < CodeBaseAddress)
	assert.True(t, IORegisterStartAddress <= IORegisterEndAddress)
}

func TestNameTableSize(t *testing.T) {
	t.Parallel()

	// Test name table calculations
	totalNameTableSize := NameTableCount * NameTableSize
	assert.Equal(t, 0x1000, totalNameTableSize)
}
