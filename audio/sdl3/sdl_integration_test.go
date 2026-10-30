//go:build sdlintegration

package sdl3

import (
	"image"
	"sync"
	"testing"
	"time"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/audio"
	"github.com/retroenv/retrogolib/gui"
	guirenderer "github.com/retroenv/retrogolib/gui/sdl3"
	"github.com/retroenv/retrogolib/input"
	sdllibrary "github.com/retroenv/retrogolib/internal/sdl3"
)

func TestSDLIntegrationConcurrentBindings(t *testing.T) {
	var workers sync.WaitGroup
	results := make(chan error, 20)
	for range 20 {
		workers.Go(func() { results <- setupLibrary() })
	}
	workers.Wait()
	close(results)
	for err := range results {
		assert.NoError(t, err)
	}
}

func TestSDLIntegrationPlayback(t *testing.T) {
	useDummyDrivers(t)
	b := newTestBackend()
	b.format.SampleRate = 48000
	b.format.Samples = 256
	playback, err := Setup(b)
	assert.NoError(t, err)
	t.Cleanup(playback.Stop)
	assert.NoError(t, playback.Start())
	awaitCallbacks(t, playback, b, b.calls.Load())
	playback.Stop()
	for err := range playback.Errors {
		assert.NoError(t, err)
	}
}

func TestSDLIntegrationStopBeforeStart(t *testing.T) {
	useDummyDrivers(t)
	for range 3 {
		playback, err := Setup(newTestBackend())
		assert.NoError(t, err)
		playback.Stop()
		assert.ErrorIs(t, playback.Start(), audio.ErrClosed)
	}
}

func TestSDLIntegrationFormats(t *testing.T) {
	useDummyDrivers(t)
	for _, format := range []audio.SampleFormat{
		audio.FormatU8, audio.FormatS8,
		audio.FormatS16LSB, audio.FormatS16MSB, audio.FormatS32LSB, audio.FormatS32MSB,
		audio.FormatF32LSB, audio.FormatF32MSB,
	} {
		b := newTestBackend()
		b.format.Format = format
		playback, err := Setup(b)
		assert.NoError(t, err)
		err = playback.Start()
		playback.Stop()
		assert.NoError(t, err)
		for err := range playback.Errors {
			assert.NoError(t, err)
		}
	}
}

func TestSDLIntegrationOpenFailureCleanup(t *testing.T) {
	useDummyDrivers(t)
	assert.NoError(t, setupLibrary())
	original := OpenAudioDeviceStream
	t.Cleanup(func() { OpenAudioDeviceStream = original })
	OpenAudioDeviceStream = func(uint32, *AudioSpec, uintptr, uintptr) uintptr { return 0 }
	playback, err := Setup(newTestBackend())
	assert.Error(t, err)
	assert.Nil(t, playback)
	var wasInit func(uint32) uint32
	assert.NoError(t, sdllibrary.LoadFunctions(map[string]any{"SDL_WasInit": &wasInit}))
	assert.Equal(t, uint32(0), wasInit(SDL_INIT_AUDIO))
	OpenAudioDeviceStream = original
	playback, err = Setup(newTestBackend())
	assert.NoError(t, err)
	playback.Stop()
}

func TestSDLIntegrationSubsystemFailure(t *testing.T) {
	useDummyDrivers(t)
	setHint(t, "SDL_AUDIO_DRIVER", "retrogolib-missing-driver")
	playback, err := Setup(newTestBackend())
	assert.Error(t, err)
	assert.Nil(t, playback)
	setHint(t, "SDL_AUDIO_DRIVER", "dummy")
	playback, err = Setup(newTestBackend())
	assert.NoError(t, err)
	playback.Stop()
}

func TestSDLIntegrationGUIAudioIsolation(t *testing.T) {
	useDummyDrivers(t)
	b := newTestBackend()
	playback, err := Setup(b)
	assert.NoError(t, err)
	t.Cleanup(playback.Stop)
	assert.NoError(t, playback.Start())
	render, cleanup, err := guirenderer.Setup(&integrationVideo{})
	assert.NoError(t, err)
	defer cleanup()
	_, err = render()
	assert.NoError(t, err)
	cleanup()
	awaitCallbacks(t, playback, b, b.calls.Load())
	playback.Stop()
	for err := range playback.Errors {
		assert.NoError(t, err)
	}

	// In the opposite order, audio shutdown must leave the GUI renderer usable.
	render, cleanup, err = guirenderer.Setup(&integrationVideo{})
	assert.NoError(t, err)
	defer cleanup()
	playback, err = Setup(newTestBackend())
	assert.NoError(t, err)
	assert.NoError(t, playback.Start())
	playback.Stop()
	_, err = render()
	assert.NoError(t, err)
}

func TestSDLIntegrationGUIFailureKeepsAudio(t *testing.T) {
	useDummyDrivers(t)
	_, cleanup, err := guirenderer.Setup(&integrationVideo{})
	assert.NoError(t, err)
	cleanup()
	original := guirenderer.CreateRenderer
	t.Cleanup(func() { guirenderer.CreateRenderer = original })
	guirenderer.CreateRenderer = func(uintptr, *byte) uintptr { return 0 }
	b := newTestBackend()
	playback, err := Setup(b)
	assert.NoError(t, err)
	t.Cleanup(playback.Stop)
	assert.NoError(t, playback.Start())
	_, _, err = guirenderer.Setup(&integrationVideo{})
	assert.Error(t, err)
	var wasInit func(uint32) uint32
	assert.NoError(t, sdllibrary.LoadFunctions(map[string]any{"SDL_WasInit": &wasInit}))
	assert.Equal(t, uint32(0), wasInit(guirenderer.SDL_INIT_VIDEO))
	assert.Equal(t, uint32(SDL_INIT_AUDIO), wasInit(SDL_INIT_AUDIO))
	awaitCallbacks(t, playback, b, b.calls.Load())
}

type integrationVideo struct{}

func (*integrationVideo) Image() *image.RGBA { return image.NewRGBA(image.Rect(0, 0, 16, 16)) }
func (*integrationVideo) Dimensions() gui.Dimensions {
	return gui.Dimensions{
		Width:       16,
		Height:      16,
		ScaleFactor: 1,
	}
}
func (*integrationVideo) WindowTitle() string { return "audio coexistence test" }
func (*integrationVideo) KeyDown(input.Key)   {}
func (*integrationVideo) KeyUp(input.Key)     {}

func awaitCallbacks(t *testing.T, playback *audio.Playback, b *testBackend, previous int64) {
	t.Helper()
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for b.calls.Load() <= previous {
		select {
		case err := <-playback.Errors:
			t.Fatalf("playback stopped unexpectedly: %v", err)
		case <-timer.C:
			t.Fatal("audio device did not consume queued frames")
		case <-ticker.C:
		}
	}
}

// useDummyDrivers selects the dummy SDL drivers for one test. SDL3 copies the
// process environment when it first reads it, so t.Setenv can have no effect.
func useDummyDrivers(t *testing.T) {
	t.Helper()
	setHint(t, "SDL_AUDIO_DRIVER", "dummy")
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
