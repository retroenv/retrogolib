package analysis

import "github.com/retroenv/retrogolib/analysis/knownbits"

// Reason distinguishes supported unknown values from refused transfers.
type Reason uint8

const (
	UnknownEffects Reason = iota
	Supported
	UnsupportedOpcode
	IllegalOperand
	InvalidInput
	UnknownMode
	UnsupportedMode
	UnsupportedVariant
)

// Flow classifies local control flow. Calls need a caller-owned summary or barrier.
type Flow uint8

const (
	FlowUnknown Flow = iota
	FlowNext
	FlowBranch
	FlowJump
	FlowCall
	FlowReturn
	FlowTrap
	FlowStop
)

// WriteRule defines the effect of a narrow resource write.
type WriteRule uint8

const (
	Replace WriteRule = iota
	Preserve
	ZeroExtend
	SignExtend
	Invalidate
)

// MemoryKind describes one semantic access. Read-modify-write uses two events.
type MemoryKind uint8

const (
	MemoryRead MemoryKind = iota
	MemoryWrite
)

// View selects a logical bit slice of target-local storage.
type View[R comparable] struct {
	Storage R
	Offset  knownbits.Width
	Bits    knownbits.Width
}

// Write replaces a view under its explicit write rule. Unknown facts still clobber.
type Write[R comparable, W knownbits.Word] struct {
	View  View[R]
	Value knownbits.Value[W]
	Rule  WriteRule
}

// Copy proves that a destination receives a source from the input snapshot.
// It carries no caller identity and does not imply later memory equivalence.
type Copy[R comparable] struct {
	Source      View[R]
	Destination View[R]
}

// MemoryEvent reports an access in target-defined logical address space.
// The caller must bind that address to bank, mapper, and physical-alias context.
// Stable is a read proof supplied by the caller. It is false for writes.
// Value describes the accessed scalar when it fits in W; otherwise it is unknown.
type MemoryEvent[W knownbits.Word, A any] struct {
	Address A
	Bytes   uint32
	Kind    MemoryKind
	Value   knownbits.Value[W]
	Stable  bool
}

// Result holds local effects and established facts. Its zero value fails closed.
// Condition is meaningful for FlowBranch and can be unknown.
type Result[R comparable, W knownbits.Word, A any] struct {
	Reason          Reason
	EffectsComplete bool
	Flow            Flow
	Condition       knownbits.Bit

	Reads  []View[R]
	Writes []Write[R, W]
	Copies []Copy[R]
	Memory []MemoryEvent[W, A]
}

// Valid checks view bounds against a declared storage width, at most 64 bits.
func (vw View[R]) Valid(storageWidth knownbits.Width) bool {
	return storageWidth > 0 && storageWidth <= 64 && vw.Bits > 0 &&
		uint16(vw.Offset)+uint16(vw.Bits) <= uint16(storageWidth)
}

// Overlaps compares slices directly. Sibling bytes do not overlap through a pair.
func (vw View[R]) Overlaps(other View[R]) bool {
	return vw.Valid(64) && other.Valid(64) && vw.Storage == other.Storage &&
		uint16(vw.Offset) < uint16(other.Offset)+uint16(other.Bits) &&
		uint16(other.Offset) < uint16(vw.Offset)+uint16(vw.Bits)
}

// Affected returns the full invalidation view under the write rule.
// Extension changes all storage bits, including bits above the input view.
func (ww Write[R, W]) Affected(width knownbits.Width) (View[R], bool) {
	part := ww.View
	if !part.Valid(width) {
		return View[R]{}, false
	}

	switch ww.Rule {
	case Replace:
		return part, part.Offset == 0 && part.Bits == width

	case Preserve, Invalidate:
		return part, true

	case ZeroExtend, SignExtend:
		return View[R]{
			Storage: part.Storage,
			Bits:    width,
		}, part.Offset == 0

	default:
		return View[R]{}, false
	}
}

