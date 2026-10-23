// Package valueanalysisprobe checks proposed scalar and register-view contracts.
package valueanalysisprobe

import (
	"testing"

	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/assert"
)

type z80Resource uint16

type dataRegister string

type vfResource uint8

func TestProbeZ80PairWritePreservesSibling(t *testing.T) {
	t.Parallel()

	parent, valid := knownbits.NewConstant[uint16](16, 0x9934)
	assert.True(t, valid)
	part, valid := knownbits.NewConstant[uint16](8, 0x12)
	assert.True(t, valid)
	upper := analysis.View[z80Resource]{
		Storage: 0,
		Offset:  8,
		Bits:    8,
	}
	lower := analysis.View[z80Resource]{
		Storage: 0,
		Bits:    8,
	}
	whole := analysis.View[z80Resource]{
		Storage: 0,
		Bits:    16,
	}
	result, valid := knownbits.Insert(16, parent, upper.Offset, upper.Bits, part)
	assert.True(t, valid)
	expected, valid := knownbits.NewConstant[uint16](16, 0x1234)
	assert.True(t, valid)
	assert.Equal(t, expected, result)
	assert.True(t, whole.Overlaps(upper))
	assert.True(t, whole.Overlaps(lower))
	assert.False(t, upper.Overlaps(lower))
}

func TestProbePartialDataRegisterWrite(t *testing.T) {
	t.Parallel()

	parent, valid := knownbits.NewConstant[uint32](32, 0x12345678)
	assert.True(t, valid)
	part, valid := knownbits.NewConstant[uint32](8, 0xaa)
	assert.True(t, valid)
	lowByte := analysis.View[dataRegister]{
		Storage: "D0",
		Bits:    8,
	}
	result, valid := knownbits.Insert(32, parent, lowByte.Offset, lowByte.Bits, part)
	assert.True(t, valid)
	expected, valid := knownbits.NewConstant[uint32](32, 0x123456aa)
	assert.True(t, valid)
	assert.Equal(t, expected, result)
	unknown, valid := knownbits.New[uint32](32, 0, 0x78)
	assert.True(t, valid)
	result, valid = knownbits.Insert(32, unknown, lowByte.Offset, lowByte.Bits, part)
	assert.True(t, valid)
	assert.Equal(t, uint32(0xaa), result.KnownOne)
	assert.Equal(t, uint32(0x55), result.KnownZero)
	lowByte.Offset = 31
	_, valid = knownbits.Insert(32, parent, lowByte.Offset, lowByte.Bits, part)
	assert.False(t, valid)
}

func TestProbeUnknownModeDoesNotSelectWidth(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		mode  knownbits.Bit
		bits  knownbits.Width
		valid bool
	}{
		{mode: knownbits.BitUnknown},
		{
			mode:  knownbits.BitClear,
			bits:  16,
			valid: true,
		},
		{
			mode:  knownbits.BitSet,
			bits:  8,
			valid: true,
		},
	} {
		bits, valid := accumulatorWidth(test.mode)
		assert.Equal(t, test.valid, valid)
		assert.Equal(t, test.bits, bits)
	}
}

func TestProbeVFIsAnOrdinaryByteResource(t *testing.T) {
	t.Parallel()

	before, valid := knownbits.NewConstant[uint8](8, 0x81)
	assert.True(t, valid)
	carry, valid := knownbits.NewConstant[uint8](8, before.KnownOne>>7)
	assert.True(t, valid)
	vf := analysis.View[vfResource]{
		Storage: 0,
		Bits:    8,
	}
	result, valid := knownbits.Insert(8, before, vf.Offset, vf.Bits, carry)
	assert.True(t, valid)
	expected, valid := knownbits.NewConstant[uint8](8, 1)
	assert.True(t, valid)
	assert.Equal(t, expected, result)
	assert.Equal(t, uint8(0x81), before.KnownOne)
}

func accumulatorWidth(mode knownbits.Bit) (knownbits.Width, bool) {
	switch mode {
	case knownbits.BitClear:
		return 16, true
	case knownbits.BitSet:
		return 8, true
	default:
		return 0, false
	}
}
