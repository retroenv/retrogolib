package analysis

import (
	"testing"

	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/assert"
)

func TestExtensionInvalidatesWholeStorage(t *testing.T) {
	t.Parallel()

	write := Write[int, uint16]{
		View: View[int]{
			Bits: 8,
		},
		Rule: ZeroExtend,
	}
	affected, valid := write.Affected(16)
	assert.True(t, valid)
	assert.Equal(t, knownbits.Width(16), affected.Bits)
	result := Result[int, uint16, uint16]{
		Reason:          Supported,
		EffectsComplete: true,
		Writes: []Write[int, uint16]{
			write,
			{
				View: View[int]{
					Offset: 8,
					Bits:   8,
				},
				Rule: Preserve,
			},
		},
	}
	assert.False(t, result.Valid(func(int) knownbits.Width { return 16 }))
}

func TestInvalidResultRecords(t *testing.T) {
	t.Parallel()

	for _, mutate := range []func(*Result[int, uint8, uint16]){
		func(rr *Result[int, uint8, uint16]) { rr.Reason = 255 },
		func(rr *Result[int, uint8, uint16]) { rr.Flow = 255 },
		func(rr *Result[int, uint8, uint16]) { rr.Condition = 255 },
		func(rr *Result[int, uint8, uint16]) { rr.Reason = UnknownEffects },
		invalidRead,
		invalidWrite,
		invalidCopy,
		invalidMemory,
	} {
		result := Result[int, uint8, uint16]{
			Reason:          Supported,
			EffectsComplete: true,
		}
		mutate(&result)
		assert.False(t, result.Valid(func(int) knownbits.Width { return 8 }))
	}
}

func TestMemoryEventValidation(t *testing.T) {
	t.Parallel()

	for _, length := range []uint32{0, 1, 2, 8, 9, ^uint32(0)} {
		event := MemoryEvent[uint8, uint16]{
			Bytes: length,
		}
		assert.True(t, validEvent(event))
		event.Value.KnownOne = 1
		assert.Equal(t, length == 1, validEvent(event))
	}
	write := MemoryEvent[uint8, uint16]{
		Kind:   MemoryWrite,
		Bytes:  1,
		Stable: true,
	}
	assert.False(t, validEvent(write))
	write.Stable = false
	write.Value.KnownOne = 1
	write.Value.KnownZero = 1
	assert.False(t, validEvent(write))
}

func invalidRead(result *Result[int, uint8, uint16]) {
	result.Reads = []View[int]{
		{
			Bits: 9,
		},
	}
}

func invalidWrite(result *Result[int, uint8, uint16]) {
	result.Writes = []Write[int, uint8]{
		{
			View: View[int]{
				Bits: 8,
			},
			Value: knownbits.Value[uint8]{
				KnownZero: 1,
				KnownOne:  1,
			},
		},
	}
}

func invalidCopy(result *Result[int, uint8, uint16]) {
	result.Copies = []Copy[int]{
		{
			Source: View[int]{
				Bits: 8,
			},
			Destination: View[int]{
				Bits: 7,
			},
		},
	}
}

func invalidMemory(result *Result[int, uint8, uint16]) {
	result.Memory = []MemoryEvent[uint8, uint16]{
		{
			Kind: 255,
		},
	}
}
