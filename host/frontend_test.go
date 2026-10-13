package host

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/audio"
	"github.com/retroenv/retrogolib/gui"
)

type testBackend struct {
	Backend
}

type frontendFixture struct {
	stage   string
	failure error

	renders  int
	cleanups int

	starts      int
	stops       int
	audioErrors chan error
}

type pacingInterruption struct {
	name string

	timeout    time.Duration
	cancelWait bool
	failAudio  bool

	wantErr error
}

func (fx *frontendFixture) setupGUI(gui.Backend) (func() (bool, error), func(), error) {
	if fx.stage == "window" {
		return nil, nil, fx.failure
	}

	return fx.render, fx.cleanup, nil
}

func (fx *frontendFixture) setupAudio(audio.Backend) (*audio.Playback, error) {
	if fx.stage == "audio setup" {
		return nil, fx.failure
	}

	return &audio.Playback{
		Start:  fx.startAudio,
		Stop:   fx.stopAudio,
		Errors: fx.audioErrors,
	}, nil
}

func (fx *frontendFixture) render() (bool, error) {
	fx.renders++
	if fx.stage == "render" {
		return false, fx.failure
	}
	if fx.stage == "audio render" {
		fx.audioErrors <- fx.failure
	}

	return fx.stage != "close", nil
}

func (fx *frontendFixture) cleanup() {
	fx.cleanups++
}

func (fx *frontendFixture) startAudio() error {
	fx.starts++
	if fx.stage == "audio start" {
		return fx.failure
	}

	return nil
}

func (fx *frontendFixture) stopAudio() {
	fx.stops++
}

func (fx *frontendFixture) failAudio() {
	fx.audioErrors <- fx.failure
}

func (pi pacingInterruption) run(t *testing.T) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), pi.timeout)
	defer cancel()

	fixture := &frontendFixture{
		failure: pi.wantErr,

		audioErrors: make(chan error, 1),
	}
	front, err := New(testBackend{}, Options{
		FrameDuration: time.Hour,

		Audio: fixture.setupAudio,
		GUI:   fixture.setupGUI,
	})
	assert.NoError(t, err)
	defer front.Close()

	// Cancellation and audio errors must interrupt an active pacing wait.
	if pi.cancelWait {
		time.AfterFunc(time.Second, cancel)
	}
	if pi.failAudio {
		time.AfterFunc(time.Second, fixture.failAudio)
	}

	start := time.Now()
	open, err := front.Step(ctx, 1)
	assert.False(t, open)
	assert.ErrorIs(t, err, pi.wantErr)
	assert.Equal(t, time.Second, time.Since(start))
	assert.Equal(t, 0, fixture.renders)
	assert.Equal(t, 1, fixture.stops)
	assert.Equal(t, 1, fixture.cleanups)

	open, err = front.Step(ctx, 2)
	assert.False(t, open)
	assert.NoError(t, err)
	front.Close()
	assert.Equal(t, 1, fixture.stops)
	assert.Equal(t, 1, fixture.cleanups)
}

func TestHeadlessLimit(t *testing.T) {
	front, err := New(nil, Options{
		Headless: true,

		FrameLimit: 3,
	})
	assert.NoError(t, err)
	open, err := front.Step(t.Context(), 2)
	assert.NoError(t, err)
	assert.True(t, open)
	open, err = front.Step(t.Context(), 3)
	assert.NoError(t, err)
	assert.False(t, open)
	open, err = front.Step(t.Context(), 4)
	assert.NoError(t, err)
	assert.False(t, open)
	front.Close()
}

func TestFrontendCleanup(t *testing.T) {
	failure := errors.New("test frontend failure")

	for _, stage := range []string{"audio setup", "audio start", "audio worker", "audio render", "render", "close"} {
		t.Run(stage, func(t *testing.T) {
			fixture := &frontendFixture{
				stage:   stage,
				failure: failure,

				audioErrors: make(chan error, 1),
			}
			wantRenders := 1

			if stage == "audio worker" {
				wantRenders = 0
			}

			front, err := New(testBackend{}, Options{
				Audio: fixture.setupAudio,
				GUI:   fixture.setupGUI,
			})
			if stage == "audio setup" || stage == "audio start" {
				assert.ErrorIs(t, err, failure)
				assert.Nil(t, front)
			} else {
				assert.NoError(t, err)
				if stage == "audio worker" {
					fixture.audioErrors <- failure
				}
				open, stepErr := front.Step(t.Context(), 1)
				assert.False(t, open)
				if stage == "close" {
					assert.NoError(t, stepErr)
				} else {
					assert.ErrorIs(t, stepErr, failure)
				}
				front.Close()
				open, stepErr = front.Step(t.Context(), 2)
				assert.False(t, open)
				assert.NoError(t, stepErr)
				assert.Equal(t, wantRenders, fixture.renders)
			}
			assert.Equal(t, 1, fixture.cleanups)
			if stage == "audio setup" {
				assert.Equal(t, 0, fixture.stops)
			} else {
				assert.Equal(t, 1, fixture.stops)
			}
		})
	}
}

func TestFrontendWindowAndValidation(t *testing.T) {
	failure := errors.New("window unavailable")
	fixture := &frontendFixture{
		stage:   "window",
		failure: failure,
	}
	_, err := New(testBackend{}, Options{
		GUI: fixture.setupGUI,
	})
	assert.ErrorIs(t, err, failure)
	_, err = New(nil, Options{})
	assert.Error(t, err)
	_, err = New(nil, Options{
		Headless: true,

		FrameDuration: -time.Second,
	})
	assert.Error(t, err)
}

