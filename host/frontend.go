// Package host supplies a frame-driven desktop frontend for console backends.
package host

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/retroenv/retrogolib/audio"
	"github.com/retroenv/retrogolib/gui"
)

// Backend supplies video, input, and audio for a console.
type Backend interface {
	audio.Backend
	gui.Backend
}

// Options selects the frontend policy.
type Options struct {
	// Headless disables window creation, audio playback, and pacing.
	Headless bool

	// FrameLimit sets the maximum completed frame count. Zero disables the limit.
	FrameLimit uint64
	// FrameDuration sets the minimum presentation interval. Zero disables pacing.
	FrameDuration time.Duration

	Audio audio.Initializer
	GUI   gui.Initializer
}

// Frontend presents frames on the caller's goroutine.
// Call New, Step, and Close on the same goroutine. SDL2 requires the main goroutine.
// Calls must not overlap. Close the frontend on all exit paths.
// Audio callbacks run on the audio worker and can overlap frontend calls.
type Frontend struct {
	options Options
	closed  bool

	render  func() (bool, error)
	cleanup func()

	playback *audio.Playback

	last time.Time
}

// New opens the selected frontend. Headless mode opens no desktop devices.
func New(backend Backend, options Options) (*Frontend, error) {
	if options.FrameDuration < 0 {
		return nil, errors.New("frame duration must not be negative")
	}

	front := &Frontend{
		options: options,
	}

	if options.Headless {
		return front, nil
	}
	if backend == nil || options.GUI == nil {
		return nil, errors.New("window backend and initializer are required")
	}

	var err error
	front.render, front.cleanup, err = options.GUI(backend)
	if err != nil {
		return nil, fmt.Errorf("opening window: %w", err)
	}

	if err := front.startAudio(backend); err != nil {
		front.Close()
		return nil, err
	}

	front.last = time.Now()
	return front, nil
}

// Step presents a completed frame. False means that the caller must stop.
// Pass the total number of completed frames, starting at 1.
// Context cancellation or an audio error interrupts the pacing wait.
// Step closes the frontend on completion, cancellation, or failure.
// Calls after closure return false with no error. Step never exits the process.
func (fr *Frontend) Step(ctx context.Context, frame uint64) (bool, error) {
	if fr.closed {
		return false, nil
	}
	if err := ctx.Err(); err != nil {
		return fr.finish(false, fmt.Errorf("presenting frame: %w", err))
	}
	if fr.options.Headless {
		return fr.finish(fr.options.FrameLimit == 0 || frame < fr.options.FrameLimit, nil)
	}
	if err := fr.audioError(); err != nil {
		return fr.finish(false, err)
	}
	if err := fr.waitFrame(ctx); err != nil {
		return fr.finish(false, err)
	}

	fr.last = time.Now()
	open, err := fr.render()
	if err == nil {
		err = fr.audioError()
	}
	if fr.options.FrameLimit > 0 && frame >= fr.options.FrameLimit {
		open = false
	}
	return fr.finish(open, err)
}

// Close stops audio and closes the window. Repeated calls have no effect.
func (fr *Frontend) Close() {
	if fr.closed {
		return
	}
	fr.closed = true

	if fr.playback != nil {
		fr.playback.Stop()
	}
	if fr.cleanup != nil {
		fr.cleanup()
	}
}

func (fr *Frontend) startAudio(backend Backend) error {
	if fr.options.Audio == nil {
		return nil
	}

	var err error
	fr.playback, err = fr.options.Audio(backend)
	if err != nil {
		return fmt.Errorf("opening audio: %w", err)
	}

	if err := fr.playback.Start(); err != nil {
		return fmt.Errorf("starting audio: %w", err)
	}
	return nil
}

func (fr *Frontend) finish(open bool, err error) (bool, error) {
	if err != nil || !open {
		fr.Close()
		return false, err
	}
	return true, nil
}

func (fr *Frontend) audioError() error {
	if fr.playback == nil {
		return nil
	}

	select {
	case err := <-fr.playback.Errors:
		if err != nil {
			return fmt.Errorf("audio playback: %w", err)
		}

	default:
	}
	return nil
}

func (fr *Frontend) waitFrame(ctx context.Context) error {
	remaining := fr.options.FrameDuration - time.Since(fr.last)

	if remaining <= 0 {
		return nil
	}

	var audioErrors <-chan error

	if fr.playback != nil {
		audioErrors = fr.playback.Errors
	}

	timer := time.NewTimer(remaining)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("waiting for frame: %w", ctx.Err())

		case err := <-audioErrors:
			if err != nil {
				return fmt.Errorf("audio playback: %w", err)
			}
			audioErrors = nil

		case <-timer.C:
			return nil
		}
	}
}
