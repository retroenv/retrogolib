package set

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestNew(t *testing.T) {
	s := New[int]()
	assert.True(t, s.IsEmpty(), "New set should be empty")
	assert.Equal(t, 0, s.Size(), "New set should have size 0")
}

func TestNewFromSlice(t *testing.T) {
	elements := []int{1, 2, 3, 2, 1}
	s := NewFromSlice(elements)

	assert.Equal(t, 3, s.Size(), "Expected size 3")
	assert.True(t, s.Contains(1), "Set should contain 1")
	assert.True(t, s.Contains(2), "Set should contain 2")
	assert.True(t, s.Contains(3), "Set should contain 3")
}

func TestAdd(t *testing.T) {
	s := New[string]()
	s.Add("hello")
	s.Add("world")
	s.Add("hello") // duplicate

	assert.Equal(t, 2, s.Size(), "Expected size 2")
	assert.True(t, s.Contains("hello"), "Set should contain hello")
	assert.True(t, s.Contains("world"), "Set should contain world")
}

func TestRemove(t *testing.T) {
	s := NewFromSlice([]int{1, 2, 3})
	s.Remove(2)

	assert.Equal(t, 2, s.Size(), "Expected size 2")
	assert.False(t, s.Contains(2), "Set should not contain removed element")
	assert.True(t, s.Contains(1), "Set should still contain 1")
	assert.True(t, s.Contains(3), "Set should still contain 3")
}

func TestContains(t *testing.T) {
	s := NewFromSlice([]int{1, 2, 3})

	assert.True(t, s.Contains(1), "Set should contain 1")
	assert.True(t, s.Contains(2), "Set should contain 2")
	assert.True(t, s.Contains(3), "Set should contain 3")
	assert.False(t, s.Contains(4), "Set should not contain 4")
}

func TestClear(t *testing.T) {
	s := NewFromSlice([]int{1, 2, 3})
	s.Clear()

	assert.True(t, s.IsEmpty(), "Set should be empty after clear")
	assert.Equal(t, 0, s.Size(), "Expected size 0")
}

func TestToSlice(t *testing.T) {
	original := []int{1, 2, 3}
	s := NewFromSlice(original)
	result := s.ToSlice()

	assert.Len(t, result, 3, "Expected slice length 3")

	resultSet := NewFromSlice(result)
	assert.True(t, s.Equal(resultSet), "ToSlice should preserve all elements")
}

func TestCopy(t *testing.T) {
	s1 := NewFromSlice([]int{1, 2, 3})
	s2 := s1.Copy()

	assert.True(t, s1.Equal(s2), "Copy should create equal set")

	s2.Add(4)
	assert.False(t, s1.Contains(4), "Modifying copy should not affect original")
	assert.False(t, s1.Equal(s2), "Modified copy should not be equal to original")
}

func TestEqual(t *testing.T) {
	s1 := NewFromSlice([]int{1, 2, 3})
	s2 := NewFromSlice([]int{3, 2, 1})
	s3 := NewFromSlice([]int{1, 2})
	s4 := NewFromSlice([]int{1, 2, 4})

	assert.True(t, s1.Equal(s2), "Sets with same elements should be equal")
	assert.False(t, s1.Equal(s3), "Sets with different sizes should not be equal")
	assert.False(t, s1.Equal(s4), "Sets with different elements should not be equal")
}

func TestUnion(t *testing.T) {
	s1 := NewFromSlice([]int{1, 2, 3})
	s2 := NewFromSlice([]int{3, 4, 5})
	result := s1.Union(s2)

	expected := NewFromSlice([]int{1, 2, 3, 4, 5})
	assert.True(t, result.Equal(expected), "Union should contain all elements from both sets")
}

func TestIntersection(t *testing.T) {
	s1 := NewFromSlice([]int{1, 2, 3, 4})
	s2 := NewFromSlice([]int{3, 4, 5, 6})
	result := s1.Intersection(s2)

	expected := NewFromSlice([]int{3, 4})
	assert.True(t, result.Equal(expected), "Intersection should contain only common elements")
}

func TestDifference(t *testing.T) {
	s1 := NewFromSlice([]int{1, 2, 3, 4})
	s2 := NewFromSlice([]int{3, 4, 5, 6})
	result := s1.Difference(s2)

	expected := NewFromSlice([]int{1, 2})
	assert.True(t, result.Equal(expected), "Difference should contain elements only in first set")
}

func TestSymmetricDifference(t *testing.T) {
	s1 := NewFromSlice([]int{1, 2, 3})
	s2 := NewFromSlice([]int{3, 4, 5})
	result := s1.SymmetricDifference(s2)

	expected := NewFromSlice([]int{1, 2, 4, 5})
	assert.True(t, result.Equal(expected), "Symmetric difference should contain elements in either set but not both")
}

func TestIsSubset(t *testing.T) {
	s1 := NewFromSlice([]int{1, 2})
	s2 := NewFromSlice([]int{1, 2, 3, 4})
	s3 := NewFromSlice([]int{1, 5})

	assert.True(t, s1.IsSubset(s2), "s1 should be subset of s2")
	assert.False(t, s1.IsSubset(s3), "s1 should not be subset of s3")
	assert.True(t, s1.IsSubset(s1), "Set should be subset of itself")
}

func TestIsSuperset(t *testing.T) {
	s1 := NewFromSlice([]int{1, 2, 3, 4})
	s2 := NewFromSlice([]int{1, 2})
	s3 := NewFromSlice([]int{1, 5})

	assert.True(t, s1.IsSuperset(s2), "s1 should be superset of s2")
	assert.False(t, s1.IsSuperset(s3), "s1 should not be superset of s3")
	assert.True(t, s1.IsSuperset(s1), "Set should be superset of itself")
}

func TestIsDisjoint(t *testing.T) {
	s1 := NewFromSlice([]int{1, 2, 3})
	s2 := NewFromSlice([]int{4, 5, 6})
	s3 := NewFromSlice([]int{3, 4, 5})

	assert.True(t, s1.IsDisjoint(s2), "s1 and s2 should be disjoint")
	assert.False(t, s1.IsDisjoint(s3), "s1 and s3 should not be disjoint")
}

func TestStringSet(t *testing.T) {
	s := New[string]()
	s.Add("hello")
	s.Add("world")

	assert.True(t, s.Contains("hello"), "String set should contain hello")
	assert.True(t, s.Contains("world"), "String set should contain world")
	assert.Equal(t, 2, s.Size(), "Expected size 2")
}

func TestEmptySetOperations(t *testing.T) {
	empty := New[int]()
	s := NewFromSlice([]int{1, 2, 3})

	assert.True(t, empty.Union(s).Equal(s), "Union with empty set should equal the other set")
	assert.True(t, empty.Intersection(s).IsEmpty(), "Intersection with empty set should be empty")
	assert.True(t, empty.Difference(s).IsEmpty(), "Difference of empty set should be empty")
	assert.True(t, empty.IsSubset(s), "Empty set should be subset of any set")
	assert.True(t, s.IsSuperset(empty), "Any set should be superset of empty set")
}
