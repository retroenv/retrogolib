package sdl3

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/audio"
)

func TestSetupRejectsInvalidFormats(t *testing.T) {
	assertInvalid := func(t *testing.T, backend audio.Backend) {
		t.Helper()
		playback, err := Setup(backend)
		assert.Error(t, err)
		assert.Nil(t, playback)
	}
	assertInvalid(t, nil)
	for _, format := range []audio.Format{
		{SampleRate: 44100, Channels: 1, Samples: -1, Format: audio.FormatS16},
		{SampleRate: 44100, Channels: 1, Samples: 0, Format: audio.FormatS16},
		{SampleRate: 44100, Channels: 257, Samples: 512, Format: audio.FormatS16},
		{SampleRate: 44100, Channels: 1, Samples: 512, Format: audio.FormatU16LSB},
		{SampleRate: 44100, Channels: 1, Samples: 512, Format: audio.FormatU16MSB},
		{SampleRate: 0, Channels: 1, Samples: 512, Format: audio.FormatS16},
		{SampleRate: 44100, Channels: 1, Samples: 512, Format: 0xffff},
	} {
		assertInvalid(t, &testBackend{format: format})
	}
}

func TestPlaybackLifecycle(t *testing.T) {
	b := newTestBackend()
	playback, _, native := newTestPlayback(t, b)
	assert.Equal(t, int64(0), b.calls.Load())
	assert.NoError(t, playback.Start())
	assert.NoError(t, playback.Start())
	playback.Stop()
	playback.Stop()
	assert.ErrorIs(t, playback.Start(), audio.ErrClosed)
	assert.Equal(t, []string{"open", "queue", "queue", "resume", "close", "release"}, native.events)
	assert.Equal(t, int64(2), b.calls.Load())
	_, ok := <-playback.Errors
	assert.False(t, ok)
}

func TestStopBeforeStart(t *testing.T) {
	playback, _, native := newTestPlayback(t, newTestBackend())
	playback.Stop()
	assert.Equal(t, []string{"open", "close", "release"}, native.events)
	assert.ErrorIs(t, playback.Start(), audio.ErrClosed)
}

func TestStopWaitsForCallback(t *testing.T) {
	b := newTestBackend()
	entered := make(chan struct{})
	proceed := make(chan struct{})
	var unblock sync.Once
	b.callback = func([]byte) { close(entered); <-proceed }
	playback, device, native := newTestPlayback(t, b)
	t.Cleanup(func() { unblock.Do(func() { close(proceed) }) })
	started := make(chan error, 1)
	go func() { started <- playback.Start() }()
	awaitSignal(t, entered)
	stopped := make(chan struct{})
	go func() { playback.Stop(); close(stopped) }()
	awaitSignal(t, device.stopRequested)
	select {
	case <-stopped:
		t.Fatal("Stop returned while the callback was active")
	default:
	}
	unblock.Do(func() { close(proceed) })
	awaitSignal(t, stopped)
	err := <-started
	assert.True(t, err == nil || errors.Is(err, audio.ErrClosed))
	assert.Equal(t, []string{"open", "close", "release"}, native.events)
}

func TestSmallBuffersPrefillBeforeResume(t *testing.T) {
	b := newTestBackend()
	b.format.SampleRate = 48000
	b.format.Samples = 256
	playback, device, native := newTestPlayback(t, b)
	assert.NoError(t, playback.Start())
	playback.Stop()
	assert.Equal(t, int64(4), b.calls.Load())
	assert.True(t, native.queued >= 48000*2*2/50)
	assert.True(t, device.interval <= 3*time.Millisecond)
	assert.Equal(t, []string{"open", "queue", "queue", "queue", "queue", "resume", "close", "release"}, native.events)
}

func TestPauseResume(t *testing.T) {
	b := newTestBackend()
	b.paused.Store(true)
	playback, _, native := newTestPlayback(t, b)
	assert.NoError(t, playback.Start())
	assert.Equal(t, int64(0), b.calls.Load())
	b.paused.Store(false)
	native.awaitEvent(t, "resume")
	b.paused.Store(true)
	native.awaitEvent(t, "pause")
	b.paused.Store(false)
	native.awaitEvent(t, "resume")
	playback.Stop()
	assert.Equal(t, int64(2), b.calls.Load())
}

func TestQueueFailure(t *testing.T) {
	for _, duringStart := range []bool{true, false} {
		name := "during playback"
		if duringStart {
			name = "during start"
		}
		t.Run(name, func(t *testing.T) {
			playback, _, native := newTestPlayback(t, newTestBackend())
			if !duringStart {
				assert.NoError(t, playback.Start())
			}
			native.mutex.Lock()
			native.failQueue = true
			native.queued = 0
			native.mutex.Unlock()
			if duringStart {
				assert.Error(t, playback.Start())
			}
			select {
			case err := <-playback.Errors:
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "queue failed")

			case <-time.After(time.Second):
				t.Fatal("playback did not report its failure")
			}
			playback.Stop()
			assert.Error(t, playback.Start())
			assert.Equal(t, []string{"close", "release"}, native.events[len(native.events)-2:])
			_, ok := <-playback.Errors
			assert.False(t, ok)
		})
	}
}

func TestConcurrentControls(t *testing.T) {
	playback, _, native := newTestPlayback(t, newTestBackend())
	var workers sync.WaitGroup
	results := make(chan error, 20)
	for range 20 {
		workers.Go(func() {
			results <- playback.Start()
			playback.Stop()
		})
	}
	workers.Wait()
	close(results)
	for err := range results {
		assert.True(t, err == nil || errors.Is(err, audio.ErrClosed))
	}
	closes := 0
	for _, event := range native.events {
		if event == "close" {
			closes++
		}
	}
	assert.Equal(t, 1, closes)
}

