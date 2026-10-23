package knownbits

import "math/bits"

// TestBit reports the bit at position, counted from zero.
func TestBit[W Word](width Width, value Value[W], position Width) (Bit, bool) {
	if !value.Valid(width) || position >= width {
		return BitUnknown, false
	}
	if value.KnownZero>>position&1 != 0 {
		return BitClear, true
	}
	if value.KnownOne>>position&1 != 0 {
		return BitSet, true
	}
	return BitUnknown, true
}

// Zero reports whether all active bits are zero.
func Zero[W Word](width Width, value Value[W]) (Bit, bool) {
	if !value.Valid(width) {
		return BitUnknown, false
	}
	if value.KnownOne != 0 {
		return BitClear, true
	}
	if _, known := value.Constant(width); known {
		return BitSet, true
	}
	return BitUnknown, true
}

// Sign reports the active high bit.
func Sign[W Word](width Width, value Value[W]) (Bit, bool) {
	return TestBit(width, value, width-1)
}

// EvenParity is set when an even number of active bits are one.
// Pass a byte slice when a CPU tests parity only in its low byte.
func EvenParity[W Word](width Width, value Value[W]) (Bit, bool) {
	if !value.Valid(width) {
		return BitUnknown, false
	}
	number, known := value.Constant(width)
	if !known {
		return BitUnknown, true
	}
	if bits.OnesCount64(uint64(number))%2 == 0 {
		return BitSet, true
	}
	return BitClear, true
}

// Equal reports concrete equality without assuming operand identity.
func Equal[W Word](width Width, left, right Value[W]) (Bit, bool) {
	if !left.Valid(width) || !right.Valid(width) {
		return BitUnknown, false
	}
	if left.KnownZero&right.KnownOne|left.KnownOne&right.KnownZero != 0 {
		return BitClear, true
	}
	_, leftKnown := left.Constant(width)
	_, rightKnown := right.Constant(width)
	if leftKnown && rightKnown {
		return BitSet, true
	}
	return BitUnknown, true
}
