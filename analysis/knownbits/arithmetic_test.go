package knownbits

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

type byteCase struct {
	fact    Value[uint8]
	numbers []uint8
}

type arithmeticOracle struct {
	zero     uint8
	one      uint8
	carry    uint8
	overflow uint8
}

func TestPartialArithmeticExhaustive(t *testing.T) {
	t.Parallel()

	cases := byteCases(4)
	for _, left := range cases {
		for _, right := range cases {
			for carry := BitUnknown; carry <= BitSet; carry++ {
				checkPartialArithmetic(t, left, right, carry, false)
				checkPartialArithmetic(t, left, right, carry, true)
			}
		}
	}
}

func TestByteConstantArithmeticExhaustive(t *testing.T) {
	t.Parallel()

	for left := range 256 {
		for right := range 256 {
			for carry := BitClear; carry <= BitSet; carry++ {
				checkByteConstants(t, uint8(left), uint8(right), carry, false)
				checkByteConstants(t, uint8(left), uint8(right), carry, true)
			}
		}
	}
}

func TestArithmeticWidthsAndBoundaries(t *testing.T) {
	t.Parallel()

	for _, width := range []Width{1, 8, 16, 24, 32, 64} {
		mask := ^uint64(0) >> (64 - width)
		maximum, _ := NewConstant(width, mask)
		one, _ := NewConstant[uint64](width, 1)
		zero, _ := NewConstant[uint64](width, 0)
		result, valid := Add(width, maximum, one, BitClear)
		assert.True(t, valid)
		assert.Equal(t, zero, result.Value)
		assert.Equal(t, BitSet, result.CarryOut)
		difference, valid := Sub(width, zero, one, BitClear)
		assert.True(t, valid)
		assert.Equal(t, maximum, difference.Value)
		assert.Equal(t, BitSet, difference.BorrowOut)
		for boundary := Width(1); boundary <= width; boundary++ {
			carry, valid := CarryAt(width, maximum, one, BitClear, boundary)
			assert.True(t, valid)
			assert.Equal(t, BitSet, carry)
			borrow, valid := BorrowAt(width, zero, one, BitClear, boundary)
			assert.True(t, valid)
			assert.Equal(t, BitSet, borrow)
		}
	}
}

func TestPartialArithmeticCarryProofs(t *testing.T) {
	t.Parallel()

	low, _ := New[uint8](8, 0x80, 0)
	one, _ := NewConstant[uint8](8, 1)
	result, valid := Add(8, low, one, BitUnknown)
	assert.True(t, valid)
	assert.Equal(t, BitClear, result.CarryOut)
	assert.Equal(t, BitUnknown, result.SignedOverflow)
	zero, _ := NewConstant[uint8](8, 0)
	result, valid = Add(8, low, zero, BitClear)
	assert.True(t, valid)
	assert.Equal(t, low, result.Value)
	assert.Equal(t, BitClear, result.SignedOverflow)
	unknown := Value[uint8]{}
	result, valid = Add(8, unknown, unknown, BitClear)
	assert.True(t, valid)
	assert.Equal(t, unknown, result.Value)
}

func checkPartialArithmetic(t *testing.T, left, right byteCase, carry Bit, subtract bool) {
	t.Helper()

	oracle := arithmeticOracle{
		zero: 15,
		one:  15,
	}
	for _, aa := range left.numbers {
		for _, bb := range right.numbers {
			for cc := range 2 {
				if carry == BitUnknown || cc == int(carry-BitClear) {
					oracle.include(4, aa, bb, cc, subtract)
				}
			}
		}
	}
	checkArithmetic(t, 4, left.fact, right.fact, carry, subtract, oracle)
}

func checkByteConstants(t *testing.T, left, right uint8, carry Bit, subtract bool) {
	t.Helper()

	leftFact, _ := NewConstant[uint8](8, left)
	rightFact, _ := NewConstant[uint8](8, right)
	oracle := arithmeticOracle{
		zero: 255,
		one:  255,
	}
	oracle.include(8, left, right, int(carry-BitClear), subtract)
	checkArithmetic(t, 8, leftFact, rightFact, carry, subtract, oracle)
}

func checkArithmetic(t *testing.T, width Width, left, right Value[uint8], carry Bit, subtract bool,
	oracle arithmeticOracle) {

	t.Helper()

	var fact Value[uint8]
	var event, overflow Bit
	var valid bool
	if subtract {
		result, ok := Sub(width, left, right, carry)
		fact, event, overflow, valid = result.Value, result.BorrowOut, result.SignedOverflow, ok
	} else {
		result, ok := Add(width, left, right, carry)
		fact, event, overflow, valid = result.Value, result.CarryOut, result.SignedOverflow, ok
	}
	assert.True(t, valid)
	assert.Equal(t, oracle.zero, fact.KnownZero)
	assert.Equal(t, oracle.one, fact.KnownOne)
	assert.Equal(t, oracleBit(oracle.carry), event)
	assert.Equal(t, oracleBit(oracle.overflow), overflow)
}

func (oo *arithmeticOracle) include(width Width, left, right uint8, carry int, subtract bool) {
	mask := (1 << width) - 1
	raw := int(left) + int(right) + carry
	signed := signedByte(width, left) + signedByte(width, right) + carry
	if subtract {
		raw = int(left) - int(right) - carry
		signed = signedByte(width, left) - signedByte(width, right) - carry
	}
	value := uint8(raw & mask)
	oo.zero &= ^value
	oo.one &= value
	if raw < 0 || raw > mask {
		oo.carry |= 2
	} else {
		oo.carry |= 1
	}
	limit := 1 << (width - 1)
	if signed < -limit || signed >= limit {
		oo.overflow |= 2
	} else {
		oo.overflow |= 1
	}
}

func signedByte(width Width, number uint8) int {
	if number&(1<<(width-1)) != 0 {
		return int(number) - (1 << width)
	}
	return int(number)
}

func oracleBit(events uint8) Bit {
	switch events {
	case 1:
		return BitClear
	case 2:
		return BitSet
	default:
		return BitUnknown
	}
}

func byteCases(width Width) []byteCase {
	count := 1
	for range width {
		count *= 3
	}
	result := make([]byteCase, 0, count)
	for code := range count {
		fact := Value[uint8]{}
		trits := code
		for position := range width {
			switch trits % 3 {
			case 1:
				fact.KnownZero |= 1 << position
			case 2:
				fact.KnownOne |= 1 << position
			}
			trits /= 3
		}
		numbers := make([]uint8, 0, 1<<(width))
		for number := range 1 << width {
			if uint8(number)&fact.KnownZero == 0 && uint8(number)&fact.KnownOne == fact.KnownOne {
				numbers = append(numbers, uint8(number))
			}
		}
		result = append(result, byteCase{
			fact:    fact,
			numbers: numbers,
		})
	}
	return result
}
