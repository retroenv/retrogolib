package m68000

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestNew(t *testing.T) {
	mem := NewBasicMemory()
	bus := NewBasicBus(mem)

	cpu, err := New(bus, WithInitialPC(0x1000), WithInitialSP(0x2000))
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x1000), cpu.PC)
	assert.Equal(t, uint32(0x2000), cpu.sp)
	assert.True(t, cpu.IsSupervisor())

	// Test nil bus error.
	cpu, err = New(nil)
	assert.Nil(t, cpu)
	assert.ErrorIs(t, err, ErrNilBus)
}

func TestNew_ResetVector(t *testing.T) {
	mem := NewBasicMemory()
	// Set reset SSP at address 0.
	mem.WriteLong(0x000000, 0x00FF0000)
	// Set reset PC at address 4.
	mem.WriteLong(0x000004, 0x00001000)

	bus := NewBasicBus(mem)
	cpu, err := New(bus)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x00FF0000), cpu.SSP)
	assert.Equal(t, uint32(0x00FF0000), cpu.sp)
	assert.Equal(t, uint32(0x00001000), cpu.PC)
}

func TestState(t *testing.T) {
	cpu := newTestCPU(t)

	cpu.D[0] = 0x12345678
	cpu.D[7] = 0xDEADBEEF
	cpu.A[0] = 0x00001000
	cpu.PC = 0x00002000
	cpu.cycles = 500

	state := cpu.State()
	assert.Equal(t, uint32(0x12345678), state.D[0])
	assert.Equal(t, uint32(0xDEADBEEF), state.D[7])
	assert.Equal(t, uint32(0x00001000), state.A[0])
	assert.Equal(t, uint32(0x00002000), state.PC)
	assert.Equal(t, uint64(500), state.Cycles)
}

func TestA7_SupervisorUserSwitch(t *testing.T) {
	cpu := newTestCPU(t)

	// Start in supervisor mode.
	cpu.sp = 0x00010000
	cpu.SSP = 0x00010000
	cpu.USP = 0x00020000

	assert.Equal(t, uint32(0x00010000), cpu.A7())

	// Switch to user mode.
	cpu.SetSR(cpu.GetSR() & ^uint16(MaskSupervisor))

	assert.Equal(t, uint32(0x00020000), cpu.A7())

	// Switch back to supervisor mode.
	cpu.SetSR(cpu.GetSR() | MaskSupervisor)
	assert.Equal(t, uint32(0x00010000), cpu.A7())
}

func TestHaltResume(t *testing.T) {
	cpu := newTestCPU(t)

	assert.False(t, cpu.Halted())
	cpu.Halt()
	assert.True(t, cpu.Halted())
	cpu.Resume()
	assert.False(t, cpu.Halted())
}

func TestCycles(t *testing.T) {
	cpu := newTestCPU(t)
	assert.Equal(t, uint64(0), cpu.Cycles())

	cpu.cycles = 42
	assert.Equal(t, uint64(42), cpu.Cycles())
}

// newTestCPU creates a CPU for testing with a basic memory/bus.
func newTestCPU(t *testing.T) *CPU {
	t.Helper()
	mem := NewBasicMemory()
	bus := NewBasicBus(mem)
	cpu, err := New(bus, WithInitialPC(0x1000), WithInitialSP(0x10000))
	assert.NoError(t, err)
	return cpu
}

// testBus implements the Bus interface for testing with custom IRQ control.
type testBus struct {
	Memory
	irqLevel uint8
}

func (b *testBus) IRQAcknowledge(level uint8) uint32 {
	return uint32(VectorAutoVector1) + uint32(level) - 1
}

func (b *testBus) IRQLevel() uint8 { return b.irqLevel }

func (b *testBus) OnReset() {}
