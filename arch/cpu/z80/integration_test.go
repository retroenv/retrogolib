package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/arch"
	"github.com/retroenv/retrogolib/assert"
)

// Z80 integration tests - instruction execution, interrupt handling, and flag operations

// =============================================================================
// Instruction Execution Integration Tests
// =============================================================================

func TestStepNOP(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory, WithSystemType(arch.GameBoy))
	assert.NoError(t, err) // Game Boy starts at 0x0100

	// Set up NOP instruction at PC
	memory.Write(0x0100, 0x00) // NOP

	initialCycles := cpu.cycles
	initialPC := cpu.PC

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error for NOP")
	assert.Equal(t, initialPC+1, cpu.PC, "PC should increment by 1 for NOP")
	assert.Equal(t, initialCycles+4, cpu.cycles, "Cycles should increment by 4 for NOP")
}

func TestStepHalt(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory, WithSystemType(arch.GameBoy))
	assert.NoError(t, err) // Game Boy starts at 0x0100

	// Set up HALT instruction at PC
	memory.Write(0x0100, 0x76) // HALT

	assert.False(t, cpu.halted, "CPU should not be halted initially")

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error for HALT")
	assert.True(t, cpu.halted, "CPU should be halted after HALT instruction")

	// Test that halted CPU just advances cycles
	initialCycles := cpu.cycles
	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error when halted")
	assert.Equal(t, initialCycles+4, cpu.cycles, "Cycles should advance when halted")
}

// =============================================================================
// Interrupt Handling Tests
// =============================================================================

func TestInterrupts(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test interrupt enable/disable
	assert.False(t, cpu.iff1, "IFF1 should be false initially")
	assert.False(t, cpu.iff2, "IFF2 should be false initially")

	// Test interrupt triggers
	cpu.TriggerIRQ()
	assert.True(t, cpu.triggerIrq, "IRQ should be triggered")

	cpu.TriggerNMI()
	assert.True(t, cpu.triggerNmi, "NMI should be triggered")
}

func TestInterruptInstructions(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory, WithSystemType(arch.GameBoy))
	assert.NoError(t, err)

	// Test DI (0xF3)
	cpu.iff1 = true
	cpu.iff2 = true
	memory.Write(0x0100, 0xF3) // DI

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.False(t, cpu.iff1, "IFF1 should be disabled")
	assert.False(t, cpu.iff2, "IFF2 should be disabled")

	// Test EI (0xFB)
	memory.Write(0x0101, 0xFB) // EI

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.True(t, cpu.iff1, "IFF1 should be enabled")
	assert.True(t, cpu.iff2, "IFF2 should be enabled")
}

// =============================================================================
// CPU Public API Tests
// =============================================================================

func TestCPUPublicAPI(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test Cycles()
	initialCycles := cpu.Cycles()
	assert.Equal(t, uint64(0), initialCycles, "Initial cycles should be 0")

	// Test register pair getters
	cpu.B = 0x12
	cpu.C = 0x34
	assert.Equal(t, uint16(0x1234), cpu.BC(), "BC should return combined B and C registers")

	cpu.D = 0x56
	cpu.E = 0x78
	assert.Equal(t, uint16(0x5678), cpu.DE(), "DE should return combined D and E registers")

	cpu.H = 0x9A
	cpu.L = 0xBC
	assert.Equal(t, uint16(0x9ABC), cpu.HL(), "HL should return combined H and L registers")

	cpu.A = 0xDE
	cpu.setFlags(0xF0)
	assert.Equal(t, uint16(0xDEF0), cpu.AF(), "AF should return combined A and flags")
}

func TestLoadProgram(t *testing.T) {
	memory := NewMemory()

	// Test LoadProgram (alias for LoadROM)
	program := []byte{0x01, 0x02, 0x03, 0x04}
	memory.LoadProgram(program)

	for i, expected := range program {
		actual := memory.Read(uint16(i))
		assert.Equal(t, expected, actual, "LoadProgram should work like LoadROM")
	}

	// Test nil program
	memory.LoadProgram(nil)
	// Should not panic or error
}

func TestCPUErrorConditions(t *testing.T) {
	// Test CPU creation with nil memory
	cpu, err := New(nil)
	assert.ErrorIs(t, err, ErrNilMemory, "Should return ErrNilMemory for nil memory")
	assert.Nil(t, cpu, "CPU should be nil when creation fails")

	// Test interrupt mode validation
	memory := NewMemory()
	cpu, err = New(memory)
	assert.NoError(t, err)

	err = cpu.SetInterruptMode(InterruptMode2)
	assert.NoError(t, err, "Should accept valid interrupt mode")

	err = cpu.SetInterruptMode(InterruptMode(3))
	assert.ErrorIs(t, err, ErrInvalidInterruptMode, "Should reject invalid interrupt mode")

	// Test valid interrupt modes
	for _, mode := range []InterruptMode{InterruptMode0, InterruptMode1, InterruptMode2} {
		err = cpu.SetInterruptMode(mode)
		assert.NoError(t, err, "Should accept valid interrupt mode")
		assert.Equal(t, mode, cpu.GetInterruptMode(), "Should return set interrupt mode")
	}
}

// =============================================================================
// Flag Operations Tests
// =============================================================================

func TestFlags_GetFlags(t *testing.T) {
	tests := []struct {
		name  string
		flags Flags
		want  uint8
	}{
		{
			name:  "all flags clear",
			flags: Flags{},
			want:  0x00,
		},
		{
			name:  "carry flag set",
			flags: Flags{C: 1},
			want:  0x01,
		},
		{
			name:  "zero flag set",
			flags: Flags{Z: 1},
			want:  0x40,
		},
		{
			name:  "sign flag set",
			flags: Flags{S: 1},
			want:  0x80,
		},
		{
			name:  "all flags set",
			flags: Flags{C: 1, N: 1, P: 1, X: 1, H: 1, Y: 1, Z: 1, S: 1},
			want:  0xFF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := &CPU{Flags: tt.flags}
			got := cpu.GetFlags()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCPU_SetFlags(t *testing.T) {
	tests := []struct {
		name  string
		input uint8
		want  Flags
	}{
		{
			name:  "all clear",
			input: 0x00,
			want:  Flags{},
		},
		{
			name:  "carry set",
			input: 0x01,
			want:  Flags{C: 1},
		},
		{
			name:  "zero set",
			input: 0x40,
			want:  Flags{Z: 1},
		},
		{
			name:  "sign set",
			input: 0x80,
			want:  Flags{S: 1},
		},
		{
			name:  "all set",
			input: 0xFF,
			want:  Flags{C: 1, N: 1, P: 1, X: 1, H: 1, Y: 1, Z: 1, S: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := &CPU{}
			cpu.setFlags(tt.input)
			assert.Equal(t, tt.want, cpu.Flags)
		})
	}
}
