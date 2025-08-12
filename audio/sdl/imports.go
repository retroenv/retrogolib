package sdl

import (
	"fmt"
	"runtime"

	"github.com/ebitengine/purego"
)

var (
	// Core SDL functions (may be shared with GUI)
	Init     func(flags uint32) int
	GetError func() string
	Quit     func()

	// Audio-specific functions
	OpenAudioDevice    func(device string, iscapture int, desired *AudioSpec, obtained *AudioSpec, allowedChanges int) AudioDeviceID
	CloseAudioDevice   func(dev AudioDeviceID)
	PauseAudioDevice   func(dev AudioDeviceID, pauseOn int)
	QueueAudio         func(dev AudioDeviceID, data []byte, length uint32) int
	GetQueuedAudioSize func(dev AudioDeviceID) uint32
	ClearQueuedAudio   func(dev AudioDeviceID)
	LockAudioDevice    func(dev AudioDeviceID)
	UnlockAudioDevice  func(dev AudioDeviceID)
	GetAudioStatus     func() int
)

var audioImports = map[string]any{
	"SDL_Init":               &Init,
	"SDL_GetError":           &GetError,
	"SDL_Quit":               &Quit,
	"SDL_OpenAudioDevice":    &OpenAudioDevice,
	"SDL_CloseAudioDevice":   &CloseAudioDevice,
	"SDL_PauseAudioDevice":   &PauseAudioDevice,
	"SDL_QueueAudio":         &QueueAudio,
	"SDL_GetQueuedAudioSize": &GetQueuedAudioSize,
	"SDL_ClearQueuedAudio":   &ClearQueuedAudio,
	"SDL_LockAudioDevice":    &LockAudioDevice,
	"SDL_UnlockAudioDevice":  &UnlockAudioDevice,
	"SDL_GetAudioStatus":     &GetAudioStatus,
}

func registerFunction(lib uintptr, name string, ptr any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("registering function '%s': %v", name, r)
		}
	}()

	purego.RegisterLibFunc(ptr, lib, name)
	return nil
}

func getSDLSystemLibrary() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return "libSDL2.dylib", nil
	case "freebsd":
		return "libSDL2.so", nil
	case "linux":
		return "libSDL2.so", nil
	case "windows":
		return "SDL2.dll", nil
	default:
		return "", fmt.Errorf("GOOS=%s is not supported", runtime.GOOS)
	}
}
