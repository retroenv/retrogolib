package analysis

import (
	"testing"

	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/assert"
)

func TestMemoryEffectsDoNotDependOnEstablishedValues(t *testing.T) {
	t.Parallel()

	for _, code := range []byte{0x46, 0x7e, 0x0a, 0x1a, 0x3a, 0xa6, 0xae, 0xb6, 0xbe} {
		checkMemoryEffects(t, code, 1, 0)
	}
	for _, code := range []byte{0x70, 0x36, 0x02, 0x12, 0x32} {
		checkMemoryEffects(t, code, 0, 1)
	}
	checkMemoryEffects(t, 0x22, 0, 2)
	checkMemoryEffects(t, 0x2a, 2, 0)
	input := fixture(0, 0xa6, 0)
	result := bufferedResult()
	Transfer(input, &result)
	after := applyWrites(t, input.Before, result)
	value, known := knownbits.Extract(16, after[AF], 8, 8)
	assert.True(t, known)
	constant, known := value.Constant(8)
	assert.True(t, known)
	assert.Equal(t, uint16(0), constant)
	assert.Len(t, result.Memory, 1)
	assert.False(t, result.Memory[0].Stable)
}

func checkMemoryEffects(t *testing.T, code byte, expectedReads, expectedWrites int) {
	t.Helper()

	input := fixture(0, code, 0)
	for _, unknown := range []bool{false, true} {
		if unknown {
			input.Before = [ResourceCount]knownbits.Value[uint16]{}
		}
		result := bufferedResult()
		Transfer(input, &result)
		assert.True(t, result.EffectsComplete)
		reads := 0
		writes := 0
		for _, event := range result.Memory {
			if event.Kind == analysis.MemoryRead {
				reads++
			} else {
				writes++
			}
		}
		assert.Equal(t, expectedReads, reads)
		assert.Equal(t, expectedWrites, writes)
	}
}
