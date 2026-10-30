// Package audio provides CGO-free PCM playback for retro console emulators.
//
// Import a renderer to register it with audio.Register. Renderers are separate
// packages and may impose additional runtime dependencies and audio format
// constraints. audio.Setup tries the registered renderers in priority order and
// returns the first one that starts.
//
// # Usage with SDL3 and SDL2
//
// Import both renderers to use SDL3 and fall back to SDL2, for example when the
// SDL3 library is not installed.
//
//	import (
//	    "github.com/retroenv/retrogolib/audio"
//	    _ "github.com/retroenv/retrogolib/audio/sdl2"
//	    _ "github.com/retroenv/retrogolib/audio/sdl3"
//	)
//
//	playback, err := audio.Setup(emulator)
//	if err != nil {
//	    return err
//	}
//	defer playback.Stop()
//	if err := playback.Start(); err != nil {
//	    return err
//	}
//
// The application should monitor playback.Errors in its event loop. A received
// error is terminal: the renderer stops playback and releases its resources.
// The channel is buffered and closes after cleanup, with no value on normal Stop.
// Stop is still safe and required after an error, and waits for all callbacks.
// Start is idempotent while playing; playback cannot restart after Stop.
//
// # Backend contract
//
// AudioFormat is read once during Setup. SampleRate is in sample frames per
// second; Samples is the number of frames supplied per AudioCallback. A stereo
// frame contains both a left and a right sample. PCM channels are interleaved.
// FormatS16, FormatS32, FormatF32 and FormatU16 always mean little-endian.
//
// AudioCallback and AudioPaused are called serially by a worker goroutine and
// can overlap the application's emulation and GUI code. Protect shared data with
// synchronization, such as a mutex or atomic.Bool for a pause flag. Callbacks
// must return promptly and must not call Start or Stop themselves.
//
// AudioCallback must fill its entire borrowed buffer without retaining it.
// Supply silence when source data runs out: zero for signed/float formats, 0x80
// for unsigned 8-bit, and 0x8000 (in the selected byte order) for unsigned 16-bit.
// Pausing preserves already queued audio. Stop discards queued audio.
package audio
