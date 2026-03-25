package m6809

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestNew(t *testing.T) {
	cpu, _ := newTestCPU(t)
	assert.Equal(t, uint16(0x8000), cpu.PC)
	assert.Equal(t, uint8(1), cpu.Flags.I)
	assert.Equal(t, uint8(1), cpu.Flags.F)
}

func TestNilMemory(t *testing.T) {
	_, err := NewMemory(nil)
	assert.Error(t, err)
}

func TestDRegister(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.A = 0x12
	cpu.B = 0x34
	assert.Equal(t, uint16(0x1234), cpu.D())

	cpu.SetD(0xABCD)
	assert.Equal(t, uint8(0xAB), cpu.A)
	assert.Equal(t, uint8(0xCD), cpu.B)
}

func TestReset(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.A = 0xFF
	cpu.B = 0xFF
	cpu.X = 0xFFFF
	cpu.Y = 0xFFFF
	cpu.Reset()
	assert.Equal(t, uint8(0), cpu.A)
	assert.Equal(t, uint8(0), cpu.B)
	assert.Equal(t, uint16(0x8000), cpu.PC)
}

func TestState(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.A = 0x42
	cpu.B = 0x10
	cpu.X = 0x1234
	cpu.Y = 0x5678
	s := cpu.State()
	assert.Equal(t, uint8(0x42), s.A)
	assert.Equal(t, uint8(0x10), s.B)
	assert.Equal(t, uint16(0x1234), s.X)
	assert.Equal(t, uint16(0x5678), s.Y)
}

func TestPushPopS(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.S = 0x0200
	cpu.pushS8(0xAB)
	assert.Equal(t, uint16(0x01FF), cpu.S)
	v := cpu.popS8()
	assert.Equal(t, uint8(0xAB), v)
	assert.Equal(t, uint16(0x0200), cpu.S)
}

func TestPushPopS16(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.S = 0x0200
	cpu.pushS16(0x1234)
	v := cpu.popS16()
	assert.Equal(t, uint16(0x1234), v)
}

func TestPushPopU(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.U = 0x0200
	cpu.pushU8(0xCD)
	assert.Equal(t, uint16(0x01FF), cpu.U)
	v := cpu.popU8()
	assert.Equal(t, uint8(0xCD), v)
	assert.Equal(t, uint16(0x0200), cpu.U)
}

func TestFlags(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.Flags.Set(0xFF)
	assert.Equal(t, uint8(0xFF), cpu.Flags.Get())

	cpu.Flags.Set(0x00)
	assert.Equal(t, uint8(0x00), cpu.Flags.Get())

	cpu.Flags.C = 1
	cpu.Flags.Z = 1
	expected := uint8(MaskCarry | MaskZero)
	assert.Equal(t, expected, cpu.Flags.Get())
}

// testMem is a simple flat 64KB memory for testing.
type testMem struct {
	data [65536]byte
}

func newTestCPU(t *testing.T) (*CPU, *testMem) {
	t.Helper()
	mem := &testMem{}
	// Set reset vector to $8000 (big-endian)
	mem.WriteWord(VectorRESET, 0x8000)
	wrapped, err := NewMemory(mem)
	assert.NoError(t, err)
	cpu, err := New(wrapped)
	assert.NoError(t, err)
	return cpu, mem
}

func (m *testMem) Read(addr uint16) uint8     { return m.data[addr] }
func (m *testMem) Write(addr uint16, v uint8) { m.data[addr] = v }
func (m *testMem) ReadWord(addr uint16) uint16 {
	// Big-endian (6809 native)
	return uint16(m.data[addr])<<8 | uint16(m.data[addr+1])
}
func (m *testMem) WriteWord(addr uint16, v uint16) {
	// Big-endian (6809 native)
	m.data[addr] = uint8(v >> 8)
	m.data[addr+1] = uint8(v)
}
