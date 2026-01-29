// Package gui provides cross-platform GUI rendering capabilities for retro console emulators.
//
// This package defines interfaces and abstractions for GUI rendering while avoiding
// CGO dependencies. The actual rendering implementation is provided by backend packages
// that implement the required interfaces.
//
// # Architecture
//
// The GUI package uses a backend interface pattern:
//   - Backend interface: Defines what the GUI system needs from the application
//   - Initializer function: Sets up the GUI renderer and returns control functions
//   - Dimensions type: Manages window sizing and scaling
//
// # Backend Interface
//
// Applications implement the Backend interface to provide:
//   - Image data for rendering
//   - Window dimensions and scaling
//   - Window title
//   - Input event handling
//
// # Basic Usage
//
//	type MyEmulator struct {
//		display *image.RGBA
//		// ... other fields
//	}
//
//	func (e *MyEmulator) Image() *image.RGBA {
//		return e.display
//	}
//
//	func (e *MyEmulator) Dimensions() gui.Dimensions {
//		return gui.Dimensions{
//			Width:       256,
//			Height:      240,
//			ScaleFactor: 2.0,
//		}
//	}
//
//	func (e *MyEmulator) WindowTitle() string {
//		return "My Retro Emulator"
//	}
//
//	// ... implement other Backend methods
//
//	// Start GUI
//	render, cleanup, err := gui.Setup(emulator)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer cleanup()
//
//	for {
//		running, err := render()
//		if err != nil {
//			log.Fatal(err)
//		}
//		if !running {
//			break
//		}
//	}
//
// # CGO-Free Operation
//
// This package achieves CGO-free operation by using pure Go libraries like
// ebitengine/purego for system interaction, making cross-compilation easier
// and reducing external dependencies.
//
// # Input Handling
//
// Input events are delivered through the Backend interface methods KeyDown
// and KeyUp, using key codes defined in the input package.
package gui
