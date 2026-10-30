// Package sdl3 provides CGO-free SDL3 audio playback for package audio.
// It requires an SDL3 shared library at runtime.
//
// Import this package for side effects to register Setup with audio.Register at
// a higher priority than the SDL2 renderer, or call Setup directly. If this
// renderer fails, audio.Setup tries the next registered renderer.
// Samples sets the callback block size and need not be a power of two.
// SDL converts the input format to the device format. Unsigned 16-bit PCM
// is not supported. Start fills the input queue before playback starts.
// Stop closes the stream and its device, then releases the audio subsystem.
// Audio and GUI use separate subsystem references.
package sdl3
