package cartridge

import (
	"errors"
	"fmt"
	"io"
	"math/bits"
)

// The loader limits each ROM area to 64 MiB to bound header-controlled allocations.
const maxROMSize = 64 * 1024 * 1024

var (
	errROMSize  = errors.New("unsupported ROM size")
	errRAMSize  = errors.New("unrepresentable NES 2.0 RAM size")
	errMetadata = errors.New("invalid cartridge metadata")
)

// RAMSizes contains RAM sizes in bytes. Zero means that the memory is absent.
type RAMSizes struct {
	PRGVolatile    int
	PRGNonvolatile int

	CHRVolatile    int
	CHRNonvolatile int
}

// NES2Metadata contains the additional NES 2.0 cartridge information.
// ConsoleData is header byte 13: Vs. System data or an extended console type.
// See https://www.nesdev.org/wiki/NES_2.0 for the field encodings.
type NES2Metadata struct {
	Submapper byte
	RAMSizes  RAMSizes

	ConsoleType     byte
	ConsoleData     byte
	ExpansionDevice byte

	MiscROMCount byte
	MiscROM      []byte
}

func (c *Cartridge) fileHeader() (header, error) {
	if c.Mapper > 0xFFF || (len(c.Trainer) != 0 && len(c.Trainer) != 512) {
		return header{}, errMetadata
	}

	hdr := header{
		Magic:       iNESFileMagic,
		NumPRG:      byte(len(c.PRG) / 16384),
		NumCHR:      byte(len(c.CHR) / 8192),
		NumRAM:      c.RAM,
		VideoFormat: c.VideoFormat,
	}
	hdr.Control1, hdr.Control2 = ControlBytes(c.Battery, byte(c.Mirror), c.Mapper, len(c.Trainer) > 0)
	if c.NES2 == nil && c.Mapper <= 0xFF && c.VideoFormat <= 1 &&
		len(c.PRG)%16384 == 0 && len(c.PRG)/16384 <= 255 &&
		len(c.CHR)%8192 == 0 && len(c.CHR)/8192 <= 255 {

		return hdr, nil
	}

	return c.nes2Header(hdr)
}

func (c *Cartridge) nes2Header(hdr header) (header, error) {
	metadata := c.NES2
	if metadata == nil {
		metadata = &NES2Metadata{}
		if c.Battery != 0 {
			metadata.RAMSizes.PRGNonvolatile = int(c.RAM) * 8192
		} else {
			metadata.RAMSizes.PRGVolatile = int(c.RAM) * 8192
		}
	}
	if metadata.Submapper > 15 || metadata.ConsoleType > 3 || c.VideoFormat > 3 ||
		metadata.MiscROMCount > 3 || metadata.ExpansionDevice > 63 ||
		(metadata.MiscROMCount == 0) != (len(metadata.MiscROM) == 0) || len(metadata.MiscROM) > maxROMSize {

		return header{}, errMetadata
	}

	prgLow, prgHigh, err := encodeROMSize(len(c.PRG), 16384)
	if err != nil {
		return header{}, fmt.Errorf("encoding PRG: %w", err)
	}

	chrLow, chrHigh, err := encodeROMSize(len(c.CHR), 8192)
	if err != nil {
		return header{}, fmt.Errorf("encoding CHR: %w", err)
	}

	ram, err := encodeRAMSizes(metadata.RAMSizes)
	if err != nil {
		return header{}, err
	}

	hdr.NumPRG, hdr.NumCHR = prgLow, chrLow
	hdr.Control2 = (hdr.Control2 & 0xF0) | 8 | metadata.ConsoleType
	hdr.NumRAM = metadata.Submapper<<4 | byte(c.Mapper>>8)
	hdr.VideoFormat = chrHigh<<4 | prgHigh
	hdr.Reserved = [6]byte{ram[0], ram[1], c.VideoFormat, metadata.ConsoleData,
		metadata.MiscROMCount, metadata.ExpansionDevice}
	return hdr, nil
}

