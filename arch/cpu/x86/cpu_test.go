package x86

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/log"
)

func TestNew(t *testing.T) {
	logger := log.NewTestLogger(t)

	tests := []struct {
		name        string
		memory      *Memory
		expectError bool
	}{
		{
			name:        "valid memory",
			memory:      createTestMemory(t, logger),
			expectError: false,
		},
		{
			name:        "nil memory",
			memory:      nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu, err := New(tt.memory)

			if tt.expectError {
				assert.ErrorIs(t, err, ErrNilMemory)
				assert.Nil(t, cpu)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cpu)
				// Basic validation that CPU was initialized
				assert.NotNil(t, cpu.Memory())
			}
		})
	}
}

func TestCPU_RegisterAccess(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory := createTestMemory(t, logger)
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test direct register access
	cpu.AX = 0x1234
	cpu.BX = 0x5678
	assert.Equal(t, uint16(0x1234), cpu.AX)
	assert.Equal(t, uint16(0x5678), cpu.BX)

	// Test segment registers
	cpu.CS = 0xF000
	cpu.DS = 0x2000
	assert.Equal(t, uint16(0xF000), cpu.CS)
	assert.Equal(t, uint16(0x2000), cpu.DS)
}

func TestCPU_SegmentedAddressing(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory := createTestMemory(t, logger)
	cpu, err := New(memory)
	assert.NoError(t, err)

	tests := []struct {
		segment  uint16
		offset   uint16
		expected uint32
	}{
		{0x0000, 0x0000, 0x00000},
		{0x1000, 0x0000, 0x10000},
		{0x0000, 0x1000, 0x01000},
		{0x1234, 0x5678, 0x179B8},          // 0x1234 << 4 + 0x5678 = 0x12340 + 0x5678
		{0xFFFF, 0x000F, 0xFFFF0 + 0x000F}, // Maximum address
	}

	for _, tt := range tests {
		result := cpu.CalculateAddress(tt.segment, tt.offset)
		assert.Equal(t, tt.expected, result)
	}
}

func TestCPU_InstructionProcessing(t *testing.T) {
	logger := log.NewTestLogger(t)
	memory := createTestMemory(t, logger)
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test that opcode table is accessible
	assert.Equal(t, 256, len(Opcodes))

	// Test that first opcode (0x00 - ADD r/m8, r8) is properly defined
	opcode := Opcodes[0x00]
	assert.NotNil(t, opcode.Instruction)
	assert.Equal(t, "add", opcode.Instruction.Name)
	assert.Equal(t, ModRMRegisterAddressing, opcode.Addressing)

	// Test a specific instruction like MOV
	movOpcode := Opcodes[0x88] // MOV r/m8, r8
	assert.NotNil(t, movOpcode.Instruction)
	assert.Equal(t, "mov", movOpcode.Instruction.Name)

	// Test that CPU can still perform address calculations (useful for assembler/disassembler)
	addr := cpu.CalculateAddress(0x1000, 0x0100)
	assert.Equal(t, uint32(0x10100), addr)
}

// createTestMemory creates a test memory instance.
func createTestMemory(t *testing.T, logger *log.Logger) *Memory {
	t.Helper()
	memory, err := NewMemory(1024*1024, logger) // 1MB
	assert.NoError(t, err)
	return memory
}
