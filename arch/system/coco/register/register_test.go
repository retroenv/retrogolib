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
	assert.Equal(t, uint16(0xFFDF), uint16(SAMTYSet))
}

func TestSAMControlAddresses(t *testing.T) {
	// TY at $FFDE/$FFDF changes the memory map, not CPU speed. Rate uses R0/R1.
	tests := []struct {
		name                     string
		clearAddress, setAddress uint16
		wantClear                uint16
	}{
		{name: "rate R0", clearAddress: SAMR0Clear, setAddress: SAMR0Set, wantClear: 0xFFD6},
		{name: "rate R1", clearAddress: SAMR1Clear, setAddress: SAMR1Set, wantClear: 0xFFD8},
		{name: "memory size M0", clearAddress: SAMM0Clear, setAddress: SAMM0Set, wantClear: 0xFFDA},
		{name: "memory size M1", clearAddress: SAMM1Clear, setAddress: SAMM1Set, wantClear: 0xFFDC},
		{name: "memory map TY", clearAddress: SAMTYClear, setAddress: SAMTYSet, wantClear: 0xFFDE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantClear, tt.clearAddress)
			assert.Equal(t, tt.wantClear+1, tt.setAddress)
		})
	}
}
