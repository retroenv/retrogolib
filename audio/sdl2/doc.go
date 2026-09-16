// Package sdl2 provides a CGO-free SDL2 audio renderer for package audio.
// It requires an SDL2 shared library at runtime, but no SDL headers or C compiler
// are needed to build.
//
// Importing this package registers Setup as audio.Setup. Applications may call
// Setup directly instead when they do not need renderer registration.
//
// Samples must be a power of two between 1 and 32768. It controls the callback
// chunk size and requests a hardware buffer size; actual latency also includes
// SDL and operating-system buffering. The renderer prefills at least 20ms and
// two hardware/callback buffers, then refills according to the chunk duration.
//
// SDL2 audio and GUI own their subsystems separately. Either can be cleaned up
// while the other remains active. Applications using SDL directly must not call
// global SDL_Quit until all library playback and GUI controls have been closed.
package sdl2
