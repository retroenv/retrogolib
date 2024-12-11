package sdl

import "github.com/retroenv/retrogolib/input"

const (
	K_UNKNOWN      = 0
	K_RETURN       = 13
	K_ESCAPE       = 27
	K_BACKSPACE    = 8
	K_TAB          = 9
	K_SPACE        = 32
	K_EXCLAIM      = 33
	K_QUOTEDBL     = 34
	K_HASH         = 35
	K_PERCENT      = 37
	K_DOLLAR       = 36
	K_AMPERSAND    = 38
	K_QUOTE        = 39
	K_LEFTPAREN    = 40
	K_RIGHTPAREN   = 41
	K_ASTERISK     = 42
	K_PLUS         = 43
	K_COMMA        = 44
	K_MINUS        = 45
	K_PERIOD       = 46
	K_SLASH        = 47
	K_0            = 48
	K_1            = 49
	K_2            = 50
	K_3            = 51
	K_4            = 52
	K_5            = 53
	K_6            = 54
	K_7            = 55
	K_8            = 56
	K_9            = 57
	K_COLON        = 58
	K_SEMICOLON    = 59
	K_LESS         = 60
	K_EQUALS       = 61
	K_GREATER      = 62
	K_QUESTION     = 63
	K_AT           = 64
	K_LEFTBRACKET  = 91
	K_BACKSLASH    = 92
	K_RIGHTBRACKET = 93
	K_CARET        = 94
	K_UNDERSCORE   = 95
	K_BACKQUOTE    = 96
	K_a            = 97
	K_b            = 98
	K_c            = 99
	K_d            = 100
	K_e            = 101
	K_f            = 102
	K_g            = 103
	K_h            = 104
	K_i            = 105
	K_j            = 106
	K_k            = 107
	K_l            = 108
	K_m            = 109
	K_n            = 110
	K_o            = 111
	K_p            = 112
	K_q            = 113
	K_r            = 114
	K_s            = 115
	K_t            = 116
	K_u            = 117
	K_v            = 118
	K_w            = 119
	K_x            = 120
	K_y            = 121
	K_z            = 122
	K_CAPSLOCK     = 1073741881
	K_F1           = 1073741882
	K_F2           = 1073741883
	K_F3           = 1073741884
	K_F4           = 1073741885
	K_F5           = 1073741886
	K_F6           = 1073741887
	K_F7           = 1073741888
	K_F8           = 1073741889
	K_F9           = 1073741890
	K_F10          = 1073741891
	K_F11          = 1073741892
	K_F12          = 1073741893
	K_F13          = 1073741928
	K_F14          = 1073741929
	K_F15          = 1073741930
	K_F16          = 1073741931
	K_F17          = 1073741932
	K_F18          = 1073741933
	K_F19          = 1073741934
	K_F20          = 1073741935
	K_F21          = 1073741936
	K_F22          = 1073741937
	K_F23          = 1073741938
	K_F24          = 1073741939
	K_PRINTSCREEN  = 1073741894
	K_SCROLLLOCK   = 1073741895
	K_PAUSE        = 1073741896
	K_INSERT       = 1073741897
	K_HOME         = 1073741898
	K_PAGEUP       = 1073741899
	K_DELETE       = 127
	K_END          = 1073741901
	K_PAGEDOWN     = 1073741902
	K_RIGHT        = 1073741903
	K_LEFT         = 1073741904
	K_DOWN         = 1073741905
	K_UP           = 1073741906
	K_NUMLOCKCLEAR = 1073741907
	K_KP_DIVIDE    = 1073741908
	K_KP_MULTIPLY  = 1073741909
	K_KP_MINUS     = 1073741910
	K_KP_PLUS      = 1073741911
	K_KP_ENTER     = 1073741912
	K_KP_1         = 1073741913
	K_KP_2         = 1073741914
	K_KP_3         = 1073741915
	K_KP_4         = 1073741916
	K_KP_5         = 1073741917
	K_KP_6         = 1073741918
	K_KP_7         = 1073741919
	K_KP_8         = 1073741920
	K_KP_9         = 1073741921
	K_KP_0         = 1073741922
	K_KP_PERIOD    = 1073741923
	K_KP_COMMA     = 1073741957
	K_KP_EQUALS    = 1073741927
	K_LCTRL        = 1073742048
	K_LSHIFT       = 1073742049
	K_LALT         = 1073742050
	K_LGUI         = 1073742051
	K_RCTRL        = 1073742052
	K_RSHIFT       = 1073742053
	K_RALT         = 1073742054
	K_RGUI         = 1073742055
)

