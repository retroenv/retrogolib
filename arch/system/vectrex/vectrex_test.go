package vectrex

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestAddressConstants(t *testing.T) {
	assert.Equal(t, 0x10000, AddressSpaceSize)
}

func TestMemoryMapRanges(t *testing.T) {
	// Cartridge ROM
	assert.Equal(t, uint16(0x0000), uint16(CartridgeStart))
	assert.Equal(t, uint16(0x7FFF), uint16(CartridgeEnd))
	assert.Equal(t, 0x8000, CartridgeMaxSize)

	// RAM
	assert.Equal(t, uint16(0xC800), uint16(RAMStart))
	assert.Equal(t, uint16(0xCBFF), uint16(RAMEnd))
	assert.Equal(t, 1024, RAMSize)
	assert.Equal(t, RAMSize, int(RAMEnd-RAMStart+1))

	// VIA
	assert.Equal(t, uint16(0xD000), uint16(VIAStart))
	assert.Equal(t, uint16(0xD00F), uint16(VIAEnd))

	// System ROM
	assert.Equal(t, uint16(0xE000), uint16(ROMStart))
	assert.Equal(t, uint16(0xFFFF), uint16(ROMEnd))
	assert.Equal(t, 0x2000, ROMSize)
}

func TestVectorAddresses(t *testing.T) {
	// All vectors must be within system ROM
	assert.True(t, ResetVector >= ROMStart)
	assert.True(t, NMIVector >= ROMStart)
	assert.True(t, SWIVector >= ROMStart)
	assert.True(t, IRQVector >= ROMStart)
	assert.True(t, FIRQVector >= ROMStart)
}

func TestCartridgeSizes(t *testing.T) {
	assert.Equal(t, 4096, CartridgeSize4K)
	assert.Equal(t, 8192, CartridgeSize8K)
	assert.Equal(t, 16384, CartridgeSize16K)
	assert.Equal(t, 32768, CartridgeSize32K)
}

func TestRAMMirrorRelationship(t *testing.T) {
	// RAM mirror starts exactly $0400 above RAM
	assert.Equal(t, uint16(RAMStart+0x0400), uint16(RAMMirrorStart))

	// Both regions have the same size
	ramSize := RAMEnd - RAMStart + 1
	mirrorSize := RAMMirrorEnd - RAMMirrorStart + 1
	assert.Equal(t, ramSize, mirrorSize)
}

func TestMemoryRegionNoOverlap(t *testing.T) {
	t.Parallel()

	regions := []struct {
		name  string
		start uint16
		end   uint16
	}{
		{"Cartridge", CartridgeStart, CartridgeEnd},
		{"RAM", RAMStart, RAMEnd},
		{"RAM mirror", RAMMirrorStart, RAMMirrorEnd},
		{"VIA", VIAStart, VIAEnd},
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
