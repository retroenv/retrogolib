package register

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestAPUConstants(t *testing.T) {
	t.Parallel()

	// Test APU register addresses are in expected range
	assert.Equal(t, 0x4000, APU_PL1_VOL)
	assert.Equal(t, 0x4001, APU_PL1_SWEEP)
	assert.Equal(t, 0x4002, APU_PL1_LO)
	assert.Equal(t, 0x4003, APU_PL1_HI)
	assert.Equal(t, 0x4017, APU_FRAME)

	// Test that APU addresses are sequential where expected
	assert.Equal(t, APU_PL1_VOL+1, APU_PL1_SWEEP)
	assert.Equal(t, APU_PL1_SWEEP+1, APU_PL1_LO)
	assert.Equal(t, APU_PL1_LO+1, APU_PL1_HI)
}

func TestControllerConstants(t *testing.T) {
	t.Parallel()

	// Test controller register addresses
	assert.Equal(t, 0x4016, JOYPAD1)
	assert.Equal(t, 0x4017, JOYPAD2)

	// Test that controller addresses are sequential
	assert.Equal(t, JOYPAD1+1, JOYPAD2)
}

func TestPPUConstants(t *testing.T) {
	t.Parallel()

	// Test PPU register addresses
	assert.Equal(t, 0x2000, PPU_CTRL)
	assert.Equal(t, 0x2001, PPU_MASK)
	assert.Equal(t, 0x2002, PPU_STATUS)
	assert.Equal(t, 0x2007, PPU_DATA)

	// Test special PPU addresses
	assert.Equal(t, 0x3f00, PALETTE_START)
	assert.Equal(t, 0x4014, OAM_DMA)

	// Test PPU addresses are in 2000-2007 range for main registers
	assert.True(t, PPU_CTRL >= 0x2000 && PPU_CTRL <= 0x2007)
	assert.True(t, PPU_DATA >= 0x2000 && PPU_DATA <= 0x2007)
}

func TestRegisterRanges(t *testing.T) {
	t.Parallel()

	// Test that register ranges don't overlap inappropriately
	// PPU registers: 0x2000-0x2007
	// APU registers: 0x4000-0x4017
	assert.True(t, PPU_DATA < APU_PL1_VOL)
	assert.True(t, PPU_CTRL < JOYPAD1)
}
