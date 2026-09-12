package cartridge

import (
	"bytes"
	"io"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestNES2MetadataRoundtrip(t *testing.T) {
	rom := make([]byte, 16+512+0x8000+0x2000+3)
	copy(rom, []byte{'N', 'E', 'S', 0x1A, 2, 1, 0xAE, 0xA9, 0x52, 0, 0x97, 0x68, 3, 0x42, 1, 0x21})
	rom[16], rom[528], rom[528+0x8000] = 0xAA, 0xBB, 0xCC
	copy(rom[len(rom)-3:], []byte{1, 2, 3})

	cart, err := LoadFile(bytes.NewReader(rom))
	assert.NoError(t, err)
	assert.Equal(t, uint16(682), cart.Mapper)
	assert.Equal(t, Mirror4, cart.Mirror)
	assert.Equal(t, byte(1), cart.Battery)
	assert.Equal(t, byte(3), cart.VideoFormat)
	assert.Equal(t, byte(0), cart.RAM)
	assert.Equal(t, byte(0xAA), cart.Trainer[0])
	assert.Equal(t, byte(0xBB), cart.PRG[0])
	assert.Equal(t, byte(0xCC), cart.CHR[0])
	assert.Equal(t, &NES2Metadata{
		Submapper: 5,
		RAMSizes: RAMSizes{
			PRGVolatile:    8192,
			PRGNonvolatile: 32768,

			CHRVolatile:    16384,
			CHRNonvolatile: 4096,
		},

		ConsoleType:     1,
		ConsoleData:     0x42,
		ExpansionDevice: 0x21,

		MiscROMCount: 1,
		MiscROM:      []byte{1, 2, 3},
	}, cart.NES2)

	var output bytes.Buffer
	assert.NoError(t, cart.Save(&output))
	assert.Equal(t, rom[:16], output.Bytes()[:16])
	assert.True(t, bytes.Equal(rom[16:], output.Bytes()[16:]))
}

func TestNES2ROMSizeRoundtrip(t *testing.T) {
	for _, test := range []struct {
		name                 string
		prg, chr             int
		lowPRG, lowCHR, high byte
	}{
		{"linear extension", 8 * 1024 * 1024, 2 * 1024 * 1024, 0, 0, 0x12},
		{"exponent", 3 * 1024, 5 * 512, 10<<2 | 1, 9<<2 | 2, 0xFF},
	} {
		t.Run(test.name, func(t *testing.T) {
			rom := make([]byte, 16+test.prg+test.chr)
			copy(rom, []byte{'N', 'E', 'S', 0x1A, test.lowPRG, test.lowCHR, 0, 8, 0, test.high})
			rom[16] = 0xAA
			if test.chr != 0 {
				rom[16+test.prg] = 0xBB
			}

			cart, err := LoadFile(bytes.NewReader(rom))
			assert.NoError(t, err)
			assert.Len(t, cart.PRG, test.prg)
			assert.Len(t, cart.CHR, test.chr)

			var output bytes.Buffer
			assert.NoError(t, cart.Save(&output))
			assert.Equal(t, rom[:16], output.Bytes()[:16])
			assert.True(t, bytes.Equal(rom[16:], output.Bytes()[16:]))
		})
	}
}

func TestNES2SizeValidation(t *testing.T) {
	low, high, err := encodeROMSize(maxROMSize, 16384)
	assert.NoError(t, err)
	assert.Equal(t, byte(26<<2), low)
	assert.Equal(t, byte(15), high)
	size, err := decodeROMSize(low, high, 16384)
	assert.NoError(t, err)
	assert.Equal(t, maxROMSize, size)

	for _, low := range []byte{63 << 2, 63<<2 | 3, 62<<2 | 2, 27 << 2, 26<<2 | 1} {
		_, err := decodeROMSize(low, 15, 16384)
		assert.ErrorIs(t, err, errROMSize)
	}
	for _, size := range []int{-1, maxROMSize + 1, 11, 0xF00 * 16384} {
		_, _, err := encodeROMSize(size, 16384)
		assert.ErrorIs(t, err, errROMSize)
	}
	for _, size := range []int{-1, 64, 129, 3 * 8192, (64 << 15) + 1} {
		_, err := encodeRAMSizes(RAMSizes{PRGVolatile: size})
		assert.ErrorIs(t, err, errRAMSize)
	}
	for shift := byte(0); shift <= 15; shift++ {
		size := decodeRAMSize(shift)
		encoded, err := encodeRAMSizes(RAMSizes{size, size, size, size})
		assert.NoError(t, err)
		assert.Equal(t, [2]byte{shift<<4 | shift, shift<<4 | shift}, encoded)
	}
}

func TestNES2TruncatedData(t *testing.T) {
	rom := make([]byte, 16+512+16384+8192)
	copy(rom, []byte{'N', 'E', 'S', 0x1A, 1, 1, 4, 8})
	for _, length := range []int{0, 15, 16, 527, 528, 528 + 16383, len(rom) - 1} {
		_, err := LoadFile(bytes.NewReader(rom[:length]))
		assert.Error(t, err)
	}
}

func TestNES2ZeroRAMAndLegacyMetadata(t *testing.T) {
	for _, format := range []byte{0, 8} {
		rom := []byte{'N', 'E', 'S', 0x1A, 0, 0, 0, format, 0, 0, 0, 0, 0, 0, 0, 0}
		cart, err := LoadFile(bytes.NewReader(rom))
		assert.NoError(t, err)
		if format == 0 {
			assert.Nil(t, cart.NES2)
		} else {
			assert.NotNil(t, cart.NES2)
			assert.Equal(t, RAMSizes{}, cart.NES2.RAMSizes)
		}
	}
}

func TestNES2SaveValidation(t *testing.T) {
	for _, cart := range []*Cartridge{
		{Mapper: 0x1000},
		{Trainer: make([]byte, 1)},
		{NES2: &NES2Metadata{Submapper: 16}},
		{NES2: &NES2Metadata{ConsoleType: 4}},
		{VideoFormat: 4},
		{NES2: &NES2Metadata{MiscROMCount: 1}},
		{NES2: &NES2Metadata{MiscROM: []byte{1}}},
		{NES2: &NES2Metadata{ExpansionDevice: 64}},
	} {
		assert.ErrorIs(t, cart.Save(io.Discard), errMetadata)
	}
	assert.ErrorIs(t, (&Cartridge{PRG: make([]byte, 11)}).Save(io.Discard), errROMSize)
}

func TestNES2PromotesLegacyRAM(t *testing.T) {
	for _, battery := range []byte{0, 1} {
		cart := &Cartridge{Mapper: 682, RAM: 1, Battery: battery}
		var output bytes.Buffer
		assert.NoError(t, cart.Save(&output))
		loaded, err := LoadFile(bytes.NewReader(output.Bytes()))
		assert.NoError(t, err)
		if battery == 0 {
			assert.Equal(t, 8192, loaded.NES2.RAMSizes.PRGVolatile)
		} else {
			assert.Equal(t, 8192, loaded.NES2.RAMSizes.PRGNonvolatile)
		}
	}
}

func TestNES2Mirroring(t *testing.T) {
	for _, test := range []struct {
		flags  byte
		mirror MirrorMode
	}{
		{0, MirrorHorizontal}, {1, MirrorVertical}, {8, Mirror4}, {9, Mirror4},
	} {
		rom := []byte{'N', 'E', 'S', 0x1A, 0, 0, test.flags, 8, 0, 0, 0, 0, 0, 0, 0, 0}
		cart, err := LoadFile(bytes.NewReader(rom))
		assert.NoError(t, err)
		assert.Equal(t, test.mirror, cart.Mirror)
	}
}

func TestNES2RejectsInvalidAreas(t *testing.T) {
	for _, rom := range [][]byte{
		{'N', 'E', 'S', 0x1A, 63<<2 | 3, 0, 0, 8, 0, 15, 0, 0, 0, 0, 0, 0},
		{'N', 'E', 'S', 0x1A, 0, 63<<2 | 3, 0, 8, 0, 0xF0, 0, 0, 0, 0, 0, 0},
		{'N', 'E', 'S', 0x1A, 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 1, 0},
	} {
		_, err := LoadFile(bytes.NewReader(rom))
		assert.ErrorIs(t, err, errROMSize)
	}
}
