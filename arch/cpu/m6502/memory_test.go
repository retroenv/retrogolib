package m6502

import (
	"testing"

	. "github.com/retroenv/retrogolib/addressing"
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
	m := NewMemory(&testMemory{})

	i := new(uint8)
	m.WriteAddressModes(1, i)
	assert.Equal(t, 1, *i)

	assert.Equal(t, 1, m.ReadAddressModes(true, i))
	assert.Equal(t, 1, m.ReadAddressModes(true, 1))
}

func TestMemoryAbsoluteInt(t *testing.T) {
	t.Parallel()
	m := NewMemory(&testMemory{})

	m.WriteAddressModes(1, 2)
	assert.Equal(t, 1, m.Read(2))
	assert.Equal(t, 1, m.ReadAddressModes(false, 2))

	m.WriteAddressModes(1, Absolute(3))
	assert.Equal(t, 1, m.Read(2))
	assert.Equal(t, 1, m.ReadAddressModes(false, Absolute(3)))
}

func TestReadWord(t *testing.T) {
	m := NewMemory(&testMemory{})
	m.Write(0, 1)
	m.Write(1, 2)
	assert.Equal(t, 0x201, m.ReadWord(0))
}

func TestReadWordBug(t *testing.T) {
	m := NewMemory(&testMemory{})
	m.Write(0x2ff, 1)
	m.Write(0x200, 2)
	assert.Equal(t, 0x201, m.ReadWordBug(0x02FF))
}

func TestWriteWord(t *testing.T) {
	m := NewMemory(&testMemory{})
	m.WriteWord(0, 0x201)
	assert.Equal(t, 0x201, m.ReadWord(0))
}
