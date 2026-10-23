package knownbits

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestBinaryBitwiseExhaustive(t *testing.T) {
	t.Parallel()

	cases := byteCases(4)
	for _, left := range cases {
		for _, right := range cases {
			for operation := range 3 {
				checkBinary(t, left, right, operation)
			}
		}
	}
}

func TestByteUnaryExhaustive(t *testing.T) {
	t.Parallel()

	for _, test := range byteCases(8) {
		complement, valid := Not(8, test.fact)
		assert.True(t, valid)
		assert.Equal(t, test.fact.KnownOne, complement.KnownZero)
		assert.Equal(t, test.fact.KnownZero, complement.KnownOne)
		for _, count := range []uint{0, 1, 3, 7, 8, 9, 64, ^uint(0)} {
			checkShifts(t, test, count)
		}
		checkPredicates(t, test)
		for carry := BitUnknown; carry <= BitSet; carry++ {
			checkCarryRotates(t, test, carry)
		}
	}
}

func checkBinary(t *testing.T, left, right byteCase, operation int) {
	t.Helper()

	zero := uint8(15)
	one := uint8(15)
	for _, aa := range left.numbers {
		for _, bb := range right.numbers {
			value := binaryConcrete(aa, bb, operation)
			zero &= ^value
			one &= value
		}
	}
	functions := []func(Width, Value[uint8], Value[uint8]) (Value[uint8], bool){And[uint8], Or[uint8], Xor[uint8]}
	fact, valid := functions[operation](4, left.fact, right.fact)
	assert.True(t, valid)
	assert.Equal(t, zero, fact.KnownZero)
	assert.Equal(t, one, fact.KnownOne)
	predicate, valid := Equal(4, left.fact, right.fact)
	assert.True(t, valid)
	events := uint8(0)
	for _, aa := range left.numbers {
		for _, bb := range right.numbers {
			if aa == bb {
				events |= 2
			} else {
				events |= 1
			}
		}
	}
	assert.Equal(t, oracleBit(events), predicate)
}

func checkShifts(t *testing.T, test byteCase, count uint) {
	t.Helper()

	functions := []func(Width, Value[uint8], uint) (Value[uint8], bool){
		ShiftLeft[uint8], ShiftRight[uint8], ShiftRightSigned[uint8], RotateLeft[uint8], RotateRight[uint8],
	}
	for operation, function := range functions {
		zero := uint8(255)
		one := uint8(255)
		for _, number := range test.numbers {
			value := shiftConcrete(number, count, operation)
			zero &= ^value
			one &= value
		}
		fact, valid := function(8, test.fact, count)
		assert.True(t, valid)
		assert.Equal(t, zero, fact.KnownZero)
		assert.Equal(t, one, fact.KnownOne)
	}
}

func checkCarryRotates(t *testing.T, test byteCase, carry Bit) {
	t.Helper()

	for _, left := range []bool{false, true} {
		zero := uint8(255)
		one := uint8(255)
		events := uint8(0)
		for _, number := range test.numbers {
			for input := range 2 {
				if carry != BitUnknown && input != int(carry-BitClear) {
					continue
				}
				value := number>>1 | uint8(input)<<7
				out := number & 1
				if left {
					value = number<<1 | uint8(input)
					out = number >> 7
				}
				zero &= ^value
				one &= value
				events |= 1 << out
			}
		}
		fact, event, valid := RotateRightCarry(8, test.fact, carry)
		if left {
			fact, event, valid = RotateLeftCarry(8, test.fact, carry)
		}
		assert.True(t, valid)
		assert.Equal(t, zero, fact.KnownZero)
		assert.Equal(t, one, fact.KnownOne)
		assert.Equal(t, oracleBit(events), event)
	}
}

func checkPredicates(t *testing.T, test byteCase) {
	t.Helper()

	var zeroEvents, signEvents, parityEvents uint8
	for _, number := range test.numbers {
		zeroEvents |= predicateEvent(number == 0)
		signEvents |= predicateEvent(number&0x80 != 0)
		parity := uint8(0)
		for position := range 8 {
			parity ^= (number >> position) & 1
		}
		parityEvents |= predicateEvent(parity == 0)
	}
	zero, valid := Zero(8, test.fact)
	assert.True(t, valid)
	assert.Equal(t, oracleBit(zeroEvents), zero)
	sign, valid := Sign(8, test.fact)
	assert.True(t, valid)
	assert.Equal(t, oracleBit(signEvents), sign)
	parity, valid := EvenParity(8, test.fact)
	assert.True(t, valid)
	assert.Equal(t, oracleBit(parityEvents), parity)
	for position := range Width(8) {
		bit, valid := TestBit(8, test.fact, position)
		assert.True(t, valid)
		events := uint8(0)
		for _, number := range test.numbers {
			events |= 1 << ((number >> position) & 1)
		}
		assert.Equal(t, oracleBit(events), bit)
	}
}

func predicateEvent(value bool) uint8 {
	if value {
		return 2
	}
	return 1
}

func binaryConcrete(left, right uint8, operation int) uint8 {
	switch operation {
	case 0:
		return left & right
	case 1:
		return left | right
	default:
		return left ^ right
	}
}

func shiftConcrete(number uint8, count uint, operation int) uint8 {
	switch operation {
	case 0:
		return number << min(count, 8)
	case 1:
		return number >> min(count, 8)
	case 2:
		return uint8(int8(number) >> min(count, 8))
	case 3:
		return number<<(count%8) | number>>(8-count%8)
	default:
		return number>>(count%8) | number<<(8-count%8)
	}
}
