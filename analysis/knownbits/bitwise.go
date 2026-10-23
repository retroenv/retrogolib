package knownbits

// Join keeps only facts common to both inputs.
func Join[W Word](width Width, left, right Value[W]) (Value[W], bool) {
	return Value[W]{
		KnownZero: left.KnownZero & right.KnownZero,
		KnownOne:  left.KnownOne & right.KnownOne,
	}, validMasks(width, left.KnownZero&left.KnownOne|right.KnownZero&right.KnownOne,
		left.KnownZero|left.KnownOne|right.KnownZero|right.KnownOne)
}

// And computes independent bit facts for left AND right.
func And[W Word](width Width, left, right Value[W]) (Value[W], bool) {
	return Value[W]{
		KnownZero: left.KnownZero | right.KnownZero,
		KnownOne:  left.KnownOne & right.KnownOne,
	}, validMasks(width, left.KnownZero&left.KnownOne|right.KnownZero&right.KnownOne,
		left.KnownZero|left.KnownOne|right.KnownZero|right.KnownOne)
}

// Or computes independent bit facts for left OR right.
func Or[W Word](width Width, left, right Value[W]) (Value[W], bool) {
	return Value[W]{
		KnownZero: left.KnownZero & right.KnownZero,
		KnownOne:  left.KnownOne | right.KnownOne,
	}, validMasks(width, left.KnownZero&left.KnownOne|right.KnownZero&right.KnownOne,
		left.KnownZero|left.KnownOne|right.KnownZero|right.KnownOne)
}

// Xor does not infer concrete operand identity from equal facts.
func Xor[W Word](width Width, left, right Value[W]) (Value[W], bool) {
	return Value[W]{
		KnownZero: left.KnownZero&right.KnownZero | left.KnownOne&right.KnownOne,
		KnownOne:  left.KnownZero&right.KnownOne | left.KnownOne&right.KnownZero,
	}, validMasks(width, left.KnownZero&left.KnownOne|right.KnownZero&right.KnownOne,
		left.KnownZero|left.KnownOne|right.KnownZero|right.KnownOne)
}

// Not complements every active bit.
func Not[W Word](width Width, value Value[W]) (Value[W], bool) {
	if !value.Valid(width) {
		return Value[W]{}, false
	}
	return Value[W]{
		KnownZero: value.KnownOne,
		KnownOne:  value.KnownZero,
	}, true
}
