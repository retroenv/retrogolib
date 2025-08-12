package sdl

// SDL Audio initialization flags.
const (
	SDL_INIT_AUDIO = 0x00000010
)

// SDL Audio Format flags.
const (
	AUDIO_U8     = 0x0008 // Unsigned 8-bit samples
	AUDIO_S8     = 0x8008 // Signed 8-bit samples
	AUDIO_U16LSB = 0x0010 // Unsigned 16-bit samples (little-endian)
	AUDIO_S16LSB = 0x8010 // Signed 16-bit samples (little-endian)
	AUDIO_U16MSB = 0x1010 // Unsigned 16-bit samples (big-endian)
	AUDIO_S16MSB = 0x9010 // Signed 16-bit samples (big-endian)
	AUDIO_S32LSB = 0x8020 // Signed 32-bit samples (little-endian)
	AUDIO_S32MSB = 0x9020 // Signed 32-bit samples (big-endian)
	AUDIO_F32LSB = 0x8120 // 32-bit floating point samples (little-endian)
	AUDIO_F32MSB = 0x9120 // 32-bit floating point samples (big-endian)

	// Platform-specific defaults
	AUDIO_U16SYS = AUDIO_U16LSB
	AUDIO_S16SYS = AUDIO_S16LSB
	AUDIO_S32SYS = AUDIO_S32LSB
	AUDIO_F32SYS = AUDIO_F32LSB
)

// SDL Audio status constants.
const (
	SDL_AUDIO_STOPPED = 0
	SDL_AUDIO_PLAYING = 1
	SDL_AUDIO_PAUSED  = 2
)

// AudioDeviceID represents an SDL audio device identifier.
type AudioDeviceID uint32

// AudioSpec represents SDL audio specification.
type AudioSpec struct {
	Freq     int32   // DSP frequency (samples per second)
	Format   uint16  // Audio data format
	Channels uint8   // Number of channels: 1 mono, 2 stereo
	Silence  uint8   // Audio buffer silence value
	Samples  uint16  // Audio buffer size in samples
	Padding  uint16  // Necessary for alignment
	Size     uint32  // Audio buffer size in bytes
	Callback uintptr // Callback function pointer (0 for queue-based)
	Userdata uintptr // User data pointer
}
