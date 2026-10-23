package analysis

import (
	"testing"

	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/arch/cpu/cpu6502"
	"github.com/retroenv/retrogolib/assert"
)

func TestArithmeticModeAdmission(t *testing.T) {
	t.Parallel()

	for _, code := range []byte{0x69, 0xe9} {
		input := fixture(code, 2, 0x99)
		input.Variant = cpu6502.VariantNMOS6502
		input.Before[P] = knownbits.Value[uint8]{}
		result := bufferedResult()
		Transfer(input, &result)
		assert.True(t, result.EffectsComplete)
		assert.Equal(t, analysis.UnknownMode, result.Reason)
		after := applyWrites(t, input.Before, result)
		assert.Equal(t, knownbits.Value[uint8]{}, after[A])
		assert.Equal(t, byte(0), (after[P].KnownZero|after[P].KnownOne)&0xc3)
		input.Before[P], _ = knownbits.NewConstant[uint8](8, 0x28)
		Transfer(input, &result)
		assert.Equal(t, analysis.UnsupportedMode, result.Reason)
		input.Variant = cpu6502.VariantNES6502
		Transfer(input, &result)
		assert.Equal(t, analysis.Supported, result.Reason)
		assert.True(t, result.EffectsComplete)
		for _, read := range result.Reads {
			assert.False(t, read.Storage == P && read.Offset == 3 && read.Bits == 1)
		}
	}
}

func TestRejectedInputClearsPreviousResult(t *testing.T) {
	t.Parallel()

	result := bufferedResult()
	input := fixture(0xaa, 1, 1)
	Transfer(input, &result)
	assert.NotEmpty(t, result.Writes)
	assert.Len(t, result.Copies, 1)
	input.Opcode.Opcode = 0x8b
	input.Opcode.Size = 2
	Transfer(input, &result)
	assert.Equal(t, analysis.UnsupportedOpcode, result.Reason)
	assert.False(t, result.EffectsComplete)
	assert.Empty(t, result.Writes)
	assert.Empty(t, result.Copies)
	input.Variant = -1
	Transfer(input, &result)
	assert.Equal(t, analysis.UnsupportedVariant, result.Reason)
	input = fixture(0x69, 1, 1)
	Transfer(input, &result)
	assert.Equal(t, analysis.IllegalOperand, result.Reason)
	input.Opcode.Size = 2
	input.Before[A].KnownZero |= input.Before[A].KnownOne
	Transfer(input, &result)
	assert.Equal(t, analysis.InvalidInput, result.Reason)
	assert.False(t, result.EffectsComplete)
}

func TestMemoryProofAndUnknownAddress(t *testing.T) {
	t.Parallel()

	input := fixture(0xa5, 2, 1)
	input.Read = func(Address) (knownbits.Value[uint8], bool) {
		value, _ := knownbits.NewConstant[uint8](8, 0xaa)
		return value, false
	}
	result := bufferedResult()
	Transfer(input, &result)
	after := applyWrites(t, input.Before, result)
	assert.Equal(t, knownbits.Value[uint8]{}, after[A])
	assert.Len(t, result.Memory, 1)
	assert.False(t, result.Memory[0].Stable)
	input = fixture(0xbd, 3, 1)
	input.Before[X] = knownbits.Value[uint8]{}
	calls := 0
	input.Read = func(Address) (knownbits.Value[uint8], bool) {
		calls++
		return knownbits.Value[uint8]{}, true
	}
	Transfer(input, &result)
	assert.Equal(t, 0, calls)
	assert.Len(t, result.Memory, 1)
	_, known := result.Memory[0].Address.Value.Constant(16)
	assert.False(t, known)
	input = fixture(0xa5, 2, 1)
	input.Read = func(Address) (knownbits.Value[uint8], bool) {
		return knownbits.Value[uint8]{
			KnownZero: 1,
			KnownOne:  1,
		}, true
	}
	Transfer(input, &result)
	assert.Equal(t, analysis.InvalidInput, result.Reason)
	assert.False(t, result.EffectsComplete)
	assert.Empty(t, result.Writes)
}

func TestBranchFlagPreservation(t *testing.T) {
	t.Parallel()

	for _, code := range []byte{0x10, 0x30, 0x50, 0x70, 0x90, 0xb0, 0xd0, 0xf0} {
		input := fixture(code, 2, 0)
		input.Before[P] = knownbits.Value[uint8]{}
		result := bufferedResult()
		Transfer(input, &result)
		assert.True(t, result.EffectsComplete)
		assert.Equal(t, analysis.FlowBranch, result.Flow)
		assert.Equal(t, knownbits.BitUnknown, result.Condition)
		assert.Empty(t, result.Writes)
		assert.Len(t, result.Reads, 1)
		position := result.Reads[0].Offset
		input.Before[P], _ = knownbits.NewConstant[uint8](8, 1<<position)
		Transfer(input, &result)
		want := knownbits.BitClear
		if code&0x20 != 0 {
			want = knownbits.BitSet
		}
		assert.Equal(t, want, result.Condition)
	}
}
