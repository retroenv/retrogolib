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
		writeLong(cpu.bus, VectorAutoVector3*4, 0x3000)
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
		assert.Equal(t, pc, readLong(cpu.bus, cpu.A7()+2))
		assert.NoError(t, cpu.Step())
		assert.Equal(t, uint32(0x3002), cpu.PC)
	}
}

func TestLevel7InterruptIsEdgeSensitive(t *testing.T) {
	// A bus that held level 7 re-entered the handler on every step. The CPU
	// must accept level 7 once per rising edge or per mask drop below 7.
	cpu := newTestCPU(t)
	bus := cpu.bus.(*BasicBus)
	writeLong(cpu.bus, VectorAutoVector7*4, 0x3000)
	for address := uint32(0x3000); address < 0x3010; address += 2 {
		cpu.bus.WriteWord(address, 0x4E71) // NOP.
	}
	bus.irqLevel = 7

	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x3000), cpu.PC)
	sp := cpu.A7()
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x3002), cpu.PC)
	assert.Equal(t, sp, cpu.A7())

	// Lowering the mask while level 7 is still asserted accepts it again.
	cpu.SetSR(MaskSupervisor | 6<<FlagIPM0)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x3000), cpu.PC)
	assert.Equal(t, sp-6, cpu.A7())
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x3002), cpu.PC)

	// A new rising edge after the bus releases level 7 is accepted once.
	bus.irqLevel = 0
	assert.NoError(t, cpu.Step())
	bus.irqLevel = 7
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x3000), cpu.PC)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x3002), cpu.PC)
}

func TestTriggerIRQRemainsPendingWhileMasked(t *testing.T) {
	// Queued requests must survive masking and wake STOP once accepted.
	cpu := newTestCPU(t)
	cpu.bus.WriteWord(cpu.PC, 0x4E71)
	writeLong(cpu.bus, VectorAutoVector3*4, 0x3000)
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
