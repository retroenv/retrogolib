// Package arch provides architecture constants and types.
package arch

import (
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

	// Z80 represents the Zilog Z80 processor used in:
	// - ZX Spectrum
	// - Amstrad CPC
	// - MSX computers
	// - Game Boy (modified Z80)
	// - Sega Master System/Game Gear
	Z80 Architecture = "z80"
)

// allSupportedArchitectures defines the single source of truth for supported architectures.
// Adding a new architecture requires updating only this slice.
var allSupportedArchitectures = []Architecture{
	CHIP8,
	M6502,
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
func FromString(s string) (Architecture, bool) {
	arch := Architecture(s)
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
