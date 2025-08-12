// Package audio implements different audio renderers for retro console emulation.
//
// The audio package provides a CGO-free interface for audio playback that can be
// used by emulators and retro console tooling. The design follows the same
// modular pattern as the GUI package, allowing different audio backends while
// maintaining a consistent API.
//
// Usage:
//
//	import (
//	    "github.com/retroenv/retrogolib/audio"
//	    _ "github.com/retroenv/retrogolib/audio/sdl" // Register SDL audio backend
//	)
//
//	// Implement the Backend interface in your emulator
//	type Emulator struct {
//	    audioBuffer []int16
//	    paused      bool
//	}
//
//	func (e *Emulator) AudioFormat() audio.Format {
//	    return audio.Format{
//	        SampleRate: 44100,
//	        Channels:   2,
//	        Samples:    512,
//	        Format:     audio.FormatS16,
//	    }
//	}
//
//	func (e *Emulator) AudioCallback(buffer []byte) {
//	    // Fill buffer with audio samples from emulator
//	}
//
//	func (e *Emulator) AudioPaused() bool {
//	    return e.paused
//	}
//
//	// Set up audio
//	audioStart, audioStop, err := audio.Setup(emulator)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer audioStop()
//
//	if err := audioStart(); err != nil {
//	    log.Fatal(err)
//	}
//
// The package is designed to work alongside the GUI package for complete
// multimedia emulation support.
package audio
