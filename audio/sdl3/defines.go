package sdl3

const (
	SDL_INIT_AUDIO                    = 0x10
	SDL_AUDIO_DEVICE_DEFAULT_PLAYBACK = 0xffffffff
)

// AudioSpec has the layout from SDL3/SDL_audio.h.
type AudioSpec struct {
	Format   uint32
	Channels int32
	Freq     int32
}
