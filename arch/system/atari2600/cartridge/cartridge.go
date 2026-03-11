// Package cartridge provides Atari 2600 ROM loading and banking scheme detection.
//
// Atari 2600 ROMs are raw binary files with no header. The banking scheme
// is determined by the ROM size: ROMs up to 4 KB fit in the ROM window
// without banking, while larger ROMs use various bank switching schemes
// triggered by accessing specific addresses.
package cartridge

import (
	"errors"
	"fmt"
	"io"

	"github.com/retroenv/retrogolib/arch/system/atari2600"
)

// BankingScheme identifies the bank switching method used by a cartridge.
type BankingScheme int

// Banking schemes, determined by ROM size.
const (
	SchemeNone BankingScheme = iota // 2 KB or 4 KB, no banking
	SchemeF8                        // 8 KB, 2 banks, triggers at $1FF8-$1FF9
	SchemeFA                        // 12 KB, 3 banks, triggers at $1FF8-$1FFA
	SchemeF6                        // 16 KB, 4 banks, triggers at $1FF6-$1FF9
	SchemeF4                        // 32 KB, 8 banks, triggers at $1FF4-$1FFB
	Scheme3F                        // 64 KB, 16 banks, Tigervision, write bank to $003F
)

// bankingSchemeNames maps banking schemes to their display names.
var bankingSchemeNames = map[BankingScheme]string{
	SchemeNone: "None",
	SchemeF8:   "F8",
	SchemeFA:   "FA",
	SchemeF6:   "F6",
	SchemeF4:   "F4",
	Scheme3F:   "3F",
}

// String returns the name of the banking scheme.
func (b BankingScheme) String() string {
	if name, ok := bankingSchemeNames[b]; ok {
		return name
	}
	return fmt.Sprintf("BankingScheme(%d)", int(b))
}

// Bank switching trigger address ranges within the ROM window.
// Bank switching is triggered by accessing (reading) these addresses.
const (
	// F8 scheme: 8 KB, 2 banks.
	F8TriggerStart = 0x1FF8
	F8TriggerEnd   = 0x1FF9

	// FA scheme: 12 KB, 3 banks.
	FATriggerStart = 0x1FF8
	FATriggerEnd   = 0x1FFA

	// F6 scheme: 16 KB, 4 banks.
	F6TriggerStart = 0x1FF6
	F6TriggerEnd   = 0x1FF9

	// F4 scheme: 32 KB, 8 banks.
	F4TriggerStart = 0x1FF4
	F4TriggerEnd   = 0x1FFB

	// 3F (Tigervision) scheme: write bank number to $003F.
	Trigger3F = 0x003F
)

// Cartridge contains an Atari 2600 cartridge ROM.
type Cartridge struct {
	ROM    []byte        // raw ROM data
	Scheme BankingScheme // detected banking scheme
	Banks  int           // number of 4 KB banks
}

// Load reads a raw Atari 2600 ROM binary and detects the banking scheme.
func Load(reader io.Reader) (*Cartridge, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("reading ROM: %w", err)
	}

	if len(data) == 0 {
		return nil, errors.New("empty ROM")
	}

	scheme, err := DetectScheme(len(data))
	if err != nil {
		return nil, err
	}

	// Pad 2 KB ROMs to 4 KB by mirroring.
	rom := data
	if len(data) == atari2600.CartridgeSize2K {
		rom = make([]byte, atari2600.CartridgeSize4K)
		copy(rom, data)
		copy(rom[atari2600.CartridgeSize2K:], data)
	}

	banks := len(rom) / atari2600.ROMWindowSize
	if banks == 0 {
		banks = 1
	}

	return &Cartridge{
		ROM:    rom,
		Scheme: scheme,
		Banks:  banks,
	}, nil
}

// DetectScheme returns the banking scheme for the given ROM size.
func DetectScheme(size int) (BankingScheme, error) {
	switch size {
	case atari2600.CartridgeSize2K, atari2600.CartridgeSize4K:
		return SchemeNone, nil
	case atari2600.CartridgeSize8K:
		return SchemeF8, nil
	case atari2600.CartridgeSize12K:
		return SchemeFA, nil
	case atari2600.CartridgeSize16K:
		return SchemeF6, nil
	case atari2600.CartridgeSize32K:
		return SchemeF4, nil
	case atari2600.CartridgeSize64K:
		return Scheme3F, nil
	default:
		return SchemeNone, fmt.Errorf("unsupported ROM size: %d bytes", size)
	}
}

// BankOffset returns the byte offset into the ROM for the given bank number.
// Returns an error if the bank number is out of range.
func (c *Cartridge) BankOffset(bank int) (int, error) {
	if bank < 0 || bank >= c.Banks {
		return 0, fmt.Errorf("bank %d out of range (0-%d)", bank, c.Banks-1)
	}
	return bank * atari2600.ROMWindowSize, nil
}

// triggerRange maps each banking scheme to its trigger address range.
var triggerRange = map[BankingScheme][2]uint16{
	SchemeF8: {F8TriggerStart, F8TriggerEnd},
	SchemeFA: {FATriggerStart, FATriggerEnd},
	SchemeF6: {F6TriggerStart, F6TriggerEnd},
	SchemeF4: {F4TriggerStart, F4TriggerEnd},
}

// TriggerBank returns the bank number selected by accessing the given address,
// or -1 if the address is not a bank switching trigger for this cartridge's scheme.
func (c *Cartridge) TriggerBank(address uint16) int {
	if c.Scheme == Scheme3F {
		// Tigervision: the bank number is written as data to $003F,
		// not encoded in the trigger address itself. Return 0 to
		// signal that the address matches; the caller must read the
		// written value to determine the actual bank.
		if address == Trigger3F {
			return 0
		}
		return -1
	}

	r, ok := triggerRange[c.Scheme]
	if ok && address >= r[0] && address <= r[1] {
		return int(address - r[0])
	}
	return -1
}
