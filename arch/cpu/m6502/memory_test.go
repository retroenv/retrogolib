package m6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

type testMemory struct {
	b [0x10000]byte
}

func (m *testMemory) Read(address uint16) uint8 {
	return m.b[address]
}

func (m *testMemory) Write(address uint16, value uint8) {
	m.b[address] = value
}

func TestMemoryImmediate(t *testing.T) {
	t.Parallel()
	m, err := NewMemory(&testMemory{})
	assert.NoError(t, err)

	i := new(uint8)
	assert.NoError(t, m.WriteAddressModes(1, i))
	assert.Equal(t, 1, *i)

	val, err := m.ReadAddressModes(true, i)
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	val, err = m.ReadAddressModes(true, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
}

func TestMemoryAbsoluteInt(t *testing.T) {
	t.Parallel()
	m, err := NewMemory(&testMemory{})
	assert.NoError(t, err)

	assert.NoError(t, m.WriteAddressModes(1, 2))
	assert.Equal(t, 1, m.Read(2))

	val, err := m.ReadAddressModes(false, 2)
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	assert.NoError(t, m.WriteAddressModes(1, Absolute(3)))
	assert.Equal(t, 1, m.Read(2))

	val, err = m.ReadAddressModes(false, Absolute(3))
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
}

func TestReadWord(t *testing.T) {
	m, err := NewMemory(&testMemory{})
	assert.NoError(t, err)
	m.Write(0, 1)
	m.Write(1, 2)
	assert.Equal(t, 0x201, m.ReadWord(0))
}

func TestReadWordBug(t *testing.T) {
	m, err := NewMemory(&testMemory{})
	assert.NoError(t, err)
	m.Write(0x2ff, 1)
	m.Write(0x200, 2)
	assert.Equal(t, 0x201, m.ReadWordBug(0x02FF))
}

func TestWriteWord(t *testing.T) {
	m, err := NewMemory(&testMemory{})
	assert.NoError(t, err)
	m.WriteWord(0, 0x201)
	assert.Equal(t, 0x201, m.ReadWord(0))
}

func TestNewMemoryValidation(t *testing.T) {
	t.Parallel()

	// Test that nil memory returns error
	_, err := NewMemory(nil)
	assert.Error(t, err, "BasicMemory cannot be nil", "Should return error for nil memory")
	assert.Contains(t, err.Error(), "BasicMemory cannot be nil", "Should contain specific error message")

	// Test that valid memory works
	mem := &testMemory{}
	memory, err := NewMemory(mem)
	assert.NoError(t, err)
	assert.NotNil(t, memory)
}
