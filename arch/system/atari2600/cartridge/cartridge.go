// Package cartridge provides Atari 2600 ROM loading and banking scheme detection.
//
// Atari 2600 ROMs are raw binary files with no header. The banking scheme
// is inferred from the ROM size. Several schemes share a size, so this is a
// heuristic rather than definitive detection. This package provides ROM metadata
// and hotspot lookup; the host implements the actual bank switching.
package cartridge

import (
	"errors"
	"fmt"
	"io"

	"github.com/retroenv/retrogolib/arch/system/atari2600"
)

// BankingScheme identifies the bank switching method used by a cartridge.
type BankingScheme int

// Supported banking schemes.
const (
	SchemeNone BankingScheme = iota // 2 KB or 4 KB, no banking
	SchemeF8                        // 8 KB, 2 banks, triggers at $1FF8-$1FF9
	SchemeFA                        // 12 KB, 3 banks, triggers at $1FF8-$1FFA
	SchemeF6                        // 16 KB, 4 banks, triggers at $1FF6-$1FF9
	SchemeF4                        // 32 KB, 8 banks, triggers at $1FF4-$1FFB
	Scheme3F                        // Tigervision: 2 KB banks, writes to $0000-$003F select the lower segment.
)

// Bank switching hotspots; 3F responds only to writes.
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

// bankingSchemeNames maps banking schemes to their display names.
var bankingSchemeNames = map[BankingScheme]string{
	SchemeNone: "None",
	SchemeF8:   "F8",
	SchemeFA:   "FA",
	SchemeF6:   "F6",
	SchemeF4:   "F4",
	Scheme3F:   "3F",
}

// triggerRange maps each banking scheme to its trigger address range.
var triggerRange = map[BankingScheme][2]uint16{
	SchemeF8: {F8TriggerStart, F8TriggerEnd},
	SchemeFA: {FATriggerStart, FATriggerEnd},
	SchemeF6: {F6TriggerStart, F6TriggerEnd},
	SchemeF4: {F4TriggerStart, F4TriggerEnd},
}

// Cartridge contains an Atari 2600 cartridge ROM.
type Cartridge struct {
	ROM    []byte        // Raw ROM data.
	Scheme BankingScheme // Inferred banking scheme.
	Banks  int           // Number of banks: 2 KB for 3F, 4 KB for other schemes.
}

// String returns the name of the banking scheme.
func (bs BankingScheme) String() string {
	if name, ok := bankingSchemeNames[bs]; ok {
		return name
	}
	return fmt.Sprintf("BankingScheme(%d)", int(bs))
}

// BankOffset returns the byte offset into the ROM for the given bank number.
// Returns an error if the bank number is out of range.
func (cart *Cartridge) BankOffset(bank int) (int, error) {
	if bank < 0 || bank >= cart.Banks {
		return 0, fmt.Errorf("bank %d out of range (0-%d)", bank, cart.Banks-1)
	}
	return bank * cart.Scheme.bankSize(), nil
}

// TriggerBank returns the bank number selected by accessing the given address,
// or -1 if the address is not a bank switching trigger for this cartridge's scheme.
// For 3F, a result of 0 identifies a write hotspot; the written data selects the bank.
func (cart *Cartridge) TriggerBank(address uint16) int {
	address &= atari2600.AddressMask
	if cart.Scheme == Scheme3F {
		if address <= Trigger3F {
			return 0
		}
		return -1
	}

	r, ok := triggerRange[cart.Scheme]
	if ok && address >= r[0] && address <= r[1] {
		return int(address - r[0])
	}
	return -1
}

func (bs BankingScheme) bankSize() int {
	if bs == Scheme3F {
		return atari2600.ROMWindowSize / 2
	}
	return atari2600.ROMWindowSize
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

	return &Cartridge{
		ROM:    rom,
		Scheme: scheme,
		Banks:  len(rom) / scheme.bankSize(),
	}, nil
}

// DetectScheme returns the default banking scheme for a supported ROM size.
// Size alone cannot distinguish every cartridge format.
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
