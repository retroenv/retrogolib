package knownbits

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestSlicesAndExtensionsExhaustive(t *testing.T) {
	t.Parallel()

	for _, test := range byteCases(4) {
		for from := Width(1); from <= 4; from++ {
			part, valid := Truncate(4, from, test.fact)
			assert.True(t, valid)
			for to := from; to <= 8; to++ {
				checkExtension(t, from, to, part, false)
				checkExtension(t, from, to, part, true)
			}
		}
		for offset := range Width(4) {
			for size := Width(1); size <= 4-offset; size++ {
				extracted, valid := Extract(4, test.fact, offset, size)
				assert.True(t, valid)
				restored, valid := Insert(4, test.fact, offset, size, extracted)
				assert.True(t, valid)
				assert.Equal(t, test.fact, restored)
			}
		}
	}
}

func TestInsertPreservesUnaffectedBits(t *testing.T) {
	t.Parallel()

	parent, _ := NewConstant[uint32](32, 0x12345678)
	part, _ := NewConstant[uint32](8, 0xaa)
	result, valid := Insert(32, parent, 0, 8, part)
	assert.True(t, valid)
	number, known := result.Constant(32)
	assert.True(t, known)
	assert.Equal(t, uint32(0x123456aa), number)
	unknownPart := Value[uint32]{}
	result, valid = Insert(32, parent, 8, 8, unknownPart)
	assert.True(t, valid)
	assert.Equal(t, parent.KnownZero&^0xff00, result.KnownZero)
	assert.Equal(t, parent.KnownOne&^0xff00, result.KnownOne)
}

func TestConcatenationAndFullWidthBoundaries(t *testing.T) {
	t.Parallel()

	high, _ := NewConstant[uint64](32, 0x12345678)
	low, _ := NewConstant[uint64](32, 0x9abcdef0)
	joined, valid := Concat(32, high, 32, low)
	assert.True(t, valid)
	number, known := joined.Constant(64)
	assert.True(t, known)
	assert.Equal(t, uint64(0x123456789abcdef0), number)
	for _, width := range []Width{1, 16, 24, 32, 64} {
		constant, _ := NewConstant[uint64](width, 1)
		for _, count := range []uint{0, 1, uint(width), uint(width) + 1, ^uint(0)} {
			shifted, valid := ShiftLeft(width, constant, count)
			assert.True(t, valid)
			assert.True(t, shifted.Valid(width))
			rotated, valid := RotateRight(width, constant, count)
			assert.True(t, valid)
			restored, valid := RotateLeft(width, rotated, count)
			assert.True(t, valid)
			assert.Equal(t, constant, restored)
		}
	}
}

func checkExtension(t *testing.T, from, to Width, value Value[uint8], signed bool) {
	t.Helper()

	mask := uint8((1 << to) - 1)
	zero := mask
	one := mask
	for number := range 1 << from {
		if uint8(number)&value.KnownZero != 0 || uint8(number)&value.KnownOne != value.KnownOne {
			continue
		}
		concrete := uint8(number)
		if signed {
			concrete = uint8(signedByte(from, concrete)) & mask
		}
		zero &= ^concrete
		one &= concrete
	}
	fact, valid := ZeroExtend(from, to, value)
	if signed {
		fact, valid = SignExtend(from, to, value)
	}
	assert.True(t, valid)
	assert.Equal(t, zero, fact.KnownZero)
	assert.Equal(t, one, fact.KnownOne)
}
