// Package sdl provides a SDL GUI renderer.
package sdl

import (
	"fmt"
	"unsafe"

	"github.com/retroenv/retrogolib/gui"
)

func Setup(backend gui.Backend) (guiRender func() (bool, error), guiCleanup func(), err error) {
	dimensions := backend.Dimensions()

	window, renderer, tex, err := setupSDL(dimensions, backend)
	if err != nil {
		return nil, nil, err
	}

	render := func() (bool, error) {
		return renderSDL(dimensions, backend, renderer, tex)
	}

	cleanup := func() {
		DestroyTexture(tex)
		DestroyRenderer(renderer)
		DestroyWindow(window)
		Quit()
	}
	return render, cleanup, nil
}

func setupSDL(dimensions gui.Dimensions, backend gui.Backend) (uintptr, uintptr, uintptr, error) {
	if err := setupLibrary(); err != nil {
		return 0, 0, 0, fmt.Errorf("setting up SDL library: %w", err)
	}

	if ret := Init(SDL_INIT_EVERYTHING); ret != 0 {
		return 0, 0, 0, fmt.Errorf("initializing SDL: %s", GetError())
	}

	height := int32(float64(dimensions.Height) * dimensions.ScaleFactor)
	width := int32(float64(dimensions.Width) * dimensions.ScaleFactor)

	window := CreateWindow(backend.WindowTitle(), SDL_WINDOWPOS_CENTERED,
		SDL_WINDOWPOS_CENTERED, width, height,
		SDL_WINDOW_SHOWN|SDL_WINDOW_ALLOW_HIGHDPI)
	if window == 0 {
		return 0, 0, 0, fmt.Errorf("creating SDL window: %s", GetError())
	}

	renderer := CreateRenderer(window, -1, SDL_RENDERER_ACCELERATED)
	if renderer == 0 {
		return 0, 0, 0, fmt.Errorf("creating SDL renderer: %s", GetError())
	}

	tex := CreateTexture(renderer, uint32(SDL_PIXELFORMAT_ABGR8888),
		SDL_TEXTUREACCESS_STREAMING, int32(dimensions.Width), int32(dimensions.Height))
	if tex == 0 {
		return 0, 0, 0, fmt.Errorf("creating SDL texture: %s", GetError())
	}

	return window, renderer, tex, nil
}

// renderSDL
func renderSDL(dimensions gui.Dimensions, backend gui.Backend, renderer uintptr, tex uintptr) (bool, error) {
	var event event
	for ret := PollEvent(&event); ret != 0; ret = PollEvent(&event) {
		switch event.Event {
		case SDL_QUIT:
			return false, nil

		case SDL_KEYDOWN:
			keyEvent := (*keyboardEvent)(unsafe.Pointer(&event.Event))
			if keyEvent.Keysym.Sym == K_ESCAPE {
				return false, nil
			}

			controllerKey, ok := keyMapping[keyEvent.Keysym.Sym]
			if ok {
				backend.KeyDown(controllerKey)
			}

		case SDL_KEYUP:
			keyEvent := (*keyboardEvent)(unsafe.Pointer(&event.Event))
			controllerKey, ok := keyMapping[keyEvent.Keysym.Sym]
			if ok {
				backend.KeyUp(controllerKey)
			}
		}
	}

	image := backend.Image()
	if ret := UpdateTexture(tex, 0, image.Pix, dimensions.Width); ret != 0 {
		return false, fmt.Errorf("updating SDL texture: %s", GetError())
	}

	if ret := RenderCopy(renderer, tex, 0, 0); ret != 0 {
		return false, fmt.Errorf("copying SDL texture: %s", GetError())
	}
	RenderPresent(renderer)

	return true, nil
}