// Reset clears the result and retains caller-owned slice capacity.
func (rr *Result[R, W, A]) Reset() {
	rr.Reason = UnknownEffects
	rr.EffectsComplete = false
	rr.Flow = FlowUnknown
	rr.Condition = knownbits.BitUnknown
	rr.Reads = rr.Reads[:0]
	rr.Writes = rr.Writes[:0]
	rr.Copies = rr.Copies[:0]
	rr.Memory = rr.Memory[:0]
}

// Valid checks output bounds, disjoint writes, and copy coverage against target widths.
// Address validity remains target-owned. Use this at validation boundaries.
func (rr *Result[R, W, A]) Valid(width func(R) knownbits.Width) bool {
	if rr.Reason > UnsupportedVariant || rr.Flow > FlowStop || rr.Condition > knownbits.BitSet {
		return false
	}
	if rr.EffectsComplete && rr.Reason == UnknownEffects {
		return false
	}

	for _, read := range rr.Reads {
		if !read.Valid(width(read.Storage)) {
			return false
		}
	}

	if !rr.validWrites(width) {
		return false
	}
	for _, proof := range rr.Copies {
		if !rr.validCopy(proof, width) {
			return false
		}
	}

	for _, event := range rr.Memory {
		if !validEvent(event) {
			return false
		}
	}

	return true
}

func (rr *Result[R, W, A]) validWrites(width func(R) knownbits.Width) bool {
	for index, write := range rr.Writes {
		if _, valid := Apply(width(write.View.Storage), knownbits.Value[W]{}, write); !valid {
			return false
		}
		affected, _ := write.Affected(width(write.View.Storage))
		for _, other := range rr.Writes[:index] {
			otherAffected, _ := other.Affected(width(other.View.Storage))
			if affected.Overlaps(otherAffected) {
				return false
			}
		}
	}

	return true
}

func (rr *Result[R, W, A]) validCopy(proof Copy[R], width func(R) knownbits.Width) bool {
	if !proof.Source.Valid(width(proof.Source.Storage)) || !proof.Destination.Valid(width(proof.Destination.Storage)) ||
		proof.Source.Bits != proof.Destination.Bits {

		return false
	}

	for _, write := range rr.Writes {
		destination := proof.Destination
		affected, _ := write.Affected(width(write.View.Storage))
		if affected.Storage == destination.Storage && affected.Offset <= destination.Offset &&
			uint16(affected.Offset)+uint16(affected.Bits) >= uint16(destination.Offset)+uint16(destination.Bits) {

			return true
		}
	}

	return false
}

// Apply computes one storage value. The caller must supply the selected storage.
// Preserve and Invalidate affect only the selected view. Extension requires offset zero.
func Apply[R comparable, W knownbits.Word](width knownbits.Width, before knownbits.Value[W],
	write Write[R, W]) (knownbits.Value[W], bool) {

	view := write.View
	if !before.Valid(width) || !view.Valid(width) || !write.Value.Valid(view.Bits) {
		return knownbits.Value[W]{}, false
	}

	switch write.Rule {
	case Replace:
		return write.Value, view.Offset == 0 && view.Bits == width && write.Value.Valid(width)

	case Preserve:
		return knownbits.Insert(width, before, view.Offset, view.Bits, write.Value)

	case Invalidate:
		return knownbits.Insert(width, before, view.Offset, view.Bits, knownbits.Value[W]{})

	case ZeroExtend:
		if view.Offset == 0 {
			return knownbits.ZeroExtend(view.Bits, width, write.Value)
		}

	case SignExtend:
		if view.Offset == 0 {
			return knownbits.SignExtend(view.Bits, width, write.Value)
		}
	}

	return knownbits.Value[W]{}, false
}

func validEvent[W knownbits.Word, A any](event MemoryEvent[W, A]) bool {
	if event.Kind > MemoryWrite || event.Kind == MemoryWrite && event.Stable {
		return false
	}
	if event.Bytes > 0 && event.Bytes <= 8 && event.Value.Valid(knownbits.Width(event.Bytes*8)) {
		return true
	}
	return event.Value == (knownbits.Value[W]{})
}
