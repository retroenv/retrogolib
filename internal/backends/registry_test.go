package backends

import (
	"errors"
	"fmt"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

var errEmpty = errors.New("empty")

func names(r *Registry[int]) []string {
	entries := r.Entries()
	result := make([]string, 0, len(entries))
	for _, e := range entries {
		result = append(result, e.Name)
	}
	return result
}

func TestRegistry_Order(t *testing.T) {
	var r Registry[int]
	r.Register("low", 10, 1)
	r.Register("high", 20, 2)
	r.Register("low2", 10, 3)

	assert.Equal(t, []string{"high", "low", "low2"}, names(&r))
	assert.Equal(t, 2, r.Entries()[0].Value)
}

func TestRegistry_Unregister(t *testing.T) {
	var r Registry[int]
	unregister := r.Register("a", 1, 1)
	r.Register("b", 2, 2)

	unregister()
	unregister()
	assert.Equal(t, []string{"b"}, names(&r))

	r.Register("a", 1, 1) // The name can be used again.
	assert.Equal(t, []string{"b", "a"}, names(&r))
}

func TestRegistry_Empty(t *testing.T) {
	var r Registry[int]
	assert.Len(t, r.Entries(), 0)
}

func TestRegistry_Panics(t *testing.T) {
	var r Registry[int]
	r.Register("a", 1, 1)

	assert.Panics(t, func() { r.Register("a", 2, 2) })
	assert.Panics(t, func() { r.Register("", 2, 2) })
}

func TestRegistry_Len(t *testing.T) {
	var r Registry[int]
	assert.Equal(t, 0, r.Len())
	r.Register("a", 1, 1)
	assert.Equal(t, 1, r.Len())
}

func TestFirst_FallsBackInPriorityOrder(t *testing.T) {
	var r Registry[error]
	r.Register("old", 10, nil)
	r.Register("unused", 1, nil)
	r.Register("new", 20, errors.New("new failed"))

	var calls []string
	try := func(err error) (string, error) {
		calls = append(calls, fmt.Sprint(err))
		return "result", err
	}

	result, info, err := First(&r, errEmpty, try)
	assert.NoError(t, err)
	assert.Equal(t, "result", result)
	assert.Equal(t, []string{"new failed", "<nil>"}, calls)
	assert.Equal(t, "old", info.Name)
	assert.Error(t, info.Skipped)
	assert.Contains(t, info.Skipped.Error(), "new: new failed")
}

func TestFirst_FirstSucceeds(t *testing.T) {
	var r Registry[int]
	r.Register("only", 1, 1)

	result, info, err := First(&r, errEmpty, func(v int) (int, error) { return v, nil })
	assert.NoError(t, err)
	assert.Equal(t, 1, result)
	assert.Equal(t, "only", info.Name)
	assert.NoError(t, info.Skipped)
}

func TestFirst_AllFail(t *testing.T) {
	var r Registry[error]
	r.Register("first", 2, errors.New("first error"))
	r.Register("second", 1, errors.New("second error"))

	result, info, err := First(&r, errEmpty, func(err error) (*int, error) {
		return new(int), err
	})
	assert.Error(t, err)
	assert.True(t, result == nil, "result should be the zero value on error")
	assert.Equal(t, "", info.Name)
	assert.Contains(t, err.Error(), "first: first error")
	assert.Contains(t, err.Error(), "second: second error")
}

func TestFirst_Empty(t *testing.T) {
	var r Registry[int]
	_, _, err := First(&r, errEmpty, func(int) (int, error) { return 0, nil })
	assert.True(t, errors.Is(err, errEmpty), "error should be errEmpty")
}
