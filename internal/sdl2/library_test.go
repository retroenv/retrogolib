package sdl2

import (
	"sync"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestSubsystemOwnership(t *testing.T) {
	counts := make(map[uint32]int)
	s := subsystems{
		init: func(flags uint32) int { counts[flags]++; return 0 },
		quit: func(flags uint32) { counts[flags]-- },
	}
	audio, err := s.acquire(0x10)
	assert.NoError(t, err)
	secondAudio, err := s.acquire(0x10)
	assert.NoError(t, err)
	video, err := s.acquire(0x20)
	assert.NoError(t, err)
	video()
	video()
	assert.Equal(t, 2, counts[0x10])
	assert.Equal(t, 0, counts[0x20])
	audio()
	assert.Equal(t, 1, counts[0x10])
	secondAudio()
	assert.Equal(t, 0, counts[0x10])
}

func TestSubsystemFailure(t *testing.T) {
	s := subsystems{
		init:     func(uint32) int { return -1 },
		getError: func() string { return "unavailable" },
	}
	release, err := s.acquire(0x10)
	assert.Nil(t, release)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unavailable")
}

func TestConcurrentSubsystemOwnership(t *testing.T) {
	count := 0
	s := subsystems{
		init: func(uint32) int { count++; return 0 },
		quit: func(uint32) { count-- },
	}
	var workers sync.WaitGroup
	for range 20 {
		workers.Go(func() {
			release, _ := s.acquire(0x10)
			release()
			release()
		})
	}
	workers.Wait()
	assert.Equal(t, 0, count)
}
