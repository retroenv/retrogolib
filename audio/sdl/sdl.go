// Package sdl provides a SDL audio renderer.
package sdl

import (
	"fmt"
	"sync"
	"time"

	"github.com/retroenv/retrogolib/audio"
)

// Ensure Setup is registered with the audio package.
func init() {
	audio.Setup = Setup
}

// audioDevice holds the SDL audio device state.
type audioDevice struct {
	deviceID AudioDeviceID
	backend  audio.Backend
	buffer   []byte
	running  bool
	mutex    sync.RWMutex
}

// Setup initializes the SDL audio subsystem and returns control functions.
func Setup(backend audio.Backend) (audioStart func() error, audioStop func(), err error) {
	if err := setupLibrary(); err != nil {
		return nil, nil, fmt.Errorf("setting up SDL audio library: %w", err)
	}

	// Initialize SDL audio subsystem if not already initialized
	if ret := Init(SDL_INIT_AUDIO); ret != 0 {
		return nil, nil, fmt.Errorf("initializing SDL audio: %s", GetError())
	}

	format := backend.AudioFormat()
	device := &audioDevice{
		backend: backend,
		buffer:  make([]byte, format.BufferSize()),
	}

	// Convert audio format to SDL format
	sdlFormat, err := convertToSDLFormat(format.Format)
	if err != nil {
		return nil, nil, fmt.Errorf("unsupported audio format: %w", err)
	}

	// Configure audio specification
	desired := AudioSpec{
		Freq:     int32(format.SampleRate),
		Format:   sdlFormat,
		Channels: uint8(format.Channels),
		Samples:  uint16(format.Samples),
		Callback: 0, // Use queue-based audio (simpler than callbacks)
	}

	var obtained AudioSpec
	device.deviceID = OpenAudioDevice("", 0, &desired, &obtained, 0)
	if device.deviceID == 0 {
		return nil, nil, fmt.Errorf("opening audio device: %s", GetError())
	}

	start := func() error {
		return device.start()
	}

	stop := func() {
		device.stop()
	}

	return start, stop, nil
}

// start begins audio playback.
func (d *audioDevice) start() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.running {
		return nil
	}

	d.running = true
	PauseAudioDevice(d.deviceID, 0) // Unpause

	// Start audio pumping goroutine
	go d.audioLoop()

	return nil
}

// stop terminates audio playback and cleans up resources.
func (d *audioDevice) stop() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if !d.running {
		return
	}

	d.running = false
	ClearQueuedAudio(d.deviceID)
	CloseAudioDevice(d.deviceID)
}

// audioLoop pumps audio data from the backend to the SDL device.
func (d *audioDevice) audioLoop() {
	// Use a ticker for consistent timing
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		d.mutex.RLock()
		running := d.running
		d.mutex.RUnlock()

		if !running {
			return
		}

		d.processAudioFrame()
	}
}

// processAudioFrame handles one frame of audio processing.
func (d *audioDevice) processAudioFrame() {
	// Check if audio should be paused
	if d.backend.AudioPaused() {
		PauseAudioDevice(d.deviceID, 1) // Pause
		return
	}

	// Ensure audio is unpaused
	PauseAudioDevice(d.deviceID, 0) // Unpause

	// Check if we need to queue more audio
	// Keep a minimum amount of audio queued to prevent underruns
	queued := GetQueuedAudioSize(d.deviceID)
	bufferSize := uint32(len(d.buffer))

	// Queue audio if we have less than 2 buffers worth queued
	if queued < bufferSize*2 {
		// Get audio data from backend
		d.backend.AudioCallback(d.buffer)

		// Queue the audio data
		if ret := QueueAudio(d.deviceID, d.buffer, bufferSize); ret != 0 {
			// SDL error occurred, but continue trying
			return
		}
	}
}

// convertToSDLFormat converts audio.SampleFormat to SDL format.
func convertToSDLFormat(format audio.SampleFormat) (uint16, error) {
	switch format {
	case audio.FormatU8:
		return AUDIO_U8, nil
	case audio.FormatS8:
		return AUDIO_S8, nil
	case audio.FormatU16LSB: // FormatU16 is an alias for this
		return AUDIO_U16LSB, nil
	case audio.FormatS16LSB: // FormatS16 is an alias for this
		return AUDIO_S16LSB, nil
	case audio.FormatU16MSB:
		return AUDIO_U16MSB, nil
	case audio.FormatS16MSB:
		return AUDIO_S16MSB, nil
	case audio.FormatS32LSB: // FormatS32 is an alias for this
		return AUDIO_S32LSB, nil
	case audio.FormatS32MSB:
		return AUDIO_S32MSB, nil
	case audio.FormatF32LSB: // FormatF32 is an alias for this
		return AUDIO_F32LSB, nil
	case audio.FormatF32MSB:
		return AUDIO_F32MSB, nil
	default:
		return 0, fmt.Errorf("unsupported format: %d", format)
	}
}
