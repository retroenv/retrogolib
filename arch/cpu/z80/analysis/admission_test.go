package analysis

import (
	"testing"

	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/arch/cpu/z80"
	"github.com/retroenv/retrogolib/assert"
)

func TestUnsupportedAndMalformedForms(t *testing.T) {
	t.Parallel()

	input := fixture(0, 0x41, 1)
	result := bufferedResult()
	Transfer(input, &result)
	assert.NotEmpty(t, result.Copies)
	input = fixture(0, 0xeb, 1)
	Transfer(input, &result)
	assert.False(t, result.EffectsComplete)
	assert.Equal(t, analysis.UnsupportedOpcode, result.Reason)
	assert.Equal(t, analysis.FlowUnknown, result.Flow)
	assert.Empty(t, result.Writes)
	assert.Empty(t, result.Copies)
	input = fixture(0, 0x3e, 1)
	input.Opcode.Size = 3
	Transfer(input, &result)
	assert.Equal(t, analysis.IllegalOperand, result.Reason)
	input = fixture(0, 0x3e, 1)
	input.Operand.KnownOne |= 0x100
	Transfer(input, &result)
	assert.Equal(t, analysis.InvalidInput, result.Reason)
	input = fixture(0, 0x3e, 1)
	input.Before[AF].KnownZero |= input.Before[AF].KnownOne
	Transfer(input, &result)
	assert.Equal(t, analysis.InvalidInput, result.Reason)
	input = fixture(z80.PrefixCB, 0, 1)
	Transfer(input, &result)
	assert.Equal(t, analysis.UnsupportedOpcode, result.Reason)
	assert.False(t, result.EffectsComplete)
}

func TestMemoryStabilityAndInvalidRead(t *testing.T) {
	t.Parallel()

	input := fixture(0, 0x7e, 1)
	input.Read = func(Address) (knownbits.Value[uint16], bool) {
		value, _ := knownbits.NewConstant[uint16](8, 0xaa)
		return value, false
	}
	result := bufferedResult()
	Transfer(input, &result)
	after := applyWrites(t, input.Before, result)
	assert.Equal(t, uint16(0), (after[AF].KnownZero|after[AF].KnownOne)&0xff00)
	assert.Len(t, result.Memory, 1)
	assert.False(t, result.Memory[0].Stable)
	input.Before[HL] = knownbits.Value[uint16]{}
	calls := 0
	input.Read = func(Address) (knownbits.Value[uint16], bool) {
		calls++
		return knownbits.Value[uint16]{}, true
	}
	Transfer(input, &result)
	assert.Equal(t, 0, calls)
	assert.Len(t, result.Memory, 1)
	input.Before[HL], _ = knownbits.NewConstant[uint16](16, 0)
	input.Read = func(Address) (knownbits.Value[uint16], bool) {
		return knownbits.Value[uint16]{
			KnownOne: 0x100,
		}, true
	}
	Transfer(input, &result)
	assert.Equal(t, analysis.InvalidInput, result.Reason)
	assert.False(t, result.EffectsComplete)
	assert.Empty(t, result.Writes)
}
