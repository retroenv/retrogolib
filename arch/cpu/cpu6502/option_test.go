package cpu6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestNewOptions(t *testing.T) {
	t.Parallel()

	hook := PreExecutionHook(func(_ *CPU, _ *Instruction, _ ...any) {})
	opts := NewOptions(nil, WithTracing(), WithPreExecutionHook(hook), WithVariant(Variant65C02))

	assert.True(t, opts.tracing)
	assert.NotNil(t, opts.preExecutionHook)
	assert.Equal(t, Variant65C02, opts.variant)
}
