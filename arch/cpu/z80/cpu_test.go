package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/arch"
	"github.com/retroenv/retrogolib/assert"
)

// Core CPU functionality tests - initialization, state management, basic operations

func TestNew(t *testing.T) {
	memory := NewBasicMemory()

	// Test default initialization
	cpu, err := New(memory)
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x0000), cpu.PC, "PC should be initialized to 0x0000 for generic system")
	assert.Equal(t, uint16(0xFFFF), cpu.SP, "SP should be initialized to 0xFFFF for generic system")
	assert.Equal(t, uint64(0), cpu.cycles, "Cycles should start at 0")
	assert.False(t, cpu.halted, "CPU should not be halted initially")
	assert.False(t, cpu.iff1, "IFF1 should be false initially")
	assert.False(t, cpu.iff2, "IFF2 should be false initially")
	assert.Equal(t, uint8(0), cpu.im, "Interrupt mode should be 0 initially")

	// Test Game Boy system initialization
	gameboyMemory := NewBasicMemory()
	gameboyCPU, err := New(gameboyMemory, WithSystemType(arch.GameBoy))
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x0100), gameboyCPU.PC, "PC should be initialized to 0x0100 for Game Boy")
	assert.Equal(t, uint16(0xFFFE), gameboyCPU.SP, "SP should be initialized to 0xFFFE for Game Boy")

	// Test error case
	cpu, err = New(nil)
	assert.Nil(t, cpu)
	assert.ErrorIs(t, err, ErrNilMemory)
}

func TestHaltedState(t *testing.T) {
	memory := NewBasicMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	assert.False(t, cpu.Halted(), "CPU should not be halted initially")

	cpu.Halt()
	assert.True(t, cpu.Halted(), "CPU should be halted after Halt()")

	cpu.Resume()
	assert.False(t, cpu.Halted(), "CPU should not be halted after Resume()")
}

func TestState(t *testing.T) {
	memory := NewBasicMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Set some register values
	cpu.A = 0x12
	cpu.B = 0x34
	cpu.SP = 0x5678
	cpu.PC = 0x9ABC
	cpu.IX = 0xDEF0
	cpu.cycles = 1000
	cpu.setFlags(0xFF)

	state := cpu.State()

	assert.Equal(t, uint8(0x12), state.A, "State A should match CPU A")
	assert.Equal(t, uint8(0x34), state.B, "State B should match CPU B")
	assert.Equal(t, uint16(0x5678), state.SP, "State SP should match CPU SP")
	assert.Equal(t, uint16(0x9ABC), state.PC, "State PC should match CPU PC")
	assert.Equal(t, uint16(0xDEF0), state.IX, "State IX should match CPU IX")
	assert.Equal(t, uint64(1000), state.Cycles, "State cycles should match CPU cycles")

	// Check flags
	assert.Equal(t, uint8(1), state.Flags.C, "State carry flag should be set")
	assert.Equal(t, uint8(1), state.Flags.Z, "State zero flag should be set")
	assert.Equal(t, uint8(1), state.Flags.S, "State sign flag should be set")
}

func TestMemoryAccess(t *testing.T) {
	memory := NewBasicMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test that CPU has access to memory
	assert.Equal(t, memory, cpu.Memory(), "CPU should return the same memory instance")

	// Test memory operations through CPU
	memory.Write(0x1000, 0x42)
	assert.Equal(t, uint8(0x42), memory.Read(0x1000), "Memory read should return written value")
}
