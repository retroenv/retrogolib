package analysis

import "github.com/retroenv/retrogolib/analysis/knownbits"

func (tt *transfer) flag(position knownbits.Width) knownbits.Bit {
	value := tt.read(view(AF, position, 1))
	bit, _ := knownbits.TestBit(1, value, 0)
	return bit
}

func (tt *transfer) writeFlag(position knownbits.Width, value knownbits.Bit) {
	part := knownbits.Value[uint16]{}
	switch value {
	case knownbits.BitClear:
		part.KnownZero = 1

	case knownbits.BitSet:
		part.KnownOne = 1
	}

	tt.write(view(AF, position, 1), part)
}

func (tt *transfer) szxy(value knownbits.Value[uint16]) {
	sign, _ := knownbits.Sign(8, value)
	zero, _ := knownbits.Zero(8, value)
	tt.writeFlag(7, sign)
	tt.writeFlag(6, zero)
	for _, position := range []knownbits.Width{3, 5} {
		bit, _ := knownbits.TestBit(8, value, position)
		tt.writeFlag(position, bit)
	}
}

func (tt *transfer) byteALU(code byte) {
	left := tt.read(view(AF, 8, 8))
	operation := (code >> 3) & 3
	var right knownbits.Value[uint16]
	self := false
	switch {
	case code&0xc0 == 0xc0:
		right = tt.input.Operand

	case code&7 == 6:
		right = tt.readMemory(tt.indirectAddress())

	default:
		source := tt.byteRegister(code&7, false)
		self = source == view(AF, 8, 8)
		right = tt.read(source)
	}

	if operation == 3 {
		tt.compare(left, right, self)
		return
	}
	value := logical(operation, left, right)
	if self && operation == 1 {
		value, _ = knownbits.NewConstant[uint16](8, 0)
	}

	tt.write(view(AF, 8, 8), value)
	tt.szxy(value)
	parity, _ := knownbits.EvenParity(8, value)
	tt.writeFlag(2, parity)
	tt.writeFlag(0, knownbits.BitClear)
	tt.writeFlag(1, knownbits.BitClear)
	half := knownbits.BitClear
	if operation == 0 {
		half = knownbits.BitSet
	}

	tt.writeFlag(4, half)
	tt.logicalCopy(operation, right, self)
}

func (tt *transfer) compare(left, right knownbits.Value[uint16], self bool) {
	difference, _ := knownbits.Sub(8, left, right, knownbits.BitClear)
	half, _ := knownbits.BorrowAt(8, left, right, knownbits.BitClear, 4)
	if self {
		difference.Value, _ = knownbits.NewConstant[uint16](8, 0)
		difference.BorrowOut = knownbits.BitClear
		difference.SignedOverflow = knownbits.BitClear
		half = knownbits.BitClear
	}

	tt.szxy(difference.Value)
	for _, position := range []knownbits.Width{3, 5} {
		bit, _ := knownbits.TestBit(8, right, position)
		tt.writeFlag(position, bit)
	}

	tt.writeFlag(0, difference.BorrowOut)
	tt.writeFlag(1, knownbits.BitSet)
	tt.writeFlag(2, difference.SignedOverflow)
	tt.writeFlag(4, half)
}

func (tt *transfer) carryControl(complement bool) {
	accumulator := tt.read(view(AF, 8, 8))
	oldXY := knownbits.Value[uint16]{}
	for _, position := range []knownbits.Width{3, 5} {
		bit := tt.read(view(AF, position, 1))
		oldXY.KnownZero |= bit.KnownZero << position
		oldXY.KnownOne |= bit.KnownOne << position
	}

	q := tt.read(view(Q, 0, 8))
	notLatch, _ := knownbits.Not(8, q)
	retained, _ := knownbits.And(8, oldXY, notLatch)
	xy, _ := knownbits.Or(8, accumulator, retained)
	for _, position := range []knownbits.Width{3, 5} {
		bit, _ := knownbits.TestBit(8, xy, position)
		tt.writeFlag(position, bit)
	}

	carry := knownbits.BitSet
	half := knownbits.BitClear
	if complement {
		half = tt.flag(0)
		carry = invert(half)
	}

	tt.writeFlag(0, carry)
	tt.writeFlag(1, knownbits.BitClear)
	tt.writeFlag(4, half)
}

func logical(operation byte, left, right knownbits.Value[uint16]) knownbits.Value[uint16] {
	switch operation {
	case 0:
		result, _ := knownbits.And(8, left, right)
		return result

	case 1:
		result, _ := knownbits.Xor(8, left, right)
		return result

	default:
		result, _ := knownbits.Or(8, left, right)
		return result
	}
}

func invert(value knownbits.Bit) knownbits.Bit {
	switch value {
	case knownbits.BitClear:
		return knownbits.BitSet

	case knownbits.BitSet:
		return knownbits.BitClear

	default:
		return knownbits.BitUnknown
	}
}
