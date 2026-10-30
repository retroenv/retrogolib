package gui

import (
	"errors"
	"fmt"
	"image"

	"github.com/retroenv/retrogolib/input"
	"github.com/retroenv/retrogolib/internal/backends"
	"github.com/retroenv/retrogolib/renderer"
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

// ErrNoRenderer is returned by Setup when no GUI renderer is registered.
var ErrNoRenderer = errors.New("no GUI renderer registered")

var renderers backends.Registry[renderer.Renderer[Initializer]]

// Register adds a GUI renderer for Setup and returns a function that removes it.
// Renderer packages call Register in their init function.
// Setup tries renderers with a higher priority first. Renderers with an equal
// priority run in registration order. Register panics when the name is empty,
// the name is already registered, or the renderer has no Setup function.
func Register(name string, priority int, r renderer.Renderer[Initializer]) (unregister func()) {
	if r.Setup == nil {
		panic("gui: renderer setup function is nil")
	}
	return renderers.Register(name, priority, r)
}

// Available reports whether at least one GUI renderer is registered.
func Available() bool {
	return renderers.Len() > 0
}

// Setup starts the first registered GUI renderer that initializes without error.
// It is SetupWithInfo without the renderer information.
func Setup(backend Backend) (guiRender func() (bool, error), guiCleanup func(), err error) {
	guiRender, guiCleanup, _, err = SetupWithInfo(backend)
	return guiRender, guiCleanup, err
}

// SetupWithInfo starts the first registered GUI renderer that initializes without
// error and returns information about it. It tries the renderers in priority
// order, for example a renderer for a newer library version before a renderer
// for an older one. This is the entrypoint for code importing this package to
// start the GUI. Call it on the main goroutine.
//
// When all renderers fail, the returned error contains the error of every
// renderer. It returns ErrNoRenderer when no renderer is registered.
func SetupWithInfo(backend Backend) (guiRender func() (bool, error), guiCleanup func(), info renderer.Info, err error) {
	result, info, err := backends.First(&renderers, ErrNoRenderer, func(r renderer.Renderer[Initializer]) (started, error) {
		return newStarted(r.Setup(backend))
	})
	if err != nil {
		return nil, nil, renderer.Info{}, fmt.Errorf("starting GUI renderer: %w", err)
	}

	return result.render, result.cleanup, info, nil
}

// started contains the functions of a renderer that started.
type started struct {
	render  func() (bool, error)
	cleanup func()
}

// newStarted combines the results of an Initializer call.
func newStarted(render func() (bool, error), cleanup func(), err error) (started, error) {
	result := started{
		render:  render,
		cleanup: cleanup,
	}

	return result, err
}
