package cpu6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestIRQRemainsPendingWhileMasked(t *testing.T) {
	t.Parallel()

	cpu := cpuTestSetup(t)
	cpu.memory.Write(cpu.PC, 0xea)
	cpu.TriggerIrq()

	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x8001), cpu.PC)
	assert.True(t, cpu.State().Interrupts.IrqTriggered)
	assert.False(t, cpu.State().Interrupts.IrqRunning)
}

func TestStepServicesIRQBeforeNextInstruction(t *testing.T) {
	t.Parallel()

	cpu := cpuTestSetup(t)
	cpu.Flags.I = 0
	cpu.memory.Write(cpu.PC, 0xea)
	startPC := cpu.PC
	cpu.TriggerIrq()

	assert.NoError(t, cpu.Step())

	assert.Equal(t, uint16(testIrqAddress), cpu.PC)
	assert.Equal(t, initialCycles+7, cpu.cycles)
	assert.Equal(t, InitialStack-3, cpu.SP)
	assert.Equal(t, startPC, stackedPC(cpu))
	status := cpu.memory.Read(StackBase + InitialStack - 2)
	assert.Equal(t, uint8(0), status&0b0001_0000)
	assert.Equal(t, uint8(0b0010_0000), status&0b0010_0000)
	assert.True(t, cpu.State().Interrupts.IrqRunning)
}

func TestNMIPrecedesIRQAndIgnoresInterruptMask(t *testing.T) {
	t.Parallel()

	memory, err := NewMemory(&testMemory{})
	assert.NoError(t, err)
	memory.WriteWord(ResetAddress, 0x8000)
	memory.WriteWord(NMIAddress, 0x9000)
	memory.WriteWord(IrqAddress, 0xa000)
	cpu := New(memory)
	cpu.TriggerIrq()
	cpu.TriggerNMI()

	assert.NoError(t, cpu.Step())

	assert.Equal(t, uint16(0x9000), cpu.PC)
	state := cpu.State()
	assert.True(t, state.Interrupts.NMIRunning)
	assert.True(t, state.Interrupts.IrqTriggered)
}

func TestInterruptAcceptsZeroVector(t *testing.T) {
	t.Parallel()

	memory, err := NewMemory(&testMemory{})
	assert.NoError(t, err)
	memory.WriteWord(ResetAddress, 0x8000)
	cpu := New(memory)
	cpu.TriggerNMI()

	assert.NoError(t, cpu.Step())

	assert.Equal(t, uint16(0), cpu.PC)
	assert.Equal(t, initialCycles+7, cpu.cycles)
	assert.Equal(t, InitialStack-3, cpu.SP)
}

func TestBRKOnlySetsBreakInStackedStatus(t *testing.T) {
	t.Parallel()

	cpu := cpuTestSetup(t)
	cpu.Flags.B = 0

	assert.NoError(t, brk(cpu))

	status := cpu.memory.Read(StackBase + InitialStack - 2)
	assert.Equal(t, uint8(0b0001_0000), status&0b0001_0000)
	assert.Equal(t, uint8(0), cpu.Flags.B)
}

func stackedPC(cpu *CPU) uint16 {
	low := uint16(cpu.memory.Read(StackBase + InitialStack - 1))
	high := uint16(cpu.memory.Read(StackBase + InitialStack))
	return high<<8 | low
}
