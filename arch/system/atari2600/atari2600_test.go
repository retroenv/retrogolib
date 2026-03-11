package atari2600

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestAddressConstants(t *testing.T) {
	assert.Equal(t, uint16(0x1FFF), uint16(AddressMask))
	assert.Equal(t, uint16(0x2000), uint16(AddressSpaceSize))
}

func TestMemoryMapRanges(t *testing.T) {
	// TIA
	assert.True(t, TIAWriteStart == 0x00)
	assert.True(t, TIAWriteEnd == 0x2C)

	// RAM
	assert.Equal(t, uint16(0x0080), uint16(RAMStart))
	assert.Equal(t, uint16(0x00FF), uint16(RAMEnd))
	assert.Equal(t, 128, RAMSize)
	assert.Equal(t, RAMSize, int(RAMEnd-RAMStart+1))

	// RAM mirror
	assert.Equal(t, uint16(0x0180), uint16(RAMMirrorStart))
	assert.Equal(t, uint16(0x01FF), uint16(RAMMirrorEnd))

	// RIOT I/O
	assert.Equal(t, uint16(0x0280), uint16(RIOTStart))
	assert.Equal(t, uint16(0x0297), uint16(RIOTEnd))

	// ROM
	assert.Equal(t, uint16(0x1000), uint16(ROMStart))
	assert.Equal(t, uint16(0x1FFF), uint16(ROMEnd))
	assert.Equal(t, 0x1000, ROMWindowSize)
}

func TestCartridgeSizes(t *testing.T) {
	assert.Equal(t, 2048, CartridgeSize2K)
	assert.Equal(t, 4096, CartridgeSize4K)
	assert.Equal(t, 8192, CartridgeSize8K)
	assert.Equal(t, 16384, CartridgeSize16K)
	assert.Equal(t, 32768, CartridgeSize32K)
	assert.Equal(t, 65536, CartridgeSize64K)
}

func TestResetVector(t *testing.T) {
	// Reset vector must be within the ROM window
	assert.True(t, ResetVector >= ROMStart)
	assert.True(t, ResetVector <= ROMEnd)
}

func TestAddressMasking(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		address  uint16
		expected uint16
	}{
		{"low address unchanged", 0x0080, 0x0080},
		{"ROM address unchanged", 0x1000, 0x1000},
		{"max address unchanged", 0x1FFF, 0x1FFF},
		{"mirror at 0x2000", 0x2000, 0x0000},
		{"mirror at 0x3000", 0x3000, 0x1000},
		{"mirror at 0x4000", 0x4000, 0x0000},
		{"mirror at 0xF000", 0xF000, 0x1000},
		{"mirror at 0xFFFC", 0xFFFC, 0x1FFC},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			masked := tt.address & AddressMask
			assert.Equal(t, tt.expected, masked)
		})
	}
}

func TestTIAMirroring(t *testing.T) {
	t.Parallel()

	// TIA decodes only the low 6 bits, so addresses mirror every 64 bytes.
	// Write register VSYNC ($00) mirrors at $40, $80, $C0, etc.
	addresses := []uint16{0x0000, 0x0040, 0x0080, 0x00C0, 0x0100}
	for _, addr := range addresses {
		masked := addr & TIAMirrorMask
		assert.Equal(t, uint16(0x0000), masked,
			"address 0x%04X should mirror to TIA register 0x00", addr)
	}

	// CXCLR ($2C) mirrors at $6C, $AC, etc.
	assert.Equal(t, uint16(0x002C), uint16(0x006C)&TIAMirrorMask)
	assert.Equal(t, uint16(0x002C), uint16(0x00AC)&TIAMirrorMask)
}

func TestRAMMirrorRelationship(t *testing.T) {
	t.Parallel()

	// RAM mirror starts exactly $100 above RAM.
	assert.Equal(t, uint16(RAMStart+0x0100), uint16(RAMMirrorStart))
	assert.Equal(t, uint16(RAMEnd+0x0100), uint16(RAMMirrorEnd))

	// Both regions have the same size.
	ramSize := RAMEnd - RAMStart + 1
	mirrorSize := RAMMirrorEnd - RAMMirrorStart + 1
	assert.Equal(t, ramSize, mirrorSize)
}

func TestMemoryRegionNoOverlap(t *testing.T) {
	t.Parallel()

	// Verify key memory regions don't overlap within the 13-bit address space.
	regions := []struct {
		name  string
		start uint16
		end   uint16
	}{
		{"TIA", TIAWriteStart, TIAWriteEnd},
		{"RAM", RAMStart, RAMEnd},
		{"RAM mirror", RAMMirrorStart, RAMMirrorEnd},
		{"RIOT", RIOTStart, RIOTEnd},
		{"ROM", ROMStart, ROMEnd},
	}

	for i := range regions {
		for j := i + 1; j < len(regions); j++ {
			a := regions[i]
			b := regions[j]
			overlaps := a.start <= b.end && b.start <= a.end
			assert.True(t, !overlaps,
				"%s (0x%04X-0x%04X) overlaps with %s (0x%04X-0x%04X)",
				a.name, a.start, a.end, b.name, b.start, b.end)
		}
	}
}
