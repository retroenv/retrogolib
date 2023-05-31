package register

import (
	. "github.com/retroenv/retrogolib/addressing"
)

// APUAddressToName maps address constants from address to name.
var APUAddressToName = map[uint16]AccessModeConstant{
	APU_PL1_VOL:    {Constant: "APU_PL1_VOL", Mode: WriteAccess},
	APU_PL1_SWEEP:  {Constant: "APU_PL1_SWEEP", Mode: WriteAccess},
	APU_PL1_LO:     {Constant: "APU_PL1_LO", Mode: WriteAccess},
	APU_PL1_HI:     {Constant: "APU_PL1_HI", Mode: WriteAccess},
	APU_PL2_VOL:    {Constant: "APU_PL2_VOL", Mode: WriteAccess},
	APU_PL2_SWEEP:  {Constant: "APU_PL2_SWEEP", Mode: WriteAccess},
	APU_PL2_LO:     {Constant: "APU_PL2_LO", Mode: WriteAccess},
	APU_PL2_HI:     {Constant: "APU_PL2_HI", Mode: WriteAccess},
	APU_TRI_LINEAR: {Constant: "APU_TRI_LINEAR", Mode: WriteAccess},
	APU_TRI_LO:     {Constant: "APU_TRI_LO", Mode: WriteAccess},
	APU_TRI_HI:     {Constant: "APU_TRI_HI", Mode: WriteAccess},
	APU_NOISE_VOL:  {Constant: "APU_NOISE_VOL", Mode: WriteAccess},
	APU_NOISE_LO:   {Constant: "APU_NOISE_LO", Mode: WriteAccess},
	APU_NOISE_HI:   {Constant: "APU_NOISE_HI", Mode: WriteAccess},
	APU_DMC_FREQ:   {Constant: "APU_DMC_FREQ", Mode: WriteAccess},
	APU_DMC_RAW:    {Constant: "APU_DMC_RAW", Mode: ReadWriteAccess},
	APU_DMC_START:  {Constant: "APU_DMC_START", Mode: ReadWriteAccess},
	APU_DMC_LEN:    {Constant: "APU_DMC_LEN", Mode: ReadWriteAccess},
	APU_SND_CHN:    {Constant: "APU_SND_CHN", Mode: ReadWriteAccess},
	APU_FRAME:      {Constant: "APU_FRAME", Mode: WriteAccess},
}

// ControllerAddressToName maps address constants from address to name.
var ControllerAddressToName = map[uint16]AccessModeConstant{
	JOYPAD1: {Constant: "JOYPAD1", Mode: ReadWriteAccess},
	JOYPAD2: {Constant: "JOYPAD2", Mode: ReadAccess},
}

// PPUAddressToName maps address constants from address to name.
var PPUAddressToName = map[uint16]AccessModeConstant{
	PPU_CTRL:   {Constant: "PPU_CTRL", Mode: WriteAccess},
	PPU_MASK:   {Constant: "PPU_MASK", Mode: WriteAccess},
	PPU_STATUS: {Constant: "PPU_STATUS", Mode: ReadAccess},
	OAM_ADDR:   {Constant: "OAM_ADDR", Mode: WriteAccess},
	OAM_DATA:   {Constant: "OAM_DATA", Mode: ReadWriteAccess},
	PPU_SCROLL: {Constant: "PPU_SCROLL", Mode: WriteAccess},
	PPU_ADDR:   {Constant: "PPU_ADDR", Mode: WriteAccess},
	PPU_DATA:   {Constant: "PPU_DATA", Mode: ReadWriteAccess},

	PALETTE_START: {Constant: "PALETTE_START", Mode: ReadWriteAccess},

	OAM_DMA: {Constant: "OAM_DMA", Mode: WriteAccess},
}
