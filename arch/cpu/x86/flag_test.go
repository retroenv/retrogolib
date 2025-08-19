package x86

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestFlags_GetMethods(t *testing.T) {
	// Test flag getter methods on a standalone Flags value
	var flags Flags

	// Test initial state (all flags should be false)
	assert.False(t, flags.GetCarry())
	assert.False(t, flags.GetZero())
	assert.False(t, flags.GetSign())
	assert.False(t, flags.GetOverflow())
	assert.False(t, flags.GetParity())

	// Test setting individual flags
	flags = MaskCarry | MaskZero
	assert.True(t, flags.GetCarry())
	assert.True(t, flags.GetZero())
	assert.False(t, flags.GetSign())

	// Test all flags set
	flags = 0xFFFF
	assert.True(t, flags.GetCarry())
	assert.True(t, flags.GetZero())
	assert.True(t, flags.GetSign())
	assert.True(t, flags.GetOverflow())
	assert.True(t, flags.GetParity())
}

func TestFlags_Constants(t *testing.T) {
	// Test that flag constants are properly defined
	assert.Equal(t, Flags(1<<0), MaskCarry)
	assert.Equal(t, Flags(1<<2), MaskParity)
	assert.Equal(t, Flags(1<<6), MaskZero)
	assert.Equal(t, Flags(1<<7), MaskSign)
	assert.Equal(t, Flags(1<<11), MaskOverflow)
}

func TestFlags_SetMethods(t *testing.T) {
	tests := []struct {
		name    string
		setFunc func(*Flags, bool)
		getFunc func(Flags) bool
		mask    Flags
	}{
		{"carry", (*Flags).SetCarry, Flags.GetCarry, MaskCarry},
		{"parity", (*Flags).SetParity, Flags.GetParity, MaskParity},
		{"auxcarry", (*Flags).SetAuxCarry, Flags.GetAuxCarry, MaskAuxCarry},
		{"zero", (*Flags).SetZero, Flags.GetZero, MaskZero},
		{"sign", (*Flags).SetSign, Flags.GetSign, MaskSign},
		{"trap", (*Flags).SetTrap, Flags.GetTrap, MaskTrap},
		{"interrupt", (*Flags).SetInterrupt, Flags.GetInterrupt, MaskInterrupt},
		{"direction", (*Flags).SetDirection, Flags.GetDirection, MaskDirection},
		{"overflow", (*Flags).SetOverflow, Flags.GetOverflow, MaskOverflow},
		{"nested", (*Flags).SetNested, Flags.GetNested, MaskNested},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var flags Flags

			// Test setting the flag
			tt.setFunc(&flags, true)
			assert.True(t, tt.getFunc(flags), "flag should be set")
			assert.Equal(t, tt.mask, flags&tt.mask, "only the specific flag should be set")

			// Test clearing the flag
			tt.setFunc(&flags, false)
			assert.False(t, tt.getFunc(flags), "flag should be cleared")
			assert.Equal(t, Flags(0), flags&tt.mask, "flag should be completely cleared")
		})
	}
}

func TestFlags_SetIOPL(t *testing.T) {
	var flags Flags

	// Test setting different IOPL levels
	flags.SetIOPL(0)
	assert.Equal(t, uint8(0), flags.GetIOPL())

	flags.SetIOPL(1)
	assert.Equal(t, uint8(1), flags.GetIOPL())

	flags.SetIOPL(2)
	assert.Equal(t, uint8(2), flags.GetIOPL())

	flags.SetIOPL(3)
	assert.Equal(t, uint8(3), flags.GetIOPL())

	// Test that values > 3 are masked to 2 bits
	flags.SetIOPL(7) // Binary: 111, should become 11 (3)
	assert.Equal(t, uint8(3), flags.GetIOPL())
}
