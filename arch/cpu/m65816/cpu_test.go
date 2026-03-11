package m65816

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestNew(t *testing.T) {
	cpu, _ := newTestCPU(t)
	assert.Equal(t, uint16(0x8000), cpu.PC)
	assert.True(t, cpu.E)
	assert.Equal(t, uint8(1), cpu.Flags.M)
	assert.Equal(t, uint8(1), cpu.Flags.X)
}

func TestNilMemory(t *testing.T) {
	_, err := NewMemory(nil)
	assert.Error(t, err)
}

func TestAccWidth(t *testing.T) {
	cpu, _ := newTestCPU(t)

	// Emulation mode: always 1
	assert.Equal(t, 1, cpu.AccWidth())

	// Native mode, M=1: still 1
	cpu.E = false
	cpu.Flags.M = 1
	assert.Equal(t, 1, cpu.AccWidth())

	// Native mode, M=0: 2
	cpu.Flags.M = 0
	assert.Equal(t, 2, cpu.AccWidth())
}

func TestIdxWidth(t *testing.T) {
	cpu, _ := newTestCPU(t)

	cpu.E = false
	cpu.Flags.X = 0
	assert.Equal(t, 2, cpu.IdxWidth())
	cpu.Flags.X = 1
	assert.Equal(t, 1, cpu.IdxWidth())
}

func TestAccumulatorAB(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.C = 0x1234
	assert.Equal(t, uint8(0x34), cpu.A())
	assert.Equal(t, uint8(0x12), cpu.B())
}

func TestReset(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.C = 0xFFFF
	cpu.X = 0xFFFF
	cpu.Y = 0xFFFF
	cpu.Reset()
	assert.Equal(t, uint16(0), cpu.C)
	assert.True(t, cpu.E)
	assert.Equal(t, uint16(0x8000), cpu.PC)
}

func TestState(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.C = 0x42
	cpu.X = 0x10
	cpu.Y = 0x20
	s := cpu.State()
	assert.Equal(t, uint16(0x42), s.C)
	assert.Equal(t, uint16(0x10), s.X)
}

func TestPushPop(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.SP = 0x01FF
	cpu.push8(0xAB)
	assert.Equal(t, uint16(0x01FE), cpu.SP)
	v := cpu.pop8()
	assert.Equal(t, uint8(0xAB), v)
	assert.Equal(t, uint16(0x01FF), cpu.SP)
}

func TestPush16Pop16(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.SP = 0x01FF
	cpu.push16(0x1234)
	v := cpu.pop16()
	assert.Equal(t, uint16(0x1234), v)
}

func TestFullPC(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.PB = 0x02
	cpu.PC = 0x8000
	assert.Equal(t, uint32(0x028000), cpu.FullPC())
}

// testMem is a simple flat 16 MB memory for testing.
type testMem struct {
	data [1 << 24]byte
}

func newTestCPU(t *testing.T) (*CPU, *testMem) {
	t.Helper()
	mem := &testMem{}
	// Set reset vector to $8000
	mem.WriteWord(VectorEmuRESET, 0x8000)
	wrapped, err := NewMemory(mem)
	assert.NoError(t, err)
	cpu, err := New(wrapped)
	assert.NoError(t, err)
	return cpu, mem
}

func (m *testMem) Read(addr uint32) uint8     { return m.data[addr&0xFFFFFF] }
func (m *testMem) Write(addr uint32, v uint8) { m.data[addr&0xFFFFFF] = v }
func (m *testMem) ReadWord(addr uint32) uint16 {
	lo := uint16(m.data[addr&0xFFFFFF])
	hi := uint16(m.data[(addr+1)&0xFFFFFF])
	return hi<<8 | lo
}
func (m *testMem) WriteWord(addr uint32, v uint16) {
	m.data[addr&0xFFFFFF] = uint8(v)
	m.data[(addr+1)&0xFFFFFF] = uint8(v >> 8)
}
