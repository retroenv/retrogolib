// Package register contains constants for Atari 2600 hardware register addresses.
package register

// TIA write registers ($00-$2C).
// These are written by the CPU to control video and audio output.
const (
	VSYNC  = 0x00 // Vertical sync set-clear
	VBLANK = 0x01 // Vertical blank set-clear
	WSYNC  = 0x02 // Wait for leading edge of horizontal blank
	RSYNC  = 0x03 // Reset horizontal sync counter
	NUSIZ0 = 0x04 // Number-size player-missile 0
	NUSIZ1 = 0x05 // Number-size player-missile 1
	COLUP0 = 0x06 // Color-luminance player 0
	COLUP1 = 0x07 // Color-luminance player 1
	COLUPF = 0x08 // Color-luminance playfield
	COLUBK = 0x09 // Color-luminance background
	CTRLPF = 0x0A // Control playfield ball size and collisions
	REFP0  = 0x0B // Reflect player 0
	REFP1  = 0x0C // Reflect player 1
	PF0    = 0x0D // Playfield register byte 0
	PF1    = 0x0E // Playfield register byte 1
	PF2    = 0x0F // Playfield register byte 2
	RESP0  = 0x10 // Reset player 0
	RESP1  = 0x11 // Reset player 1
	RESM0  = 0x12 // Reset missile 0
	RESM1  = 0x13 // Reset missile 1
	RESBL  = 0x14 // Reset ball
	AUDC0  = 0x15 // Audio control 0
	AUDC1  = 0x16 // Audio control 1
	AUDF0  = 0x17 // Audio frequency 0
	AUDF1  = 0x18 // Audio frequency 1
	AUDV0  = 0x19 // Audio volume 0
	AUDV1  = 0x1A // Audio volume 1
	GRP0   = 0x1B // Graphics player 0
	GRP1   = 0x1C // Graphics player 1
	ENAM0  = 0x1D // Graphics enable missile 0
	ENAM1  = 0x1E // Graphics enable missile 1
	ENABL  = 0x1F // Graphics enable ball
	HMP0   = 0x20 // Horizontal motion player 0
	HMP1   = 0x21 // Horizontal motion player 1
	HMM0   = 0x22 // Horizontal motion missile 0
	HMM1   = 0x23 // Horizontal motion missile 1
	HMBL   = 0x24 // Horizontal motion ball
	VDELP0 = 0x25 // Vertical delay player 0
	VDELP1 = 0x26 // Vertical delay player 1
	VDELBL = 0x27 // Vertical delay ball
	RESMP0 = 0x28 // Reset missile 0 to player 0
	RESMP1 = 0x29 // Reset missile 1 to player 1
	HMOVE  = 0x2A // Apply horizontal motion
	HMCLR  = 0x2B // Clear horizontal motion registers
	CXCLR  = 0x2C // Clear collision latches
)

// TIAWriteCount is the number of TIA write registers.
const TIAWriteCount = 45

// TIA read registers ($00-$0D).
// These are read by the CPU to check collisions and input state.
// Only the upper bits are valid; lower bits return undefined values.
const (
	CXM0P  = 0x00 // Read collision M0-P1, M0-P0 (bits 7-6)
	CXM1P  = 0x01 // Read collision M1-P0, M1-P1 (bits 7-6)
	CXP0FB = 0x02 // Read collision P0-PF, P0-BL (bits 7-6)
	CXP1FB = 0x03 // Read collision P1-PF, P1-BL (bits 7-6)
	CXM0FB = 0x04 // Read collision M0-PF, M0-BL (bits 7-6)
	CXM1FB = 0x05 // Read collision M1-PF, M1-BL (bits 7-6)
	CXBLPF = 0x06 // Read collision BL-PF (bit 7 only)
	CXPPMM = 0x07 // Read collision P0-P1, M0-M1 (bits 7-6)
	INPT0  = 0x08 // Read paddle 0 input (bit 7)
	INPT1  = 0x09 // Read paddle 1 input (bit 7)
	INPT2  = 0x0A // Read paddle 2 input (bit 7)
	INPT3  = 0x0B // Read paddle 3 input (bit 7)
	INPT4  = 0x0C // Read joystick 0 trigger (bit 7)
	INPT5  = 0x0D // Read joystick 1 trigger (bit 7)
)

// TIAReadCount is the number of TIA read registers.
const TIAReadCount = 14

// TIAWriteNames maps TIA write register addresses to their names.
var TIAWriteNames = map[uint16]string{
	VSYNC:  "VSYNC",
	VBLANK: "VBLANK",
	WSYNC:  "WSYNC",
	RSYNC:  "RSYNC",
	NUSIZ0: "NUSIZ0",
	NUSIZ1: "NUSIZ1",
	COLUP0: "COLUP0",
	COLUP1: "COLUP1",
	COLUPF: "COLUPF",
	COLUBK: "COLUBK",
	CTRLPF: "CTRLPF",
	REFP0:  "REFP0",
	REFP1:  "REFP1",
	PF0:    "PF0",
	PF1:    "PF1",
	PF2:    "PF2",
	RESP0:  "RESP0",
	RESP1:  "RESP1",
	RESM0:  "RESM0",
	RESM1:  "RESM1",
	RESBL:  "RESBL",
	AUDC0:  "AUDC0",
	AUDC1:  "AUDC1",
	AUDF0:  "AUDF0",
	AUDF1:  "AUDF1",
	AUDV0:  "AUDV0",
	AUDV1:  "AUDV1",
	GRP0:   "GRP0",
	GRP1:   "GRP1",
	ENAM0:  "ENAM0",
	ENAM1:  "ENAM1",
	ENABL:  "ENABL",
	HMP0:   "HMP0",
	HMP1:   "HMP1",
	HMM0:   "HMM0",
	HMM1:   "HMM1",
	HMBL:   "HMBL",
	VDELP0: "VDELP0",
	VDELP1: "VDELP1",
	VDELBL: "VDELBL",
	RESMP0: "RESMP0",
	RESMP1: "RESMP1",
	HMOVE:  "HMOVE",
	HMCLR:  "HMCLR",
	CXCLR:  "CXCLR",
}

// TIAReadNames maps TIA read register addresses to their names.
var TIAReadNames = map[uint16]string{
	CXM0P:  "CXM0P",
	CXM1P:  "CXM1P",
	CXP0FB: "CXP0FB",
	CXP1FB: "CXP1FB",
	CXM0FB: "CXM0FB",
	CXM1FB: "CXM1FB",
	CXBLPF: "CXBLPF",
	CXPPMM: "CXPPMM",
	INPT0:  "INPT0",
	INPT1:  "INPT1",
	INPT2:  "INPT2",
	INPT3:  "INPT3",
	INPT4:  "INPT4",
	INPT5:  "INPT5",
}
