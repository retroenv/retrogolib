package cpu6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestNewRejectsNilMemory(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() { New(nil) })
}

func TestStateIncludesUnusedFlag(t *testing.T) {
	t.Parallel()

	cpu := cpuTestSetup(t)
	cpu.Flags.U = 1

	assert.Equal(t, cpu.Flags, cpu.State().Flags)
}

func TestResetClearsTransientExecutionState(t *testing.T) {
	t.Parallel()

	cpu := cpuTestSetup(t)
	cpu.branchTaken = true
	cpu.TraceStep = TraceStep{PC: 0x1234, OpcodeOperands: []byte{0xea}}

	cpu.Reset()

	assert.False(t, cpu.branchTaken)
	assert.Equal(t, TraceStep{}, cpu.TraceStep)
}

func TestBranchValidatesTarget(t *testing.T) {
	t.Parallel()

	cpu := cpuTestSetup(t)

	assert.ErrorIs(t, bcc(cpu), ErrMissingParameter)
	assert.ErrorIs(t, bcc(cpu, "invalid"), ErrInvalidParameterType)
}
