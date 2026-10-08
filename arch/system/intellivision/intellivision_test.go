package intellivision

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestMemoryMap(t *testing.T) {
	assert.Equal(t, 0x10000, AddressSpaceWords)
	assert.Equal(t, 0xFFFF, OpenBusValue)

	regions := []struct {
		name       string
		start, end int
		wantStart  int
		wantEnd    int
	}{
		{"STIC", STICStart, STICEnd, 0x0000, 0x003F},
		{"scratchpad RAM", ScratchpadRAMStart, ScratchpadRAMEnd, 0x0100, 0x01EF},
		{"PSG", PSGStart, PSGEnd, 0x01F0, 0x01FF},
		{"system RAM", SystemRAMStart, SystemRAMEnd, 0x0200, 0x035F},
		{"Executive ROM", ExecutiveROMStart, ExecutiveROMEnd, 0x1000, 0x1FFF},
		{"graphics ROM", GraphicsROMStart, GraphicsROMEnd, 0x3000, 0x37FF},
		{"graphics RAM", GraphicsRAMStart, GraphicsRAMEnd, 0x3800, 0x39FF},
	}

	for i, region := range regions {
		t.Run(region.name, func(t *testing.T) {
			assert.Equal(t, region.wantStart, region.start)
			assert.Equal(t, region.wantEnd, region.end)
			assert.True(t, region.start <= region.end)
			assert.True(t, region.end < AddressSpaceWords)
		})
		if i > 0 {
			assert.True(t, regions[i-1].end < region.start)
		}
	}

	assert.Equal(t, ExecutiveROMStart, ResetAddress)
	assert.Equal(t, ExecutiveROMStart+4, InterruptAddress)
}
