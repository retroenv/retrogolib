// Package codedatalog reads and writes raw FCEUX NES Code/Data Logger files.
// https://fceux.com/web/help/CodeDataLogger.html
package codedatalog

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/retroenv/retrogolib/arch/system/nes/cartridge"
)

// PrgFlag defines flags for the PRG ROM.
type PrgFlag byte

// ChrFlag defines flags for the CHR ROM.
type ChrFlag byte

const (
	Code              PrgFlag = 0b0000_0001
	Data              PrgFlag = 0b0000_0010
	RomBankMappedLow  PrgFlag = 0b0000_0100
	RomBankMappedHigh PrgFlag = 0b0000_1000
	IndirectCode      PrgFlag = 0b0001_0000
	IndirectData      PrgFlag = 0b0010_0000
	PCMAudio          PrgFlag = 0b0100_0000
	// SubEntryPoint is a legacy Mesen annotation. FCEUX exports omit it.
	SubEntryPoint PrgFlag = 0b1000_0000

	DrawnOnScreen        ChrFlag = 0b0000_0001
	ReadProgrammatically ChrFlag = 0b0000_0010
)

// ErrUnsupportedFormat identifies a CDL variant that this package cannot read.
var ErrUnsupportedFormat = errors.New("unsupported cdl format")

// prgMappingShift moves CPU address bits A14-A13 to CDL bits 3-2.
const prgMappingShift = 11

const mesen2Signature = "CDLv2"

var prgFlagNames = []flagName[PrgFlag]{
	{
		flag: Code,
		name: "code",
	},
	{
		flag: Data,
		name: "data",
	},
	{
		flag: RomBankMappedLow,
		name: "rom_bank_mapped_low",
	},
	{
		flag: RomBankMappedHigh,
		name: "rom_bank_mapped_high",
	},
	{
		flag: IndirectCode,
		name: "indirect_code",
	},
	{
		flag: IndirectData,
		name: "indirect_data",
	},
	{
		flag: PCMAudio,
		name: "pcm_audio",
	},
	{
		flag: SubEntryPoint,
		name: "sub_entry_point",
	},
}

var chrFlagNames = []flagName[ChrFlag]{
	{
		flag: DrawnOnScreen,
		name: "drawn_on_screen",
	},
	{
		flag: ReadProgrammatically,
		name: "read_programmatically",
	},
}

// Names returns the names of the set PRG flags in bit order.
func (f PrgFlag) Names() []string {
	return names(f, prgFlagNames)
}

// Names returns the names of the set CHR flags in bit order.
// Bits without a defined name are omitted.
func (f ChrFlag) Names() []string {
	return names(f, chrFlagNames)
}

// WithMapping replaces the FCEUX mapping bits with CPU address bits A14-A13.
// These bits identify the last accessed window at $8000, $A000, $C000, or $E000.
// https://fceux.com/web/help/CodeDataLogger.html
func (f PrgFlag) WithMapping(address uint16) PrgFlag {
	mask := RomBankMappedLow | RomBankMappedHigh
	return f&^mask | PrgFlag(address>>prgMappingShift)&mask
}

// LoadFile returns the PRG flags of a raw FCEUX CDL file.
// It accepts a PRG-only file or a PRG mask followed by a complete CHR mask.
// It rejects Mesen 2 headers because their flag assignments are different.
// https://github.com/SourMesen/Mesen2/blob/master/Core/Debugger/CodeDataLogger.cpp
func LoadFile(cart *cartridge.Cartridge, reader io.Reader) ([]PrgFlag, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("reading cartridge data: %w", err)
	}

	if len(data) != len(cart.PRG) && len(data) != len(cart.PRG)+len(cart.CHR) {
		if bytes.HasPrefix(data, []byte(mesen2Signature)) {
			return nil, fmt.Errorf("mesen 2 CDLv2 file: %w", ErrUnsupportedFormat)
		}
		return nil, errors.New("invalid cdl file size for PRG and CHR ROM")
	}

	prgFlags := make([]PrgFlag, len(cart.PRG))
	for i := range len(cart.PRG) {
		prgFlags[i] = PrgFlag(data[i])
	}

	return prgFlags, nil
}

// Write writes the PRG flags followed by the CHR flags in raw FCEUX CDL format.
// It clears the reserved PRG bit used for legacy Mesen subroutine entries and
// the undefined CHR bits. The input slices are not changed.
func Write(destination io.Writer, prg []PrgFlag, chr []ChrFlag) error {
	if err := writeFlags(destination, prg, "PRG", ^SubEntryPoint); err != nil {
		return err
	}
	return writeFlags(destination, chr, "CHR", DrawnOnScreen|ReadProgrammatically)
}

type flagName[T ~byte] struct {
	flag T
	name string
}

func names[T ~byte](value T, labels []flagName[T]) []string {
	var result []string
	for _, label := range labels {
		if value&label.flag != 0 {
			result = append(result, label.name)
		}
	}

	return result
}

func writeFlags[T ~byte](destination io.Writer, flags []T, space string, mask T) error {
	if len(flags) == 0 {
		return nil
	}
	data := make([]byte, len(flags))
	for index, flag := range flags {
		data[index] = byte(flag & mask)
	}

	written, err := destination.Write(data)
	if err == nil && written != len(data) {
		err = io.ErrShortWrite
	}

	if err != nil {
		return fmt.Errorf("writing %s flags: %w", space, err)
	}
	return nil
}
