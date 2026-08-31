package sm83

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestNewOptions(t *testing.T) {
	t.Parallel()

	hook := PreExecutionHook(func(_ *CPU, _ uint8, _ ...any) {})
	opts := newOptions(nil, WithTracing(), WithPreExecutionHook(hook), WithInitialPC(0x1234), WithInitialSP(0x5678))

	assert.True(t, opts.tracing)
	assert.NotNil(t, opts.preExecutionHook)
	assert.Equal(t, uint16(0x1234), opts.initialPC)
	assert.Equal(t, uint16(0x5678), opts.initialSP)
}
