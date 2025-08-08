package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestNew(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory)

	assert.Equal(t, uint16(0x0100), cpu.PC, "PC should be initialized to 0x0100")
	assert.Equal(t, uint16(InitialStack), cpu.SP, "SP should be initialized correctly")
	assert.Equal(t, uint64(0), cpu.cycles, "Cycles should start at 0")
	assert.False(t, cpu.halted, "CPU should not be halted initially")
	assert.False(t, cpu.iff1, "IFF1 should be false initially")
	assert.False(t, cpu.iff2, "IFF2 should be false initially")
	assert.Equal(t, uint8(0), cpu.im, "Interrupt mode should be 0 initially")
}

func TestRegisterPairs(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory)

	// Test BC register pair
	cpu.B = 0x12
	cpu.C = 0x34
	assert.Equal(t, uint16(0x1234), cpu.BC(), "BC register pair should return correct value")

	cpu.setBC(0x5678)
	assert.Equal(t, uint8(0x56), cpu.B, "B should be set correctly")
	assert.Equal(t, uint8(0x78), cpu.C, "C should be set correctly")

	// Test DE register pair
	cpu.D = 0xAB
	cpu.E = 0xCD
	assert.Equal(t, uint16(0xABCD), cpu.DE(), "DE register pair should return correct value")

	cpu.setDE(0xEF01)
	assert.Equal(t, uint8(0xEF), cpu.D, "D should be set correctly")
	assert.Equal(t, uint8(0x01), cpu.E, "E should be set correctly")

	// Test HL register pair
	cpu.H = 0x23
	cpu.L = 0x45
	assert.Equal(t, uint16(0x2345), cpu.HL(), "HL register pair should return correct value")

	cpu.setHL(0x6789)
	assert.Equal(t, uint8(0x67), cpu.H, "H should be set correctly")
	assert.Equal(t, uint8(0x89), cpu.L, "L should be set correctly")

	// Test AF register pair
	cpu.A = 0xF0
	cpu.setFlags(0x0F)
	assert.Equal(t, uint16(0xF00F), cpu.AF(), "AF register pair should return correct value")

	cpu.setAF(0x1E2D)
	assert.Equal(t, uint8(0x1E), cpu.A, "A should be set correctly")
	assert.Equal(t, uint8(0x2D), cpu.GetFlags(), "Flags should be set correctly")
}

func TestStackOperations(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory)

	// Test push and pop byte
	originalSP := cpu.SP
	cpu.push(0x42)
	assert.Equal(t, originalSP-1, cpu.SP, "SP should decrement after push")
	assert.Equal(t, uint8(0x42), memory.Read(cpu.SP), "Value should be stored at SP")

	value := cpu.pop()
	assert.Equal(t, uint8(0x42), value, "Popped value should match pushed value")
	assert.Equal(t, originalSP, cpu.SP, "SP should return to original value")

	// Test push and pop 16-bit word
	cpu.push16(0x1234)
	assert.Equal(t, originalSP-2, cpu.SP, "SP should decrement by 2 after push16")

	word := cpu.pop16()
	assert.Equal(t, uint16(0x1234), word, "Popped word should match pushed word")
	assert.Equal(t, originalSP, cpu.SP, "SP should return to original value")
}

func TestExchange(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory)

	// Set main registers
	cpu.A = 0x11
	cpu.B = 0x22
	cpu.C = 0x33
	cpu.D = 0x44
	cpu.E = 0x55
	cpu.H = 0x66
	cpu.L = 0x77
	cpu.setFlags(0x88)

	// Set alternate registers
	cpu.A_ = 0xAA
	cpu.B_ = 0xBB
	cpu.C_ = 0xCC
	cpu.D_ = 0xDD
	cpu.E_ = 0xEE
	cpu.H_ = 0xFF
	cpu.L_ = 0x00
	cpu.Flags_.C = 1
	cpu.Flags_.Z = 1

	// Test exchange
	cpu.exchange()

	// Check that main and alternate registers are swapped
	assert.Equal(t, uint8(0xAA), cpu.A, "A should be swapped")
	assert.Equal(t, uint8(0xBB), cpu.B, "B should be swapped")
	assert.Equal(t, uint8(0xCC), cpu.C, "C should be swapped")
	assert.Equal(t, uint8(0xDD), cpu.D, "D should be swapped")
	assert.Equal(t, uint8(0xEE), cpu.E, "E should be swapped")
	assert.Equal(t, uint8(0xFF), cpu.H, "H should be swapped")
	assert.Equal(t, uint8(0x00), cpu.L, "L should be swapped")

	assert.Equal(t, uint8(0x11), cpu.A_, "A' should be swapped")
	assert.Equal(t, uint8(0x22), cpu.B_, "B' should be swapped")
	assert.Equal(t, uint8(0x33), cpu.C_, "C' should be swapped")
	assert.Equal(t, uint8(0x44), cpu.D_, "D' should be swapped")
	assert.Equal(t, uint8(0x55), cpu.E_, "E' should be swapped")
	assert.Equal(t, uint8(0x66), cpu.H_, "H' should be swapped")
	assert.Equal(t, uint8(0x77), cpu.L_, "L' should be swapped")
}

func TestExchangeAF(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory)

	// Set main AF
	cpu.A = 0x12
	cpu.setFlags(0x34)

	// Set alternate AF
	cpu.A_ = 0x56
	cpu.Flags_.C = 1
	cpu.Flags_.Z = 1
	cpu.Flags_.S = 1

	// Test AF exchange
	cpu.exchangeAF()

	// Check that AF is swapped
	assert.Equal(t, uint8(0x56), cpu.A, "A should be swapped")
	assert.Equal(t, uint8(0x12), cpu.A_, "A' should be swapped")

	// Flags should be swapped
	assert.Equal(t, uint8(1), cpu.Flags.C, "C flag should be from alternate")
	assert.Equal(t, uint8(1), cpu.Flags.Z, "Z flag should be from alternate")
	assert.Equal(t, uint8(1), cpu.Flags.S, "S flag should be from alternate")
}

func TestHaltedState(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory)

	assert.False(t, cpu.Halted(), "CPU should not be halted initially")

	cpu.Halt()
	assert.True(t, cpu.Halted(), "CPU should be halted after Halt()")

	cpu.Resume()
	assert.False(t, cpu.Halted(), "CPU should not be halted after Resume()")
}

func TestState(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory)

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

func TestInterrupts(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory)

	// Test interrupt enable/disable
	assert.False(t, cpu.iff1, "IFF1 should be false initially")
	assert.False(t, cpu.iff2, "IFF2 should be false initially")

	// Test interrupt triggers
	cpu.TriggerIRQ()
	assert.True(t, cpu.triggerIrq, "IRQ should be triggered")

	cpu.TriggerNMI()
	assert.True(t, cpu.triggerNmi, "NMI should be triggered")
}

func TestMemoryAccess(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory)

	// Test that CPU has access to memory
	assert.Equal(t, memory, cpu.Memory(), "CPU should return the same memory instance")

	// Test memory operations through CPU
	memory.Write(0x1000, 0x42)
	assert.Equal(t, uint8(0x42), memory.Read(0x1000), "Memory read should return written value")
}
