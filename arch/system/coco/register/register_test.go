package register

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestPIA0RegisterAddresses(t *testing.T) {
	for addr := range PIA0Names {
		assert.True(t, addr >= 0xFF00 && addr <= 0xFF03,
			"PIA0 register 0x%04X out of range", addr)
	}
}

func TestPIA0RegisterCompleteness(t *testing.T) {
	assert.Len(t, PIA0Names, 4)
}

func TestPIA1RegisterAddresses(t *testing.T) {
	for addr := range PIA1Names {
		assert.True(t, addr >= 0xFF20 && addr <= 0xFF23,
			"PIA1 register 0x%04X out of range", addr)
	}
}

func TestPIA1RegisterCompleteness(t *testing.T) {
	assert.Len(t, PIA1Names, 4)
}

func TestSAMRegisterRange(t *testing.T) {
	// SAM registers span $FFC0-$FFDF (32 bytes)
	assert.Equal(t, uint16(0xFFC0), uint16(SAMV0Clear))
	assert.Equal(t, uint16(0xFFDF), uint16(SAMRateSet))
}
