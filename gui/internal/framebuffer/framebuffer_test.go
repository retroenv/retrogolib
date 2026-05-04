package framebuffer

import (
	"image"
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/gui"
)

func TestRGBABytes(t *testing.T) {
	dimensions := gui.Dimensions{
		ScaleFactor: 2,
		Height:      2,
		Width:       2,
	}
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))

	pixels, err := RGBABytes(dimensions, img)
	assert.NoError(t, err)
	assert.Equal(t, img.Pix, pixels)
}

func TestRGBAPointer(t *testing.T) {
	dimensions := gui.Dimensions{
		ScaleFactor: 2,
		Height:      2,
		Width:       2,
	}
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))

	ptr, err := RGBAPointer(dimensions, img)
	assert.NoError(t, err)
	assert.NotEqual(t, uintptr(0), ptr)
}

func TestRGBABytesRejectsInvalidImages(t *testing.T) {
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
			_, err := RGBABytes(dimensions, tt.img)
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
			err := ValidateDimensions(tt.dimensions)
			assert.ErrorContains(t, err, tt.want)
		})
	}
}
