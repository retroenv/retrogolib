package analysis

import (
	"testing"

	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/z80"
	"github.com/retroenv/retrogolib/assert"
)

func TestCarryControlPreservesFlagsAndClobbersLatch(t *testing.T) {
	t.Parallel()

	for _, code := range []byte{0x37, 0x3f} {
		for _, latch := range []uint16{0, 0xff} {
			input := fixture(0, code, 0)
			input.Before[AF], _ = knownbits.NewConstant[uint16](16, 0xed)
			input.Before[Q], _ = knownbits.NewConstant[uint16](8, latch)
			result := bufferedResult()
			Transfer(input, &result)
			after := applyWrites(t, input.Before, result)
			expected := uint16(0xc4)
			if latch == 0 {
				expected |= 0x28
			}
			if code == 0x37 {
				expected |= 1
			} else {
				expected |= 0x10
			}
			assert.Equal(t, expected, after[AF].KnownOne)
			assert.Equal(t, knownbits.Value[uint16]{}, after[Q])
		}
	}
}

func TestCompareIdentityAndBorrowPolarity(t *testing.T) {
	t.Parallel()

	input := fixture(0, 0xbf, 0)
	input.Before[AF], _ = knownbits.New[uint16](16, 0xff, 0)
	result := bufferedResult()
	Transfer(input, &result)
	after := applyWrites(t, input.Before, result)
	assert.Equal(t, uint16(0x42), after[AF].KnownOne)
	assert.Equal(t, uint16(0x95), after[AF].KnownZero)
	input.Opcode.Opcode = 0xb8
	input.Before[BC] = knownbits.Value[uint16]{}
	Transfer(input, &result)
	after = applyWrites(t, input.Before, result)
	zero, valid := knownbits.TestBit(16, after[AF], 6)
	assert.True(t, valid)
	assert.Equal(t, knownbits.BitUnknown, zero)
	input = fixture(0, 0xfe, 0)
	input.Before[AF], _ = knownbits.NewConstant[uint16](16, 0x100)
	input.Operand, _ = knownbits.NewConstant[uint16](8, 0x28)
	Transfer(input, &result)
	after = applyWrites(t, input.Before, result)
	assert.Equal(t, uint16(1), after[AF].KnownOne&1)
	assert.Equal(t, uint16(0x28), after[AF].KnownOne&0x28)
	input = fixture(0, 0xa8, 0)
	input.Before[AF] = knownbits.Value[uint16]{}
	input.Before[BC] = knownbits.Value[uint16]{}
	Transfer(input, &result)
	after = applyWrites(t, input.Before, result)
	assert.Equal(t, uint16(0), (after[AF].KnownZero|after[AF].KnownOne)&0xff00)
}

func TestRefreshReadAndLoadOverride(t *testing.T) {
	t.Parallel()

	input := fixture(z80.PrefixED, 0x5f, 0)
	input.Before[R], _ = knownbits.NewConstant[uint16](8, 0xff)
	input.Before[IFF2], _ = knownbits.NewConstant[uint16](1, 1)
	result := bufferedResult()
	Transfer(input, &result)
	after := applyWrites(t, input.Before, result)
	assert.Equal(t, uint16(0x81), after[R].KnownOne)
	assert.Equal(t, uint16(0x81), after[AF].KnownOne>>8)
	assert.Equal(t, uint16(4), after[AF].KnownOne&4)
	assert.Empty(t, result.Copies)
	input = fixture(z80.PrefixED, 0x4f, 0x23)
	input.Before[R], _ = knownbits.NewConstant[uint16](8, 0xff)
	Transfer(input, &result)
	after = applyWrites(t, input.Before, result)
	assert.Equal(t, uint16(0x23), after[R].KnownOne)
	assert.True(t, result.Valid(StorageWidth))
}
