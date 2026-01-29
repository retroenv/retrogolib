package arch

import (
	"strings"

	"github.com/retroenv/retrogolib/set"
)

// System represents a complete retro computing system.
// This is separate from CPU architecture and handles system-specific
// concerns like executable format, system calls, and runtime constraints.
type System string

// Supported systems.
const (
	// CHIP8System represents the Chip-8 virtual machine system.
	// Used on COSMAC VIP, Telmac 1800, ETI 660, and modern emulators.
	CHIP8System System = "chip8"

	// DOS represents MS-DOS and compatible systems.
	DOS System = "dos"

	// GameBoy represents the Nintendo Game Boy handheld system.
	// Includes original Game Boy, Game Boy Pocket, and Game Boy Color compatibility.
	GameBoy System = "gameboy"

	// Generic represents a generic system without specific hardware quirks.
	// Can be used for any CPU architecture when no system-specific behavior is needed.
	Generic System = "generic"

	// NES represents the Nintendo Entertainment System (Famicom).
	NES System = "nes"

	// ZXSpectrum represents the Sinclair ZX Spectrum home computer series.
	ZXSpectrum System = "zx-spectrum"
)

// allSupportedSystems defines the single source of truth for supported systems.
// Adding a new system requires updating only this slice.
var allSupportedSystems = []System{
	CHIP8System,
	DOS,
	GameBoy,
	Generic,
	NES,
	ZXSpectrum,
}

// supportedSystemsSet provides O(1) lookup performance for system validation.
var supportedSystemsSet = set.NewFromSlice(allSupportedSystems)

// String returns the string representation of the system.
func (s System) String() string {
	return string(s)
}

// IsValid returns true if the system is supported.
func (s System) IsValid() bool {
	return supportedSystemsSet.Contains(s)
}

// SystemFromString creates a System from a string.
// Returns the system and true if valid, or empty System and false if invalid.
// The comparison is case-insensitive.
func SystemFromString(s string) (System, bool) {
	sys := System(strings.ToLower(s))
	if sys.IsValid() {
		return sys, true
	}
	return "", false
}

// SupportedSystems returns a slice of all supported systems.
func SupportedSystems() []System {
	// Return a copy to prevent external mutation
	result := make([]System, len(allSupportedSystems))
	copy(result, allSupportedSystems)
	return result
}
