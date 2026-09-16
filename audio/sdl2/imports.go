package sdl2

import (
	"sync"

	sdllibrary "github.com/retroenv/retrogolib/internal/sdl2"
)

var (
	GetError           func() string
	OpenAudioDevice    func(device *byte, iscapture int, desired, obtained *AudioSpec, allowedChanges int) AudioDeviceID
	CloseAudioDevice   func(dev AudioDeviceID)
	PauseAudioDevice   func(dev AudioDeviceID, pauseOn int)
	QueueAudio         func(dev AudioDeviceID, data []byte, length uint32) int
	GetQueuedAudioSize func(dev AudioDeviceID) uint32
)

var setupLibrary = sync.OnceValue(func() error {
	return sdllibrary.LoadFunctions(map[string]any{
		"SDL_GetError":           &GetError,
		"SDL_OpenAudioDevice":    &OpenAudioDevice,
		"SDL_CloseAudioDevice":   &CloseAudioDevice,
		"SDL_PauseAudioDevice":   &PauseAudioDevice,
		"SDL_QueueAudio":         &QueueAudio,
		"SDL_GetQueuedAudioSize": &GetQueuedAudioSize,
	})
})
