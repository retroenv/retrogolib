package knownbits

// ShiftLeft fills low bits with zero. Counts at least width produce zero.
func ShiftLeft[W Word](width Width, value Value[W], count uint) (Value[W], bool) {
	mask, valid := widthMask[W](width)
	if !valid || !value.Valid(width) {
		return Value[W]{}, false
	}

	count = min(count, uint(width))
	fill := mask & ^(mask << count)

	return New(width, value.KnownZero<<count|W(fill), value.KnownOne<<count)
}

// ShiftRight fills high bits with zero. Counts at least width produce zero.
func ShiftRight[W Word](width Width, value Value[W], count uint) (Value[W], bool) {
	mask, valid := widthMask[W](width)
	if !valid || !value.Valid(width) {
		return Value[W]{}, false
	}

	count = min(count, uint(width))
	fill := mask & ^(mask >> count)

	return Value[W]{
		KnownZero: value.KnownZero>>count | W(fill),
		KnownOne:  value.KnownOne >> count,
	}, true
}

// ShiftRightSigned fills high bits with the original sign bit.
// Counts at least width produce only copies of that bit.
func ShiftRightSigned[W Word](width Width, value Value[W], count uint) (Value[W], bool) {
	mask, valid := widthMask[W](width)
	if !valid || !value.Valid(width) {
		return Value[W]{}, false
	}

	count = min(count, uint(width))
	fill := W(mask & ^(mask >> count))
	result := Value[W]{
		KnownZero: value.KnownZero >> count,
		KnownOne:  value.KnownOne >> count,
	}
	if value.KnownZero>>(width-1)&1 != 0 {
		result.KnownZero |= fill
	}

	if value.KnownOne>>(width-1)&1 != 0 {
		result.KnownOne |= fill
	}

	return result, true
}

// RotateLeft rotates within width. Counts are reduced modulo width.
func RotateLeft[W Word](width Width, value Value[W], count uint) (Value[W], bool) {
	if !value.Valid(width) {
		return Value[W]{}, false
	}

	count %= uint(width)
	other := uint(width) - count

	return New(width, value.KnownZero<<count|value.KnownZero>>other, value.KnownOne<<count|value.KnownOne>>other)
}

// RotateRight rotates within width. Counts are reduced modulo width.
func RotateRight[W Word](width Width, value Value[W], count uint) (Value[W], bool) {
	if !value.Valid(width) {
		return Value[W]{}, false
	}

	count %= uint(width)
	other := uint(width) - count

	return New(width, value.KnownZero>>count|value.KnownZero<<other, value.KnownOne>>count|value.KnownOne<<other)
}

// RotateLeftCarry rotates one bit through a separate carry input.
// The outgoing carry is the original high bit.
func RotateLeftCarry[W Word](width Width, value Value[W], carry Bit) (Value[W], Bit, bool) {
	if !value.Valid(width) || carry > BitSet {
		return Value[W]{}, BitUnknown, false
	}

	outgoing, _ := TestBit(width, value, width-1)
	result, _ := New(width, value.KnownZero<<1, value.KnownOne<<1)

	return fillCarry(result, 0, carry), outgoing, true
}

// RotateRightCarry rotates one bit through a separate carry input.
// The outgoing carry is the original low bit.
func RotateRightCarry[W Word](width Width, value Value[W], carry Bit) (Value[W], Bit, bool) {
	if !value.Valid(width) || carry > BitSet {
		return Value[W]{}, BitUnknown, false
	}

	outgoing, _ := TestBit(width, value, 0)
	result := Value[W]{
		KnownZero: value.KnownZero >> 1,
		KnownOne:  value.KnownOne >> 1,
	}

	return fillCarry(result, width-1, carry), outgoing, true
}

func fillCarry[W Word](value Value[W], position Width, carry Bit) Value[W] {
	switch carry {
	case BitClear:
		value.KnownZero |= W(1) << position
	case BitSet:
		value.KnownOne |= W(1) << position
	}

	return value
}
