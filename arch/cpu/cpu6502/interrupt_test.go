package cpu6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestIRQLineMaskingAndLiveVector(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)
	cpu.PC = 0x8000
	cpu.Flags.I = 1
	cpu.SetIRQ(true)

	assert.False(t, cpu.CheckInterrupts())
	assert.True(t, cpu.State().Interrupts.IrqTriggered)

	const vector = 0x9234
	cpu.memory.WriteWord(IrqAddress, vector)
	cpu.Flags.I = 0
	cycles := cpu.Cycles()
	stackPointer := cpu.SP

	assert.True(t, cpu.CheckInterrupts())
	assert.Equal(t, uint16(vector), cpu.PC)
	assert.Equal(t, cycles+7, cpu.Cycles())
	assert.Equal(t, byte(0), cpu.memory.Read(StackBase+uint16(stackPointer-2))&0b0001_0000)
	assert.True(t, cpu.State().Interrupts.IrqTriggered)
	assert.False(t, cpu.CheckInterrupts())

	cpu.SetIRQ(false)
	assert.False(t, cpu.State().Interrupts.IrqTriggered)
}

func TestNMIPriorityAndLiveVector(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)
	cpu.Flags.I = 0
	cpu.SetIRQ(true)
	cpu.TriggerNMI()

	const vector = 0x9678
	cpu.memory.WriteWord(NMIAddress, vector)

	assert.True(t, cpu.CheckInterrupts())
	assert.Equal(t, uint16(vector), cpu.PC)
	assert.False(t, cpu.State().Interrupts.NMITriggered)
	assert.True(t, cpu.State().Interrupts.IrqTriggered)
}

func TestTriggerIRQQueuesOneRequest(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)
	cpu.Flags.I = 0
	cpu.TriggerIrq()

	assert.True(t, cpu.CheckInterrupts())
	cpu.Flags.I = 0
	assert.False(t, cpu.CheckInterrupts())
}

func TestBRKUsesLiveVectorAndSoftwareStatus(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)
	cpu.PC = 0x8000
	cpu.Flags.B = 0
	stackPointer := cpu.SP

	const vector = 0x9456
	cpu.memory.WriteWord(IrqAddress, vector)

	assert.NoError(t, brk(cpu))
	assert.Equal(t, uint16(vector), cpu.PC)
	assert.Equal(t, byte(0b0001_0000), cpu.memory.Read(StackBase+uint16(stackPointer-2))&0b0001_0000)
	assert.Equal(t, uint8(0), cpu.Flags.B)
}

func TestStallCyclesAccumulate(t *testing.T) {
	t.Parallel()
	cpu := cpuTestSetup(t)
	cpu.PC = 0x8000
	cpu.memory.Write(cpu.PC, 0xEA)
	cpu.Flags.I = 0
	cpu.SetIRQ(true)
	cycles := cpu.Cycles()

	cpu.StallCycles(2)
	cpu.StallCycles(1)
	assert.False(t, cpu.CheckInterrupts())
	cpu.SetIRQ(false)

	for range 3 {
		assert.NoError(t, cpu.Step())
	}

	assert.Equal(t, uint16(0x8000), cpu.PC)
	assert.Equal(t, cycles+3, cpu.Cycles())

	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x8001), cpu.PC)
	assert.Equal(t, cycles+5, cpu.Cycles())
}
