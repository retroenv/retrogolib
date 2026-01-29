package gui_test

import (
	"image"
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/gui"
	"github.com/retroenv/retrogolib/input"
)

// MockBackend implements the gui.Backend interface for testing
type MockBackend struct {
	img        *image.RGBA
	dimensions gui.Dimensions
	title      string
	downKeys   []input.Key
	upKeys     []input.Key
}

func NewMockBackend() *MockBackend {
	return &MockBackend{
		img: image.NewRGBA(image.Rect(0, 0, 256, 240)),
		dimensions: gui.Dimensions{
			Width:       256,
			Height:      240,
			ScaleFactor: 2.0,
		},
		title:    "Test Window",
		downKeys: make([]input.Key, 0),
		upKeys:   make([]input.Key, 0),
	}
}

func (m *MockBackend) Image() *image.RGBA {
	return m.img
}

func (m *MockBackend) Dimensions() gui.Dimensions {
	return m.dimensions
}

func (m *MockBackend) WindowTitle() string {
	return m.title
}

func (m *MockBackend) KeyDown(key input.Key) {
	m.downKeys = append(m.downKeys, key)
}

func (m *MockBackend) KeyUp(key input.Key) {
	m.upKeys = append(m.upKeys, key)
}

func TestDimensions_Properties(t *testing.T) {
	dims := gui.Dimensions{
		Width:       320,
		Height:      240,
		ScaleFactor: 1.5,
	}

	assert.Equal(t, 320, dims.Width, "Width should be set correctly")
	assert.Equal(t, 240, dims.Height, "Height should be set correctly")
	assert.Equal(t, 1.5, dims.ScaleFactor, "ScaleFactor should be set correctly")
}

func TestDimensions_ZeroValues(t *testing.T) {
	dims := gui.Dimensions{}

	assert.Equal(t, 0, dims.Width, "Default width should be 0")
	assert.Equal(t, 0, dims.Height, "Default height should be 0")
	assert.Equal(t, 0.0, dims.ScaleFactor, "Default scale factor should be 0.0")
}

func TestBackend_Interface(t *testing.T) {
	backend := NewMockBackend()

	// Test Image method
	img := backend.Image()
	assert.NotNil(t, img, "Image should not be nil")
	assert.Equal(t, 256, img.Bounds().Dx(), "Image width should be 256")
	assert.Equal(t, 240, img.Bounds().Dy(), "Image height should be 240")

	// Test Dimensions method
	dims := backend.Dimensions()
	assert.Equal(t, 256, dims.Width, "Dimensions width should be 256")
	assert.Equal(t, 240, dims.Height, "Dimensions height should be 240")
	assert.Equal(t, 2.0, dims.ScaleFactor, "Scale factor should be 2.0")

	// Test WindowTitle method
	title := backend.WindowTitle()
	assert.Equal(t, "Test Window", title, "Window title should match")
}

func TestBackend_KeyEvents(t *testing.T) {
	backend := NewMockBackend()

	// Test KeyDown
	backend.KeyDown(input.A)
	backend.KeyDown(input.Space)
	backend.KeyDown(input.Enter)

	assert.Equal(t, 3, len(backend.downKeys), "Should have recorded 3 key down events")
	assert.Equal(t, input.A, backend.downKeys[0], "First key down should be A")
	assert.Equal(t, input.Space, backend.downKeys[1], "Second key down should be Space")
	assert.Equal(t, input.Enter, backend.downKeys[2], "Third key down should be Enter")

	// Test KeyUp
	backend.KeyUp(input.A)
	backend.KeyUp(input.Escape)

	assert.Equal(t, 2, len(backend.upKeys), "Should have recorded 2 key up events")
	assert.Equal(t, input.A, backend.upKeys[0], "First key up should be A")
	assert.Equal(t, input.Escape, backend.upKeys[1], "Second key up should be Escape")
}

