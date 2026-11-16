package cpu6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

const (
	executionTestPC     = 0x8000
	executionTestVector = 0x9000
	executionTestStore  = 0x0200
	executionTestValue  = 0x31
)

type executionRecorder struct {
	cycles       []BusCycle
	events       []StackEvent
	instructions int
}

type executionInterruptTest struct {
	name string

	mask   byte
	opcode byte
	source InterruptSource

	wantImmediate bool
}

func TestExecutionStateRoundTrip(t *testing.T) {
	cpu := cpuTestSetup(t)
	cpu.irqLine = true
	cpu.irqPolled = true
	cpu.nmiSample = true
	cpu.jammed = true
	cpu.jamCycle = 2
	cpu.stallCycles = 17
	saved := cpu.SaveExecutionState()
	cpu.Reset()
	assert.NoError(t, cpu.RestoreExecutionState(saved))
	assert.Equal(t, saved, cpu.SaveExecutionState())
	assert.True(t, cpu.Jammed())
	invalid := saved
	invalid.State.Flags.C = 2
	assert.Error(t, cpu.RestoreExecutionState(invalid))
	assert.Equal(t, saved, cpu.SaveExecutionState())
}

func TestRuntimeHooksCanBeRemoved(t *testing.T) {
	cpu := cpuTestSetup(t)
	var pre, stack int
	cpu.SetPreExecutionHook(func(*CPU, *Instruction, ...any) { pre++ })
	cpu.SetStackEventHook(func(StackEvent) { stack++ })
	cpu.Memory().Write(0, 0x48)
	cpu.PC = 0
	assert.NoError(t, cpu.Step())
	assert.Equal(t, 1, pre)
	assert.Equal(t, 1, stack)
	cpu.SetPreExecutionHook(nil)
	cpu.SetStackEventHook(nil)
	assert.Nil(t, cpu.PreExecutionHook())
	assert.Nil(t, cpu.StackEventHook())
	cpu.PC = 0
	assert.NoError(t, cpu.Step())
	assert.Equal(t, 1, pre)
	assert.Equal(t, 1, stack)
}

func TestConstructorHooksCanBeReplaced(t *testing.T) {
	for _, cycleMode := range []bool{false, true} {
		t.Run(executionModeName(cycleMode), func(t *testing.T) {
			bus := &testMemory{}
			copy(bus.b[executionTestPC:], []byte{0x48, 0xa9, executionTestValue, 0x8d, 0x00, 0x02})
			var initial, replacement executionRecorder
			cpu := newExecutionTestCPU(t, bus, &initial, cycleMode)
			assert.Nil(t, cpu.opts.preExecutionHook)
			assert.Nil(t, cpu.opts.stackEventHook)

			for range 3 {
				assert.NoError(t, cpu.Step())
			}

			assert.Equal(t, 3, initial.instructions)
			assert.Len(t, initial.events, 1)
			cpu.SetPreExecutionHook(replacement.recordInstruction)
			cpu.SetStackEventHook(replacement.recordStack)
			cpu.PC = executionTestPC
			for range 3 {
				assert.NoError(t, cpu.Step())
			}

			assert.Equal(t, 3, initial.instructions)
			assert.Len(t, initial.events, 1)
			assert.Equal(t, 3, replacement.instructions)
			assert.Len(t, replacement.events, 1)
			cpu.SetPreExecutionHook(nil)
			cpu.SetStackEventHook(nil)
			cpu.PC = executionTestPC
			for range 3 {
				assert.NoError(t, cpu.Step())
			}

			assert.Equal(t, 3, replacement.instructions)
			assert.Len(t, replacement.events, 1)
		})
	}
}

func TestExecutionStateResumesInstructions(t *testing.T) {
	for _, cycleMode := range []bool{false, true} {
		t.Run(executionModeName(cycleMode), func(t *testing.T) {
			bus := &testMemory{}
			copy(bus.b[executionTestPC:], []byte{
				0xa9, executionTestValue, 0x48, 0xa9, 0x00, 0x68, 0x8d, 0x00, 0x02, 0xe8,
			})
			var original, replay executionRecorder
			cpu := newExecutionTestCPU(t, bus, &original, cycleMode)
			assert.NoError(t, cpu.Step())
			assert.NoError(t, cpu.Step())
			cpu.StallCycles(2)
			saved := cpu.SaveExecutionState()
			replayBus := *bus
			restored := newExecutionTestCPU(t, &replayBus, &replay, cycleMode)
			assert.NoError(t, restored.RestoreExecutionState(saved))
			original = executionRecorder{}

			for range 6 {
				assert.NoError(t, cpu.Step())
				assert.NoError(t, restored.Step())
				assert.Equal(t, cpu.SaveExecutionState(), restored.SaveExecutionState())
			}

			assert.Equal(t, byte(executionTestValue), restored.A)
			assert.Equal(t, byte(1), restored.X)
			assert.Equal(t, byte(executionTestValue), replayBus.b[executionTestStore])
			assert.Equal(t, original, replay)
			assert.Equal(t, bus.b, replayBus.b)
		})
	}
}

