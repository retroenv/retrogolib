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

var sdlKeyMapping = map[sdl.Keycode]input.Key{
	sdl.K_UNKNOWN:      input.Unknown,
	sdl.K_RETURN:       input.Enter,
	sdl.K_ESCAPE:       input.Escape,
	sdl.K_BACKSPACE:    input.Backspace,
	sdl.K_TAB:          input.Tab,
	sdl.K_SPACE:        input.Space,
	sdl.K_COMMA:        input.Comma,
	sdl.K_MINUS:        input.Minus,
	sdl.K_SLASH:        input.Slash,
	sdl.K_0:            input.Key0,
	sdl.K_1:            input.Key1,
	sdl.K_2:            input.Key2,
	sdl.K_3:            input.Key3,
	sdl.K_4:            input.Key4,
	sdl.K_5:            input.Key5,
	sdl.K_6:            input.Key6,
	sdl.K_7:            input.Key7,
	sdl.K_8:            input.Key8,
	sdl.K_9:            input.Key9,
	sdl.K_COLON:        input.Semicolon,
	sdl.K_SEMICOLON:    input.Semicolon,
	sdl.K_EQUALS:       input.Equal,
	sdl.K_LEFTBRACKET:  input.LeftBracket,
	sdl.K_BACKSLASH:    input.Backslash,
	sdl.K_RIGHTBRACKET: input.RightBracket,
	sdl.K_a:            input.A,
	sdl.K_b:            input.B,
	sdl.K_c:            input.C,
	sdl.K_d:            input.D,
	sdl.K_e:            input.E,
	sdl.K_f:            input.F,
	sdl.K_g:            input.G,
	sdl.K_h:            input.H,
	sdl.K_i:            input.I,
	sdl.K_j:            input.J,
	sdl.K_k:            input.K,
	sdl.K_l:            input.L,
	sdl.K_m:            input.M,
	sdl.K_n:            input.N,
	sdl.K_o:            input.O,
	sdl.K_p:            input.P,
	sdl.K_q:            input.Q,
	sdl.K_r:            input.R,
	sdl.K_s:            input.S,
	sdl.K_t:            input.T,
	sdl.K_u:            input.U,
	sdl.K_v:            input.V,
	sdl.K_w:            input.W,
	sdl.K_x:            input.X,
	sdl.K_y:            input.Y,
	sdl.K_z:            input.Z,
	sdl.K_F1:           input.F1,
	sdl.K_F2:           input.F2,
	sdl.K_F3:           input.F3,
	sdl.K_F4:           input.F4,
	sdl.K_F5:           input.F5,
	sdl.K_F6:           input.F6,
	sdl.K_F7:           input.F7,
	sdl.K_F8:           input.F8,
	sdl.K_F9:           input.F9,
	sdl.K_F10:          input.F10,
	sdl.K_F11:          input.F11,
	sdl.K_F12:          input.F12,
	sdl.K_PRINTSCREEN:  input.PrintScreen,
	sdl.K_SCROLLLOCK:   input.ScrollLock,
	sdl.K_PAUSE:        input.Pause,
	sdl.K_INSERT:       input.Insert,
	sdl.K_HOME:         input.Home,
	sdl.K_PAGEUP:       input.PageUp,
	sdl.K_DELETE:       input.Delete,
	sdl.K_END:          input.End,
	sdl.K_PAGEDOWN:     input.PageDown,
	sdl.K_RIGHT:        input.Right,
	sdl.K_LEFT:         input.Left,
	sdl.K_DOWN:         input.Down,
	sdl.K_UP:           input.Up,
	sdl.K_KP_DIVIDE:    input.KPDivide,
	sdl.K_KP_MULTIPLY:  input.KPMultiply,
	sdl.K_KP_MINUS:     input.KPSubtract,
	sdl.K_KP_PLUS:      input.KPAdd,
	sdl.K_KP_ENTER:     input.KPEnter,
	sdl.K_KP_1:         input.KP1,
	sdl.K_KP_2:         input.KP2,
	sdl.K_KP_3:         input.KP3,
	sdl.K_KP_4:         input.KP4,
	sdl.K_KP_5:         input.KP5,
	sdl.K_KP_6:         input.KP6,
	sdl.K_KP_7:         input.KP7,
	sdl.K_KP_8:         input.KP8,
	sdl.K_KP_9:         input.KP9,
	sdl.K_KP_0:         input.KP0,
	sdl.K_KP_PERIOD:    input.KPDecimal,
	sdl.K_F13:          input.F13,
	sdl.K_F14:          input.F14,
	sdl.K_F15:          input.F15,
	sdl.K_F16:          input.F16,
	sdl.K_F17:          input.F17,
	sdl.K_F18:          input.F18,
	sdl.K_F19:          input.F19,
	sdl.K_F20:          input.F20,
	sdl.K_F21:          input.F21,
	sdl.K_F22:          input.F22,
	sdl.K_F23:          input.F23,
	sdl.K_F24:          input.F24,
	sdl.K_LCTRL:        input.LeftControl,
	sdl.K_LSHIFT:       input.LeftShift,
	sdl.K_LALT:         input.LeftAlt,
	sdl.K_LGUI:         input.LeftSuper,
	sdl.K_RCTRL:        input.RightControl,
	sdl.K_RSHIFT:       input.RightShift,
	sdl.K_RALT:         input.RightAlt,
	sdl.K_RGUI:         input.RightSuper,
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

	data := unsafe.Pointer(&image.Pix[0])
	if err := tex.Update(nil, data, dimensions.Width); err != nil {
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
