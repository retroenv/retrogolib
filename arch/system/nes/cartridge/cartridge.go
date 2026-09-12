// Package cartridge provides .nes ROM loading and saving.
package cartridge

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Cartridge contains a NES cartridge content.
type Cartridge struct {
	PRG     []byte // PRG-ROM banks
	CHR     []byte // CHR-ROM banks
	RAM     byte   // Legacy iNES PRG-RAM size in 8 KiB banks; NES 2.0 uses NES2.RAMSizes.
	Trainer []byte

	// NES2 is nil for legacy iNES files. NES 2.0 RAM sizes are explicit, including zero.
	NES2 *NES2Metadata

	Mapper      uint16     // mapper type
	Mirror      MirrorMode // mirroring mode
	Battery     byte       // battery present
	VideoFormat byte       // 0 NTSC, 1 PAL, 2 multiple-region, 3 Dendy (NES 2.0).
}

// New returns a new cartridge.
func New() *Cartridge {
	return &Cartridge{
		PRG:     make([]byte, 0x8000),
		CHR:     make([]byte, 0x2000),
		Mapper:  0,
		Mirror:  MirrorVertical,
		Battery: 0,
	}
}

// Save writes iNES, or NES 2.0 when the cartridge requires its fields or size encodings.
func (c *Cartridge) Save(writer io.Writer) error {
	header, err := c.fileHeader()
	if err != nil {
		return err
	}

	if err := binary.Write(writer, binary.LittleEndian, header); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}

	if len(c.Trainer) > 0 {
		if err := binary.Write(writer, binary.LittleEndian, c.Trainer); err != nil {
			return fmt.Errorf("writing trainer: %w", err)
		}
	}

	if err := binary.Write(writer, binary.LittleEndian, c.PRG); err != nil {
		return fmt.Errorf("writing PRG: %w", err)
	}

	if len(c.CHR) > 0 {
		if err := binary.Write(writer, binary.LittleEndian, c.CHR); err != nil {
			return fmt.Errorf("writing CHR: %w", err)
		}
	}

	if c.NES2 != nil {
		if _, err := writer.Write(c.NES2.MiscROM); err != nil {
			return fmt.Errorf("writing miscellaneous ROM: %w", err)
		}
	}

	return nil
}