func TestExecutionStateResumesInterruptPolling(t *testing.T) {
	for _, test := range []executionInterruptTest{
		{
			name: "IRQ",

			opcode: 0xea,
			source: InterruptIRQ,

			wantImmediate: true,
		},
		{
			name: "IRQ after CLI",

			mask:   1,
			opcode: 0x58,
			source: InterruptIRQ,
		},
		{
			name: "NMI",

			mask:   1,
			opcode: 0xea,
			source: InterruptNMI,

			wantImmediate: true,
		},
	} {
		t.Run(test.name, test.run)
	}
}

func TestExecutionStateResumesKIL(t *testing.T) {
	bus := &testMemory{}
	bus.b[executionTestPC] = 0x02
	var original, replay executionRecorder
	cpu := newExecutionTestCPU(t, bus, &original, true)
	assert.NoError(t, cpu.Step())
	assert.NoError(t, cpu.Step())
	cpu.TriggerNMI()
	saved := cpu.SaveExecutionState()
	replayBus := *bus
	restored := newExecutionTestCPU(t, &replayBus, &replay, true)
	assert.NoError(t, restored.RestoreExecutionState(saved))
	original = executionRecorder{}
	assert.True(t, restored.Jammed())

	for range 3 {
		assert.False(t, cpu.CheckInterrupts())
		assert.False(t, restored.CheckInterrupts())
		assert.NoError(t, cpu.Step())
		assert.NoError(t, restored.Step())
		assert.Equal(t, cpu.SaveExecutionState(), restored.SaveExecutionState())
	}

	assert.Equal(t, []BusCycle{
		{
			Address: 0xfffe,
			Kind:    BusCycleDummy,
		},
		{
			Address: 0xfffe,
			Kind:    BusCycleDummy,
		},
		{
			Address: 0xffff,
			Kind:    BusCycleDummy,
		},
	}, replay.cycles)
	assert.Equal(t, uint16(executionTestPC+1), restored.PC)
	assert.Equal(t, original, replay)
	assert.Equal(t, bus.b, replayBus.b)
	assert.True(t, restored.Jammed())
}

func TestRuntimeHookChangesDuringExecution(t *testing.T) {
	cpu := cpuTestSetup(t)
	cpu.Memory().Write(0, 0x4c)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for range 1000 {
			_ = cpu.Step()
		}
	}()
	for range 1000 {
		cpu.SetPreExecutionHook(func(*CPU, *Instruction, ...any) {})
		cpu.SetPreExecutionHook(nil)
	}
	<-done
}

func (r *executionRecorder) recordCycle(cycle BusCycle) bool {
	r.cycles = append(r.cycles, cycle)
	return false
}

func (r *executionRecorder) recordInstruction(*CPU, *Instruction, ...any) {
	r.instructions++
}

func (r *executionRecorder) recordStack(event StackEvent) {
	r.events = append(r.events, event)
}

func (test executionInterruptTest) run(t *testing.T) {
	t.Helper()
	bus := &testMemory{}
	copy(bus.b[executionTestPC:], []byte{test.opcode, 0xea})
	bus.b[executionTestVector] = 0x40 // RTI returns to the saved program.
	var original, replay executionRecorder
	cpu := newExecutionTestCPU(t, bus, &original, true)
	cpu.memory.WriteWord(IrqAddress, executionTestVector)
	cpu.memory.WriteWord(NMIAddress, executionTestVector)
	cpu.Flags.I = test.mask
	if test.source == InterruptNMI {
		cpu.TriggerNMI()
	} else {
		cpu.SetIRQ(true)
	}

	assert.NoError(t, cpu.Step())
	saved := cpu.SaveExecutionState()
	replayBus := *bus
	restored := newExecutionTestCPU(t, &replayBus, &replay, true)
	assert.NoError(t, restored.RestoreExecutionState(saved))
	original = executionRecorder{}
	assert.Equal(t, test.wantImmediate, cpu.CheckInterrupts())
	assert.Equal(t, test.wantImmediate, restored.CheckInterrupts())
	if !test.wantImmediate {
		assert.NoError(t, cpu.Step())
		assert.NoError(t, restored.Step())
		assert.True(t, cpu.CheckInterrupts())
		assert.True(t, restored.CheckInterrupts())
	}

	assert.Equal(t, uint16(executionTestVector), restored.PC)
	assert.Equal(t, cpu.SaveExecutionState(), restored.SaveExecutionState())
	assert.NoError(t, cpu.Step())
	assert.NoError(t, restored.Step())
	assert.Equal(t, cpu.SaveExecutionState(), restored.SaveExecutionState())
	assert.Equal(t, original, replay)
	assert.Equal(t, bus.b, replayBus.b)
}

func executionModeName(cycleMode bool) string {
	if cycleMode {
		return "bus cycles"
	}

	return "instructions"
}

func newExecutionTestCPU(t *testing.T, bus *testMemory, recorder *executionRecorder, cycleMode bool) *CPU {
	t.Helper()
	memory, err := NewMemory(bus)
	assert.NoError(t, err)
	options := []Option{WithPreExecutionHook(recorder.recordInstruction), WithStackEventHook(recorder.recordStack)}
	if cycleMode {
		options = append(options, WithCycleHook(recorder.recordCycle))
	}

	cpu := newTestCPU(t, memory, options...)
	cpu.PC = executionTestPC
	return cpu
}
