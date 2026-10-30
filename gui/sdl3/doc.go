// Package sdl3 provides a CGO-free SDL3 GUI renderer for package gui.
// It requires an SDL3 shared library at runtime.
//
// Import this package for side effects to select gui.Setup, or call Setup
// directly. Import only one GUI renderer for side effects.
// Call Setup, render, and cleanup on the main goroutine.
// The renderer scales the frame with nearest-pixel sampling.
// Cleanup releases the video subsystem. Audio and GUI use separate subsystem
// references.
//
// To test the SDL2 and SDL3 backends against the installed shared libraries,
// run make test-sdl.
package sdl3
