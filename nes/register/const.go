// Package register contains constants that represent special memory register addresses.
package register

// APU (Audio Processing Unit) constants
const (
	APU_PL1_VOL    = 0x4000
	APU_PL1_SWEEP  = 0x4001
	APU_PL1_LO     = 0x4002
	APU_PL1_HI     = 0x4003
	APU_PL2_VOL    = 0x4004
	APU_PL2_SWEEP  = 0x4005
	APU_PL2_LO     = 0x4006
	APU_PL2_HI     = 0x4007
	APU_TRI_LINEAR = 0x4008
	APU_TRI_LO     = 0x400A
	APU_TRI_HI     = 0x400B
	APU_NOISE_VOL  = 0x400C
	APU_NOISE_LO   = 0x400E
	APU_NOISE_HI   = 0x400F
	APU_DMC_FREQ   = 0x4010
	APU_DMC_RAW    = 0x4011
	APU_DMC_START  = 0x4012
	APU_DMC_LEN    = 0x4013
	APU_SND_CHN    = 0x4015
	APU_FRAME      = 0x4017
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
