package sdl2

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/retroenv/retrogolib/audio"
	sdllibrary "github.com/retroenv/retrogolib/internal/sdl2"
)

// Setup opens the default output device without starting playback.
// Samples must be a power of two between 1 and 32768 sample frames.
func Setup(backend audio.Backend) (*audio.Playback, error) {
	if backend == nil {
		return nil, errors.New("audio backend is nil")
	}
	format := backend.AudioFormat()
	if err := validateFormat(format); err != nil {
		return nil, fmt.Errorf("validating audio format: %w", err)
	}
	if err := setupLibrary(); err != nil {
		return nil, fmt.Errorf("setting up SDL audio library: %w", err)
	}

	release, err := sdllibrary.Acquire(SDL_INIT_AUDIO)
	if err != nil {
		return nil, fmt.Errorf("acquiring SDL audio: %w", err)
	}
	device, err := openDevice(backend, format, release)
	if err != nil {
		release()
		return nil, err
	}
	return device.playback(), nil
}

type audioDevice struct {
	deviceID    AudioDeviceID
	backend     audio.Backend
	buffer      []byte
	targetBytes uint32
	interval    time.Duration
	paused      bool
	release     func()

	startRequests chan chan error
	stopRequested chan struct{}
	done          chan struct{}
	errors        chan error
	stopOnce      sync.Once
	terminalError error // Read only after done is closed.
}

func init() {
	audio.Setup = Setup
}

func (d *audioDevice) playback() *audio.Playback {
	d.startRequests = make(chan chan error)
	d.stopRequested = make(chan struct{})
	d.done = make(chan struct{})
	d.errors = make(chan error, 1)
	go d.run()
	return &audio.Playback{
		Start: d.start, Stop: d.stop, Errors: d.errors,
	}
}

func (d *audioDevice) start() error {
	reply := make(chan error, 1)
	select {
	case <-d.done:
		return d.startError()
	case d.startRequests <- reply:
	}
	select {
	case err := <-reply:
		return err
	case <-d.done:
		return d.startError()
	}
}

func (d *audioDevice) startError() error {
	if d.terminalError != nil {
		return d.terminalError
	}
	return audio.ErrClosed
}

func (d *audioDevice) stop() {
	d.stopOnce.Do(func() { close(d.stopRequested) })
	<-d.done
}

func (d *audioDevice) stopping() bool {
	select {
	case <-d.stopRequested:
		return true
	default:
		return false
	}
}

func (d *audioDevice) run() {
	d.terminalError = d.pump()
	CloseAudioDevice(d.deviceID)
	d.release()
	if d.terminalError != nil {
		d.errors <- d.terminalError
	}
	close(d.errors)
	close(d.done)
}

func (d *audioDevice) pump() error {
	var ticker *time.Ticker
	var ticks <-chan time.Time
	defer func() {
		if ticker != nil {
			ticker.Stop()
		}
	}()
	for {
		if d.stopping() {
			return nil
		}
		select {
		case <-d.stopRequested:
			return nil
		case reply := <-d.startRequests:
			if ticks != nil {
				reply <- nil
				continue
			}
			err := d.processAudioFrame()
			reply <- err
			if err != nil {
				return err
			}
			ticker = time.NewTicker(d.interval)
			ticks = ticker.C
		case <-ticks:
			if err := d.processAudioFrame(); err != nil {
				return err
			}
		}
	}
}

func (d *audioDevice) processAudioFrame() error {
	if d.stopping() {
		return nil
	}
	// SDL error messages belong to the OS thread that made the failing call.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if d.backend.AudioPaused() {
		if !d.paused {
			PauseAudioDevice(d.deviceID, 1)
			d.paused = true
		}
		return nil
	}
	if err := d.refill(); err != nil {
		return err
	}
	if d.paused && !d.stopping() {
		PauseAudioDevice(d.deviceID, 0)
		d.paused = false
	}
	return nil
}

func (d *audioDevice) refill() error {
	queued := GetQueuedAudioSize(d.deviceID)
	bufferSize := uint32(len(d.buffer))
	// Use a snapshot to bound the work even if the device drains while we fill.
	for queued < d.targetBytes {
		if d.stopping() {
			return nil
		}
		d.backend.AudioCallback(d.buffer)
		if d.stopping() {
			return nil
		}
		if QueueAudio(d.deviceID, d.buffer, bufferSize) != 0 {
			return fmt.Errorf("queueing SDL audio: %s", GetError())
		}
		queued += bufferSize
	}
	return nil
}

func openDevice(backend audio.Backend, format audio.Format, release func()) (*audioDevice, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	desired := AudioSpec{
		Freq: int32(format.SampleRate), Format: uint16(format.Format),
		Channels: uint8(format.Channels), Samples: uint16(format.Samples),
	}
	var obtained AudioSpec
	id := OpenAudioDevice(nil, 0, &desired, &obtained, 0)
	if id == 0 {
		return nil, fmt.Errorf("opening audio device: %s", GetError())
	}
	if obtained.Freq != desired.Freq || obtained.Format != desired.Format ||
		obtained.Channels != desired.Channels || obtained.Samples == 0 {

		CloseAudioDevice(id)
		return nil, errors.New("SDL returned an incompatible audio format")
	}

	// Cover at least 20ms and two hardware/callback buffers, including small blocks.
	frames := max(int64(format.SampleRate)/50, 2*int64(max(format.Samples, int(obtained.Samples))))
	interval := time.Second * time.Duration(format.Samples) / time.Duration(format.SampleRate) / 2
	return &audioDevice{
		deviceID: id, backend: backend, buffer: make([]byte, format.BufferSize()),
		targetBytes: uint32(frames * int64(format.Channels*format.BytesPerSample())),
		interval:    max(time.Millisecond, min(5*time.Millisecond, interval)),
		paused:      true, release: release,
	}, nil
}

func validateFormat(format audio.Format) error {
	if err := format.Validate(); err != nil {
		return fmt.Errorf("invalid PCM format: %w", err)
	}
	if format.Samples > 32768 || format.Samples&(format.Samples-1) != 0 {
		return fmt.Errorf("SDL2 sample frames must be a power of two between 1 and 32768, got %d", format.Samples)
	}
	return nil
}
