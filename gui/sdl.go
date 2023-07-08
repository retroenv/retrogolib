//go:build !nesgo && !nogui && sdl

package gui

import (
	"fmt"
	"unsafe"

	"github.com/retroenv/retrogolib/input"
	"github.com/veandco/go-sdl2/sdl"
)

func init() {
	Setup = setupSDLGui
}

// TODO: complete
var sdlKeyMapping = map[sdl.Keycode]input.Key{
	sdl.K_UP:        input.Up,
	sdl.K_DOWN:      input.Down,
	sdl.K_LEFT:      input.Left,
	sdl.K_RIGHT:     input.Right,
	sdl.K_z:         input.A,
	sdl.K_x:         input.B,
	sdl.K_RETURN:    input.Enter,
	sdl.K_BACKSPACE: input.Backspace,
}

func setupSDLGui(backend Backend) (guiRender func() (bool, error), guiCleanup func(), err error) {
	dimensions := backend.Dimensions()

	window, renderer, tex, err := setupSDL(dimensions, backend)
	if err != nil {
		return nil, nil, err
	}

	render := func() (bool, error) {
		return renderSDL(dimensions, backend, renderer, tex)
	}

	cleanup := func() {
		_ = tex.Destroy()
		_ = renderer.Destroy()
		_ = window.Destroy()
		sdl.Quit()
	}
	return render, cleanup, nil
}

func setupSDL(dimensions Dimensions, backend Backend) (*sdl.Window, *sdl.Renderer, *sdl.Texture, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, nil, nil, fmt.Errorf("initializing SDL: %w", err)
	}

	height := int32(float64(dimensions.Height) * dimensions.ScaleFactor)
	width := int32(float64(dimensions.Width) * dimensions.ScaleFactor)

	window, err := sdl.CreateWindow(backend.WindowTitle(), sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED, width, height,
		sdl.WINDOW_SHOWN|sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("creating SDL window: %w", err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("creating SDL renderer: %w", err)
	}

	tex, err := renderer.CreateTexture(uint32(sdl.PIXELFORMAT_ABGR8888),
		sdl.TEXTUREACCESS_STREAMING, int32(dimensions.Width), int32(dimensions.Height))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("creating SDL texture: %w", err)
	}

	return window, renderer, tex, nil
}

func renderSDL(dimensions Dimensions, backend Backend, renderer *sdl.Renderer, tex *sdl.Texture) (bool, error) {
	running := true

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch et := event.(type) {
		case *sdl.QuitEvent:
			running = false
			break

		case *sdl.KeyboardEvent:
			if et.Type == sdl.KEYDOWN && et.Keysym.Sym == sdl.K_ESCAPE {
				running = false
				break
			}
			onSDLKey(backend, et)
		}
	}

	image := backend.Image()
	if err := tex.Update(nil, unsafe.Pointer(&image.Pix), dimensions.Width); err != nil {
		return false, err
	}

	if err := renderer.Copy(tex, nil, nil); err != nil {
		return false, err
	}
	renderer.Present()

	return running, nil
}

func onSDLKey(backend Backend, event *sdl.KeyboardEvent) {
	controllerKey, ok := sdlKeyMapping[event.Keysym.Sym]
	if !ok {
		return
	}

	switch event.Type {
	case sdl.KEYDOWN:
		backend.KeyDown(controllerKey)

	case sdl.KEYUP:
		backend.KeyUp(controllerKey)
	}
}
