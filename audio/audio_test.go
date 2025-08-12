package audio

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestFormat_BytesPerSample(t *testing.T) {
	tests := []struct {
		name     string
		format   SampleFormat
		expected int
	}{
		{"8-bit unsigned", FormatU8, 1},
		{"8-bit signed", FormatS8, 1},
		{"16-bit little-endian unsigned", FormatU16LSB, 2},
		{"16-bit little-endian signed", FormatS16LSB, 2},
		{"16-bit big-endian unsigned", FormatU16MSB, 2},
		{"16-bit big-endian signed", FormatS16MSB, 2},
		{"32-bit little-endian signed", FormatS32LSB, 4},
		{"32-bit big-endian signed", FormatS32MSB, 4},
		{"32-bit little-endian float", FormatF32LSB, 4},
		{"32-bit big-endian float", FormatF32MSB, 4},
		{"unknown format defaults to 2", SampleFormat(0xFFFF), 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := Format{Format: tt.format}
			result := format.BytesPerSample()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormat_BufferSize(t *testing.T) {
	tests := []struct {
		name     string
		format   Format
		expected int
	}{
		{
			name: "mono 16-bit 512 samples",
			format: Format{
				SampleRate: 44100,
				Channels:   1,
				Samples:    512,
				Format:     FormatS16,
			},
			expected: 1024, // 512 * 1 * 2
		},
		{
			name: "stereo 16-bit 512 samples",
			format: Format{
				SampleRate: 44100,
				Channels:   2,
				Samples:    512,
				Format:     FormatS16,
			},
			expected: 2048, // 512 * 2 * 2
		},
		{
			name: "stereo 32-bit float 1024 samples",
			format: Format{
				SampleRate: 48000,
				Channels:   2,
				Samples:    1024,
				Format:     FormatF32,
			},
			expected: 8192, // 1024 * 2 * 4
		},
		{
			name: "mono 8-bit 256 samples",
			format: Format{
				SampleRate: 22050,
				Channels:   1,
				Samples:    256,
				Format:     FormatU8,
			},
			expected: 256, // 256 * 1 * 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.format.BufferSize()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSampleFormatConstants(t *testing.T) {
	// Test that platform-specific defaults are set correctly
	assert.Equal(t, FormatU16LSB, FormatU16)
	assert.Equal(t, FormatS16LSB, FormatS16)
	assert.Equal(t, FormatS32LSB, FormatS32)
	assert.Equal(t, FormatF32LSB, FormatF32)
}

// MockBackend implements Backend interface for testing.
type MockBackend struct {
	format       Format
	callbackData []byte
	paused       bool
}

func NewMockBackend(format Format) *MockBackend {
	return &MockBackend{
		format:       format,
		callbackData: make([]byte, format.BufferSize()),
	}
}

func (m *MockBackend) AudioFormat() Format {
	return m.format
}

func (m *MockBackend) AudioCallback(buffer []byte) {
	copy(buffer, m.callbackData)
}

func (m *MockBackend) AudioPaused() bool {
	return m.paused
}

func (m *MockBackend) SetPaused(paused bool) {
	m.paused = paused
}

func (m *MockBackend) SetCallbackData(data []byte) {
	m.callbackData = make([]byte, len(data))
	copy(m.callbackData, data)
}

func TestMockBackend(t *testing.T) {
	format := Format{
		SampleRate: 44100,
		Channels:   2,
		Samples:    512,
		Format:     FormatS16,
	}

	backend := NewMockBackend(format)

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
}

func BenchmarkFormat_BytesPerSample(b *testing.B) {
	format := Format{Format: FormatS16}

	b.ResetTimer()
	for range b.N {
		_ = format.BytesPerSample()
	}
}

func BenchmarkFormat_BufferSize(b *testing.B) {
	format := Format{
		SampleRate: 44100,
		Channels:   2,
		Samples:    512,
		Format:     FormatS16,
	}

	b.ResetTimer()
	for range b.N {
		_ = format.BufferSize()
	}
}

func BenchmarkMockBackend_AudioCallback(b *testing.B) {
	format := Format{
		SampleRate: 44100,
		Channels:   2,
		Samples:    512,
		Format:     FormatS16,
	}

	backend := NewMockBackend(format)
	buffer := make([]byte, format.BufferSize())

	b.ResetTimer()
	for range b.N {
		backend.AudioCallback(buffer)
	}
}
