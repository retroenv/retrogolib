package cpu65816

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestNewOptions(t *testing.T) {
	t.Parallel()

	hook := PreExecutionHook(func(_ *CPU, _ *Instruction, _ ...any) {})
	opts := newOptions(nil, WithTracing(), WithPreExecutionHook(hook))

	assert.True(t, opts.tracing)
	assert.NotNil(t, opts.preExecutionHook)
}
