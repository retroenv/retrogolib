package sdl3

import (
	"sync"

	sdllibrary "github.com/retroenv/retrogolib/internal/sdl3"
)

var (
	// GetError returns the last error message.
	GetError func() string

	// CreateWindow creates a window.
	CreateWindow func(title string, w, h int32, flags uint64) uintptr
	// DestroyWindow destroys a window.
	DestroyWindow func(window uintptr)

	// CreateRenderer creates a renderer.
	CreateRenderer func(window uintptr, name *byte) uintptr
	// RenderTexture copies a portion of the texture to the rendering target.
	RenderTexture func(renderer uintptr, texture uintptr, srcRect uintptr, dstRect uintptr) bool
	// RenderPresent updates the screen with any rendering performed.
	RenderPresent func(renderer uintptr) bool
	// DestroyRenderer destroys a renderer.
	DestroyRenderer func(renderer uintptr)

	// CreateTexture creates a texture.
	CreateTexture func(renderer uintptr, format uint32, access int, w, h int32) uintptr
	// UpdateTexture updates the given texture rectangle with new pixel data.
	UpdateTexture func(texture uintptr, rect uintptr, pixels []byte, pitch int32) bool
	// SetTextureScaleMode sets the scale mode for texture copy operations.
	SetTextureScaleMode func(texture uintptr, scaleMode int32) bool
	// DestroyTexture destroys a texture.
	DestroyTexture func(texture uintptr)

	// PollEvent polls for currently pending events.
	PollEvent func(event *event) bool
)

var imports = map[string]any{
	"SDL_CreateRenderer":      &CreateRenderer,
	"SDL_CreateTexture":       &CreateTexture,
	"SDL_CreateWindow":        &CreateWindow,
	"SDL_DestroyRenderer":     &DestroyRenderer,
	"SDL_DestroyTexture":      &DestroyTexture,
	"SDL_DestroyWindow":       &DestroyWindow,
	"SDL_GetError":            &GetError,
	"SDL_PollEvent":           &PollEvent,
	"SDL_RenderTexture":       &RenderTexture,
	"SDL_RenderPresent":       &RenderPresent,
	"SDL_SetTextureScaleMode": &SetTextureScaleMode,
	"SDL_UpdateTexture":       &UpdateTexture,
}

var setupLibrary = sync.OnceValue(func() error {
	return sdllibrary.LoadFunctions(imports)
})
