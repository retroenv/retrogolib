// Package sdl2 coordinates SDL2 library loading and subsystem ownership.
package sdl2

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/retroenv/retrogolib/internal/dynlib"
)

// LoadFunctions registers a backend's bindings. Each backend must call it once.
func LoadFunctions(imports map[string]any) error {
	var names []string
	switch runtime.GOOS {
	case "darwin":
		names = []string{"libSDL2.dylib"}
	case "linux":
		names = []string{"libSDL2-2.0.so.0", "libSDL2.so"}
	case "freebsd":
		names = []string{"libSDL2.so"}
	case "windows":
		names = []string{"SDL2.dll"}
	default:
		return fmt.Errorf("SDL2 is not supported on GOOS=%s", runtime.GOOS)
	}
	var lastErr error
	for _, name := range names {
		if _, err := dynlib.LoadFunctions(name, imports); err != nil {
			lastErr = err
			continue
		}
		return nil
	}
	return fmt.Errorf("loading SDL2: %w", lastErr)
}

// Acquire initializes the requested subsystems and returns an idempotent release.
// Video callers must keep their goroutine on the main OS thread through release.
func Acquire(flags uint32) (func(), error) {
	if err := loadCore(); err != nil {
		return nil, err
	}
	return lifecycle.acquire(flags)
}

type subsystems struct {
	mutex    sync.Mutex
	init     func(uint32) int
	quit     func(uint32)
	getError func() string
}

func (s *subsystems) acquire(flags uint32) (func(), error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if s.init(flags) != 0 {
		return nil, fmt.Errorf("initializing SDL2 subsystems %#x: %s", flags, s.getError())
	}
	return sync.OnceFunc(func() {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		s.quit(flags)
	}), nil
}

var lifecycle subsystems

var loadCore = sync.OnceValue(func() error {
	return LoadFunctions(map[string]any{
		"SDL_InitSubSystem": &lifecycle.init,
		"SDL_QuitSubSystem": &lifecycle.quit,
		"SDL_GetError":      &lifecycle.getError,
	})
})
