package analysis

import (
	"testing"

	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/assert"
)

func TestUnknownPCDoesNotInventReturnBytes(t *testing.T) {
	t.Parallel()

	input := fixture(0x20, 3, 0)
	input.PC = knownbits.Value[uint16]{}
	result := bufferedResult()
	Transfer(input, &result)
	assert.True(t, result.EffectsComplete)
	assert.Equal(t, analysis.FlowCall, result.Flow)
	assert.Len(t, result.Memory, 2)
	assert.Equal(t, knownbits.Value[uint8]{}, result.Memory[0].Value)
	assert.Equal(t, knownbits.Value[uint8]{}, result.Memory[1].Value)
	input.PC = knownbits.Value[uint16]{
		KnownZero: 1,
		KnownOne:  1,
	}
	Transfer(input, &result)
	assert.Equal(t, analysis.InvalidInput, result.Reason)
	assert.False(t, result.EffectsComplete)
}
