package register

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestVIARegisterAddresses(t *testing.T) {
	for addr := range VIANames {
		assert.True(t, addr >= 0xD000 && addr <= 0xD00F,
			"VIA register 0x%04X out of range", addr)
	}
}

func TestVIARegisterCompleteness(t *testing.T) {
	assert.Len(t, VIANames, VIARegisterCount)
}

func TestButtonBits(t *testing.T) {
	// Verify button bits don't overlap
	allBits := ButtonRight | ButtonLeft | ButtonDown | ButtonUp
	assert.Equal(t, uint8(0x0F), allBits)
}

func TestIRQBits(t *testing.T) {
	// Timer 1 is the most commonly used interrupt
	assert.Equal(t, uint8(0x40), uint8(IRQTimer1))
	assert.Equal(t, uint8(0x80), uint8(IRQAny))
}
