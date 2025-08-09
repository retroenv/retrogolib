package x86

import (
	"fmt"
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/log"
)

func TestParity(t *testing.T) {
	tests := []struct {
		value    uint8
		expected bool
	}{
		{0x00, true},  // 0 bits set (even)
		{0x01, false}, // 1 bit set (odd)
		{0x03, true},  // 2 bits set (even)
		{0x07, false}, // 3 bits set (odd)
		{0x0F, true},  // 4 bits set (even)
		{0x1F, false}, // 5 bits set (odd)
		{0x3F, true},  // 6 bits set (even)
		{0x7F, false}, // 7 bits set (odd)
		{0xFF, true},  // 8 bits set (even)
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("0x%02X", tt.value), func(t *testing.T) {
			result := parity(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParityAllValues(t *testing.T) {
	// Test all possible byte values to ensure lookup table is correct
	for i := range 256 {
		value := uint8(i)
		expected := computeParityByBitCount(value)
		result := parity(value)
		if result != expected {
			t.Errorf("parity(0x%02X): expected %v, got %v", value, expected, result)
		}
	}
}

// computeParityByBitCount computes parity by counting bits (reference implementation)
func computeParityByBitCount(value uint8) bool {
	count := 0
	for i := range 8 {
		if (value & (1 << i)) != 0 {
			count++
		}
	}
	return count%2 == 0
}

func TestFlagGettersAndSetters(t *testing.T) {
	tests := []struct {
		name     string
		setFlag  func(*CPU, bool)
		getFlag  func(Flags) bool
		flagMask Flags
	}{
		{"carry", (*CPU).SetCarry, Flags.GetCarry, MaskCarry},
		{"zero", (*CPU).SetZero, Flags.GetZero, MaskZero},
		{"sign", (*CPU).SetSign, Flags.GetSign, MaskSign},
		{"overflow", (*CPU).SetOverflow, Flags.GetOverflow, MaskOverflow},
		{"parity", (*CPU).SetParity, Flags.GetParity, MaskParity},
		{"auxcarry", (*CPU).SetAuxCarry, Flags.GetAuxCarry, MaskAuxCarry},
		{"interrupt", (*CPU).SetInterrupt, Flags.GetInterrupt, MaskInterrupt},
		{"direction", (*CPU).SetDirection, Flags.GetDirection, MaskDirection},
		{"trap", (*CPU).SetTrap, Flags.GetTrap, MaskTrap},
		{"nested", (*CPU).SetNested, Flags.GetNested, MaskNested},
	}

	logger := log.NewTestLogger(t)
	memory, err := NewMemory(64*1024, logger)
	assert.NoError(t, err)

	cpu, err := New(memory)
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test setting flag to true
			tt.setFlag(cpu, true)
			assert.True(t, tt.getFlag(cpu.Flags))
			assert.True(t, (cpu.Flags&tt.flagMask) != 0)

			// Test setting flag to false
			tt.setFlag(cpu, false)
			assert.False(t, tt.getFlag(cpu.Flags))
			assert.True(t, (cpu.Flags&tt.flagMask) == 0)
		})
	}
}

func TestSetSZP8(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(64*1024, logger)
	assert.NoError(t, err)

	cpu, err := New(memory)
	assert.NoError(t, err)

	tests := []struct {
		value  uint8
		sign   bool
		zero   bool
		parity bool
	}{
		{0x00, false, true, true},   // zero, even parity
		{0x01, false, false, false}, // positive, odd parity
		{0x80, true, false, false},  // negative, odd parity (1 bit set)
		{0xFF, true, false, true},   // negative, even parity (8 bits set)
		{0x0F, false, false, true},  // positive, even parity
		{0x07, false, false, false}, // positive, odd parity
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("0x%02X", tt.value), func(t *testing.T) {
			cpu.SetSZP8(tt.value)
			assert.Equal(t, tt.sign, cpu.Flags.GetSign())
			assert.Equal(t, tt.zero, cpu.Flags.GetZero())
			assert.Equal(t, tt.parity, cpu.Flags.GetParity())
		})
	}
}

func TestSetSZP16(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory, err := NewMemory(64*1024, logger)
	assert.NoError(t, err)

	cpu, err := New(memory)
	assert.NoError(t, err)

	tests := []struct {
		value  uint16
		sign   bool
		zero   bool
		parity bool // parity only considers low byte
	}{
		{0x0000, false, true, true},   // zero, even parity (low byte 0x00)
		{0x0001, false, false, false}, // positive, odd parity (low byte 0x01)
		{0x8000, true, false, true},   // negative, even parity (low byte 0x00)
		{0xFFFF, true, false, true},   // negative, even parity (low byte 0xFF)
		{0x1234, false, false, false}, // positive, odd parity (low byte 0x34 has 3 bits set)
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("0x%04X", tt.value), func(t *testing.T) {
			cpu.SetSZP16(tt.value)
			assert.Equal(t, tt.sign, cpu.Flags.GetSign())
			assert.Equal(t, tt.zero, cpu.Flags.GetZero())
			assert.Equal(t, tt.parity, cpu.Flags.GetParity())
		})
	}
}
