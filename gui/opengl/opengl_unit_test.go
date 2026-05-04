package opengl

import (
	"image"
	"testing"
	"unsafe"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/gui"
	"github.com/retroenv/retrogolib/input"
)

func TestCString(t *testing.T) {
	raw := []byte("GLFW failed\x00ignored")

	assert.Equal(t, "GLFW failed", cString((*byte)(unsafe.Pointer(&raw[0]))))
	assert.Equal(t, "", cString(nil))
}

func TestRgbaPixels(t *testing.T) {
	dimensions := gui.Dimensions{
		ScaleFactor: 2,
		Height:      2,
		Width:       2,
	}
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))

	ptr, err := rgbaPixels(dimensions, img)
	assert.NoError(t, err)
	assert.NotEqual(t, uintptr(0), ptr)
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
	assert.Equal(t, input.Escape, keyMapping[GLFW_KEY_ESCAPE])
	assert.Equal(t, input.A, keyMapping[GLFW_KEY_A])
	assert.Equal(t, input.KPEnter, keyMapping[GLFW_KEY_KP_ENTER])
	assert.Equal(t, input.Menu, keyMapping[GLFW_KEY_MENU])
}
