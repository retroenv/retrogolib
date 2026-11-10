package cpu6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestFlagsByte(t *testing.T) {
	for raw := range 256 {
		value := byte(raw)
		flags := Flags{
			C: value & 1,
			Z: value >> 1 & 1,
			I: value >> 2 & 1,
			D: value >> 3 & 1,
			B: value >> 4 & 1,
			U: value >> 5 & 1,
			V: value >> 6 & 1,
			N: value >> 7 & 1,
		}
		assert.Equal(t, value, flags.Byte())
	}
}

func TestFlagsByteUsesSnapshot(t *testing.T) {
	cpu := cpuTestSetup(t)
	cpu.setFlags(0xa5)
	snapshot := cpu.State().Flags
	cpu.setFlags(0)

	assert.Equal(t, byte(0xa5), snapshot.Byte())
	assert.Equal(t, byte(0), cpu.GetFlags())
}
