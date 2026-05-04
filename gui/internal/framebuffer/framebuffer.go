package framebuffer

import (
	"errors"
	"fmt"
	"image"
	"unsafe"

	"github.com/retroenv/retrogolib/gui"
)

// BytesPerPixel is the number of bytes in an RGBA frame-buffer pixel.
const BytesPerPixel = 4

// RGBABytes validates and returns RGBA pixel bytes for dimensions.
func RGBABytes(dimensions gui.Dimensions, img *image.RGBA) ([]byte, error) {
	if err := ValidateDimensions(dimensions); err != nil {
		return nil, err
	}
	if img == nil {
		return nil, errors.New("image is nil")
	}
	if len(img.Pix) == 0 {
		return nil, errors.New("image has no pixel data")
	}
	if img.Stride != dimensions.Width*BytesPerPixel {
		return nil, fmt.Errorf("image stride %d does not match expected stride %d",
			img.Stride, dimensions.Width*BytesPerPixel)
	}

	minLen := (dimensions.Height-1)*img.Stride + dimensions.Width*BytesPerPixel
	if len(img.Pix) < minLen {
		return nil, fmt.Errorf("image pixel data has length %d, need at least %d", len(img.Pix), minLen)
	}
	return img.Pix, nil
}

// RGBAPointer validates and returns an RGBA pixel pointer for dimensions.
func RGBAPointer(dimensions gui.Dimensions, img *image.RGBA) (uintptr, error) {
	pixels, err := RGBABytes(dimensions, img)
	if err != nil {
		return 0, err
	}
	return uintptr(unsafe.Pointer(&pixels[0])), nil
}

// ValidateDimensions validates renderer frame-buffer dimensions.
func ValidateDimensions(dimensions gui.Dimensions) error {
	if dimensions.Width <= 0 {
		return fmt.Errorf("width must be positive, got %d", dimensions.Width)
	}
	if dimensions.Height <= 0 {
		return fmt.Errorf("height must be positive, got %d", dimensions.Height)
	}
	if !(dimensions.ScaleFactor > 0) {
		return fmt.Errorf("scale factor must be positive, got %f", dimensions.ScaleFactor)
	}
	return nil
}
