package coco

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestAddressConstants(t *testing.T) {
	assert.Equal(t, 0x10000, AddressSpaceSize)
}

func TestMemoryMapRanges(t *testing.T) {
	// RAM
	assert.Equal(t, uint16(0x0000), uint16(RAMStart))
	assert.Equal(t, uint16(0x7FFF), uint16(RAMEnd))
	assert.Equal(t, 32768, RAMSize32K)

	// Extended BASIC ROM
	assert.Equal(t, uint16(0x8000), uint16(ExtendedBASICStart))
	assert.Equal(t, uint16(0x9FFF), uint16(ExtendedBASICEnd))
	assert.Equal(t, 0x2000, ExtendedBASICSize)

	// Color BASIC ROM
	assert.Equal(t, uint16(0xA000), uint16(ColorBASICStart))
	assert.Equal(t, uint16(0xBFFF), uint16(ColorBASICEnd))
	assert.Equal(t, 0x2000, ColorBASICSize)

	// Cartridge ROM
	assert.Equal(t, uint16(0xC000), uint16(CartridgeStart))
	assert.Equal(t, uint16(0xFEFF), uint16(CartridgeEnd))
}

func TestIOAddressRanges(t *testing.T) {
	assert.Equal(t, uint16(0xFF00), uint16(PIA0Start))
	assert.Equal(t, uint16(0xFF03), uint16(PIA0End))
	assert.Equal(t, uint16(0xFF20), uint16(PIA1Start))
	assert.Equal(t, uint16(0xFF23), uint16(PIA1End))
	assert.Equal(t, uint16(0xFFC0), uint16(SAMStart))
	assert.Equal(t, uint16(0xFFDF), uint16(SAMEnd))
}

func TestVectorAddresses(t *testing.T) {
	assert.Equal(t, uint16(0xFFFE), uint16(ResetVector))
	assert.Equal(t, uint16(0xFFFC), uint16(NMIVector))
	assert.Equal(t, uint16(0xFFFA), uint16(SWIVector))
	assert.Equal(t, uint16(0xFFF8), uint16(IRQVector))
	assert.Equal(t, uint16(0xFFF6), uint16(FIRQVector))
}

func TestCartridgeSizes(t *testing.T) {
	assert.Equal(t, 8192, CartridgeSize8K)
	assert.Equal(t, 16384, CartridgeSize16K)
	assert.Equal(t, 32768, CartridgeSize32K)
}

func TestMemoryRegionNoOverlap(t *testing.T) {
	t.Parallel()

	regions := []struct {
		name  string
		start uint16
		end   uint16
	}{
		{"RAM", RAMStart, RAMEnd},
		{"Extended BASIC", ExtendedBASICStart, ExtendedBASICEnd},
		{"Color BASIC", ColorBASICStart, ColorBASICEnd},
		{"PIA0", PIA0Start, PIA0End},
		{"PIA1", PIA1Start, PIA1End},
		{"SAM", SAMStart, SAMEnd},
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
