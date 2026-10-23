package analysis

import (
	"testing"

	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/arch/cpu/analysis"
	"github.com/retroenv/retrogolib/assert"
)

func TestMemoryEffectsDoNotDependOnEstablishedValues(t *testing.T) {
	t.Parallel()

	for _, code := range []byte{0x65, 0xe5, 0x05, 0x45, 0x25, 0xa5, 0xa6, 0xa4, 0xa7} {
		checkMemoryEffects(t, code, 1, 0)
	}
	for _, code := range []byte{0x85, 0x86, 0x84, 0x87} {
		checkMemoryEffects(t, code, 0, 1)
	}
	for _, code := range []byte{0x06, 0x46, 0x26, 0x66, 0xe6, 0xc6, 0x07, 0x27, 0x47, 0x67, 0xc7, 0xe7} {
		checkMemoryEffects(t, code, 1, 1)
	}
	input := fixture(0x25, 2, 0)
	result := bufferedResult()
	Transfer(input, &result)
	after := applyWrites(t, input.Before, result)
	value, known := after[A].Constant(8)
	assert.True(t, known)
	assert.Equal(t, byte(0), value)
	assert.Len(t, result.Memory, 1)
	assert.False(t, result.Memory[0].Stable)
}

func checkMemoryEffects(t *testing.T, code byte, expectedReads, expectedWrites int) {
	t.Helper()

	input := fixture(code, 2, 0)
	for _, unknown := range []bool{false, true} {
		if unknown {
			input.Before = [ResourceCount]knownbits.Value[uint8]{}
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