func loadNES2(reader io.Reader, hdr header) (*Cartridge, error) {
	prgSize, err := decodeROMSize(hdr.NumPRG, hdr.VideoFormat&15, 16384)
	if err != nil {
		return nil, fmt.Errorf("decoding PRG: %w", err)
	}

	chrSize, err := decodeROMSize(hdr.NumCHR, hdr.VideoFormat>>4, 8192)
	if err != nil {
		return nil, fmt.Errorf("decoding CHR: %w", err)
	}

	cart := &Cartridge{
		Mapper:      uint16(hdr.Control1>>4) | uint16(hdr.Control2&0xF0) | uint16(hdr.NumRAM&15)<<8,
		Mirror:      MirrorMode(hdr.Control1 & verticalMirroringFlag),
		Battery:     (hdr.Control1 & batteryFlag) >> 1,
		VideoFormat: hdr.Reserved[2] & 3,
		NES2:        decodeNES2Metadata(hdr),
	}
	if hdr.Control1&fourScreenFlag != 0 {
		cart.Mirror = Mirror4
	}

	if err := readNES2Data(reader, cart, hdr.Control1&trainerFlag != 0, prgSize, chrSize); err != nil {
		return nil, err
	}
	return cart, nil
}

func decodeNES2Metadata(hdr header) *NES2Metadata {
	return &NES2Metadata{
		Submapper: hdr.NumRAM >> 4,
		RAMSizes: RAMSizes{
			PRGVolatile:    decodeRAMSize(hdr.Reserved[0] & 15),
			PRGNonvolatile: decodeRAMSize(hdr.Reserved[0] >> 4),

			CHRVolatile:    decodeRAMSize(hdr.Reserved[1] & 15),
			CHRNonvolatile: decodeRAMSize(hdr.Reserved[1] >> 4),
		},

		ConsoleType:     hdr.Control2 & 3,
		ConsoleData:     hdr.Reserved[3],
		ExpansionDevice: hdr.Reserved[5] & 63,
		MiscROMCount:    hdr.Reserved[4] & 3,
	}
}

func readNES2Data(reader io.Reader, cart *Cartridge, trainer bool, prgSize, chrSize int) error {
	areas := []struct {
		name string
		size int
		data *[]byte
	}{
		{name: "trainer", data: &cart.Trainer},
		{name: "PRG", size: prgSize, data: &cart.PRG},
		{name: "CHR", size: chrSize, data: &cart.CHR},
	}
	if trainer {
		areas[0].size = 512
	}

	for _, area := range areas {
		if area.size == 0 {
			continue
		}
		*area.data = make([]byte, area.size)
		if _, err := io.ReadFull(reader, *area.data); err != nil {
			return fmt.Errorf("reading %s: %w", area.name, err)
		}
	}

	if cart.NES2.MiscROMCount != 0 {
		data, err := io.ReadAll(io.LimitReader(reader, maxROMSize+1))
		if err != nil {
			return fmt.Errorf("reading miscellaneous ROM: %w", err)
		}
		if len(data) == 0 || len(data) > maxROMSize {
			return errROMSize
		}
		cart.NES2.MiscROM = data
	}
	return nil
}

func decodeROMSize(low, high byte, unit int) (int, error) {
	size := (uint64(high)<<8 | uint64(low)) * uint64(unit)
	if high == 15 {
		// NES 2.0 uses 2^E * (2*M+1) bytes when the high nibble is 15.
		// Check the exponent before multiplication to prevent uint64 overflow.
		exponent, multiplier := low>>2, uint64(low&3)*2+1
		if exponent > 26 {
			return 0, errROMSize
		}
		size = (uint64(1) << exponent) * multiplier
	}
	if size > maxROMSize {
		return 0, errROMSize
	}
	return int(size), nil
}

func encodeROMSize(size, unit int) (byte, byte, error) {
	if size < 0 || size > maxROMSize {
		return 0, 0, errROMSize
	}
	// A high nibble of 15 selects the exponent form, so linear counts stop at 0xEFF.
	if size%unit == 0 && size/unit < 0xF00 {
		return byte(size / unit), byte((size / unit) >> 8), nil
	}

	exponent := bits.TrailingZeros(uint(size))
	multiplier := size >> exponent
	if multiplier > 7 {
		return 0, 0, errROMSize
	}
	return byte(exponent<<2 | (multiplier-1)/2), 15, nil
}

func decodeRAMSize(shift byte) int {
	if shift == 0 {
		return 0
	}
	return 64 << shift
}

func encodeRAMSizes(sizes RAMSizes) ([2]byte, error) {
	var result [2]byte

	for i, size := range []int{sizes.PRGVolatile, sizes.PRGNonvolatile, sizes.CHRVolatile, sizes.CHRNonvolatile} {
		if size == 0 {
			continue
		}
		if size < 128 || size > 64<<15 || size&(size-1) != 0 {
			return result, errRAMSize
		}
		result[i/2] |= byte(bits.TrailingZeros(uint(size))-6) << (4 * (i % 2))
	}
	return result, nil
}
