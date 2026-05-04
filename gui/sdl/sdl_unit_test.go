package sdl

import (
	"image"
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/gui"
	"github.com/retroenv/retrogolib/input"
)

func TestRgbaPixels(t *testing.T) {
	dimensions := gui.Dimensions{
		ScaleFactor: 2,
		Height:      2,
		Width:       2,
	}
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))

	pixels, err := rgbaPixels(dimensions, img)
	assert.NoError(t, err)
	assert.Equal(t, img.Pix, pixels)
}

func TestRgbaPixelsRejectsInvalidImages(t *testing.T) {
	dimensions := gui.Dimensions{
		ScaleFactor: 2,
		Height:      2,
		Width:       2,
	}

	tests := []struct {
		name string
		img  *image.RGBA
		want string
	}{
		{
			name: "nil",
			want: "image is nil",
		},
		{
			name: "empty pix",
			img:  image.NewRGBA(image.Rect(0, 0, 0, 0)),
			want: "image has no pixel data",
		},
		{
			name: "wrong stride",
			img: &image.RGBA{
				Pix:    make([]byte, 24),
				Stride: 12,
				Rect:   image.Rect(0, 0, 2, 2),
			},
			want: "image stride",
		},
		{
			name: "too short",
			img: &image.RGBA{
				Pix:    make([]byte, 15),
				Stride: 8,
				Rect:   image.Rect(0, 0, 2, 2),
			},
			want: "image pixel data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := rgbaPixels(dimensions, tt.img)
			assert.ErrorContains(t, err, tt.want)
		})
	}
}

func TestValidateDimensions(t *testing.T) {
	tests := []struct {
		name       string
		dimensions gui.Dimensions
		want       string
	}{
		{
			name: "width",
			dimensions: gui.Dimensions{
				ScaleFactor: 2,
				Height:      1,
			},
			want: "width must be positive",
		},
		{
			name: "height",
			dimensions: gui.Dimensions{
				ScaleFactor: 2,
				Width:       1,
			},
			want: "height must be positive",
		},
		{
			name: "scale",
			dimensions: gui.Dimensions{
				Height: 1,
				Width:  1,
			},
			want: "scale factor must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDimensions(tt.dimensions)
			assert.ErrorContains(t, err, tt.want)
		})
	}
}

func TestKeyMapping(t *testing.T) {
	assert.Equal(t, input.Apostrophe, keyMapping[K_QUOTE])
	assert.Equal(t, input.Period, keyMapping[K_PERIOD])
	assert.Equal(t, input.CapsLock, keyMapping[K_CAPSLOCK])
	assert.Equal(t, input.NumLock, keyMapping[K_NUMLOCKCLEAR])
	assert.Equal(t, input.KPEqual, keyMapping[K_KP_EQUALS])
	assert.Equal(t, input.Menu, keyMapping[K_APPLICATION])
}

func TestCleanupSDL(t *testing.T) {
	var calls []string

	originalDestroyTexture := DestroyTexture
	originalDestroyRenderer := DestroyRenderer
	originalDestroyWindow := DestroyWindow
	originalQuit := Quit
	defer func() {
		DestroyTexture = originalDestroyTexture
		DestroyRenderer = originalDestroyRenderer
		DestroyWindow = originalDestroyWindow
		Quit = originalQuit
	}()

	DestroyTexture = func(uintptr) {
		calls = append(calls, "texture")
	}
	DestroyRenderer = func(uintptr) {
		calls = append(calls, "renderer")
	}
	DestroyWindow = func(uintptr) {
		calls = append(calls, "window")
	}
	Quit = func() {
		calls = append(calls, "quit")
	}

	cleanupSDL(1, 2, 3)

	assert.Equal(t, []string{"texture", "renderer", "window", "quit"}, calls)
}
