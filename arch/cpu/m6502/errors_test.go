package m6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

type errorTestMemory struct {
	b [0x10000]byte
}

func (m *errorTestMemory) Read(address uint16) uint8 {
	return m.b[address]
}

func (m *errorTestMemory) Write(address uint16, value uint8) {
	m.b[address] = value
}

func TestErrorConstants(t *testing.T) {
	basicMem := &errorTestMemory{}
	memory, err := NewMemory(basicMem)
	assert.NoError(t, err)
	cpu := New(memory)

	// Test missing parameter error
	err = jsr(cpu) // No parameters
	assert.ErrorIs(t, err, ErrMissingParameter, "Should return ErrMissingParameter")

	// Test invalid parameter type error
	err = jsr(cpu, "invalid") // Wrong type
	assert.ErrorIs(t, err, ErrInvalidParameterType, "Should return ErrInvalidParameterType")

	// Test unsupported addressing mode from memory operations
	_, err = memory.ReadAddressModes(true, "invalid_param")
	assert.ErrorIs(t, err, ErrUnsupportedAddressingMode, "Should return ErrUnsupportedAddressingMode")

	// Test invalid register type
	_, err = memory.indirectMemoryPointer(IndirectResolved(0x1000), "invalid_register")
	assert.ErrorIs(t, err, ErrInvalidRegisterType, "Should return ErrInvalidRegisterType")

	// Test unknown opcode error with detailed error message
	cpu.PC = 0x1234
	basicMem.b[0x1234] = 0x02 // Opcode 0x02 is not defined (has nil instruction)
	err = cpu.Step()
	assert.ErrorIs(t, err, ErrUnknownOpcode, "Should return ErrUnknownOpcode for invalid opcode")
	assert.Contains(t, err.Error(), "0x02", "Error should contain opcode byte")
	assert.Contains(t, err.Error(), "0x1234", "Error should contain PC address")
}

// TestUnsupportedAddressingModeError tests the error message format for unsupported addressing modes.
func TestUnsupportedAddressingModeError(t *testing.T) {
	basicMem := &errorTestMemory{}
	memory, err := NewMemory(basicMem)
	assert.NoError(t, err)
	cpu := New(memory)

	// Test unsupported addressing mode error message format
	invalidMode := AddressingMode(99) // Use an invalid addressing mode
	params, opcodes, pageCrossed, err := readOpParams(cpu, invalidMode)
	assert.ErrorIs(t, err, ErrUnsupportedAddressingMode, "Should return ErrUnsupportedAddressingMode")
	assert.Contains(t, err.Error(), "0x63", "Error should contain properly formatted mode value (99 = 0x63)")
	assert.Nil(t, params, "Params should be nil on error")
	assert.Nil(t, opcodes, "Opcodes should be nil on error")
	assert.False(t, pageCrossed, "PageCrossed should be false on error")
}
