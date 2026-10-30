package knownbits

// Truncate keeps the low to bits. It cannot increase the width.
func Truncate[W Word](from, to Width, value Value[W]) (Value[W], bool) {
	if !value.Valid(from) || to > from {
		return Value[W]{}, false
	}
	return New(to, value.KnownZero, value.KnownOne)
}

// Extract selects size bits at offset, counted from the least significant bit.
// Bit positions do not depend on memory byte order.
func Extract[W Word](width Width, value Value[W], offset, size Width) (Value[W], bool) {
	if !value.Valid(width) || uint16(offset)+uint16(size) > uint16(width) {
		return Value[W]{}, false
	}
	return New(size, value.KnownZero>>offset, value.KnownOne>>offset)
}

// Insert replaces size bits at offset and preserves every other parent bit.
// The parent and part use the same Word storage type.
func Insert[W Word](width Width, parent Value[W], offset, size Width, part Value[W]) (Value[W], bool) {
	mask, valid := widthMask[W](size)
	if !valid || !parent.Valid(width) || !part.Valid(size) || uint16(offset)+uint16(size) > uint16(width) {
		return Value[W]{}, false
	}
	slice := W(mask) << offset

	return Value[W]{
		KnownZero: parent.KnownZero&^slice | part.KnownZero<<offset,
		KnownOne:  parent.KnownOne&^slice | part.KnownOne<<offset,
	}, true
}

// Concat puts high above low. The combined width must fit in W.
func Concat[W Word](highWidth Width, high Value[W], lowWidth Width, low Value[W]) (Value[W], bool) {
	total := uint16(highWidth) + uint16(lowWidth)
	if !high.Valid(highWidth) || !low.Valid(lowWidth) || total > 64 {
		return Value[W]{}, false
	}
	return New(Width(total), high.KnownZero<<lowWidth|low.KnownZero, high.KnownOne<<lowWidth|low.KnownOne)
}

// ZeroExtend proves that each added high bit is zero.
func ZeroExtend[W Word](from, to Width, value Value[W]) (Value[W], bool) {
	mask, valid := widthMask[W](to)
	if !valid || !value.Valid(from) || to < from {
		return Value[W]{}, false
	}
	oldMask, _ := widthMask[W](from)

	return Value[W]{
		KnownZero: value.KnownZero | W(mask&^oldMask),
		KnownOne:  value.KnownOne,
	}, true
}

// SignExtend copies the known sign bit into each added high bit.
func SignExtend[W Word](from, to Width, value Value[W]) (Value[W], bool) {
	mask, valid := widthMask[W](to)
	if !valid || !value.Valid(from) || to < from {
		return Value[W]{}, false
	}

	oldMask, _ := widthMask[W](from)
	high := W(mask &^ oldMask)
	if value.KnownZero>>(from-1)&1 != 0 {
		value.KnownZero |= high
	}

	if value.KnownOne>>(from-1)&1 != 0 {
		value.KnownOne |= high
	}

	return value, true
}
