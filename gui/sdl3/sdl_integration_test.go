//go:build sdlintegration

package sdl3

import (
	"image"
	"image/color"
	"runtime"
	"testing"
	"unsafe"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/gui"
	"github.com/retroenv/retrogolib/input"
	sdllibrary "github.com/retroenv/retrogolib/internal/sdl3"
)

func TestSDLIntegrationEvents(t *testing.T) {
	useDummyDrivers(t)
	b := &eventBackend{}
	render, cleanup, err := Setup(b)
	assert.NoError(t, err)
	defer cleanup()
	var push func(*event) bool
	assert.NoError(t, sdllibrary.LoadFunctions(map[string]any{"SDL_PushEvent": &push}))
	for _, kind := range []uint32{SDL_EVENT_KEY_DOWN, SDL_EVENT_KEY_UP} {
		var ev event
		key := (*keyboardEvent)(unsafe.Pointer(&ev))
		key.Type = kind
		key.Key = K_a
		key.Down = kind == SDL_EVENT_KEY_DOWN
		assert.True(t, push(&ev))
	}
	running, err := render()
	assert.NoError(t, err)
	assert.True(t, running)
	assert.Equal(t, []input.Key{input.A}, b.down)
	assert.Equal(t, []input.Key{input.A}, b.up)
	ev := event{Type: SDL_EVENT_QUIT}
	assert.True(t, push(&ev))
	running, err = render()
	assert.NoError(t, err)
	assert.False(t, running)
	cleanup()
}

// TestSDLIntegrationFramePixels reads back the presented frame. It checks the
// channel order and the nearest-pixel scale at the edge between two pixels.
func TestSDLIntegrationFramePixels(t *testing.T) {
	useDummyDrivers(t)
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	assert.NoError(t, setupLibrary())
	release, err := sdllibrary.Acquire(SDL_INIT_VIDEO)
	assert.NoError(t, err)
	defer release()

	b := &frameBackend{img: image.NewRGBA(image.Rect(0, 0, 2, 1))}
	b.img.SetRGBA(0, 0, color.RGBA{
		R: 0xff,
		G: 0x40,
		A: 0xff,
	})
	b.img.SetRGBA(1, 0, color.RGBA{
		G: 0x80,
		B: 0xff,
		A: 0xff,
	})
	dimensions := b.Dimensions()
	window, renderer, tex, err := setupSDL(dimensions, b)
	assert.NoError(t, err)
	defer cleanupSDL(window, renderer, tex)
	running, err := renderSDL(dimensions, b, renderer, tex)
	assert.NoError(t, err)
	assert.True(t, running)

	var readPixels func(renderer uintptr) uintptr
	var readSurfacePixel func(surface uintptr, x, y int32, r, g, b, a *uint8) bool
	var destroySurface func(surface uintptr)
	assert.NoError(t, sdllibrary.LoadFunctions(map[string]any{
		"SDL_RenderReadPixels": &readPixels,
		"SDL_ReadSurfacePixel": &readSurfacePixel,
		"SDL_DestroySurface":   &destroySurface,
	}))
	surface := readPixels(renderer)
	assert.True(t, surface != 0, GetError())
	defer destroySurface(surface)

	for _, pixel := range []struct {
		x    int32
		want color.RGBA
	}{
		{0, b.img.RGBAAt(0, 0)},
		{3, b.img.RGBAAt(0, 0)},
		{4, b.img.RGBAAt(1, 0)},
		{7, b.img.RGBAAt(1, 0)},
	} {
		var got color.RGBA
		assert.True(t, readSurfacePixel(surface, pixel.x, 0, &got.R, &got.G, &got.B, &got.A), GetError())
		assert.Equal(t, pixel.want, got, "pixel %d", pixel.x)
	}
}

type frameBackend struct {
	setupBackend
	img *image.RGBA
}

func (b *frameBackend) Image() *image.RGBA { return b.img }
func (*frameBackend) Dimensions() gui.Dimensions {
	return gui.Dimensions{
		Width:       2,
		Height:      1,
		ScaleFactor: 4,
	}
}

type eventBackend struct {
	down, up []input.Key
}

func (*eventBackend) Image() *image.RGBA { return image.NewRGBA(image.Rect(0, 0, 16, 16)) }
func (*eventBackend) Dimensions() gui.Dimensions {
	return gui.Dimensions{
		Width:       16,
		Height:      16,
		ScaleFactor: 1,
	}
}
func (*eventBackend) WindowTitle() string     { return "SDL3 event test" }
func (b *eventBackend) KeyDown(key input.Key) { b.down = append(b.down, key) }
func (b *eventBackend) KeyUp(key input.Key)   { b.up = append(b.up, key) }

// useDummyDrivers selects the dummy SDL drivers for one test. SDL3 copies the
// process environment when it first reads it, so t.Setenv can have no effect.
func useDummyDrivers(t *testing.T) {
	t.Helper()
	setHint(t, "SDL_VIDEO_DRIVER", "dummy")
	setHint(t, "SDL_RENDER_DRIVER", "software")
}

// setHint overrides an SDL hint until the test ends.
func setHint(t *testing.T, name, value string) {
	t.Helper()
	const hintOverride = 2 // SDL_HINT_OVERRIDE
	var setHintWithPriority func(name, value string, priority int32) bool
	var resetHint func(name string) bool
	assert.NoError(t, sdllibrary.LoadFunctions(map[string]any{
		"SDL_SetHintWithPriority": &setHintWithPriority,
		"SDL_ResetHint":           &resetHint,
	}))
	assert.True(t, setHintWithPriority(name, value, hintOverride))
	t.Cleanup(func() { resetHint(name) })
}
