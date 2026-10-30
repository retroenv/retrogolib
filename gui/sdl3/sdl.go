package sdl3

import (
	"fmt"
	"runtime"
	"sync"
	"unsafe"

	"github.com/retroenv/retrogolib/gui"
	"github.com/retroenv/retrogolib/gui/internal/framebuffer"
	"github.com/retroenv/retrogolib/internal/backends"
	sdllibrary "github.com/retroenv/retrogolib/internal/sdl3"
	"github.com/retroenv/retrogolib/renderer"
)

// Setup initializes the SDL library and returns a render and cleanup function.
// Call Setup, render, and cleanup on the same main goroutine. Cleanup is idempotent.
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

	if err := setupLibrary(); err != nil {
		return nil, nil, fmt.Errorf("setting up SDL library: %w", err)
	}
	release, err := sdllibrary.Acquire(SDL_INIT_VIDEO)
	if err != nil {
		return nil, nil, fmt.Errorf("acquiring SDL video: %w", err)
	}
	window, renderer, tex, err := setupSDL(dimensions, backend)
	if err != nil {
		release()
		return nil, nil, err
	}

	render := func() (bool, error) {
		return renderSDL(dimensions, backend, renderer, tex)
	}

	cleanup := sync.OnceFunc(func() {
		cleanupSDL(window, renderer, tex)
		release()
		runtime.UnlockOSThread()
	})

	return render, cleanup, nil
}

func init() {
	gui.Register("sdl3", backends.PrioritySDL3, renderer.Renderer[gui.Initializer]{
		Setup: Setup,
	})
}

// setupSDL initializes the SDL library and creates the window, renderer, and texture.
func setupSDL(dimensions gui.Dimensions, backend gui.Backend) (uintptr, uintptr, uintptr, error) {
	height := int32(float64(dimensions.Height) * dimensions.ScaleFactor)
	width := int32(float64(dimensions.Width) * dimensions.ScaleFactor)

	window := CreateWindow(backend.WindowTitle(), width, height, SDL_WINDOW_HIGH_PIXEL_DENSITY)
	if window == 0 {
		return 0, 0, 0, fmt.Errorf("creating SDL window: %s", GetError())
	}

	renderer := CreateRenderer(window, nil)
	if renderer == 0 {
		err := fmt.Errorf("creating SDL renderer: %s", GetError())
		cleanupSDL(window, 0, 0)
		return 0, 0, 0, err
	}

	tex := CreateTexture(renderer, uint32(SDL_PIXELFORMAT_ABGR8888),
		SDL_TEXTUREACCESS_STREAMING, int32(dimensions.Width), int32(dimensions.Height))
	if tex == 0 {
		err := fmt.Errorf("creating SDL texture: %s", GetError())
		cleanupSDL(window, renderer, 0)
		return 0, 0, 0, err
	}
	// SDL3 textures use linear scaling by default. Nearest scaling keeps pixel edges sharp.
	if !SetTextureScaleMode(tex, SDL_SCALEMODE_NEAREST) {
		err := fmt.Errorf("setting SDL texture scale mode: %s", GetError())
		cleanupSDL(window, renderer, tex)
		return 0, 0, 0, err
	}

	return window, renderer, tex, nil
}

// renderSDL renders the image to the SDL window.
func renderSDL(dimensions gui.Dimensions, backend gui.Backend, renderer, tex uintptr) (bool, error) {
	var ev event
	for ret := PollEvent(&ev); ret; ret = PollEvent(&ev) {
		switch ev.Type {
		case SDL_EVENT_QUIT:
			return false, nil

		case SDL_EVENT_KEY_DOWN:
			keyEvent := (*keyboardEvent)(unsafe.Pointer(&ev))
			if keyEvent.Key == K_ESCAPE {
				return false, nil
			}

			controllerKey, ok := keyMapping[keyEvent.Key]
			if ok {
				backend.KeyDown(controllerKey)
			}

		case SDL_EVENT_KEY_UP:
			keyEvent := (*keyboardEvent)(unsafe.Pointer(&ev))
			controllerKey, ok := keyMapping[keyEvent.Key]
			if ok {
				backend.KeyUp(controllerKey)
			}
		}
	}

	pixels, err := framebuffer.RGBABytes(dimensions, backend.Image())
	if err != nil {
		return false, fmt.Errorf("getting image pixels: %w", err)
	}

	if ret := UpdateTexture(tex, 0, pixels, int32(dimensions.Width*framebuffer.BytesPerPixel)); !ret {
		return false, fmt.Errorf("updating SDL texture: %s", GetError())
	}

	if ret := RenderTexture(renderer, tex, 0, 0); !ret {
		return false, fmt.Errorf("copying SDL texture: %s", GetError())
	}
	if !RenderPresent(renderer) {
		return false, fmt.Errorf("presenting SDL frame: %s", GetError())
	}

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
}
