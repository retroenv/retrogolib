package sdl

const (
	SDL_INIT_TIMER          = 0x00000001
	SDL_INIT_AUDIO          = 0x00000010
	SDL_INIT_VIDEO          = 0x00000020
	SDL_INIT_JOYSTICK       = 0x00000200
	SDL_INIT_HAPTIC         = 0x00001000
	SDL_INIT_GAMECONTROLLER = 0x00002000
	SDL_INIT_EVENTS         = 0x00004000
	SDL_INIT_SENSOR         = 0x00008000
	SDL_INIT_NOPARACHUTE    = 0x00100000
	SDL_INIT_EVERYTHING     = SDL_INIT_TIMER | SDL_INIT_AUDIO | SDL_INIT_VIDEO | SDL_INIT_EVENTS | SDL_INIT_JOYSTICK | SDL_INIT_HAPTIC | SDL_INIT_GAMECONTROLLER | SDL_INIT_SENSOR

	SDL_WINDOWPOS_CENTERED = 0x2FFF0000

	SDL_WINDOW_FULLSCREEN         = 0x00000001
	SDL_WINDOW_OPENGL             = 0x00000002
	SDL_WINDOW_SHOWN              = 0x00000004
	SDL_WINDOW_HIDDEN             = 0x00000008
	SDL_WINDOW_BORDERLESS         = 0x00000010
	SDL_WINDOW_RESIZABLE          = 0x00000020
	SDL_WINDOW_MINIMIZED          = 0x00000040
	SDL_WINDOW_MAXIMIZED          = 0x00000080
	SDL_WINDOW_MOUSE_GRABBED      = 0x00000100
	SDL_WINDOW_INPUT_FOCUS        = 0x00000200
	SDL_WINDOW_MOUSE_FOCUS        = 0x00000400
	SDL_WINDOW_FULLSCREEN_DESKTOP = SDL_WINDOW_FULLSCREEN | 0x00001000
	SDL_WINDOW_FOREIGN            = 0x00000800
	SDL_WINDOW_ALLOW_HIGHDPI      = 0x00002000

	SDL_RENDERER_SOFTWARE      = 0x00000001
	SDL_RENDERER_ACCELERATED   = 0x00000002
	SDL_RENDERER_PRESENTVSYNC  = 0x00000004
	SDL_RENDERER_TARGETTEXTURE = 0x00000008

	SDL_PIXELFORMAT_ABGR8888 = 0x16762004

	SDL_TEXTUREACCESS_STREAMING = 1
)

// events
const (
	SDL_QUIT            = 0x100
	SDL_DISPLAYEVENT    = 0x150
	SDL_WINDOWEVENT     = 0x200
	SDL_KEYDOWN         = 0x300
	SDL_KEYUP           = 0x301
	SDL_MOUSEMOTION     = 0x400
	SDL_MOUSEBUTTONDOWN = 0x401
	SDL_MOUSEBUTTONUP   = 0x402
	SDL_MOUSEWHEEL      = 0x403
	SDL_LASTEVENT       = 0x1FFF
)

type event struct {
	Type uint32
	_    [64]byte
}

type keyboardEvent struct {
	Type      uint32 // KEYDOWN, KEYUP
	Timestamp uint32 // timestamp of the event
	WindowID  uint32 // the window with keyboard focus, if any
	State     uint8  // PRESSED, RELEASED
	Repeat    uint8  // non-zero if this is a key repeat
	_         uint8  // padding
	_         uint8  // padding
	Keysym    keySym // Keysym representing the key that was pressed or released
}

type scancode uint32
type keycode int32

type keySym struct {
	Scancode scancode // SDL physical key code
	Sym      keycode  // SDL virtual key code
	Mod      uint16   // current key modifiers
	_        uint32   // unused
}