func TestBackend_KeyEventSequence(t *testing.T) {
	backend := NewMockBackend()

	// Simulate typical key press sequence
	backend.KeyDown(input.A)
	backend.KeyUp(input.A)
	backend.KeyDown(input.B)
	backend.KeyUp(input.B)

	assert.Equal(t, 2, len(backend.downKeys), "Should have 2 key down events")
	assert.Equal(t, 2, len(backend.upKeys), "Should have 2 key up events")

	assert.Equal(t, input.A, backend.downKeys[0], "First down should be A")
	assert.Equal(t, input.B, backend.downKeys[1], "Second down should be B")
	assert.Equal(t, input.A, backend.upKeys[0], "First up should be A")
	assert.Equal(t, input.B, backend.upKeys[1], "Second up should be B")
}

func TestInitializer_Type(t *testing.T) {
	// Test that Initializer type can be assigned
	init := gui.Initializer(func(backend gui.Backend) (func() (bool, error), func(), error) {
		return func() (bool, error) { return true, nil }, func() {}, nil
	})

	assert.NotNil(t, init, "Initializer should be assignable")

	// Test that it can be called
	backend := NewMockBackend()
	render, cleanup, err := init(backend)

	assert.NoError(t, err, "Initializer should not return error")
	assert.NotNil(t, render, "Render function should not be nil")
	assert.NotNil(t, cleanup, "Cleanup function should not be nil")

	// Test that render function works
	running, err := render()
	assert.NoError(t, err, "Render function should not return error")
	assert.True(t, running, "Render function should return true")
}

func TestSetup_GlobalVariable(t *testing.T) {
	// Test that Setup variable can be assigned
	originalSetup := gui.Setup
	defer func() {
		gui.Setup = originalSetup
	}()

	gui.Setup = func(backend gui.Backend) (func() (bool, error), func(), error) {
		return func() (bool, error) { return false, nil }, func() {}, nil
	}

	assert.NotNil(t, gui.Setup, "Setup should be assignable")

	// Test that assigned Setup can be used
	backend := NewMockBackend()
	render, cleanup, err := gui.Setup(backend)

	assert.NoError(t, err, "Setup should not return error")
	assert.NotNil(t, render, "Render function should not be nil")
	assert.NotNil(t, cleanup, "Cleanup function should not be nil")

	// Test render function behavior
	running, err := render()
	assert.NoError(t, err, "Render function should not return error")
	assert.False(t, running, "Render function should return false")
}

func TestDimensions_Equality(t *testing.T) {
	dims1 := gui.Dimensions{Width: 320, Height: 240, ScaleFactor: 2.0}
	dims2 := gui.Dimensions{Width: 320, Height: 240, ScaleFactor: 2.0}
	dims3 := gui.Dimensions{Width: 640, Height: 480, ScaleFactor: 1.0}

	assert.Equal(t, dims1, dims2, "Identical dimensions should be equal")
	assert.NotEqual(t, dims1, dims3, "Different dimensions should not be equal")
}

func TestBackend_MultipleKeyEvents(t *testing.T) {
	backend := NewMockBackend()

	// Test rapid key events
	keys := []input.Key{
		input.A, input.B, input.C, input.Key1, input.Key2, input.Key3,
		input.F1, input.F2, input.Space, input.Enter, input.Escape,
	}

	// Press all keys
	for _, key := range keys {
		backend.KeyDown(key)
	}

	// Release all keys
	for _, key := range keys {
		backend.KeyUp(key)
	}

	assert.Equal(t, len(keys), len(backend.downKeys), "All key down events should be recorded")
	assert.Equal(t, len(keys), len(backend.upKeys), "All key up events should be recorded")

	// Verify order is preserved
	for i, key := range keys {
		assert.Equal(t, key, backend.downKeys[i], "Key down order should be preserved")
		assert.Equal(t, key, backend.upKeys[i], "Key up order should be preserved")
	}
}