func TestOpenDeviceFailure(t *testing.T) {
	native := mockAudio(t)
	native.failOpen = true
	b := newTestBackend()
	device, err := openDevice(b, b.format, func() {})
	assert.Nil(t, device)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "queue failed")
}

type testBackend struct {
	format   audio.Format
	paused   atomic.Bool
	calls    atomic.Int64
	callback func([]byte)
}

func (b *testBackend) AudioFormat() audio.Format { return b.format }
func (b *testBackend) AudioPaused() bool         { return b.paused.Load() }
func (b *testBackend) AudioCallback(buffer []byte) {
	b.calls.Add(1)
	if b.callback != nil {
		b.callback(buffer)
	}
}

type mockSDL struct {
	mutex         sync.Mutex
	events        []string
	notifications chan string
	queued        uint32
	failQueue     bool
	failOpen      bool
}

func (s *mockSDL) event(event string) {
	s.events = append(s.events, event)
	s.notifications <- event
}

func (s *mockSDL) awaitEvent(t *testing.T, event string) {
	t.Helper()
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	for {
		select {
		case got := <-s.notifications:
			if got == event {
				return
			}

		case <-timer.C:
			t.Fatalf("timed out waiting for %s", event)
		}
	}
}

func newTestBackend() *testBackend {
	return &testBackend{format: audio.Format{
		SampleRate: 44100,
		Channels:   2,
		Samples:    512,
		Format:     audio.FormatS16,
	}}
}

func newTestPlayback(t *testing.T, b *testBackend) (*audio.Playback, *audioDevice, *mockSDL) {
	t.Helper()
	native := mockAudio(t)
	d, err := openDevice(b, b.format, func() {
		native.mutex.Lock()
		defer native.mutex.Unlock()
		native.event("release")
	})
	assert.NoError(t, err)
	playback := d.playback()
	t.Cleanup(playback.Stop)
	return playback, d, native
}

func mockAudio(t *testing.T) *mockSDL {
	t.Helper()
	open := OpenAudioDeviceStream
	closeDevice := DestroyAudioStream
	pause := PauseAudioStreamDevice
	resume := ResumeAudioStreamDevice
	queue := PutAudioStreamData
	size := GetAudioStreamQueued
	getError := GetError
	t.Cleanup(func() {
		OpenAudioDeviceStream = open
		DestroyAudioStream = closeDevice
		PauseAudioStreamDevice = pause
		ResumeAudioStreamDevice = resume
		PutAudioStreamData = queue
		GetAudioStreamQueued = size
		GetError = getError
	})
	s := &mockSDL{notifications: make(chan string, 100)}
	mockOpenAudio(t, s)
	mockAudioQueue(s)
	return s
}

func mockOpenAudio(t *testing.T, s *mockSDL) {
	t.Helper()
	OpenAudioDeviceStream = func(device uint32, desired *AudioSpec, callback, userdata uintptr) uintptr {
		assert.Equal(t, uint32(0xffffffff), device)
		assert.Equal(t, uint32(audio.FormatS16), desired.Format)
		assert.Equal(t, int32(2), desired.Channels)
		assert.Equal(t, uintptr(0), callback)
		assert.Equal(t, uintptr(0), userdata)
		s.event("open")
		if s.failOpen {
			return 0
		}
		return 42
	}
}

func mockAudioQueue(s *mockSDL) {
	DestroyAudioStream = func(uintptr) {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.event("close")
	}
	PauseAudioStreamDevice = func(uintptr) bool {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.event("pause")
		return true
	}
	ResumeAudioStreamDevice = func(uintptr) bool {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.event("resume")
		return true
	}
	PutAudioStreamData = func(_ uintptr, _ []byte, size int32) bool {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		if s.failQueue {
			return false
		}
		s.queued += uint32(size)
		s.event("queue")
		return true
	}
	GetAudioStreamQueued = func(uintptr) int32 {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		return int32(s.queued)
	}
	GetError = func() string { return "queue failed" }
}

func awaitSignal(t *testing.T, signal <-chan struct{}) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for worker")
	}
}

func TestStreamControlFailures(t *testing.T) {
	for _, stage := range []string{"query", "pause", "resume"} {
		t.Run(stage, func(t *testing.T) {
			b := newTestBackend()
			playback, device, native := newTestPlayback(t, b)
			switch stage {
			case "query":
				GetAudioStreamQueued = func(uintptr) int32 { return -1 }

			case "pause":
				b.paused.Store(true)
				device.paused = false
				PauseAudioStreamDevice = func(uintptr) bool { return false }

			case "resume":
				ResumeAudioStreamDevice = func(uintptr) bool { return false }
			}
			assert.Error(t, playback.Start())
			playback.Stop()
			assert.Error(t, <-playback.Errors)
			assert.Equal(t, []string{"close", "release"}, native.events[len(native.events)-2:])
		})
	}
}

func TestNonPowerOfTwoBlock(t *testing.T) {
	b := newTestBackend()
	b.format.Samples = 513
	assert.NoError(t, validateFormat(b.format))
	playback, _, _ := newTestPlayback(t, b)
	assert.NoError(t, playback.Start())
	playback.Stop()
	assert.Equal(t, int64(2), b.calls.Load())
}