func TestEnvironment(t *testing.T) {
	t.Setenv("RETROHOST_HEADLESS", "")
	t.Setenv("RETROHOST_FRAMES", "")

	for _, value := range []string{"", "0", "false", "1", "true"} {
		t.Setenv("RETROHOST_HEADLESS", value)
		options, err := FromEnvironment(Options{})
		assert.NoError(t, err)
		assert.Equal(t, value == "1" || value == "true", options.Headless)
	}
	t.Setenv("RETROHOST_HEADLESS", "invalid")
	_, err := FromEnvironment(Options{})
	assert.Error(t, err)
	t.Setenv("RETROHOST_HEADLESS", "true")

	for _, value := range []string{"0", "-1", "invalid", "18446744073709551616"} {
		t.Setenv("RETROHOST_FRAMES", value)
		_, err = FromEnvironment(Options{})
		assert.Error(t, err)
	}
	t.Setenv("RETROHOST_FRAMES", "65")
	options, err := FromEnvironment(Options{})
	assert.NoError(t, err)
	assert.Equal(t, uint64(65), options.FrameLimit)
}

func TestWindowFrameLimit(t *testing.T) {
	fixture := &frontendFixture{}
	front, err := New(testBackend{}, Options{
		FrameLimit: 2,

		GUI: fixture.setupGUI,
	})
	assert.NoError(t, err)
	open, err := front.Step(t.Context(), 1)
	assert.NoError(t, err)
	assert.True(t, open)
	open, err = front.Step(t.Context(), 2)
	assert.NoError(t, err)
	assert.False(t, open)
	front.Close()
	assert.Equal(t, 1, fixture.cleanups)
}

func TestFrontendQueuedAudioError(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		failure := errors.New("audio worker failed")
		fixture := &frontendFixture{
			audioErrors: make(chan error, 1),
		}
		front, err := New(testBackend{}, Options{
			FrameDuration: time.Hour,

			Audio: fixture.setupAudio,
			GUI:   fixture.setupGUI,
		})
		assert.NoError(t, err)
		defer front.Close()

		// A queued audio error must stop the frontend before the pacing wait.
		fixture.audioErrors <- failure
		start := time.Now()
		open, err := front.Step(t.Context(), 1)
		assert.False(t, open)
		assert.ErrorIs(t, err, failure)
		assert.Equal(t, time.Duration(0), time.Since(start))
		assert.Equal(t, 0, fixture.renders)
		assert.Equal(t, 1, fixture.stops)
		assert.Equal(t, 1, fixture.cleanups)
	})
}

func TestFrontendPacingInterrupted(t *testing.T) {
	failure := errors.New("audio worker failed")
	tests := []pacingInterruption{
		{
			name: "cancellation",

			timeout:    time.Hour,
			cancelWait: true,

			wantErr: context.Canceled,
		},
		{
			name: "deadline",

			timeout: time.Second,

			wantErr: context.DeadlineExceeded,
		},
		{
			name: "audio error",

			timeout:   time.Hour,
			failAudio: true,

			wantErr: failure,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, test.run)
		})
	}
}

func TestFrontendCanceledContext(t *testing.T) {
	for _, headless := range []bool{false, true} {
		name := "window"

		if headless {
			name = "headless"
		}

		t.Run(name, func(t *testing.T) {
			fixture := &frontendFixture{}
			front, err := New(testBackend{}, Options{
				Headless: headless,

				Audio: fixture.setupAudio,
				GUI:   fixture.setupGUI,
			})
			assert.NoError(t, err)
			defer front.Close()

			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			open, err := front.Step(ctx, 1)
			assert.False(t, open)
			assert.ErrorIs(t, err, context.Canceled)
			assert.Equal(t, 0, fixture.renders)
			wantCleanups := 1

			if headless {
				wantCleanups = 0
			}

			assert.Equal(t, wantCleanups, fixture.stops)
			assert.Equal(t, wantCleanups, fixture.cleanups)
		})
	}
}

func TestFrontendPacing(t *testing.T) {
	for _, channel := range []string{"disabled", "nil", "open", "closed"} {
		t.Run(channel, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				fixture := &frontendFixture{}

				if channel == "open" || channel == "closed" {
					fixture.audioErrors = make(chan error)
				}
				if channel == "closed" {
					close(fixture.audioErrors)
				}

				options := Options{
					FrameDuration: 10 * time.Millisecond,

					GUI: fixture.setupGUI,
				}
				if channel != "disabled" {
					options.Audio = fixture.setupAudio
				}

				front, err := New(testBackend{}, options)
				assert.NoError(t, err)
				defer front.Close()

				for frame := uint64(1); frame <= 2; frame++ {
					start := time.Now()
					open, err := front.Step(t.Context(), frame)
					assert.True(t, open)
					assert.NoError(t, err)
					assert.Equal(t, options.FrameDuration, time.Since(start))
				}
				assert.Equal(t, 2, fixture.renders)
			})
		})
	}
}

func TestHeadlessSkipsDevicesAndPacing(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		fixture := &frontendFixture{
			stage:   "window",
			failure: errors.New("window unavailable"),
		}
		front, err := New(testBackend{}, Options{
			Headless: true,

			FrameDuration: time.Hour,

			Audio: fixture.setupAudio,
			GUI:   fixture.setupGUI,
		})
		assert.NoError(t, err)
		defer front.Close()

		start := time.Now()
		open, err := front.Step(t.Context(), 1)
		assert.True(t, open)
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(0), time.Since(start))
		front.Close()
		assert.Equal(t, 0, fixture.renders)
		assert.Equal(t, 0, fixture.starts)
		assert.Equal(t, 0, fixture.stops)
		assert.Equal(t, 0, fixture.cleanups)
	})
}
