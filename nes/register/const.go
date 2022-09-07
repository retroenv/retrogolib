// Package register contains constants that represent special memory register addresses.
package register

// APU (Audio Processing Unit) constants
const (
	SQ1_VOL       = 0x4000
	SQ1_SWEEP     = 0x4001
	SQ1_LO        = 0x4002
	SQ1_HI        = 0x4003
	SQ2_VOL       = 0x4004
	SQ2_SWEEP     = 0x4005
	SQ2_LO        = 0x4006
	SQ2_HI        = 0x4007
	TRI_LINEAR    = 0x4008
	TRI_LO        = 0x400A
	TRI_HI        = 0x400B
	NOISE_VOL     = 0x400C
	NOISE_LO      = 0x400E
	NOISE_HI      = 0x400F
	APU_DMC_CTRL  = 0x4010
	APU_CHAN_CTRL = 0x4015
	APU_FRAME     = 0x4017
)

// Controller constants
const (
	JOYPAD1 = 0x4016
	JOYPAD2 = 0x4017
)

// PPU constants
const (
	PPU_CTRL   = 0x2000
	PPU_MASK   = 0x2001
	PPU_STATUS = 0x2002
	OAM_ADDR   = 0x2003
	OAM_DATA   = 0x2004
	PPU_SCROLL = 0x2005
	PPU_ADDR   = 0x2006
	PPU_DATA   = 0x2007

	PALETTE_START = 0x3f00

	OAM_DMA = 0x4014
)
