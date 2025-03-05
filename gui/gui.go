// Package gui implements different GUIs renderers.
package gui

import (
	"image"

	"github.com/retroenv/retrogolib/input"
)

// Dimensions contains settings for the window dimensions of the rendered window.
type Dimensions struct {
	ScaleFactor float64

	Height int
	Width  int
}

// Backend is an interface that gets implemented by the backend using the selected GUI.
type Backend interface {
	Image() *image.RGBA
	Dimensions() Dimensions
	WindowTitle() string

	KeyDown(key input.Key)
	KeyUp(key input.Key)
}

// Initializer defines a setup function for the selected GUI renderer.
type Initializer func(backend Backend) (guiRender func() (bool, error), guiCleanup func(), err error)

// Setup will be set by the chosen and imported GUI renderer.
// This function is the entrypoint for code importing this package to start the GUI.
var Setup Initializer
