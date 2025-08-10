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
