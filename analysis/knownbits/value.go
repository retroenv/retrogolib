package knownbits

// Bit describes a predicate whose truth can be unknown.
type Bit uint8

const (
	BitUnknown Bit = iota
	BitClear
	BitSet
)

// Width is the number of active bits, from one through the Word size.
type Width uint8

// Word selects compact unsigned storage. Named unsigned types are supported.
type Word interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Value stores independent bit facts. Its fields must not overlap.
// Bits outside the active width must be clear.
type Value[W Word] struct {
	KnownZero W
	KnownOne  W
}

// New clears masks outside width and rejects contradictions inside width.
func New[W Word](width Width, zero, one W) (Value[W], bool) {
	mask, valid := widthMask[W](width)
	if !valid || uint64(zero&one)&mask != 0 {
		return Value[W]{}, false
	}
	return Value[W]{
		KnownZero: zero & W(mask),
		KnownOne:  one & W(mask),
	}, true
}

// NewConstant keeps the low width bits of number.
func NewConstant[W Word](width Width, number W) (Value[W], bool) {
	return New(width, ^number, number)
}

// Valid reports whether the width and masks are canonical.
func (vv Value[W]) Valid(width Width) bool {
	return validMasks(width, vv.KnownZero&vv.KnownOne, vv.KnownZero|vv.KnownOne)
}

// Constant returns the number only when every active bit is known.
func (vv Value[W]) Constant(width Width) (W, bool) {
	mask, valid := widthMask[W](width)
	if !valid || !vv.Valid(width) || uint64(vv.KnownZero|vv.KnownOne) != mask {
		return 0, false
	}
	return vv.KnownOne, true
}

// UnsignedRange returns the smallest and largest permitted numbers.
func (vv Value[W]) UnsignedRange(width Width) (W, W, bool) {
	mask, valid := widthMask[W](width)
	if !valid || !vv.Valid(width) {
		return 0, 0, false
	}
	return vv.KnownOne, ^vv.KnownZero & W(mask), true
}

func widthMask[W Word](width Width) (uint64, bool) {
	if width == 0 || ^W(0)>>(width-1) == 0 {
		return 0, false
	}
	return ^uint64(0) >> (64 - width), true
}

// validMasks checks both storage width and active width without building a mask.
func validMasks[W Word](width Width, conflicts, known W) bool {
	return width > 0 && ^W(0)>>(width-1) != 0 && conflicts == 0 && known>>(width-1) <= 1
}

func invertBit(value Bit) Bit {
	switch value {
	case BitClear:
		return BitSet
	case BitSet:
		return BitClear
	default:
		return BitUnknown
	}
}

func bitBounds(value Bit) (uint64, uint64) {
	switch value {
	case BitClear:
		return 0, 0
	case BitSet:
		return 1, 1
	default:
		return 0, 1
	}
}

func bitFromBounds(low, high uint64) Bit {
	if low != high {
		return BitUnknown
	}
	if low == 0 {
		return BitClear
	}
	return BitSet
}
