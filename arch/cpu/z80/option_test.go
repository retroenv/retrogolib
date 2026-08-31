package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/arch"
	"github.com/retroenv/retrogolib/assert"
)

type optionTestIOHandler struct{}

func (optionTestIOHandler) ReadPort(_ uint8) uint8 { return 0 }
func (optionTestIOHandler) WritePort(_, _ uint8)   {}

func TestNewOptions(t *testing.T) {
	t.Parallel()

	hook := PreExecutionHook(func(_ *CPU, _ uint8, _ ...any) {})
	handler := optionTestIOHandler{}
	opts := newOptions(nil, WithTracing(), WithPreExecutionHook(hook), WithIOHandler(handler),
		WithSystemType(arch.GameBoy), WithInitialPC(0x1234), WithInitialSP(0x5678))

	assert.True(t, opts.tracing)
	assert.NotNil(t, opts.preExecutionHook)
	assert.Equal(t, handler, opts.ioHandler)
	assert.Equal(t, arch.GameBoy, opts.systemType)
	assert.Equal(t, uint16(0x1234), opts.initialPC)
	assert.Equal(t, uint16(0x5678), opts.initialSP)
}
