package register

import "github.com/retroenv/retrogolib/arch/cpu/m6502"

// APUAddressToName maps address constants from address to name.
var APUAddressToName = map[uint16]m6502.AccessModeConstant{
	APU_PL1_VOL:    {Constant: "APU_PL1_VOL", Mode: m6502.WriteAccess},
	APU_PL1_SWEEP:  {Constant: "APU_PL1_SWEEP", Mode: m6502.WriteAccess},
	APU_PL1_LO:     {Constant: "APU_PL1_LO", Mode: m6502.WriteAccess},
	APU_PL1_HI:     {Constant: "APU_PL1_HI", Mode: m6502.WriteAccess},
	APU_PL2_VOL:    {Constant: "APU_PL2_VOL", Mode: m6502.WriteAccess},
	APU_PL2_SWEEP:  {Constant: "APU_PL2_SWEEP", Mode: m6502.WriteAccess},
	APU_PL2_LO:     {Constant: "APU_PL2_LO", Mode: m6502.WriteAccess},
	APU_PL2_HI:     {Constant: "APU_PL2_HI", Mode: m6502.WriteAccess},
	APU_TRI_LINEAR: {Constant: "APU_TRI_LINEAR", Mode: m6502.WriteAccess},
	APU_TRI_LO:     {Constant: "APU_TRI_LO", Mode: m6502.WriteAccess},
	APU_TRI_HI:     {Constant: "APU_TRI_HI", Mode: m6502.WriteAccess},
	APU_NOISE_VOL:  {Constant: "APU_NOISE_VOL", Mode: m6502.WriteAccess},
	APU_NOISE_LO:   {Constant: "APU_NOISE_LO", Mode: m6502.WriteAccess},
	APU_NOISE_HI:   {Constant: "APU_NOISE_HI", Mode: m6502.WriteAccess},
	APU_DMC_FREQ:   {Constant: "APU_DMC_FREQ", Mode: m6502.WriteAccess},
	APU_DMC_RAW:    {Constant: "APU_DMC_RAW", Mode: m6502.ReadWriteAccess},
	APU_DMC_START:  {Constant: "APU_DMC_START", Mode: m6502.ReadWriteAccess},
	APU_DMC_LEN:    {Constant: "APU_DMC_LEN", Mode: m6502.ReadWriteAccess},
	APU_SND_CHN:    {Constant: "APU_SND_CHN", Mode: m6502.ReadWriteAccess},
	APU_FRAME:      {Constant: "APU_FRAME", Mode: m6502.WriteAccess},
}

// ControllerAddressToName maps address constants from address to name.
var ControllerAddressToName = map[uint16]m6502.AccessModeConstant{
	JOYPAD1: {Constant: "JOYPAD1", Mode: m6502.ReadWriteAccess},
	JOYPAD2: {Constant: "JOYPAD2", Mode: m6502.ReadAccess},
}

// PPUAddressToName maps address constants from address to name.
var PPUAddressToName = map[uint16]m6502.AccessModeConstant{
	PPU_CTRL:   {Constant: "PPU_CTRL", Mode: m6502.WriteAccess},
	PPU_MASK:   {Constant: "PPU_MASK", Mode: m6502.WriteAccess},
	PPU_STATUS: {Constant: "PPU_STATUS", Mode: m6502.ReadAccess},
	OAM_ADDR:   {Constant: "OAM_ADDR", Mode: m6502.WriteAccess},
	OAM_DATA:   {Constant: "OAM_DATA", Mode: m6502.ReadWriteAccess},
	PPU_SCROLL: {Constant: "PPU_SCROLL", Mode: m6502.WriteAccess},
	PPU_ADDR:   {Constant: "PPU_ADDR", Mode: m6502.WriteAccess},
	PPU_DATA:   {Constant: "PPU_DATA", Mode: m6502.ReadWriteAccess},

	PALETTE_START: {Constant: "PALETTE_START", Mode: m6502.ReadWriteAccess},

	OAM_DMA: {Constant: "OAM_DMA", Mode: m6502.WriteAccess},
}
