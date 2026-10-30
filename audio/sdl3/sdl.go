package sdl3

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/retroenv/retrogolib/audio"
	sdllibrary "github.com/retroenv/retrogolib/internal/sdl3"
)

// Setup opens the default output device without starting playback.
// Samples sets the callback block size. SDL selects the device buffer size.
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
	stream      uintptr
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
		Start:  d.start,
		Stop:   d.stop,
		Errors: d.errors,
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
	DestroyAudioStream(d.stream)
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
			if !PauseAudioStreamDevice(d.stream) {
				return fmt.Errorf("pausing SDL audio: %s", GetError())
			}
			d.paused = true
		}

		return nil
	}
	if err := d.refill(); err != nil {
		return err
	}
	if d.paused && !d.stopping() {
		if !ResumeAudioStreamDevice(d.stream) {
			return fmt.Errorf("resuming SDL audio: %s", GetError())
		}
		d.paused = false
	}

	return nil
}

func (d *audioDevice) refill() error {
	size := GetAudioStreamQueued(d.stream)
	if size < 0 {
		return fmt.Errorf("querying SDL audio: %s", GetError())
	}
	queued := uint32(size)
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
		if !PutAudioStreamData(d.stream, d.buffer, int32(bufferSize)) {
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
		Freq:     int32(format.SampleRate),
		Format:   uint32(format.Format),
		Channels: int32(format.Channels),
	}
	stream := OpenAudioDeviceStream(SDL_AUDIO_DEVICE_DEFAULT_PLAYBACK, &desired, 0, 0)
	if stream == 0 {
		return nil, fmt.Errorf("opening audio stream: %s", GetError())
	}
	// Keep at least 20 ms and two callback blocks in the input queue.
	frames := max(int64(format.SampleRate)/50, 2*int64(format.Samples))
	interval := time.Second * time.Duration(format.Samples) / time.Duration(format.SampleRate) / 2

	return &audioDevice{
		stream:      stream,
		backend:     backend,
		buffer:      make([]byte, format.BufferSize()),
		targetBytes: uint32(frames * int64(format.Channels*format.BytesPerSample())),
		interval:    max(time.Millisecond, min(5*time.Millisecond, interval)),
		paused:      true,
		release:     release,
	}, nil
}

func validateFormat(format audio.Format) error {
	if err := format.Validate(); err != nil {
		return fmt.Errorf("invalid PCM format: %w", err)
	}
	if format.Format == audio.FormatU16LSB || format.Format == audio.FormatU16MSB {
		return errors.New("SDL3 does not support unsigned 16-bit PCM")
	}
	if int64(format.BufferSize()) > math.MaxInt32/2 {
		return errors.New("SDL3 audio buffer size is too large")
	}
	return nil
}
