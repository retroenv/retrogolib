package sdl

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/audio"
)

func TestConvertToSDLFormat(t *testing.T) {
	tests := []struct {
		name        string
		format      audio.SampleFormat
		expected    uint16
		expectError bool
	}{
		{"8-bit unsigned", audio.FormatU8, AUDIO_U8, false},
		{"8-bit signed", audio.FormatS8, AUDIO_S8, false},
		{"16-bit little-endian unsigned", audio.FormatU16LSB, AUDIO_U16LSB, false},
		{"16-bit little-endian signed", audio.FormatS16LSB, AUDIO_S16LSB, false},
		{"16-bit big-endian unsigned", audio.FormatU16MSB, AUDIO_U16MSB, false},
		{"16-bit big-endian signed", audio.FormatS16MSB, AUDIO_S16MSB, false},
		{"32-bit little-endian signed", audio.FormatS32LSB, AUDIO_S32LSB, false},
		{"32-bit big-endian signed", audio.FormatS32MSB, AUDIO_S32MSB, false},
		{"32-bit little-endian float", audio.FormatF32LSB, AUDIO_F32LSB, false},
		{"32-bit big-endian float", audio.FormatF32MSB, AUDIO_F32MSB, false},
		{"16-bit platform default unsigned", audio.FormatU16, AUDIO_U16LSB, false},
		{"16-bit platform default signed", audio.FormatS16, AUDIO_S16LSB, false},
		{"32-bit platform default signed", audio.FormatS32, AUDIO_S32LSB, false},
		{"32-bit platform default float", audio.FormatF32, AUDIO_F32LSB, false},
		{"unsupported format", audio.SampleFormat(0xFFFF), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertToSDLFormat(tt.format)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestAudioSpec(t *testing.T) {
	spec := AudioSpec{
		Freq:     44100,
		Format:   AUDIO_S16SYS,
		Channels: 2,
		Samples:  512,
		Callback: 0,
	}

	assert.Equal(t, int32(44100), spec.Freq)
	assert.Equal(t, uint16(AUDIO_S16SYS), spec.Format)
	assert.Equal(t, uint8(2), spec.Channels)
	assert.Equal(t, uint16(512), spec.Samples)
	assert.Equal(t, uintptr(0), spec.Callback)
}

func TestAudioDeviceID(t *testing.T) {
	var id AudioDeviceID = 42
	assert.Equal(t, AudioDeviceID(42), id)
}

// MockAudioBackend implements audio.Backend for testing.
type MockAudioBackend struct {
	format       audio.Format
	callbackData []byte
	paused       bool
	callCount    int
}

func NewMockAudioBackend(format audio.Format) *MockAudioBackend {
	return &MockAudioBackend{
		format:       format,
		callbackData: make([]byte, format.BufferSize()),
	}
}

func (m *MockAudioBackend) AudioFormat() audio.Format {
	return m.format
}

func (m *MockAudioBackend) AudioCallback(buffer []byte) {
	m.callCount++
	copy(buffer, m.callbackData)
}

func (m *MockAudioBackend) AudioPaused() bool {
	return m.paused
}

func (m *MockAudioBackend) SetPaused(paused bool) {
	m.paused = paused
}

func (m *MockAudioBackend) SetCallbackData(data []byte) {
	m.callbackData = make([]byte, len(data))
	copy(m.callbackData, data)
}

func (m *MockAudioBackend) GetCallCount() int {
	return m.callCount
}

func TestMockAudioBackend(t *testing.T) {
	format := audio.Format{
		SampleRate: 44100,
		Channels:   2,
		Samples:    512,
		Format:     audio.FormatS16,
	}

	backend := NewMockAudioBackend(format)

	// Test AudioFormat
	result := backend.AudioFormat()
	assert.Equal(t, format, result)

	// Test AudioPaused
	assert.False(t, backend.AudioPaused())
	backend.SetPaused(true)
	assert.True(t, backend.AudioPaused())

	// Test AudioCallback
	testData := []byte{0x01, 0x02, 0x03, 0x04}
	backend.SetCallbackData(testData)

	buffer := make([]byte, len(testData))
	backend.AudioCallback(buffer)
	assert.Equal(t, testData, buffer)
	assert.Equal(t, 1, backend.GetCallCount())
}

func TestAudioDevice(t *testing.T) {
	format := audio.Format{
		SampleRate: 44100,
		Channels:   2,
		Samples:    512,
		Format:     audio.FormatS16,
	}

	backend := NewMockAudioBackend(format)
	device := &audioDevice{
		backend: backend,
		buffer:  make([]byte, format.BufferSize()),
	}

	// Test initial state
	assert.False(t, device.running)

	// Mock SDL functions to avoid requiring actual SDL
	originalPauseAudioDevice := PauseAudioDevice
	originalClearQueuedAudio := ClearQueuedAudio
	originalCloseAudioDevice := CloseAudioDevice

	defer func() {
		PauseAudioDevice = originalPauseAudioDevice
		ClearQueuedAudio = originalClearQueuedAudio
		CloseAudioDevice = originalCloseAudioDevice
	}()

	PauseAudioDevice = func(dev AudioDeviceID, pauseOn int) {
		// Mock implementation
	}
	ClearQueuedAudio = func(dev AudioDeviceID) {
		// Mock implementation
	}
	CloseAudioDevice = func(dev AudioDeviceID) {
		// Mock implementation
	}

	// Test multiple starts (should be safe)
	device.mutex.Lock()
	device.running = true
	device.mutex.Unlock()

	err := device.start()
	assert.NoError(t, err)

	// Test stop
	device.stop()
	assert.False(t, device.running)
}

func TestAudioDeviceProcessFrame(t *testing.T) {
	format := audio.Format{
		SampleRate: 44100,
		Channels:   2,
		Samples:    512,
		Format:     audio.FormatS16,
	}

	backend := NewMockAudioBackend(format)
	device := &audioDevice{
		backend: backend,
		buffer:  make([]byte, format.BufferSize()),
		running: true,
	}

	// Mock the SDL functions to avoid requiring actual SDL
	originalQueueAudio := QueueAudio
	originalGetQueuedAudioSize := GetQueuedAudioSize
	originalPauseAudioDevice := PauseAudioDevice

	defer func() {
		QueueAudio = originalQueueAudio
		GetQueuedAudioSize = originalGetQueuedAudioSize
		PauseAudioDevice = originalPauseAudioDevice
	}()

	queueCalled := false
	pauseCalled := false

	QueueAudio = func(dev AudioDeviceID, data []byte, length uint32) int {
		queueCalled = true
		return 0 // Success
	}

	GetQueuedAudioSize = func(dev AudioDeviceID) uint32 {
		return 0 // Empty queue, should trigger queueing
	}

	PauseAudioDevice = func(dev AudioDeviceID, pauseOn int) {
		pauseCalled = true
	}

	// Test when not paused
	backend.SetPaused(false)
	device.processAudioFrame()

	assert.True(t, queueCalled)
	assert.True(t, pauseCalled) // Should unpause
	assert.True(t, backend.GetCallCount() > 0)

	// Reset
	queueCalled = false
	pauseCalled = false
	backend = NewMockAudioBackend(format)
	device.backend = backend

	// Test when paused
	backend.SetPaused(true)
	device.processAudioFrame()

	assert.False(t, queueCalled) // Should not queue when paused
	assert.True(t, pauseCalled)  // Should pause
}

func BenchmarkConvertToSDLFormat(b *testing.B) {
	format := audio.FormatS16

	b.ResetTimer()
	for range b.N {
		_, _ = convertToSDLFormat(format)
	}
}

func BenchmarkAudioCallback(b *testing.B) {
	format := audio.Format{
		SampleRate: 44100,
		Channels:   2,
		Samples:    512,
		Format:     audio.FormatS16,
	}

	backend := NewMockAudioBackend(format)
	buffer := make([]byte, format.BufferSize())

	b.ResetTimer()
	for range b.N {
		backend.AudioCallback(buffer)
	}
}

func BenchmarkProcessAudioFrame(b *testing.B) {
	format := audio.Format{
		SampleRate: 44100,
		Channels:   2,
		Samples:    512,
		Format:     audio.FormatS16,
	}

	backend := NewMockAudioBackend(format)
	device := &audioDevice{
		backend: backend,
		buffer:  make([]byte, format.BufferSize()),
		running: true,
	}

	// Mock SDL functions for benchmarking
	QueueAudio = func(dev AudioDeviceID, data []byte, length uint32) int {
		return 0
	}
	GetQueuedAudioSize = func(dev AudioDeviceID) uint32 {
		return 0
	}
	PauseAudioDevice = func(dev AudioDeviceID, pauseOn int) {
		// No-op
	}

	b.ResetTimer()
	for range b.N {
		device.processAudioFrame()
	}
}
