package sdl3

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
	defer func() {
		DestroyTexture = originalDestroyTexture
		DestroyRenderer = originalDestroyRenderer
		DestroyWindow = originalDestroyWindow
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

	cleanupSDL(1, 2, 3)

	assert.Equal(t, []string{"texture", "renderer", "window"}, calls)
}

func TestSetupFailureCleanup(t *testing.T) {
	originalWindow := CreateWindow
	originalRenderer := CreateRenderer
	originalTexture := CreateTexture
	originalScaleMode := SetTextureScaleMode
	originalDestroyWindow := DestroyWindow
	originalDestroyRenderer := DestroyRenderer
	originalDestroyTexture := DestroyTexture
	originalError := GetError
	t.Cleanup(func() {
		CreateWindow = originalWindow
		CreateRenderer = originalRenderer
		CreateTexture = originalTexture
		SetTextureScaleMode = originalScaleMode
		DestroyWindow = originalDestroyWindow
		DestroyRenderer = originalDestroyRenderer
		DestroyTexture = originalDestroyTexture
		GetError = originalError
	})
	for _, tc := range []struct {
		stage     string
		destroyed []string
	}{
		{"window", nil},
		{"renderer", []string{"window"}},
		{"texture", []string{"renderer", "window"}},
		{"scale mode", []string{"texture", "renderer", "window"}},
	} {
		t.Run(tc.stage, func(t *testing.T) {
			message := "original error"
			var destroyed []string
			destroy := func(name string) func(uintptr) {
				return func(uintptr) { destroyed = append(destroyed, name); message = "cleanup error" }
			}
			create := func(name string, handle uintptr) uintptr {
				if tc.stage == name {
					return 0
				}
				return handle
			}
			GetError = func() string { return message }
			DestroyWindow = destroy("window")
			DestroyRenderer = destroy("renderer")
			DestroyTexture = destroy("texture")
			CreateWindow = func(string, int32, int32, uint64) uintptr { return create("window", 1) }
			CreateRenderer = func(uintptr, *byte) uintptr { return create("renderer", 2) }
			CreateTexture = func(uintptr, uint32, int, int32, int32) uintptr { return create("texture", 3) }
			SetTextureScaleMode = func(_ uintptr, scaleMode int32) bool {
				assert.Equal(t, int32(SDL_SCALEMODE_NEAREST), scaleMode)
				return false
			}
			_, _, _, err := setupSDL(gui.Dimensions{
				Width:       16,
				Height:      16,
				ScaleFactor: 1,
			}, setupBackend{})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "original error")
			assert.Equal(t, tc.destroyed, destroyed)
		})
	}
}

type setupBackend struct{ gui.Backend }

func (setupBackend) WindowTitle() string { return "setup test" }
