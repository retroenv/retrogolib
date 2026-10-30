package sdl3

// Values and layouts follow SDL3/SDL_init.h, SDL_video.h, and SDL_events.h.
const (
	SDL_INIT_VIDEO                = 0x20
	SDL_WINDOW_HIGH_PIXEL_DENSITY = 0x2000
	SDL_PIXELFORMAT_ABGR8888      = 0x16762004
	SDL_TEXTUREACCESS_STREAMING   = 1
	SDL_SCALEMODE_NEAREST         = 0
	SDL_EVENT_QUIT                = 0x100
	SDL_EVENT_KEY_DOWN            = 0x300
	SDL_EVENT_KEY_UP              = 0x301
)

type event struct {
	Type uint32
	_    uint32
	_    [15]uint64
}

type keycode uint32

type keyboardEvent struct {
	Type      uint32
	Reserved  uint32
	Timestamp uint64
	WindowID  uint32
	Which     uint32
	Scancode  uint32
	Key       keycode
	Mod       uint16
	Raw       uint16
	Down      bool
	Repeat    bool
}
