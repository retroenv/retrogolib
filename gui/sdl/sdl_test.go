//go:build manualtest

package sdl

import (
	"image"
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/gui"
	"github.com/retroenv/retrogolib/input"
)

func TestSetupGoSDL(t *testing.T) {
	b := &backend{}
	render, cleanup, err := Setup(b)
	assert.NoError(t, err)
	_, err = render()
	assert.NoError(t, err)
	cleanup()
}

const height = 240
const width = 256

type backend struct {
	img *image.RGBA
}

func (b *backend) Image() *image.RGBA {
	if b.img == nil {
		b.img = image.NewRGBA(image.Rect(0, 0, width, height))
	}
	return b.img
}

func (b *backend) Dimensions() gui.Dimensions {
	return gui.Dimensions{
		ScaleFactor: 2.0,
		Height:      height,
		Width:       width,
	}
}

func (b *backend) WindowTitle() string {
	return "unit-test"
}

func (b *backend) KeyDown(key input.Key) {
}

func (b *backend) KeyUp(key input.Key) {
}
