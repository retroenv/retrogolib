package sdl2

// SDL_INIT_AUDIO selects the audio subsystem.
const SDL_INIT_AUDIO = 0x00000010

// AudioDeviceID represents an SDL2 audio device identifier.
type AudioDeviceID uint32

// AudioSpec mirrors the native SDL2 audio specification layout.
type AudioSpec struct {
	Freq     int32
	Format   uint16
	Channels uint8
	Silence  uint8
	Samples  uint16
	Padding  uint16
	Size     uint32
	Callback uintptr
	Userdata uintptr
}
