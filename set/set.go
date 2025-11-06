package set

// Set represents a generic set data structure using a map with struct{} values.
type Set[T comparable] map[T]struct{}

// New creates a new empty set.
func New[T comparable]() Set[T] {
	return make(Set[T])
}

// NewFromSlice creates a new set from a slice of elements.
func NewFromSlice[T comparable](elements []T) Set[T] {
	s := make(Set[T], len(elements))
	for _, element := range elements {
		s[element] = struct{}{}
	}
	return s
}

// Add adds an element to the set.
func (s Set[T]) Add(element T) {
	s[element] = struct{}{}
}

// Remove removes an element from the set.
func (s Set[T]) Remove(element T) {
	delete(s, element)
}

// Contains checks if an element exists in the set.
func (s Set[T]) Contains(element T) bool {
	_, exists := s[element]
	return exists
}

// Size returns the number of elements in the set.
func (s Set[T]) Size() int {
	return len(s)
}

// IsEmpty returns true if the set is empty.
func (s Set[T]) IsEmpty() bool {
	return len(s) == 0
}

// Clear removes all elements from the set.
func (s Set[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

// ToSlice returns all elements as a slice.
func (s Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s))
	for element := range s {
		result = append(result, element)
	}
	return result
}

// Copy creates a deep copy of the set.
func (s Set[T]) Copy() Set[T] {
	result := make(Set[T], len(s))
	for element := range s {
		result[element] = struct{}{}
	}
	return result
}

// Equal checks if two sets are equal (contain the same elements).
func (s Set[T]) Equal(other Set[T]) bool {
	if len(s) != len(other) {
		return false
	}

	for element := range s {
		if !other.Contains(element) {
			return false
		}
	}

	return true
}

// Union returns a new set containing all elements from both sets.
func (s Set[T]) Union(other Set[T]) Set[T] {
	result := make(Set[T], len(s)+len(other))

	for element := range s {
		result[element] = struct{}{}
	}

	for element := range other {
		result[element] = struct{}{}
	}

	return result
}

// Intersection returns a new set containing elements that exist in both sets.
func (s Set[T]) Intersection(other Set[T]) Set[T] {
	result := make(Set[T])

	smaller, larger := s, other
	if len(other) < len(s) {
		smaller, larger = other, s
	}

	for element := range smaller {
		if larger.Contains(element) {
			result[element] = struct{}{}
		}
	}

	return result
}

// Difference returns a new set containing elements that are in this set but not in the other.
func (s Set[T]) Difference(other Set[T]) Set[T] {
	result := make(Set[T])

	for element := range s {
		if !other.Contains(element) {
			result[element] = struct{}{}
		}
	}

	return result
}

// SymmetricDifference returns a new set containing elements that are in either set but not in both.
func (s Set[T]) SymmetricDifference(other Set[T]) Set[T] {
	result := make(Set[T])

	for element := range s {
		if !other.Contains(element) {
			result[element] = struct{}{}
		}
	}

	for element := range other {
		if !s.Contains(element) {
			result[element] = struct{}{}
		}
	}

	return result
}

// IsSubset checks if this set is a subset of the other set.
func (s Set[T]) IsSubset(other Set[T]) bool {
	if len(s) > len(other) {
		return false
	}

	for element := range s {
		if !other.Contains(element) {
			return false
		}
	}

	return true
}

// IsSuperset checks if this set is a superset of the other set.
func (s Set[T]) IsSuperset(other Set[T]) bool {
	return other.IsSubset(s)
}

// IsDisjoint checks if this set has no elements in common with the other set.
func (s Set[T]) IsDisjoint(other Set[T]) bool {
	smaller, larger := s, other
	if len(other) < len(s) {
		smaller, larger = other, s
	}

	for element := range smaller {
		if larger.Contains(element) {
			return false
		}
	}

	return true
}

// ForEach applies a function to each element in the set.
func (s Set[T]) ForEach(fn func(T)) {
	for element := range s {
		fn(element)
	}
}
