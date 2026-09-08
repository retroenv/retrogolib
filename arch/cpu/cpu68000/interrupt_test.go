package cpu68000

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestTriggerIRQServicesOneStep(t *testing.T) {
	// TriggerIRQ was a no-op; bus interrupts also ran the first handler opcode
	// in the dispatch step. Both request paths must stop at the handler entry.
	for _, fromBus := range []bool{false, true} {
		cpu := newTestCPU(t)
		cpu.SetSR(MaskSupervisor)
		cpu.bus.WriteLong(VectorAutoVector3*4, 0x3000)
		cpu.bus.WriteWord(0x3000, 0x4E71) // NOP.
		pc, sp := cpu.PC, cpu.A7()
		if fromBus {
			cpu.bus.(*BasicBus).irqLevel = 3
		} else {
			cpu.TriggerIRQ(3)
		}
		assert.NoError(t, cpu.Step())
		assert.Equal(t, uint32(0x3000), cpu.PC)
		assert.Equal(t, uint64(44), cpu.Cycles())
		assert.Equal(t, sp-6, cpu.A7())
		assert.Equal(t, pc, cpu.bus.ReadLong(cpu.A7()+2))
		assert.NoError(t, cpu.Step())
		assert.Equal(t, uint32(0x3002), cpu.PC)
	}
}

func TestTriggerIRQRemainsPendingWhileMasked(t *testing.T) {
	// Queued requests must survive masking and wake STOP once accepted.
	cpu := newTestCPU(t)
	cpu.bus.WriteWord(cpu.PC, 0x4E71)
	cpu.bus.WriteLong(VectorAutoVector3*4, 0x3000)
	cpu.TriggerIRQ(3)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint64(4), cpu.Cycles())
	cpu.SetSR(MaskSupervisor)
	cpu.stopped = true
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x3000), cpu.PC)
	assert.False(t, cpu.stopped)
	assert.Equal(t, uint64(48), cpu.Cycles())
}
