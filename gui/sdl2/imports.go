package sdl2

import (
	"sync"

	sdllibrary "github.com/retroenv/retrogolib/internal/sdl2"
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
	"SDL_CreateRenderer":  &CreateRenderer,
	"SDL_CreateTexture":   &CreateTexture,
	"SDL_CreateWindow":    &CreateWindow,
	"SDL_DestroyRenderer": &DestroyRenderer,
	"SDL_DestroyTexture":  &DestroyTexture,
	"SDL_DestroyWindow":   &DestroyWindow,
	"SDL_GetError":        &GetError,
	"SDL_Init":            &Init,
	"SDL_PollEvent":       &PollEvent,
	"SDL_Quit":            &Quit,
	"SDL_RenderCopy":      &RenderCopy,
	"SDL_RenderPresent":   &RenderPresent,
	"SDL_UpdateTexture":   &UpdateTexture,
}

var setupLibrary = sync.OnceValue(func() error {
	return sdllibrary.LoadFunctions(imports)
})
