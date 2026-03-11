package cartridge

import (
	"bytes"
	"testing"

	"github.com/retroenv/retrogolib/arch/system/atari2600"
	"github.com/retroenv/retrogolib/assert"
)

func TestDetectScheme(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		size   int
		scheme BankingScheme
	}{
		{"2K no banking", atari2600.CartridgeSize2K, SchemeNone},
		{"4K no banking", atari2600.CartridgeSize4K, SchemeNone},
		{"8K F8", atari2600.CartridgeSize8K, SchemeF8},
		{"12K FA", atari2600.CartridgeSize12K, SchemeFA},
		{"16K F6", atari2600.CartridgeSize16K, SchemeF6},
		{"32K F4", atari2600.CartridgeSize32K, SchemeF4},
		{"64K 3F", atari2600.CartridgeSize64K, Scheme3F},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			scheme, err := DetectScheme(tt.size)
			assert.NoError(t, err)
			assert.Equal(t, tt.scheme, scheme)
		})
	}
}

func TestDetectSchemeInvalidSize(t *testing.T) {
	t.Parallel()

	invalidSizes := []int{0, 1000, 3000, 5000, 100000}
	for _, size := range invalidSizes {
		_, err := DetectScheme(size)
		assert.ErrorContains(t, err, "unsupported ROM size")
	}
}

func TestLoad2K(t *testing.T) {
	t.Parallel()

	rom := make([]byte, atari2600.CartridgeSize2K)
	rom[0] = 0xEA // NOP marker

	cart, err := Load(bytes.NewReader(rom))
	assert.NoError(t, err)
	assert.Equal(t, SchemeNone, cart.Scheme)
	assert.Equal(t, 1, cart.Banks)

	// 2K ROM is mirrored to fill 4K.
	assert.Len(t, cart.ROM, atari2600.CartridgeSize4K)
	assert.Equal(t, byte(0xEA), cart.ROM[0])
	assert.Equal(t, byte(0xEA), cart.ROM[atari2600.CartridgeSize2K])
}

func TestLoad4K(t *testing.T) {
	t.Parallel()

	rom := make([]byte, atari2600.CartridgeSize4K)
	rom[0] = 0x4C // JMP marker

	cart, err := Load(bytes.NewReader(rom))
	assert.NoError(t, err)
	assert.Equal(t, SchemeNone, cart.Scheme)
	assert.Equal(t, 1, cart.Banks)
	assert.Len(t, cart.ROM, atari2600.CartridgeSize4K)
	assert.Equal(t, byte(0x4C), cart.ROM[0])
}

func TestLoadBankedROMs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		size   int
		scheme BankingScheme
		banks  int
	}{
		{"8K F8", atari2600.CartridgeSize8K, SchemeF8, 2},
		{"12K FA", atari2600.CartridgeSize12K, SchemeFA, 3},
		{"16K F6", atari2600.CartridgeSize16K, SchemeF6, 4},
		{"32K F4", atari2600.CartridgeSize32K, SchemeF4, 8},
		{"64K 3F", atari2600.CartridgeSize64K, Scheme3F, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rom := make([]byte, tt.size)
			rom[0] = 0xAA // marker in bank 0

			cart, err := Load(bytes.NewReader(rom))
			assert.NoError(t, err)
			assert.Equal(t, tt.scheme, cart.Scheme)
			assert.Equal(t, tt.banks, cart.Banks)
			assert.Len(t, cart.ROM, tt.size)
		})
	}
}

func TestLoadEmpty(t *testing.T) {
	t.Parallel()

	_, err := Load(bytes.NewReader([]byte{}))
	assert.ErrorContains(t, err, "empty ROM")
}

func TestLoadInvalidSize(t *testing.T) {
	t.Parallel()

	rom := make([]byte, 5000) // not a standard size
	_, err := Load(bytes.NewReader(rom))
	assert.ErrorContains(t, err, "unsupported ROM size")
}

func TestBankOffset(t *testing.T) {
	t.Parallel()

	rom := make([]byte, atari2600.CartridgeSize16K)
	cart, err := Load(bytes.NewReader(rom))
	assert.NoError(t, err)

	// Bank 0 starts at offset 0.
	offset, err := cart.BankOffset(0)
	assert.NoError(t, err)
	assert.Equal(t, 0, offset)

	// Bank 1 starts at 4K.
	offset, err = cart.BankOffset(1)
	assert.NoError(t, err)
	assert.Equal(t, atari2600.ROMWindowSize, offset)

	// Bank 3 starts at 12K.
	offset, err = cart.BankOffset(3)
	assert.NoError(t, err)
	assert.Equal(t, 3*atari2600.ROMWindowSize, offset)

	// Bank 4 is out of range.
	_, err = cart.BankOffset(4)
	assert.ErrorContains(t, err, "out of range")

	// Negative bank is out of range.
	_, err = cart.BankOffset(-1)
	assert.ErrorContains(t, err, "out of range")
}

