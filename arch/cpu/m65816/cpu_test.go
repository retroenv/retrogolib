package m65816

import (
	"testing"
)

// testMem is a simple flat 16 MB memory for testing.
type testMem struct {
	data [1 << 24]byte
}

func (m *testMem) ReadByte(addr uint32) uint8     { return m.data[addr&0xFFFFFF] }
func (m *testMem) WriteByte(addr uint32, v uint8) { m.data[addr&0xFFFFFF] = v }
func (m *testMem) ReadWord(addr uint32) uint16 {
	lo := uint16(m.data[addr&0xFFFFFF])
	hi := uint16(m.data[(addr+1)&0xFFFFFF])
	return hi<<8 | lo
}
func (m *testMem) WriteWord(addr uint32, v uint16) {
	m.data[addr&0xFFFFFF] = uint8(v)
	m.data[(addr+1)&0xFFFFFF] = uint8(v >> 8)
}

func newTestCPU(t *testing.T) (*CPU, *testMem) {
	t.Helper()
	mem := &testMem{}
	// Set reset vector to $8000
	mem.WriteWord(VectorEmuRESET, 0x8000)
	wrapped, err := NewMemory(mem)
	if err != nil {
		t.Fatal(err)
	}
	cpu, err := New(wrapped)
	if err != nil {
		t.Fatal(err)
	}
	return cpu, mem
}

func TestNew(t *testing.T) {
	cpu, _ := newTestCPU(t)
	if cpu.PC != 0x8000 {
		t.Errorf("expected PC=0x8000, got 0x%04X", cpu.PC)
	}
	if !cpu.E {
		t.Error("expected emulation mode (E=true)")
	}
	if cpu.Flags.M != 1 {
		t.Error("expected M=1 in emulation mode")
	}
	if cpu.Flags.X != 1 {
		t.Error("expected X=1 in emulation mode")
	}
}

func TestNilMemory(t *testing.T) {
	_, err := NewMemory(nil)
	if err == nil {
		t.Error("expected error for nil memory")
	}
}

func TestAccWidth(t *testing.T) {
	cpu, _ := newTestCPU(t)

	// Emulation mode: always 1
	if cpu.AccWidth() != 1 {
		t.Errorf("emulation mode AccWidth: want 1, got %d", cpu.AccWidth())
	}

	// Native mode, M=1: still 1
	cpu.E = false
	cpu.Flags.M = 1
	if cpu.AccWidth() != 1 {
		t.Errorf("native M=1 AccWidth: want 1, got %d", cpu.AccWidth())
	}

	// Native mode, M=0: 2
	cpu.Flags.M = 0
	if cpu.AccWidth() != 2 {
		t.Errorf("native M=0 AccWidth: want 2, got %d", cpu.AccWidth())
	}
}

func TestIdxWidth(t *testing.T) {
	cpu, _ := newTestCPU(t)

	cpu.E = false
	cpu.Flags.X = 0
	if cpu.IdxWidth() != 2 {
		t.Errorf("native X=0 IdxWidth: want 2, got %d", cpu.IdxWidth())
	}
	cpu.Flags.X = 1
	if cpu.IdxWidth() != 1 {
		t.Errorf("native X=1 IdxWidth: want 1, got %d", cpu.IdxWidth())
	}
}

func TestAccumulatorAB(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.C = 0x1234
	if cpu.A() != 0x34 {
		t.Errorf("A() = 0x%02X, want 0x34", cpu.A())
	}
	if cpu.B() != 0x12 {
		t.Errorf("B() = 0x%02X, want 0x12", cpu.B())
	}
}

func TestReset(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.C = 0xFFFF
	cpu.X = 0xFFFF
	cpu.Y = 0xFFFF
	cpu.Reset()
	if cpu.C != 0 {
		t.Errorf("after Reset C=%04X, want 0", cpu.C)
	}
	if !cpu.E {
		t.Error("after Reset should be in emulation mode")
	}
	if cpu.PC != 0x8000 {
		t.Errorf("after Reset PC=0x%04X, want 0x8000", cpu.PC)
	}
}

func TestState(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.C = 0x42
	cpu.X = 0x10
	cpu.Y = 0x20
	s := cpu.State()
	if s.C != 0x42 {
		t.Errorf("State.C = %d, want 66", s.C)
	}
	if s.X != 0x10 {
		t.Errorf("State.X = %d, want 16", s.X)
	}
}

func TestPushPop(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.SP = 0x01FF
	cpu.push8(0xAB)
	if cpu.SP != 0x01FE {
		t.Errorf("after push8 SP=%04X, want $01FE", cpu.SP)
	}
	v := cpu.pop8()
	if v != 0xAB {
		t.Errorf("pop8 = 0x%02X, want 0xAB", v)
	}
	if cpu.SP != 0x01FF {
		t.Errorf("after pop8 SP=%04X, want $01FF", cpu.SP)
	}
}

func TestPush16Pop16(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.SP = 0x01FF
	cpu.push16(0x1234)
	v := cpu.pop16()
	if v != 0x1234 {
		t.Errorf("pop16 = 0x%04X, want 0x1234", v)
	}
}

func TestFullPC(t *testing.T) {
	cpu, _ := newTestCPU(t)
	cpu.PB = 0x02
	cpu.PC = 0x8000
	if cpu.FullPC() != 0x028000 {
		t.Errorf("FullPC = 0x%06X, want 0x028000", cpu.FullPC())
	}
}
