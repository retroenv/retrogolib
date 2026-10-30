package sdl3

import (
	"sync"

	sdllibrary "github.com/retroenv/retrogolib/internal/sdl3"
)

var (
	GetError                func() string
	OpenAudioDeviceStream   func(device uint32, spec *AudioSpec, callback, userdata uintptr) uintptr
	DestroyAudioStream      func(stream uintptr)
	PauseAudioStreamDevice  func(stream uintptr) bool
	ResumeAudioStreamDevice func(stream uintptr) bool
	PutAudioStreamData      func(stream uintptr, data []byte, length int32) bool
	GetAudioStreamQueued    func(stream uintptr) int32
)

var setupLibrary = sync.OnceValue(func() error {
	return sdllibrary.LoadFunctions(map[string]any{
		"SDL_GetError":                &GetError,
		"SDL_OpenAudioDeviceStream":   &OpenAudioDeviceStream,
		"SDL_DestroyAudioStream":      &DestroyAudioStream,
		"SDL_PauseAudioStreamDevice":  &PauseAudioStreamDevice,
		"SDL_ResumeAudioStreamDevice": &ResumeAudioStreamDevice,
		"SDL_PutAudioStreamData":      &PutAudioStreamData,
		"SDL_GetAudioStreamQueued":    &GetAudioStreamQueued,
	})
})
