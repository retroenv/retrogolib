// Package arch provides architecture constants and types.
package arch

import (
	"strings"

	"github.com/retroenv/retrogolib/set"
)

// Architecture represents a target CPU architecture.
type Architecture string

// Supported CPU architectures and virtual machines.
const (
	// CHIP8 represents the Chip-8 virtual machine (not a CPU architecture).
	// The Chip-8 is an interpreted programming language used on:
	// - COSMAC VIP
	// - Telmac 1800
	// - ETI 660
	// - Modern homebrew systems and emulators
	CHIP8 Architecture = "chip8"

	// M6502 represents the MOS Technology 6502 processor used in:
	// - Apple II series
	// - Commodore 64/128/VIC-20
	// - Atari 2600/5200/8-bit computers
	// - Nintendo Entertainment System (NES/Famicom)
	M6502 Architecture = "6502"

	// M65C02 represents the WDC 65C02 processor used in:
	// - Apple IIe/IIc
	// - Atari Lynx
	// - TurboGrafx-16/PC Engine
	M65C02 Architecture = "65c02"

	// M65816 represents the WDC 65C816 (65816) processor used in:
	// - Super Nintendo Entertainment System (SNES/Super Famicom)
	// - Apple IIGS
	M65816 Architecture = "65816"

	// M6809 represents the Motorola 6809 processor used in:
	// - TRS-80 Color Computer (CoCo)
	// - Vectrex
	// - Dragon 32/64
	// - Williams arcade hardware (Defender, Robotron, Joust)
	M6809 Architecture = "6809"

	// M68000 represents the Motorola 68000 processor used in:
	// - Sega Genesis/Mega Drive
	// - Commodore Amiga
	// - Atari ST
	// - Apple Macintosh (original)
	M68000 Architecture = "m68000"

	// X86 represents the Intel x86 processor family used in:
	// - IBM PC and compatibles
	// - MS-DOS and compatible operating systems
	// - Early Windows systems
	// - Embedded x86 systems
	X86 Architecture = "x86"

	// SM83 represents the Sharp SM83 (LR35902) processor used in:
	// - Nintendo Game Boy
	// - Nintendo Game Boy Color
	SM83 Architecture = "sm83"

	// Z80 represents the Zilog Z80 processor used in:
	// - ZX Spectrum
	// - Amstrad CPC
	// - MSX computers
	// - Sega Master System/Game Gear
	Z80 Architecture = "z80"
)

// allSupportedArchitectures defines the single source of truth for supported architectures.
// Adding a new architecture requires updating only this slice.
var allSupportedArchitectures = []Architecture{
	CHIP8,
	M6502,
	M65C02,
	M65816,
	M6809,
	M68000,
	SM83,
	X86,
	Z80,
}

// supportedArchitecturesSet provides O(1) lookup performance for IsValid().
var supportedArchitecturesSet = set.NewFromSlice(allSupportedArchitectures)

// String returns the string representation of the architecture.
func (a Architecture) String() string {
	return string(a)
}

// IsValid returns true if the architecture is supported.
func (a Architecture) IsValid() bool {
	return supportedArchitecturesSet.Contains(a)
}

// FromString creates an Architecture from a string.
// Returns the architecture and true if valid, or empty Architecture and false if invalid.
// The comparison is case-insensitive.
func FromString(s string) (Architecture, bool) {
	arch := Architecture(strings.ToLower(s))
	if arch.IsValid() {
		return arch, true
	}
	return "", false
}

// SupportedArchitectures returns a slice of all supported architectures.
func SupportedArchitectures() []Architecture {
	// Return a copy to prevent external mutation
	result := make([]Architecture, len(allSupportedArchitectures))
	copy(result, allSupportedArchitectures)
	return result
}
