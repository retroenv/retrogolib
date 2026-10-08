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
	// Player 1 uses the low nibble and player 2 the high nibble of the PSG port.
	player1 := Player1Button1 | Player1Button2 | Player1Button3 | Player1Button4
	player2 := Player2Button1 | Player2Button2 | Player2Button3 | Player2Button4
	assert.Equal(t, uint8(0x0F), uint8(player1))
	assert.Equal(t, uint8(0xF0), uint8(player2))
	assert.Equal(t, uint8(0x0E), uint8(PSGIOPortA))
}

func TestPortBBits(t *testing.T) {
	// Each port B signal must occupy one distinct bit.
	bits := []uint8{
		PortBSwitch, PortBSel0, PortBSel1, PortBBC1,
		PortBBDIR, PortBCompare, PortBCart, PortBRamp,
	}

	var all uint8
	for i, bit := range bits {
		assert.Equal(t, uint8(1<<i), bit)
		all |= bit
	}
	assert.Equal(t, uint8(0xFF), all)
}

func TestIRQBits(t *testing.T) {
	// Timer 1 is the most commonly used interrupt
	assert.Equal(t, uint8(0x40), uint8(IRQTimer1))
	assert.Equal(t, uint8(0x80), uint8(IRQAny))
}
