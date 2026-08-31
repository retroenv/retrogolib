package cpu68000

import (
	"testing"

	"github.com/retroenv/retrogolib/arch"
	"github.com/retroenv/retrogolib/assert"
)

func TestNewOptions(t *testing.T) {
	t.Parallel()

	opts := newOptions(nil, WithTracing(), WithSystemType(arch.Generic), WithInitialPC(0x1234), WithInitialSP(0x5678))

	assert.True(t, opts.tracing)
	assert.Equal(t, arch.Generic, opts.systemType)
	assert.Equal(t, uint32(0x1234), opts.initialPC)
	assert.Equal(t, uint32(0x5678), opts.initialSP)
}
