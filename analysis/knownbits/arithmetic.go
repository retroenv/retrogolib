package knownbits

import "math/bits"

// AddResult contains modular sum facts and mathematical events.
type AddResult[W Word] struct {
	Value          Value[W]
	CarryOut       Bit
	SignedOverflow Bit
}

// SubResult contains modular difference facts and mathematical events.
// BorrowOut is set when unsigned subtraction needs a borrow.
type SubResult[W Word] struct {
	Value          Value[W]
	BorrowOut      Bit
	SignedOverflow Bit
}

// Add computes left + right + carryIn at width.
// Inputs have independent unknown bits, even when their facts are equal.
func Add[W Word](width Width, left, right Value[W], carryIn Bit) (AddResult[W], bool) {
	if !left.Valid(width) || !right.Valid(width) || carryIn > BitSet {
		return AddResult[W]{}, false
	}
	_, leftKnown := left.Constant(width)
	_, rightKnown := right.Constant(width)
	if leftKnown && rightKnown && carryIn != BitUnknown {
		return addConstants(width, left.KnownOne, right.KnownOne, carryIn), true
	}
	return addPartial(width, left, right, carryIn), true
}

// Sub computes left - right - borrowIn at width.
// A CPU with a no-borrow carry flag must invert that flag at the boundary.
func Sub[W Word](width Width, left, right Value[W], borrowIn Bit) (SubResult[W], bool) {
	if !left.Valid(width) || !right.Valid(width) || borrowIn > BitSet {
		return SubResult[W]{}, false
	}
	complement, _ := Not(width, right)
	addition, _ := Add(width, left, complement, invertBit(borrowIn))
	return SubResult[W]{
		Value:          addition.Value,
		BorrowOut:      invertBit(addition.CarryOut),
		SignedOverflow: addition.SignedOverflow,
	}, true
}

// CarryAt reports the carry after boundary low bits of addition.
// The boundary must be from one through width.
func CarryAt[W Word](width Width, left, right Value[W], carryIn Bit, boundary Width) (Bit, bool) {
	if !left.Valid(width) || !right.Valid(width) || boundary == 0 || boundary > width {
		return BitUnknown, false
	}
	lowLeft, _ := Truncate(width, boundary, left)
	lowRight, _ := Truncate(width, boundary, right)
	result, valid := Add(boundary, lowLeft, lowRight, carryIn)
	return result.CarryOut, valid
}

// BorrowAt reports the borrow after boundary low bits of subtraction.
// The boundary must be from one through width.
func BorrowAt[W Word](width Width, left, right Value[W], borrowIn Bit, boundary Width) (Bit, bool) {
	if !left.Valid(width) || !right.Valid(width) || boundary == 0 || boundary > width {
		return BitUnknown, false
	}
	lowLeft, _ := Truncate(width, boundary, left)
	lowRight, _ := Truncate(width, boundary, right)
	result, valid := Sub(boundary, lowLeft, lowRight, borrowIn)
	return result.BorrowOut, valid
}

func addConstants[W Word](width Width, left, right W, carryIn Bit) AddResult[W] {
	carry, _ := bitBounds(carryIn)
	sum, carryOut := bits.Add64(uint64(left), uint64(right), carry)
	if width < 64 {
		carryOut = sum >> width
	}
	value, _ := NewConstant(width, W(sum))
	sign := W(1) << (width - 1)
	overflow := BitClear
	if ^(left^right)&(left^value.KnownOne)&sign != 0 {
		overflow = BitSet
	}
	return AddResult[W]{
		Value:          value,
		CarryOut:       bitFromBounds(carryOut, carryOut),
		SignedOverflow: overflow,
	}
}

func addPartial[W Word](width Width, left, right Value[W], carryIn Bit) AddResult[W] {
	carryMin, carryMax := bitBounds(carryIn)
	result := AddResult[W]{}
	for position := range width {
		leftMin := uint64(left.KnownOne>>position) & 1
		leftMax := uint64(^left.KnownZero>>position) & 1
		rightMin := uint64(right.KnownOne>>position) & 1
		rightMax := uint64(^right.KnownZero>>position) & 1
		if position == width-1 {
			result.SignedOverflow = additionOverflow(leftMin, leftMax, rightMin, rightMax, carryMin, carryMax)
		}
		low := leftMin + rightMin + carryMin
		high := leftMax + rightMax + carryMax
		if low == high {
			if low&1 == 0 {
				result.Value.KnownZero |= W(1) << position
			} else {
				result.Value.KnownOne |= W(1) << position
			}
		}
		carryMin = low >> 1
		carryMax = high >> 1
	}
	result.CarryOut = bitFromBounds(carryMin, carryMax)
	return result
}

func additionOverflow(leftMin, leftMax, rightMin, rightMax, carryMin, carryMax uint64) Bit {
	low := uint64(1)
	high := uint64(0)
	for left := leftMin; left <= leftMax; left++ {
		for right := rightMin; right <= rightMax; right++ {
			for carry := carryMin; carry <= carryMax; carry++ {
				event := carry ^ ((left + right + carry) >> 1)
				low = min(low, event)
				high = max(high, event)
			}
		}
	}
	return bitFromBounds(low, high)
}
