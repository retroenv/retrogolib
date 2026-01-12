package cartridge

import (
	"bytes"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func testRom() []byte {
	prg := make([]byte, 2*16384)
	prg[0] = 0x80 // marker

	chr := make([]byte, 8192)
	chr[0] = 0x81 // marker

	b := make([]byte, 0, 16+len(prg)+len(chr))
	b = append(b, iNESFileMagic[:]...)
	b = append(b, []byte{2, 1, 1, 0, 0}...)       // prg, chr, control 1, control 2, ram
	b = append(b, []byte{0, 0, 0, 0, 0, 0, 0}...) // reserved/padding
	b = append(b, prg...)
	b = append(b, chr...)

	return b
}

func TestLoadFile(t *testing.T) {
	rom := testRom()
	reader := bytes.NewReader(rom)

	cart, err := LoadFile(reader)
	assert.NoError(t, err)

	assert.Equal(t, 0x80, cart.PRG[0])
	assert.Equal(t, 0x81, cart.CHR[0])
	assert.Equal(t, 0, cart.Mapper)
	assert.Equal(t, 1, cart.Mirror)
	assert.Equal(t, 0, cart.Battery)
}

func TestCartridgeSave(t *testing.T) {
	c := New()
	c.PRG[0] = 0x80 // marker
	c.CHR[0] = 0x81 // marker

	var buf bytes.Buffer
	assert.NoError(t, c.Save(&buf))

	rom := testRom()
	b := buf.Bytes()
	assert.Equal(t, rom, b)
}

func TestLoadBuffer(t *testing.T) {
	// test a small buffer
	buf := []byte{0x60}
	reader := bytes.NewReader(buf)

	cart, err := LoadBuffer(reader)
	assert.NoError(t, err)
	assert.Equal(t, 0x60, cart.PRG[0])

	// test a large buffer
	buf = make([]byte, 16384+1000)
	buf[0] = 0x60
	reader = bytes.NewReader(buf)

	cart, err = LoadBuffer(reader)
	assert.NoError(t, err)
	assert.Equal(t, 0x60, cart.PRG[0])
}

func TestInvalidROM(t *testing.T) {
	t.Parallel()

	// Test empty file
	reader := bytes.NewReader([]byte{})
	_, err := LoadFile(reader)
	assert.ErrorContains(t, err, "reading header")

	// Test invalid magic
	invalidRom := make([]byte, 0, 16)
	invalidRom = append(invalidRom, []byte{0x4E, 0x45, 0x53, 0x00}...) // Invalid magic
	invalidRom = append(invalidRom, []byte{2, 1, 1, 0, 0}...)
	invalidRom = append(invalidRom, make([]byte, 7)...) // padding

	reader = bytes.NewReader(invalidRom)
	_, err = LoadFile(reader)
	assert.ErrorContains(t, err, "magic")
}

func TestCartridgeProperties(t *testing.T) {
	t.Parallel()

	// Test cartridge with different properties
	rom := make([]byte, 0, 4+5+7+16384)
	rom = append(rom, iNESFileMagic[:]...)
	rom = append(rom, []byte{1, 0, 0x02, 0x08, 0}...) // prg=1, chr=0, mapper=0, mirror=vertical, battery=1
	rom = append(rom, make([]byte, 7)...)             // padding
	rom = append(rom, make([]byte, 16384)...)         // PRG ROM

	reader := bytes.NewReader(rom)
	cart, err := LoadFile(reader)
	assert.NoError(t, err)

	assert.Equal(t, 0, cart.Mapper)
	assert.Equal(t, 0, cart.Mirror) // 0 means horizontal mirroring in the control byte
	assert.Equal(t, 1, cart.Battery)
	assert.Equal(t, 16384, len(cart.PRG))
	assert.Equal(t, 0, len(cart.CHR)) // No CHR ROM
}

func TestNewCartridge(t *testing.T) {
	t.Parallel()

	cart := New()
	assert.NotNil(t, cart)
	assert.Equal(t, 0, cart.Mapper)
	assert.Equal(t, MirrorVertical, cart.Mirror) // Default is vertical
	assert.Equal(t, 0, cart.Battery)

	// Check default sizes
	assert.Equal(t, 32768, len(cart.PRG)) // Default PRG size
	assert.Equal(t, 8192, len(cart.CHR))  // Default CHR size
}

func TestSaveLoadRoundtrip(t *testing.T) {
	t.Parallel()

	// Create original cartridge with specific data
	original := New()
	original.Mapper = 1
	original.Mirror = MirrorHorizontal
	original.Battery = 1

	// Add some test data
	original.PRG[0] = 0xAA
	original.PRG[100] = 0xBB
	original.CHR[0] = 0xCC
	original.CHR[50] = 0xDD

	// Save to buffer
	var buf bytes.Buffer
	assert.NoError(t, original.Save(&buf))

	// Load from buffer
	reader := bytes.NewReader(buf.Bytes())
	loaded, err := LoadFile(reader)
	assert.NoError(t, err)

	// Verify properties match
	assert.Equal(t, original.Mapper, loaded.Mapper)
	assert.Equal(t, original.Mirror, loaded.Mirror)
	assert.Equal(t, original.Battery, loaded.Battery)

	// Verify data matches
	assert.Equal(t, original.PRG[0], loaded.PRG[0])
	assert.Equal(t, original.PRG[100], loaded.PRG[100])
	assert.Equal(t, original.CHR[0], loaded.CHR[0])
	assert.Equal(t, original.CHR[50], loaded.CHR[50])
}
