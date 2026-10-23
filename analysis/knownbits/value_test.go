package knownbits

import (
	"testing"
	"unsafe"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/set"
)

func TestLayoutAndNamedStorage(t *testing.T) {
	t.Parallel()

	type address uint32
	value, valid := NewConstant[address](24, 0x123456)
	assert.True(t, valid)
	assert.Equal(t, address(0x123456), value.KnownOne)
	assert.Equal(t, uintptr(2), unsafe.Sizeof(Value[uint8]{}))
	assert.Equal(t, uintptr(1), unsafe.Alignof(Value[uint8]{}))
	assert.Equal(t, uintptr(16), unsafe.Sizeof(Value[uint64]{}))
	facts := set.New[Value[uint8]]()
	facts.Add(Value[uint8]{})
	assert.True(t, facts.Contains(Value[uint8]{}))
}

func TestConstructionAndQueries(t *testing.T) {
	t.Parallel()

	for width := Width(1); width <= 64; width++ {
		mask := ^uint64(0) >> (64 - width)
		value, valid := NewConstant(width, ^uint64(0))
		assert.True(t, valid)
		assert.True(t, value.Valid(width))
		number, known := value.Constant(width)
		assert.True(t, known)
		assert.Equal(t, mask, number)
		low, high, valid := (Value[uint64]{}).UnsignedRange(width)
		assert.True(t, valid)
		assert.Equal(t, uint64(0), low)
		assert.Equal(t, mask, high)
	}
	value, valid := New[uint16](8, 0xff00, 0xff80)
	assert.True(t, valid)
	assert.Equal(t, uint16(0x80), value.KnownOne)
	assert.Equal(t, uint16(0), value.KnownZero)
	low, high, valid := value.UnsignedRange(8)
	assert.True(t, valid)
	assert.Equal(t, uint16(0x80), low)
	assert.Equal(t, uint16(0xff), high)
	_, known := value.Constant(8)
	assert.False(t, known)
}

func TestInvalidConstructionAndFacts(t *testing.T) {
	t.Parallel()

	for _, width := range []Width{0, 9, 64, 255} {
		_, valid := NewConstant[uint8](width, 0)
		assert.False(t, valid)
		assert.False(t, (Value[uint8]{}).Valid(width))
	}
	_, valid := New[uint8](8, 1, 1)
	assert.False(t, valid)
	for _, value := range []Value[uint16]{
		{
			KnownZero: 1,
			KnownOne:  1,
		},
		{
			KnownOne: 0x100,
		},
	} {
		assert.False(t, value.Valid(8))
		_, known := value.Constant(8)
		assert.False(t, known)
		_, _, valid := value.UnsignedRange(8)
		assert.False(t, valid)
	}
}

func TestJoinLawsAndIndependentOperands(t *testing.T) {
	t.Parallel()

	cases := byteCases(3)
	for _, left := range cases {
		idempotent, valid := Join(3, left.fact, left.fact)
		assert.True(t, valid)
		assert.Equal(t, left.fact, idempotent)
		for _, right := range cases {
			forward, _ := Join(3, left.fact, right.fact)
			reverse, _ := Join(3, right.fact, left.fact)
			assert.Equal(t, forward, reverse)
			for _, third := range cases {
				first, _ := Join(3, forward, third.fact)
				pair, _ := Join(3, right.fact, third.fact)
				second, _ := Join(3, left.fact, pair)
				assert.Equal(t, first, second)
			}
		}
	}
	unknown := Value[uint8]{}
	equal, valid := Equal(8, unknown, unknown)
	assert.True(t, valid)
	assert.Equal(t, BitUnknown, equal)
	xor, valid := Xor(8, unknown, unknown)
	assert.True(t, valid)
	assert.Equal(t, unknown, xor)
}
