package analysis

import (
	"testing"

	"github.com/retroenv/retrogolib/analysis/knownbits"
	"github.com/retroenv/retrogolib/assert"
)

func TestViewsAndWriteRules(t *testing.T) {
	t.Parallel()

	parent, _ := knownbits.NewConstant[uint32](32, 0x12345678)
	part, _ := knownbits.NewConstant[uint32](8, 0xaa)
	write := Write[int, uint32]{
		View: View[int]{
			Storage: 0,
			Bits:    8,
		},
		Value: part,
		Rule:  Preserve,
	}
	result, valid := Apply(32, parent, write)
	assert.True(t, valid)
	assert.Equal(t, uint32(0x123456aa), result.KnownOne)
	write.Rule = ZeroExtend
	result, valid = Apply(32, parent, write)
	assert.True(t, valid)
	assert.Equal(t, uint32(0xaa), result.KnownOne)
	assert.Equal(t, uint32(0xffffff55), result.KnownZero)
	write.Rule = SignExtend
	result, valid = Apply(32, parent, write)
	assert.True(t, valid)
	assert.Equal(t, uint32(0xffffffaa), result.KnownOne)
	write.Rule = Invalidate
	result, valid = Apply(32, parent, write)
	assert.True(t, valid)
	assert.Equal(t, uint32(0x12345600), result.KnownOne)
	assert.Equal(t, uint32(0xedcba900), result.KnownZero)
	write.Rule = Replace
	_, valid = Apply(32, parent, write)
	assert.False(t, valid)
	write.View.Bits = 32
	write.Value = parent
	result, valid = Apply(32, parent, write)
	assert.True(t, valid)
	assert.Equal(t, parent, result)
}

func TestInvalidViewsAndWrites(t *testing.T) {
	t.Parallel()

	for _, part := range []View[int]{
		{},
		{
			Offset: 255,
			Bits:   2,
		},
		{
			Offset: 31,
			Bits:   2,
		},
	} {
		assert.False(t, part.Valid(32))
	}
	left := View[int]{
		Bits: 8,
	}
	right := View[int]{
		Offset: 8,
		Bits:   8,
	}
	whole := View[int]{
		Bits: 16,
	}
	assert.False(t, left.Overlaps(right))
	assert.True(t, left.Overlaps(whole))
	right.Storage = 1
	assert.False(t, right.Overlaps(whole))
	for _, rule := range []WriteRule{ZeroExtend, SignExtend, 255} {
		_, valid := Apply(16, knownbits.Value[uint16]{}, Write[int, uint16]{
			View: View[int]{
				Offset: 8,
				Bits:   8,
			},
			Rule: rule,
		})
		assert.False(t, valid)
	}
}

func TestResultResetAndValidation(t *testing.T) {
	t.Parallel()

	result := Result[int, uint8, uint16]{
		Reason:          Supported,
		EffectsComplete: true,
		Flow:            FlowNext,
		Writes: []Write[int, uint8]{
			{
				View: View[int]{
					Bits: 8,
				},
				Rule: Replace,
			},
		},
		Copies: []Copy[int]{
			{
				Source: View[int]{
					Storage: 1,
					Bits:    8,
				},
				Destination: View[int]{
					Bits: 8,
				},
			},
		},
	}
	width := func(int) knownbits.Width { return 8 }
	assert.True(t, result.Valid(width))
	result.Writes = append(result.Writes, result.Writes[0])
	assert.False(t, result.Valid(width))
	result.Writes = result.Writes[:1]
	result.Copies[0].Destination.Storage = 2
	assert.False(t, result.Valid(width))
	capacity := cap(result.Writes)
	result.Reset()
	assert.Equal(t, capacity, cap(result.Writes))
	assert.Empty(t, result.Writes)
	assert.Empty(t, result.Copies)
	assert.False(t, result.EffectsComplete)
	assert.Equal(t, UnknownEffects, result.Reason)
}
