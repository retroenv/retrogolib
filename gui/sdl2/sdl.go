// Package sdl2 provides a SDL2 GUI renderer.
package sdl2

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/retroenv/retrogolib/gui"
	"github.com/retroenv/retrogolib/gui/internal/framebuffer"
)

// Setup initializes the SDL library and returns a render and cleanup function.
func Setup(backend gui.Backend) (guiRender func() (bool, error), guiCleanup func(), err error) {
	dimensions := backend.Dimensions()
	if err := framebuffer.ValidateDimensions(dimensions); err != nil {
		return nil, nil, fmt.Errorf("validating dimensions: %w", err)
	}

	runtime.LockOSThread()
	defer func() {
		if err != nil {
			runtime.UnlockOSThread()
		}
	}()

	window, renderer, tex, err := setupSDL(dimensions, backend)
	if err != nil {
		return nil, nil, err
	}

	render := func() (bool, error) {
		return renderSDL(dimensions, backend, renderer, tex)
	}

	cleanup := func() {
		cleanupSDL(window, renderer, tex)
	}
	return render, cleanup, nil
}

// setupSDL initializes the SDL library and creates the window, renderer, and texture.
func setupSDL(dimensions gui.Dimensions, backend gui.Backend) (uintptr, uintptr, uintptr, error) {
	if err := framebuffer.ValidateDimensions(dimensions); err != nil {
		return 0, 0, 0, fmt.Errorf("validating dimensions: %w", err)
	}

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
		Quit()
		return 0, 0, 0, fmt.Errorf("creating SDL window: %s", GetError())
	}

	renderer := CreateRenderer(window, -1, SDL_RENDERER_ACCELERATED)
	if renderer == 0 {
		cleanupSDL(window, 0, 0)
		return 0, 0, 0, fmt.Errorf("creating SDL renderer: %s", GetError())
	}

	tex := CreateTexture(renderer, uint32(SDL_PIXELFORMAT_ABGR8888),
		SDL_TEXTUREACCESS_STREAMING, int32(dimensions.Width), int32(dimensions.Height))
	if tex == 0 {
		cleanupSDL(window, renderer, 0)
		return 0, 0, 0, fmt.Errorf("creating SDL texture: %s", GetError())
	}

	return window, renderer, tex, nil
}

// renderSDL renders the image to the SDL window.
func renderSDL(dimensions gui.Dimensions, backend gui.Backend, renderer uintptr, tex uintptr) (bool, error) {
	var ev event
	for ret := PollEvent(&ev); ret != 0; ret = PollEvent(&ev) {
		switch ev.Type {
		case SDL_QUIT:
			return false, nil

		case SDL_KEYDOWN:
			keyEvent := (*keyboardEvent)(unsafe.Pointer(&ev))
			if keyEvent.Keysym.Sym == K_ESCAPE {
				return false, nil
			}

			controllerKey, ok := keyMapping[keyEvent.Keysym.Sym]
			if ok {
				backend.KeyDown(controllerKey)
			}

		case SDL_KEYUP:
			keyEvent := (*keyboardEvent)(unsafe.Pointer(&ev))
			controllerKey, ok := keyMapping[keyEvent.Keysym.Sym]
			if ok {
				backend.KeyUp(controllerKey)
			}
		}
	}

	pixels, err := framebuffer.RGBABytes(dimensions, backend.Image())
	if err != nil {
		return false, fmt.Errorf("getting image pixels: %w", err)
	}

	if ret := UpdateTexture(tex, 0, pixels, dimensions.Width*framebuffer.BytesPerPixel); ret != 0 {
		return false, fmt.Errorf("updating SDL texture: %s", GetError())
	}

	if ret := RenderCopy(renderer, tex, 0, 0); ret != 0 {
		return false, fmt.Errorf("copying SDL texture: %s", GetError())
	}
	RenderPresent(renderer)

	return true, nil
}

func cleanupSDL(window, renderer, tex uintptr) {
	if tex != 0 {
		DestroyTexture(tex)
	}
	if renderer != 0 {
		DestroyRenderer(renderer)
	}
	if window != 0 {
		DestroyWindow(window)
	}
	Quit()
}
