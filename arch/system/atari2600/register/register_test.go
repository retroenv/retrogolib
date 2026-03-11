package register

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestTIAWriteRegisterCompleteness(t *testing.T) {
	assert.Len(t, TIAWriteNames, TIAWriteCount)
}

func TestTIAReadRegisterCompleteness(t *testing.T) {
	assert.Len(t, TIAReadNames, TIAReadCount)
}

func TestTIAWriteRegisterAddresses(t *testing.T) {
	// Verify all write registers are in the valid range
	for addr := range TIAWriteNames {
		assert.True(t, addr <= 0x2C,
			"TIA write register 0x%02X out of range", addr)
	}
}

func TestTIAReadRegisterAddresses(t *testing.T) {
	// Verify all read registers are in the valid range
	for addr := range TIAReadNames {
		assert.True(t, addr <= 0x0D,
			"TIA read register 0x%02X out of range", addr)
	}
}

func TestRIOTRegisterAddresses(t *testing.T) {
	// Verify all RIOT registers are in the valid range
	for addr := range RIOTNames {
		assert.True(t, addr >= 0x0280 && addr <= 0x0297,
			"RIOT register 0x%04X out of range", addr)
	}
}

func TestRIOTRegisterCompleteness(t *testing.T) {
	// 10 RIOT registers defined
	assert.Len(t, RIOTNames, 10)
}

func TestConsoleSwitchBits(t *testing.T) {
	// Verify switch bits don't overlap
	allBits := SwitchReset | SwitchSelect | SwitchBW | SwitchP0Diff | SwitchP1Diff
	assert.Equal(t, uint8(0xCB), allBits)
}

func TestJoystickBits(t *testing.T) {
	// Player 0 uses upper nibble
	p0bits := Joy0Right | Joy0Left | Joy0Down | Joy0Up
	assert.Equal(t, uint8(0xF0), p0bits)

	// Player 1 uses lower nibble
	p1bits := Joy1Right | Joy1Left | Joy1Down | Joy1Up
	assert.Equal(t, uint8(0x0F), p1bits)

	// No overlap
	assert.Equal(t, uint8(0x00), p0bits&p1bits)
}
