package set

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestBitSetBasics(t *testing.T) {
	b := NewBitSet()
	assert.True(t, b.IsEmpty())
	assert.Equal(t, 0, b.Size())

	b.Add(5)
	b.Add(10)
	assert.False(t, b.IsEmpty())
	assert.Equal(t, 2, b.Size())
	assert.True(t, b.Contains(5))
	assert.True(t, b.Contains(10))
	assert.False(t, b.Contains(0))

	b.Remove(5)
	assert.False(t, b.Contains(5))
	assert.Equal(t, 1, b.Size())

	b.Clear()
	assert.True(t, b.IsEmpty())
}

func TestBitSetBoundaries(t *testing.T) {
	b := NewBitSet()
	b.Add(0)
	b.Add(63)

	assert.Equal(t, 2, b.Size())
	assert.Equal(t, 0, b.Min())
	assert.Equal(t, 63, b.Max())
}

func TestBitSetFromSlice(t *testing.T) {
	b := NewBitSetFromSlice([]int{1, 3, 7, 10})

	assert.Equal(t, 4, b.Size())
	assert.True(t, b.Contains(1))
	assert.True(t, b.Contains(10))
	assert.False(t, b.Contains(0))
}

func TestBitSetToSlice(t *testing.T) {
	b := NewBitSetFromSlice([]int{3, 1, 7})
	slice := b.ToSlice()

	assert.Len(t, slice, 3)
	assert.Equal(t, []int{1, 3, 7}, slice)
}

func TestBitSetEqual(t *testing.T) {
	b1 := NewBitSetFromSlice([]int{1, 3})
	b2 := NewBitSetFromSlice([]int{1, 3})
	b3 := NewBitSetFromSlice([]int{1, 5})

	assert.True(t, b1.Equal(&b2))
	assert.False(t, b1.Equal(&b3))
}

func TestBitSetCopy(t *testing.T) {
	b1 := NewBitSetFromSlice([]int{1, 3})
	b2 := b1.Copy()
	b2.Add(5)

	assert.Equal(t, 2, b1.Size())
	assert.Equal(t, 3, b2.Size())
}

func TestBitSetUnion(t *testing.T) {
	b1 := NewBitSetFromSlice([]int{1, 3})
	b2 := NewBitSetFromSlice([]int{3, 5})
	result := b1.Union(&b2)

	assert.Equal(t, 3, (&result).Size())
	slice := (&result).ToSlice()
	assert.Equal(t, []int{1, 3, 5}, slice)
}

func TestBitSetIntersection(t *testing.T) {
	b1 := NewBitSetFromSlice([]int{1, 3, 5})
	b2 := NewBitSetFromSlice([]int{3, 5, 7})
	result := b1.Intersection(&b2)

	assert.Equal(t, 2, (&result).Size())
	assert.True(t, (&result).Contains(3))
	assert.True(t, (&result).Contains(5))
}

func TestBitSetDifference(t *testing.T) {
	b1 := NewBitSetFromSlice([]int{1, 3, 5})
	b2 := NewBitSetFromSlice([]int{3, 7})
	result := b1.Difference(&b2)

	assert.Equal(t, 2, (&result).Size())
	slice := (&result).ToSlice()
	assert.Equal(t, []int{1, 5}, slice)
}

func TestBitSetIsSubset(t *testing.T) {
	b1 := NewBitSetFromSlice([]int{1, 3})
	b2 := NewBitSetFromSlice([]int{1, 3, 5})

	assert.True(t, b1.IsSubset(&b2))
	assert.False(t, b2.IsSubset(&b1))
}

func TestBitSetIsDisjoint(t *testing.T) {
	b1 := NewBitSetFromSlice([]int{1, 3})
	b2 := NewBitSetFromSlice([]int{5, 7})
	b3 := NewBitSetFromSlice([]int{3, 5})

	assert.True(t, b1.IsDisjoint(&b2))
	assert.False(t, b1.IsDisjoint(&b3))
}

func TestBitSetMinMax(t *testing.T) {
	b := NewBitSet()
	assert.Equal(t, -1, b.Min())
	assert.Equal(t, -1, b.Max())

	b = NewBitSetFromSlice([]int{10, 5, 20})
	assert.Equal(t, 5, b.Min())
	assert.Equal(t, 20, b.Max())
}

func TestBitSetPopFirst(t *testing.T) {
	b := NewBitSet()
	assert.Equal(t, -1, b.PopFirst())

	b = NewBitSetFromSlice([]int{10, 5, 20})
	assert.Equal(t, 5, b.PopFirst())
	assert.Equal(t, 2, b.Size())
	assert.Equal(t, 10, b.PopFirst())
	assert.Equal(t, 20, b.PopFirst())
	assert.True(t, b.IsEmpty())
}

func TestBitSetForEach(t *testing.T) {
	b := NewBitSetFromSlice([]int{1, 3, 7})

	var visited []int
	b.ForEach(func(n int) {
		visited = append(visited, n)
	})

	assert.Equal(t, []int{1, 3, 7}, visited)
}

func TestBitSetString(t *testing.T) {
	b := NewBitSet()
	assert.Equal(t, "BitSet{}", b.String())

	b = NewBitSetFromSlice([]int{1})
	assert.Equal(t, "BitSet{1}", b.String())

	b = NewBitSetFromSlice([]int{1, 3, 5})
	assert.Equal(t, "BitSet{1, 3, 5}", b.String())
}

func TestBitSetAllValues(t *testing.T) {
	b := NewBitSet()
	for i := range 64 {
		b.Add(i)
	}

	assert.Equal(t, 64, b.Size())
	assert.Equal(t, 0, b.Min())
	assert.Equal(t, 63, b.Max())
}

func TestBitSetPanics(t *testing.T) {
	b := NewBitSet()

	assert.Panics(t, func() { b.Add(-1) })
	assert.Panics(t, func() { b.Add(64) })
	assert.Panics(t, func() { b.Remove(-1) })
	assert.Panics(t, func() { b.Remove(64) })
	assert.Panics(t, func() { b.Contains(-1) })
	assert.Panics(t, func() { b.Contains(64) })
}
