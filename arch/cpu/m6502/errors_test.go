package m6502

import (
	"errors"
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
	memory := NewMemory(basicMem)
	cpu := New(memory)

	// Test missing parameter error
	err := jsr(cpu) // No parameters
	assert.True(t, errors.Is(err, ErrMissingParameter), "Should return ErrMissingParameter")

	// Test invalid parameter type error
	err = jsr(cpu, "invalid") // Wrong type
	assert.True(t, errors.Is(err, ErrInvalidParameterType), "Should return ErrInvalidParameterType")

	// Test unsupported addressing mode from memory operations
	_, err = memory.ReadAddressModes(true, "invalid_param")
	assert.True(t, errors.Is(err, ErrUnsupportedAddressingMode), "Should return ErrUnsupportedAddressingMode")

	// Test invalid register type
	_, err = memory.indirectMemoryPointer(IndirectResolved(0x1000), "invalid_register")
	assert.True(t, errors.Is(err, ErrInvalidRegisterType), "Should return ErrInvalidRegisterType")
}
