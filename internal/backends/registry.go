package backends

import (
	"errors"
	"fmt"
	"slices"
	"sync"

	"github.com/retroenv/retrogolib/renderer"
)

// Priorities of the SDL renderers. A newer SDL major version runs before an older one.
// The audio and GUI renderers use the same values.
const (
	PrioritySDL2 = 10
	PrioritySDL3 = 20
)

// Entry is a registered value with the name of its renderer.
type Entry[T any] struct {
	Name  string
	Value T
}

type entry[T any] struct {
	Entry[T]

	priority int
	sequence uint64
}

// Registry stores named values. The zero value is ready to use.
// All methods are safe for concurrent use.
type Registry[T any] struct {
	mutex    sync.RWMutex
	entries  []entry[T]
	sequence uint64
}

// Register adds a value and returns a function that removes it again.
// A higher priority runs first. Equal priorities run in registration order.
// It panics when the name is empty or already registered.
func (r *Registry[T]) Register(name string, priority int, value T) (unregister func()) {
	if name == "" {
		panic("backends: renderer name is empty")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, e := range r.entries {
		if e.Name == name {
			panic(fmt.Sprintf("backends: renderer %q is already registered", name))
		}
	}

	r.sequence++
	added := entry[T]{
		Entry: Entry[T]{
			Name:  name,
			Value: value,
		},
		priority: priority,
		sequence: r.sequence,
	}

	// Insert before the first entry with a lower priority. This keeps the
	// registration order for equal priorities.
	index := slices.IndexFunc(r.entries, func(e entry[T]) bool {
		return e.priority < priority
	})
	if index < 0 {
		index = len(r.entries)
	}

	r.entries = slices.Insert(r.entries, index, added)

	var once sync.Once

	return func() {
		once.Do(func() {
			r.mutex.Lock()
			defer r.mutex.Unlock()
			r.entries = slices.DeleteFunc(r.entries, func(e entry[T]) bool {
				return e.sequence == added.sequence
			})
		})
	}
}

// Len returns the number of registered values.
func (r *Registry[T]) Len() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.entries)
}

// Entries returns a snapshot of the registered values in run order.
func (r *Registry[T]) Entries() []Entry[T] {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]Entry[T], len(r.entries))
	for i, e := range r.entries {
		result[i] = e.Entry
	}

	return result
}

// LibraryVersion returns the library version of a renderer. It returns an
// empty string when the renderer has no version function.
func LibraryVersion[S any](r renderer.Renderer[S]) string {
	if r.Version == nil {
		return ""
	}

	return r.Version()
}

// First calls try for each registered value in run order and returns the
// result of the first call that succeeds. When all calls fail, the returned
// error contains the error of every call, prefixed with the entry name.
// It returns errEmpty when no value is registered. libraryVersion returns the
// library version of a value and may be nil.
func First[T, R any](r *Registry[T], errEmpty error, try func(T) (R, error),
	libraryVersion func(T) string) (R, renderer.Info, error) {

	var zero R

	entries := r.Entries()
	if len(entries) == 0 {
		return zero, renderer.Info{}, errEmpty
	}

	var errs []error
	for _, e := range entries {
		result, err := try(e.Value)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", e.Name, err))
			continue
		}

		info := renderer.Info{
			Name:    e.Name,
			Skipped: errors.Join(errs...),
		}
		if libraryVersion != nil {
			info.Version = libraryVersion(e.Value)
		}

		return result, info, nil
	}

	return zero, renderer.Info{}, errors.Join(errs...)
}
