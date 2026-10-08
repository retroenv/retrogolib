package cpu6502

import (
	"sync"
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
	cpu.TraceStep = TraceStep{
		PC:             0x1234,
		OpcodeOperands: []byte{0xea},
	}

	cpu.Reset()

	assert.False(t, cpu.branchTaken)
	assert.Equal(t, TraceStep{}, cpu.TraceStep)
}

// TestConcurrentInterruptRequestsAndStalls verifies the documented contract
// that interrupt requests and stalls can come from another goroutine. Run it
// with the race detector to check the synchronization.
func TestConcurrentInterruptRequestsAndStalls(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts []Option
	}{
		{name: "instruction mode"},
		{name: "bus-cycle mode", opts: []Option{WithCycleHook(func(BusCycle) bool { return false })}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			memory, err := NewMemory(&testMemory{})
			assert.NoError(t, err)
			memory.WriteWord(ResetAddress, 0x8000)
			memory.WriteWord(NMIAddress, 0x8000)
			memory.WriteWord(IrqAddress, 0x8000)
			for address := uint16(0x8000); address < 0x8100; address++ {
				memory.Write(address, 0xea)
			}
			cpu := New(memory, append(tt.opts, WithVariant(VariantNMOS6502))...)

			const iterations = 2000
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				for range iterations {
					cpu.TriggerNMI()
					cpu.TriggerIrq()
					cpu.SetIRQ(true)
					cpu.SetIRQ(false)
					cpu.StallCycles(1)
				}
			}()
			for range iterations {
				cpu.CheckInterrupts()
				assert.NoError(t, cpu.Step())
			}
			wg.Wait()

			assert.True(t, cpu.Cycles() > initialCycles+iterations)
		})
	}
}

func TestBranchValidatesTarget(t *testing.T) {
	t.Parallel()

	cpu := cpuTestSetup(t)

	assert.ErrorIs(t, bcc(cpu), ErrMissingParameter)
	assert.ErrorIs(t, bcc(cpu, "invalid"), ErrInvalidParameterType)
}
