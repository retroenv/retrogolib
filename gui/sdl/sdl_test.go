package sdl

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/input"
)

func TestKeyMapping(t *testing.T) {
	assert.Equal(t, input.Apostrophe, keyMapping[K_QUOTE])
	assert.Equal(t, input.Period, keyMapping[K_PERIOD])
	assert.Equal(t, input.CapsLock, keyMapping[K_CAPSLOCK])
	assert.Equal(t, input.NumLock, keyMapping[K_NUMLOCKCLEAR])
	assert.Equal(t, input.KPEqual, keyMapping[K_KP_EQUALS])
	assert.Equal(t, input.Menu, keyMapping[K_APPLICATION])
}

func TestCleanupSDL(t *testing.T) {
	var calls []string

	originalDestroyTexture := DestroyTexture
	originalDestroyRenderer := DestroyRenderer
	originalDestroyWindow := DestroyWindow
	originalQuit := Quit
	defer func() {
		DestroyTexture = originalDestroyTexture
		DestroyRenderer = originalDestroyRenderer
		DestroyWindow = originalDestroyWindow
		Quit = originalQuit
	}()

	DestroyTexture = func(uintptr) {
		calls = append(calls, "texture")
	}
	DestroyRenderer = func(uintptr) {
		calls = append(calls, "renderer")
	}
	DestroyWindow = func(uintptr) {
		calls = append(calls, "window")
	}
	Quit = func() {
		calls = append(calls, "quit")
	}

	cleanupSDL(1, 2, 3)

	assert.Equal(t, []string{"texture", "renderer", "window", "quit"}, calls)
}
