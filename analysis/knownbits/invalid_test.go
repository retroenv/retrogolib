package knownbits

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestOperationsRejectMalformedFacts(t *testing.T) {
	t.Parallel()

	for _, width := range []Width{0, 4, 8, 9} {
		for _, value := range []Value[uint8]{
			{
				KnownZero: 1,
				KnownOne:  1,
			},
			{
				KnownOne: 0x80,
			},
		} {
			if value.Valid(width) {
				continue
			}
			checkMalformedFacts(t, width, value)
		}
	}
}

func TestOperationsRejectInvalidParameters(t *testing.T) {
	t.Parallel()

	value := Value[uint8]{}
	_, valid := Extract(8, value, 255, 2)
	assert.False(t, valid)
	_, valid = Extract(8, value, 0, 0)
	assert.False(t, valid)
	_, valid = Insert(8, value, 7, 2, value)
	assert.False(t, valid)
	_, valid = Concat(8, value, 8, value)
	assert.False(t, valid)
	_, valid = Truncate(4, 8, value)
	assert.False(t, valid)
	_, valid = ZeroExtend(8, 4, value)
	assert.False(t, valid)
	_, valid = SignExtend(8, 4, value)
	assert.False(t, valid)
	_, valid = TestBit(8, value, 8)
	assert.False(t, valid)
	_, valid = CarryAt(8, value, value, BitClear, 0)
	assert.False(t, valid)
	_, valid = BorrowAt(8, value, value, BitClear, 9)
	assert.False(t, valid)
	for _, bit := range []Bit{3, 255} {
		_, valid = Add(8, value, value, bit)
		assert.False(t, valid)
		_, valid = Sub(8, value, value, bit)
		assert.False(t, valid)
		_, _, valid = RotateLeftCarry(8, value, bit)
		assert.False(t, valid)
		_, _, valid = RotateRightCarry(8, value, bit)
		assert.False(t, valid)
		_, valid = CarryAt(8, value, value, bit, 4)
		assert.False(t, valid)
		_, valid = BorrowAt(8, value, value, bit, 4)
		assert.False(t, valid)
	}
}

func checkMalformedFacts(t *testing.T, width Width, value Value[uint8]) {
	t.Helper()

	for _, function := range []func(Width, Value[uint8], Value[uint8]) (Value[uint8], bool){
		Join[uint8], And[uint8], Or[uint8], Xor[uint8],
	} {
		_, valid := function(width, value, Value[uint8]{})
		assert.False(t, valid)
		_, valid = function(width, Value[uint8]{}, value)
		assert.False(t, valid)
	}
	_, _, valid := Meet(width, value, Value[uint8]{})
	assert.False(t, valid)
	_, _, valid = Meet(width, Value[uint8]{}, value)
	assert.False(t, valid)
	_, valid = Contains(width, value, 0)
	assert.False(t, valid)
	_, valid = Not(width, value)
	assert.False(t, valid)
	for _, function := range []func(Width, Value[uint8], uint) (Value[uint8], bool){
		ShiftLeft[uint8], ShiftRight[uint8], ShiftRightSigned[uint8], RotateLeft[uint8], RotateRight[uint8],
	} {
		_, valid = function(width, value, 1)
		assert.False(t, valid)
	}
	_, valid = Add(width, value, value, BitUnknown)
	assert.False(t, valid)
	_, valid = Sub(width, value, value, BitUnknown)
	assert.False(t, valid)
	_, valid = Equal(width, value, value)
	assert.False(t, valid)
	for _, predicate := range []func(Width, Value[uint8]) (Bit, bool){Zero[uint8], Sign[uint8], EvenParity[uint8]} {
		_, valid = predicate(width, value)
		assert.False(t, valid)
	}
	_, _, valid = RotateLeftCarry(width, value, BitClear)
	assert.False(t, valid)
	_, _, valid = RotateRightCarry(width, value, BitClear)
	assert.False(t, valid)
	_, valid = Truncate(width, 1, value)
	assert.False(t, valid)
	_, valid = Extract(width, value, 0, 1)
	assert.False(t, valid)
	_, valid = Insert(width, value, 0, 1, Value[uint8]{})
	assert.False(t, valid)
	_, valid = ZeroExtend(width, 8, value)
	assert.False(t, valid)
	_, valid = SignExtend(width, 8, value)
	assert.False(t, valid)
}
