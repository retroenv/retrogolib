package knownbits

// Join keeps only facts common to both inputs.
func Join[W Word](width Width, left, right Value[W]) (Value[W], bool) {
	result := Value[W]{
		KnownZero: left.KnownZero & right.KnownZero,
		KnownOne:  left.KnownOne & right.KnownOne,
	}
	conflicts := left.KnownZero&left.KnownOne | right.KnownZero&right.KnownOne
	known := left.KnownZero | left.KnownOne | right.KnownZero | right.KnownOne

	return result, validMasks(width, conflicts, known)
}

// Meet combines two sets of facts about one concrete value.
// Feasible is false when the facts contradict. The result is then unspecified.
func Meet[W Word](width Width, left, right Value[W]) (Value[W], bool, bool) {
	if !left.Valid(width) || !right.Valid(width) {
		return Value[W]{}, false, false
	}
	result := Value[W]{
		KnownZero: left.KnownZero | right.KnownZero,
		KnownOne:  left.KnownOne | right.KnownOne,
	}
	return result, result.KnownZero&result.KnownOne == 0, true
}

// And computes independent bit facts for left AND right.
func And[W Word](width Width, left, right Value[W]) (Value[W], bool) {
	result := Value[W]{
		KnownZero: left.KnownZero | right.KnownZero,
		KnownOne:  left.KnownOne & right.KnownOne,
	}
	conflicts := left.KnownZero&left.KnownOne | right.KnownZero&right.KnownOne
	known := left.KnownZero | left.KnownOne | right.KnownZero | right.KnownOne

	return result, validMasks(width, conflicts, known)
}

// Or computes independent bit facts for left OR right.
func Or[W Word](width Width, left, right Value[W]) (Value[W], bool) {
	result := Value[W]{
		KnownZero: left.KnownZero & right.KnownZero,
		KnownOne:  left.KnownOne | right.KnownOne,
	}
	conflicts := left.KnownZero&left.KnownOne | right.KnownZero&right.KnownOne
	known := left.KnownZero | left.KnownOne | right.KnownZero | right.KnownOne

	return result, validMasks(width, conflicts, known)
}

// Xor does not infer concrete operand identity from equal facts.
func Xor[W Word](width Width, left, right Value[W]) (Value[W], bool) {
	result := Value[W]{
		KnownZero: left.KnownZero&right.KnownZero | left.KnownOne&right.KnownOne,
		KnownOne:  left.KnownZero&right.KnownOne | left.KnownOne&right.KnownZero,
	}
	conflicts := left.KnownZero&left.KnownOne | right.KnownZero&right.KnownOne
	known := left.KnownZero | left.KnownOne | right.KnownZero | right.KnownOne

	return result, validMasks(width, conflicts, known)
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
