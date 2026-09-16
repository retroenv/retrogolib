//go:build sdlintegration

package sdl2

import (
	"image"
	"sync"
	"testing"
	"time"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/audio"
	"github.com/retroenv/retrogolib/gui"
	guirenderer "github.com/retroenv/retrogolib/gui/sdl2"
	"github.com/retroenv/retrogolib/input"
	sdllibrary "github.com/retroenv/retrogolib/internal/sdl2"
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
	t.Setenv("SDL_AUDIODRIVER", "dummy")
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
	t.Setenv("SDL_AUDIODRIVER", "dummy")
	for range 3 {
		playback, err := Setup(newTestBackend())
		assert.NoError(t, err)
		playback.Stop()
		assert.ErrorIs(t, playback.Start(), audio.ErrClosed)
	}
}

func TestSDLIntegrationFormats(t *testing.T) {
	t.Setenv("SDL_AUDIODRIVER", "dummy")
	for _, format := range []audio.SampleFormat{
		audio.FormatU8, audio.FormatS8, audio.FormatU16LSB, audio.FormatU16MSB,
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
	t.Setenv("SDL_AUDIODRIVER", "dummy")
	assert.NoError(t, setupLibrary())
	original := OpenAudioDevice
	t.Cleanup(func() { OpenAudioDevice = original })
	OpenAudioDevice = func(*byte, int, *AudioSpec, *AudioSpec, int) AudioDeviceID { return 0 }
	playback, err := Setup(newTestBackend())
	assert.Error(t, err)
	assert.Nil(t, playback)
	var wasInit func(uint32) uint32
	assert.NoError(t, sdllibrary.LoadFunctions(map[string]any{"SDL_WasInit": &wasInit}))
	assert.Equal(t, uint32(0), wasInit(SDL_INIT_AUDIO))
	OpenAudioDevice = original
	playback, err = Setup(newTestBackend())
	assert.NoError(t, err)
	playback.Stop()
}

func TestSDLIntegrationSubsystemFailure(t *testing.T) {
	t.Setenv("SDL_AUDIODRIVER", "retrogolib-missing-driver")
	playback, err := Setup(newTestBackend())
	assert.Error(t, err)
	assert.Nil(t, playback)
	t.Setenv("SDL_AUDIODRIVER", "dummy")
	playback, err = Setup(newTestBackend())
	assert.NoError(t, err)
	playback.Stop()
}

func TestSDLIntegrationGUIAudioIsolation(t *testing.T) {
	t.Setenv("SDL_AUDIODRIVER", "dummy")
	t.Setenv("SDL_VIDEODRIVER", "dummy")
	t.Setenv("SDL_RENDER_DRIVER", "software")
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
	t.Setenv("SDL_AUDIODRIVER", "dummy")
	t.Setenv("SDL_VIDEODRIVER", "dummy")
	t.Setenv("SDL_RENDER_DRIVER", "software")
	_, cleanup, err := guirenderer.Setup(&integrationVideo{})
	assert.NoError(t, err)
	cleanup()
	original := guirenderer.CreateRenderer
	t.Cleanup(func() { guirenderer.CreateRenderer = original })
	guirenderer.CreateRenderer = func(uintptr, int, uint32) uintptr { return 0 }
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
		Width: 16, Height: 16, ScaleFactor: 1,
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
