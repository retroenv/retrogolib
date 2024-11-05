package cartridge

import (
	"bytes"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func testRom() []byte {
	b := []byte{iNESFileMagic[0], iNESFileMagic[1], iNESFileMagic[2], iNESFileMagic[3]}
	b = append(b, []byte{2, 1, 1, 0, 0}...)       // prg, chr, control 1, control 2, ram
	b = append(b, []byte{0, 0, 0, 0, 0, 0, 0}...) // reserved/padding

	prg := make([]byte, 2*16384)
	prg[0] = 0x80 // marker
	b = append(b, prg...)

	chr := make([]byte, 8192)
	chr[0] = 0x81 // marker
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