func TestTriggerBankF8(t *testing.T) {
	t.Parallel()

	rom := make([]byte, atari2600.CartridgeSize8K)
	cart, err := Load(bytes.NewReader(rom))
	assert.NoError(t, err)

	assert.Equal(t, 0, cart.TriggerBank(0x1FF8))
	assert.Equal(t, 1, cart.TriggerBank(0x1FF9))
	assert.Equal(t, -1, cart.TriggerBank(0x1FF7)) // below range
	assert.Equal(t, -1, cart.TriggerBank(0x1FFA)) // above range
	assert.Equal(t, -1, cart.TriggerBank(0x1000)) // not a trigger
}

func TestTriggerBankFA(t *testing.T) {
	t.Parallel()

	rom := make([]byte, atari2600.CartridgeSize12K)
	cart, err := Load(bytes.NewReader(rom))
	assert.NoError(t, err)

	assert.Equal(t, 0, cart.TriggerBank(0x1FF8))
	assert.Equal(t, 1, cart.TriggerBank(0x1FF9))
	assert.Equal(t, 2, cart.TriggerBank(0x1FFA))
	assert.Equal(t, -1, cart.TriggerBank(0x1FFB))
}

func TestTriggerBankF6(t *testing.T) {
	t.Parallel()

	rom := make([]byte, atari2600.CartridgeSize16K)
	cart, err := Load(bytes.NewReader(rom))
	assert.NoError(t, err)

	assert.Equal(t, 0, cart.TriggerBank(0x1FF6))
	assert.Equal(t, 1, cart.TriggerBank(0x1FF7))
	assert.Equal(t, 2, cart.TriggerBank(0x1FF8))
	assert.Equal(t, 3, cart.TriggerBank(0x1FF9))
	assert.Equal(t, -1, cart.TriggerBank(0x1FF5)) // below range
	assert.Equal(t, -1, cart.TriggerBank(0x1FFA)) // above range
}

func TestTriggerBankF4(t *testing.T) {
	t.Parallel()

	rom := make([]byte, atari2600.CartridgeSize32K)
	cart, err := Load(bytes.NewReader(rom))
	assert.NoError(t, err)

	for i := range 8 {
		assert.Equal(t, i, cart.TriggerBank(uint16(0x1FF4+i)))
	}
	assert.Equal(t, -1, cart.TriggerBank(0x1FF3)) // below range
	assert.Equal(t, -1, cart.TriggerBank(0x1FFC)) // above range (reset vector)
}

func TestTriggerBank3F(t *testing.T) {
	t.Parallel()

	rom := make([]byte, atari2600.CartridgeSize64K)
	cart, err := Load(bytes.NewReader(rom))
	assert.NoError(t, err)

	assert.Equal(t, 0, cart.TriggerBank(0x003F))
	assert.Equal(t, -1, cart.TriggerBank(0x003E)) // not the trigger
	assert.Equal(t, -1, cart.TriggerBank(0x0040)) // not the trigger
}

func TestTriggerBankNone(t *testing.T) {
	t.Parallel()

	rom := make([]byte, atari2600.CartridgeSize4K)
	cart, err := Load(bytes.NewReader(rom))
	assert.NoError(t, err)

	// No banking scheme: all addresses return -1.
	assert.Equal(t, -1, cart.TriggerBank(0x1FF8))
	assert.Equal(t, -1, cart.TriggerBank(0x003F))
	assert.Equal(t, -1, cart.TriggerBank(0x1000))
}

func TestBankingSchemeString(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "None", SchemeNone.String())
	assert.Equal(t, "F8", SchemeF8.String())
	assert.Equal(t, "FA", SchemeFA.String())
	assert.Equal(t, "F6", SchemeF6.String())
	assert.Equal(t, "F4", SchemeF4.String())
	assert.Equal(t, "3F", Scheme3F.String())
	assert.Equal(t, "BankingScheme(99)", BankingScheme(99).String())
}
