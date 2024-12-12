package sdl

import (
	"fmt"
	"runtime"

	"github.com/ebitengine/purego"
)

var (
	// Init initializes the SDL library.
	Init func(flags uint32) int
	// GetError returns the last error message.
	GetError func() string
	// Quit quits the SDL library.
	Quit func()

	// CreateWindow creates a window.
	CreateWindow func(title string, x, y, w, h int32, flags uint32) uintptr
	// DestroyWindow destroys a window.
	DestroyWindow func(window uintptr)

	// CreateRenderer creates a renderer.
	CreateRenderer func(window uintptr, index int, flags uint32) uintptr
	// RenderCopy copies a portion of the texture to the rendering target.
	RenderCopy func(renderer uintptr, texture uintptr, srcRect uintptr, dstRect uintptr) int
	// RenderPresent updates the screen with any rendering performed.
	RenderPresent func(renderer uintptr)
	// DestroyRenderer destroys a renderer.
	DestroyRenderer func(renderer uintptr)

	// CreateTexture creates a texture.
	CreateTexture func(renderer uintptr, format uint32, access int, w, h int32) uintptr
	// UpdateTexture updates the given texture rectangle with new pixel data.
	UpdateTexture func(texture uintptr, rect uintptr, pixels []byte, pitch int) int
	// DestroyTexture destroys a texture.
	DestroyTexture func(texture uintptr)

	// PollEvent polls for currently pending events.
	PollEvent func(event *event) int
)

var imports = map[string]any{
	"SDL_Init":     &Init,
	"SDL_GetError": &GetError,
	"SDL_Quit":     &Quit,

	"SDL_CreateWindow":  &CreateWindow,
	"SDL_DestroyWindow": &DestroyWindow,

	"SDL_CreateRenderer":  &CreateRenderer,
	"SDL_RenderCopy":      &RenderCopy,
	"SDL_RenderPresent":   &RenderPresent,
	"SDL_DestroyRenderer": &DestroyRenderer,

	"SDL_CreateTexture":  &CreateTexture,
	"SDL_UpdateTexture":  &UpdateTexture,
	"SDL_DestroyTexture": &DestroyTexture,

	"SDL_PollEvent": &PollEvent,
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