var keyMapping = map[keycode]input.Key{
	K_UNKNOWN:      input.Unknown,
	K_RETURN:       input.Enter,
	K_ESCAPE:       input.Escape,
	K_BACKSPACE:    input.Backspace,
	K_TAB:          input.Tab,
	K_SPACE:        input.Space,
	K_COMMA:        input.Comma,
	K_MINUS:        input.Minus,
	K_SLASH:        input.Slash,
	K_0:            input.Key0,
	K_1:            input.Key1,
	K_2:            input.Key2,
	K_3:            input.Key3,
	K_4:            input.Key4,
	K_5:            input.Key5,
	K_6:            input.Key6,
	K_7:            input.Key7,
	K_8:            input.Key8,
	K_9:            input.Key9,
	K_COLON:        input.Semicolon,
	K_SEMICOLON:    input.Semicolon,
	K_EQUALS:       input.Equal,
	K_LEFTBRACKET:  input.LeftBracket,
	K_BACKSLASH:    input.Backslash,
	K_RIGHTBRACKET: input.RightBracket,
	K_a:            input.A,
	K_b:            input.B,
	K_c:            input.C,
	K_d:            input.D,
	K_e:            input.E,
	K_f:            input.F,
	K_g:            input.G,
	K_h:            input.H,
	K_i:            input.I,
	K_j:            input.J,
	K_k:            input.K,
	K_l:            input.L,
	K_m:            input.M,
	K_n:            input.N,
	K_o:            input.O,
	K_p:            input.P,
	K_q:            input.Q,
	K_r:            input.R,
	K_s:            input.S,
	K_t:            input.T,
	K_u:            input.U,
	K_v:            input.V,
	K_w:            input.W,
	K_x:            input.X,
	K_y:            input.Y,
	K_z:            input.Z,
	K_F1:           input.F1,
	K_F2:           input.F2,
	K_F3:           input.F3,
	K_F4:           input.F4,
	K_F5:           input.F5,
	K_F6:           input.F6,
	K_F7:           input.F7,
	K_F8:           input.F8,
	K_F9:           input.F9,
	K_F10:          input.F10,
	K_F11:          input.F11,
	K_F12:          input.F12,
	K_PRINTSCREEN:  input.PrintScreen,
	K_SCROLLLOCK:   input.ScrollLock,
	K_PAUSE:        input.Pause,
	K_INSERT:       input.Insert,
	K_HOME:         input.Home,
	K_PAGEUP:       input.PageUp,
	K_DELETE:       input.Delete,
	K_END:          input.End,
	K_PAGEDOWN:     input.PageDown,
	K_RIGHT:        input.Right,
	K_LEFT:         input.Left,
	K_DOWN:         input.Down,
	K_UP:           input.Up,
	K_KP_DIVIDE:    input.KPDivide,
	K_KP_MULTIPLY:  input.KPMultiply,
	K_KP_MINUS:     input.KPSubtract,
	K_KP_PLUS:      input.KPAdd,
	K_KP_ENTER:     input.KPEnter,
	K_KP_1:         input.KP1,
	K_KP_2:         input.KP2,
	K_KP_3:         input.KP3,
	K_KP_4:         input.KP4,
	K_KP_5:         input.KP5,
	K_KP_6:         input.KP6,
	K_KP_7:         input.KP7,
	K_KP_8:         input.KP8,
	K_KP_9:         input.KP9,
	K_KP_0:         input.KP0,
	K_KP_PERIOD:    input.KPDecimal,
	K_F13:          input.F13,
	K_F14:          input.F14,
	K_F15:          input.F15,
	K_F16:          input.F16,
	K_F17:          input.F17,
	K_F18:          input.F18,
	K_F19:          input.F19,
	K_F20:          input.F20,
	K_F21:          input.F21,
	K_F22:          input.F22,
	K_F23:          input.F23,
	K_F24:          input.F24,
	K_LCTRL:        input.LeftControl,
	K_LSHIFT:       input.LeftShift,
	K_LALT:         input.LeftAlt,
	K_LGUI:         input.LeftSuper,
	K_RCTRL:        input.RightControl,
	K_RSHIFT:       input.RightShift,
	K_RALT:         input.RightAlt,
	K_RGUI:         input.RightSuper,
}
