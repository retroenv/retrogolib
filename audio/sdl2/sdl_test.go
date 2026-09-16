package sdl2

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
		{SampleRate: 44100, Channels: 1, Samples: 65536, Format: audio.FormatS16},
		{SampleRate: 44100, Channels: 1, Samples: 513, Format: audio.FormatS16},
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
	entered, proceed := make(chan struct{}), make(chan struct{})
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

func TestObtainedFormat(t *testing.T) {
	native := mockAudio(t)
	b := newTestBackend()
	native.samples = 2048
	device, err := openDevice(b, b.format, func() {})
	assert.NoError(t, err)
	assert.Equal(t, uint32(4096*4), device.targetBytes)
	native.incompatible = true
	_, err = openDevice(b, b.format, func() {})
	assert.Error(t, err)
	assert.Equal(t, []string{"open", "open", "close"}, native.events)
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
	incompatible  bool
	samples       uint16
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
		SampleRate: 44100, Channels: 2, Samples: 512, Format: audio.FormatS16,
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
	open, closeDevice, pause := OpenAudioDevice, CloseAudioDevice, PauseAudioDevice
	queue, size, getError := QueueAudio, GetQueuedAudioSize, GetError
	t.Cleanup(func() {
		OpenAudioDevice, CloseAudioDevice, PauseAudioDevice = open, closeDevice, pause
		QueueAudio, GetQueuedAudioSize, GetError = queue, size, getError
	})
	s := &mockSDL{notifications: make(chan string, 100)}
	mockOpenAudio(t, s)
	mockAudioQueue(s)
	return s
}

func mockOpenAudio(t *testing.T, s *mockSDL) {
	t.Helper()
	OpenAudioDevice = func(name *byte, capture int, desired, obtained *AudioSpec, changes int) AudioDeviceID {
		assert.Nil(t, name)
		assert.Equal(t, 0, capture)
		assert.Equal(t, 0, changes)
		s.event("open")
		if s.failOpen {
			return 0
		}
		*obtained = *desired
		if s.samples != 0 {
			obtained.Samples = s.samples
		}
		if s.incompatible {
			obtained.Freq++
		}
		return 42
	}
}

func mockAudioQueue(s *mockSDL) {
	CloseAudioDevice = func(AudioDeviceID) {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.event("close")
	}
	PauseAudioDevice = func(_ AudioDeviceID, pause int) {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		if pause == 0 {
			s.event("resume")
		} else {
			s.event("pause")
		}
	}
	QueueAudio = func(_ AudioDeviceID, _ []byte, size uint32) int {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		if s.failQueue {
			return -1
		}
		s.queued += size
		s.event("queue")
		return 0
	}
	GetQueuedAudioSize = func(AudioDeviceID) uint32 {
		s.mutex.Lock()
		defer s.mutex.Unlock()
		return s.queued
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
