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
	assert.Equal(t, uint16(0xD7FF), uint16(VIARegionEnd))
	assert.Equal(t, uint16(0xD800), uint16(DualSelectStart))
	assert.Equal(t, uint16(0xDFFF), uint16(DualSelectEnd))

	// System ROM
	assert.Equal(t, uint16(0xE000), uint16(ROMStart))
	assert.Equal(t, uint16(0xFFFF), uint16(ROMEnd))
	assert.Equal(t, 0x2000, ROMSize)
}

func TestVectorAddresses(t *testing.T) {
	// The 6809 vector table occupies the last 16 bytes of the system ROM.
	assert.Equal(t, uint16(0xFFFE), uint16(ResetVector))
	assert.Equal(t, uint16(0xFFFC), uint16(NMIVector))
	assert.Equal(t, uint16(0xFFFA), uint16(SWIVector))
	assert.Equal(t, uint16(0xFFF8), uint16(IRQVector))
	assert.Equal(t, uint16(0xFFF6), uint16(FIRQVector))
	assert.Equal(t, uint16(0xFFF4), uint16(SWI2Vector))
	assert.Equal(t, uint16(0xFFF2), uint16(SWI3Vector))
	assert.True(t, SWI3Vector >= ROMStart)
}

func TestVIAMirrorMask(t *testing.T) {
	// Every address in the VIA region selects one of the 16 registers.
	for addr := uint16(VIAStart); addr <= VIARegionEnd; addr++ {
		offset := addr & VIAMirrorMask
		assert.True(t, uint16(VIAStart)+offset <= VIAEnd)
	}
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
		{"VIA", VIAStart, VIARegionEnd},
		{"Dual select", DualSelectStart, DualSelectEnd},
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
