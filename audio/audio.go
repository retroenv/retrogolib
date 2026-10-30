package audio

import (
	"errors"
	"fmt"
	"math"

	"github.com/retroenv/retrogolib/internal/backends"
	"github.com/retroenv/retrogolib/renderer"
)

// SampleFormat represents different audio sample formats.
type SampleFormat uint16

// Common audio sample formats.
const (
	FormatU8     SampleFormat = 0x0008 // Unsigned 8-bit samples
	FormatS8     SampleFormat = 0x8008 // Signed 8-bit samples
	FormatU16LSB SampleFormat = 0x0010 // Unsigned 16-bit samples (little-endian)
	FormatS16LSB SampleFormat = 0x8010 // Signed 16-bit samples (little-endian)
	FormatU16MSB SampleFormat = 0x1010 // Unsigned 16-bit samples (big-endian)
	FormatS16MSB SampleFormat = 0x9010 // Signed 16-bit samples (big-endian)
	FormatS32LSB SampleFormat = 0x8020 // Signed 32-bit samples (little-endian)
	FormatS32MSB SampleFormat = 0x9020 // Signed 32-bit samples (big-endian)
	FormatF32LSB SampleFormat = 0x8120 // 32-bit floating point samples (little-endian)
	FormatF32MSB SampleFormat = 0x9120 // 32-bit floating point samples (big-endian)

	// Default formats use little-endian byte order on every platform.
	FormatU16 = FormatU16LSB
	FormatS16 = FormatS16LSB
	FormatS32 = FormatS32LSB
	FormatF32 = FormatF32LSB
)

// Format contains audio format specifications.
type Format struct {
	SampleRate int          // Sample frames per second (e.g., 44100, 48000)
	Channels   int          // Interleaved channels, from 1 (mono) to 8
	Samples    int          // Sample frames per callback, each containing Channels samples
	Format     SampleFormat // Sample format (e.g., FormatS16, FormatF32)
}

// Backend is an interface that gets implemented by the backend using the selected audio renderer.
type Backend interface {
	// AudioFormat returns the desired audio format specification.
	AudioFormat() Format

	// AudioCallback fills the buffer with audio samples.
	// Calls are serial, on a playback worker, and may overlap application code.
	// Fill the entire buffer, including format-appropriate silence on underrun.
	// The buffer has Format.BufferSize() bytes and must not be retained.
	// Callbacks must return promptly and must not call playback Start or Stop.
	AudioCallback(buffer []byte)

	// AudioPaused returns whether audio should be paused, preserving queued data.
	// It runs on the playback worker; synchronize access to shared application state.
	AudioPaused() bool
}

// Initializer sets up a renderer without starting playback.
type Initializer func(backend Backend) (*Playback, error)

// Playback controls one audio device. Its controls are safe for concurrent calls.
// Do not modify these fields. Always call Stop, even if Start was never called.
type Playback struct {
	// Start prefills the queue and starts playback. Repeated calls are harmless.
	// It returns ErrClosed after Stop, or the terminal playback error after failure.
	Start func() error
	// Stop closes the device and waits for all callbacks to finish. It is idempotent.
	Stop func()
	// Errors receives at most one terminal error, then closes after cleanup.
	// Normal Stop closes it without an error. Playback never blocks on a reader.
	Errors <-chan error
}

// BytesPerSample returns the number of bytes per sample for the format.
func (f Format) BytesPerSample() int {
	switch f.Format {
	case FormatU8, FormatS8:
		return 1
	case FormatU16LSB, FormatS16LSB, FormatU16MSB, FormatS16MSB:
		return 2
	case FormatS32LSB, FormatS32MSB, FormatF32LSB, FormatF32MSB:
		return 4
	default:
		return 0
	}
}

// BufferSize returns the callback size in bytes, or zero for an invalid format.
func (f Format) BufferSize() int {
	if f.Validate() != nil {
		return 0
	}
	return f.Samples * f.Channels * f.BytesPerSample()
}

// Validate checks the PCM format and ensures its buffer size cannot overflow.
// Renderers may impose additional hardware-specific limits.
func (f Format) Validate() error {
	if f.SampleRate <= 0 || f.SampleRate > math.MaxInt32 {
		return fmt.Errorf("sample rate must be between 1 and %d, got %d", math.MaxInt32, f.SampleRate)
	}
	if f.Channels < 1 || f.Channels > 8 {
		return fmt.Errorf("channels must be between 1 and 8, got %d", f.Channels)
	}
	if f.Samples <= 0 {
		return fmt.Errorf("sample frames must be positive, got %d", f.Samples)
	}
	bytesPerSample := f.BytesPerSample()
	if bytesPerSample == 0 {
		return fmt.Errorf("unsupported sample format: %#x", f.Format)
	}
	maxBytes := min(uint64(^uint(0)>>1), uint64(math.MaxUint32))
	if uint64(f.Samples) > maxBytes/uint64(f.Channels*bytesPerSample) {
		return errors.New("audio buffer size overflows")
	}
	return nil
}

// ErrClosed is returned when starting a playback that has been stopped.
var ErrClosed = errors.New("audio playback is closed")

// ErrNoRenderer is returned by Setup when no audio renderer is registered.
var ErrNoRenderer = errors.New("no audio renderer registered")

var renderers backends.Registry[renderer.Renderer[Initializer]]

// Register adds an audio renderer for Setup and returns a function that removes it.
// Renderer packages call Register in their init function.
// Setup tries renderers with a higher priority first. Renderers with an equal
// priority run in registration order. Register panics when the name is empty,
// the name is already registered, or the renderer has no Setup function.
func Register(name string, priority int, r renderer.Renderer[Initializer]) (unregister func()) {
	if r.Setup == nil {
		panic("audio: renderer setup function is nil")
	}
	return renderers.Register(name, priority, r)
}

// Available reports whether at least one audio renderer is registered.
func Available() bool {
	return renderers.Len() > 0
}

// Setup opens the first registered audio renderer that initializes without error.
// It is SetupWithInfo without the renderer information.
func Setup(backend Backend) (*Playback, error) {
	playback, _, err := SetupWithInfo(backend)
	return playback, err
}

// SetupWithInfo opens the first registered audio renderer that initializes
// without error and returns information about it. It tries the renderers in
// priority order, for example a renderer for a newer library version before a
// renderer for an older one. This is the entrypoint for code importing this
// package to start audio.
//
// When all renderers fail, the returned error contains the error of every
// renderer. It returns ErrNoRenderer when no renderer is registered.
func SetupWithInfo(backend Backend) (*Playback, renderer.Info, error) {
	playback, info, err := backends.First(&renderers, ErrNoRenderer, func(r renderer.Renderer[Initializer]) (*Playback, error) {
		return r.Setup(backend)
	})
	if err != nil {
		return nil, renderer.Info{}, fmt.Errorf("opening audio renderer: %w", err)
	}

	return playback, info, nil
}
