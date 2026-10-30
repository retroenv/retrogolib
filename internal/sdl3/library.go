// Package sdl3 coordinates SDL3 library loading and subsystem ownership.
package sdl3

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
		names = []string{"libSDL3.dylib"}
	case "linux":
		names = []string{"libSDL3.so.0", "libSDL3.so"}
	case "freebsd":
		names = []string{"libSDL3.so"}
	case "windows":
		names = []string{"SDL3.dll"}
	default:
		return fmt.Errorf("SDL3 is not supported on GOOS=%s", runtime.GOOS)
	}

	var lastErr error
	for _, name := range names {
		if _, err := dynlib.LoadFunctions(name, imports); err != nil {
			lastErr = err
			continue
		}
		return nil
	}

	return fmt.Errorf("loading SDL3: %w", lastErr)
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
	init     func(uint32) bool
	quit     func(uint32)
	getError func() string
}

func (s *subsystems) acquire(flags uint32) (func(), error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if !s.init(flags) {
		return nil, fmt.Errorf("initializing SDL3 subsystems %#x: %s", flags, s.getError())
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
