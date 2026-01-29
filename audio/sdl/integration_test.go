//go:build manualtest

package sdl

import (
	"encoding/binary"
	"math"
	"testing"
	"time"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/audio"
)

// TestPlayMelody plays a short melody through the SDL audio backend.
// Run with: go test -tags manualtest -run TestPlayMelody -v ./audio/sdl/
func TestPlayMelody(t *testing.T) {
	b := newMelodyBackend()

	audioStart, audioStop, err := Setup(b)
	assert.NoError(t, err)
	defer audioStop()

	err = audioStart()
	assert.NoError(t, err)

	// Play melody for its full duration.
	time.Sleep(b.totalDuration())
}

// melodyBackend implements audio.Backend and synthesizes a simple melody.
type melodyBackend struct {
	format  audio.Format
	phase   float64
	sample  int
	notes   []note
	paused  bool
}

type note struct {
	frequency float64       // Hz, 0 = rest
	duration  time.Duration // how long to play
}

func newMelodyBackend() *melodyBackend {
	return &melodyBackend{
		format: audio.Format{
			SampleRate: 44100,
			Channels:   1,
			Samples:    1024,
			Format:     audio.FormatS16,
		},
		// "Ode to Joy" opening phrase (simplified).
		notes: []note{
			{329.63, 400 * time.Millisecond}, // E4
			{329.63, 400 * time.Millisecond}, // E4
			{349.23, 400 * time.Millisecond}, // F4
			{392.00, 400 * time.Millisecond}, // G4
			{392.00, 400 * time.Millisecond}, // G4
			{349.23, 400 * time.Millisecond}, // F4
			{329.63, 400 * time.Millisecond}, // E4
			{293.66, 400 * time.Millisecond}, // D4
			{261.63, 400 * time.Millisecond}, // C4
			{261.63, 400 * time.Millisecond}, // C4
			{293.66, 400 * time.Millisecond}, // D4
			{329.63, 400 * time.Millisecond}, // E4
			{329.63, 600 * time.Millisecond}, // E4 (dotted)
			{293.66, 200 * time.Millisecond}, // D4 (short)
			{293.66, 800 * time.Millisecond}, // D4 (held)
		},
	}
}

func (b *melodyBackend) totalDuration() time.Duration {
	var total time.Duration
	for _, n := range b.notes {
		total += n.duration
	}
	// Add a small buffer for audio queue to drain.
	return total + 200*time.Millisecond
}

func (b *melodyBackend) AudioFormat() audio.Format {
	return b.format
}

func (b *melodyBackend) AudioPaused() bool {
	return b.paused
}

func (b *melodyBackend) AudioCallback(buffer []byte) {
	samplesPerBuffer := len(buffer) / b.format.BytesPerSample()

	for i := range samplesPerBuffer {
		freq := b.currentFrequency()
		var sample int16
		if freq > 0 {
			// Sine wave with simple amplitude envelope to reduce clicks.
			sample = int16(16000 * math.Sin(b.phase))
			b.phase += 2 * math.Pi * freq / float64(b.format.SampleRate)
			if b.phase > 2*math.Pi {
				b.phase -= 2 * math.Pi
			}
		} else {
			b.phase = 0
		}

		binary.LittleEndian.PutUint16(buffer[i*2:], uint16(sample))
		b.sample++
	}
}

// currentFrequency returns the frequency for the current sample position.
func (b *melodyBackend) currentFrequency() float64 {
	elapsed := 0
	for _, n := range b.notes {
		noteSamples := int(n.duration.Seconds() * float64(b.format.SampleRate))
		if b.sample < elapsed+noteSamples {
			return n.frequency
		}
		elapsed += noteSamples
	}
	return 0 // silence after melody ends
}
