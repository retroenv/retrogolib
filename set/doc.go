// Package set provides a generic, type-safe set data structure implementation.
//
// The set is implemented using Go's map with struct{} values for memory efficiency
// and leverages Go generics for compile-time type safety. It supports all standard
// set operations including union, intersection, difference, and subset testing.
//
// # Basic Usage
//
//	// Create a new set
//	numbers := set.New[int]()
//	numbers.Add(1)
//	numbers.Add(2)
//	numbers.Add(3)
//
//	// Create from slice
//	colors := set.NewFromSlice([]string{"red", "green", "blue"})
//
//	// Check membership
//	if numbers.Contains(2) {
//		fmt.Println("Set contains 2")
//	}
//
// # Set Operations
//
//	s1 := set.NewFromSlice([]int{1, 2, 3})
//	s2 := set.NewFromSlice([]int{3, 4, 5})
//
//	union := s1.Union(s2)           // {1, 2, 3, 4, 5}
//	intersection := s1.Intersection(s2) // {3}
//	difference := s1.Difference(s2)     // {1, 2}
//
//	// Test relationships
//	subset := set.NewFromSlice([]int{1, 2})
//	fmt.Println(subset.IsSubset(s1)) // true
//
// # Performance
//
// All basic operations (Add, Remove, Contains) are O(1) average case.
// Set operations like Union and Intersection are O(n+m) where n and m
// are the sizes of the input sets.
//
// # Memory Efficiency
//
// The implementation uses map[T]struct{} internally, where struct{} has
// zero size, making the memory overhead minimal compared to alternatives
// like map[T]bool.
//
// # Thread Safety
//
// Sets are not thread-safe. If concurrent access is required, external
// synchronization must be provided by the caller.
//
// # Type Constraints
//
// The set can hold any comparable type, which includes:
//   - Basic types: int, string, bool, float64, etc.
//   - Arrays and structs containing only comparable types
//   - Pointers and channels
//   - Interface types with comparable dynamic values
//
// Note: Slices, maps, and functions are not comparable and cannot be used
// as set elements.
package set
