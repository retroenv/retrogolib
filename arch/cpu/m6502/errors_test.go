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
}
