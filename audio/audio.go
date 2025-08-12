package audio

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

	// Platform-specific defaults
	FormatU16 = FormatU16LSB
	FormatS16 = FormatS16LSB
	FormatS32 = FormatS32LSB
	FormatF32 = FormatF32LSB
)

// Format contains audio format specifications.
type Format struct {
	SampleRate int          // Samples per second (e.g., 44100, 48000)
	Channels   int          // Number of channels: 1 (mono) or 2 (stereo)
	Samples    int          // Buffer size in samples (e.g., 512, 1024)
	Format     SampleFormat // Sample format (e.g., FormatS16, FormatF32)
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
		return 2 // Default to 16-bit
	}
}

// BufferSize returns the total buffer size in bytes for the format.
func (f Format) BufferSize() int {
	return f.Samples * f.Channels * f.BytesPerSample()
}

// Backend is an interface that gets implemented by the backend using the selected audio renderer.
type Backend interface {
	// AudioFormat returns the desired audio format specification.
	AudioFormat() Format

	// AudioCallback fills the buffer with audio samples.
	// Called by the audio system when more audio data is needed.
	// The buffer size will match the format's BufferSize().
	AudioCallback(buffer []byte)

	// AudioPaused returns whether audio should be paused.
	AudioPaused() bool
}

// Initializer defines a setup function for the selected audio renderer.
type Initializer func(backend Backend) (audioStart func() error, audioStop func(), err error)

// Setup will be set by the chosen and imported audio renderer.
// This function is the entrypoint for code importing this package to start audio.
var Setup Initializer
