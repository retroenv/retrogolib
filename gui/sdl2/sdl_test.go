package sdl2

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/gui"
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

	assert.Equal(t, []string{"texture", "renderer", "window"}, calls)
}

func TestSetupFailureCleanup(t *testing.T) {
	originalWindow, originalRenderer, originalTexture := CreateWindow, CreateRenderer, CreateTexture
	originalDestroyWindow, originalDestroyRenderer := DestroyWindow, DestroyRenderer
	originalError := GetError
	t.Cleanup(func() {
		CreateWindow, CreateRenderer, CreateTexture = originalWindow, originalRenderer, originalTexture
		DestroyWindow, DestroyRenderer = originalDestroyWindow, originalDestroyRenderer
		GetError = originalError
	})
	for _, stage := range []string{"window", "renderer", "texture"} {
		t.Run(stage, func(t *testing.T) {
			message := "original error"
			var destroyed []string
			GetError = func() string { return message }
			DestroyWindow = func(uintptr) { destroyed = append(destroyed, "window"); message = "cleanup error" }
			DestroyRenderer = func(uintptr) { destroyed = append(destroyed, "renderer"); message = "cleanup error" }
			CreateWindow = func(string, int32, int32, int32, int32, uint32) uintptr {
				if stage == "window" {
					return 0
				}
				return 1
			}
			CreateRenderer = func(uintptr, int, uint32) uintptr {
				if stage == "renderer" {
					return 0
				}
				return 2
			}
			CreateTexture = func(uintptr, uint32, int, int32, int32) uintptr { return 0 }
			_, _, _, err := setupSDL(gui.Dimensions{
				Width: 16, Height: 16, ScaleFactor: 1,
			}, setupBackend{})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "original error")
			switch stage {
			case "window":
				assert.Len(t, destroyed, 0)
			case "renderer":
				assert.Equal(t, []string{"window"}, destroyed)
			case "texture":
				assert.Equal(t, []string{"renderer", "window"}, destroyed)
			}
		})
	}
}

type setupBackend struct{ gui.Backend }

func (setupBackend) WindowTitle() string { return "setup test" }
