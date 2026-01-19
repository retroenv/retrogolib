package set

import (
	"math/bits"
	"strconv"
	"strings"
)

// BitSet represents a set of integers in range [0, 63] using a single uint64.
// This provides zero-allocation set operations for small integer sets where
// the maximum value is known to be < 64 (like block IDs, register indices, etc).
//
// BitSet is a value type (not a pointer) and can be safely copied.
type BitSet uint64

// NewBitSet creates a new empty BitSet.
func NewBitSet() BitSet {
	return 0
}

// NewBitSetFromSlice creates a BitSet from a slice of integers.
// Values >= 64 will cause a panic.
func NewBitSetFromSlice(values []int) BitSet {
	var b BitSet
	for _, v := range values {
		b.Add(v)
	}
	return b
}

// Add adds an integer to the set.
// Panics if n < 0 or n >= 64.
func (b *BitSet) Add(n int) {
	if n < 0 || n >= 64 {
		panic("BitSet: value must be in range [0, 63]")
	}
	*b |= 1 << n
}

// Remove removes an integer from the set.
// Panics if n < 0 or n >= 64.
func (b *BitSet) Remove(n int) {
	if n < 0 || n >= 64 {
		panic("BitSet: value must be in range [0, 63]")
	}
	*b &^= 1 << n
}

// Contains checks if an integer exists in the set.
// Panics if n < 0 or n >= 64.
func (b *BitSet) Contains(n int) bool {
	if n < 0 || n >= 64 {
		panic("BitSet: value must be in range [0, 63]")
	}
	return (*b)&(1<<n) != 0
}

// Size returns the number of elements in the set.
func (b *BitSet) Size() int {
	return bits.OnesCount64(uint64(*b))
}

// IsEmpty returns true if the set is empty.
func (b *BitSet) IsEmpty() bool {
	return *b == 0
}

// Clear removes all elements from the set.
func (b *BitSet) Clear() {
	*b = 0
}

// ToSlice returns all elements as a slice in ascending order.
func (b *BitSet) ToSlice() []int {
	result := make([]int, 0, b.Size())
	for i := range 64 {
		if (*b)&(1<<i) != 0 {
			result = append(result, i)
		}
	}
	return result
}

// Copy creates a copy of the set.
// Since BitSet is a value type, this is just a simple assignment.
func (b *BitSet) Copy() BitSet {
	return *b
}

// Equal checks if two sets are equal.
func (b *BitSet) Equal(other *BitSet) bool {
	return *b == *other
}

// Union returns a new set containing all elements from both sets.
func (b *BitSet) Union(other *BitSet) BitSet {
	return (*b) | (*other)
}

// Intersection returns a new set containing elements that exist in both sets.
func (b *BitSet) Intersection(other *BitSet) BitSet {
	return (*b) & (*other)
}

// Difference returns a new set containing elements that are in this set but not in the other.
func (b *BitSet) Difference(other *BitSet) BitSet {
	return (*b) &^ (*other)
}

// IsSubset checks if this set is a subset of the other set.
func (b *BitSet) IsSubset(other *BitSet) bool {
	return ((*b) & (*other)) == *b
}

// IsDisjoint checks if this set has no elements in common with the other set.
func (b *BitSet) IsDisjoint(other *BitSet) bool {
	return ((*b) & (*other)) == 0
}

// ForEach applies a function to each element in the set in ascending order.
func (b *BitSet) ForEach(fn func(int)) {
	for i := range 64 {
		if (*b)&(1<<i) != 0 {
			fn(i)
		}
	}
}

// Min returns the minimum element in the set.
// Returns -1 if the set is empty.
func (b *BitSet) Min() int {
	if *b == 0 {
		return -1
	}
	return bits.TrailingZeros64(uint64(*b))
}

// Max returns the maximum element in the set.
// Returns -1 if the set is empty.
func (b *BitSet) Max() int {
	if *b == 0 {
		return -1
	}
	return 63 - bits.LeadingZeros64(uint64(*b))
}

// PopFirst removes and returns the minimum element.
// Returns -1 if the set is empty.
func (b *BitSet) PopFirst() int {
	if *b == 0 {
		return -1
	}
	n := bits.TrailingZeros64(uint64(*b))
	*b &^= 1 << n
	return n
}

// String returns a string representation of the set for debugging.
// Format: "BitSet{1, 3, 5}" or "BitSet{}" if empty.
func (b *BitSet) String() string {
	if *b == 0 {
		return "BitSet{}"
	}

	var sb strings.Builder
	sb.WriteString("BitSet{")

	first := true
	for i := range 64 {
		if (*b)&(1<<i) != 0 {
			if !first {
				sb.WriteString(", ")
			}
			sb.WriteString(strconv.Itoa(i))
			first = false
		}
	}

	sb.WriteByte('}')
	return sb.String()
}
