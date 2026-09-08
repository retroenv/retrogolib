package sm83

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestStepInterruptDispatch(t *testing.T) {
	// Interrupt entry used to execute the vector's first instruction immediately.
	cpu, mem := newStepCPU(t, 0x00)
	cpu.ime = true
	mem.Write(AddrIE, IntVBlank|IntTimer)
	mem.Write(AddrIF, IntVBlank|IntTimer)
	mem.Write(VectorVBlank, 0x3C) // INC A.
	assert.NoError(t, cpu.Step())
	assert.Equal(t, VectorVBlank, cpu.PC)
	assert.Equal(t, uint64(5), cpu.Cycles())
	assert.Equal(t, uint8(0), cpu.A)
	assert.Equal(t, uint16(0x100), mem.ReadWord(cpu.SP))
	assert.Equal(t, IntTimer, mem.Read(AddrIF))
	assert.False(t, cpu.ime)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint8(1), cpu.A)
}

func TestEIDelay(t *testing.T) {
	// A DI immediately after EI must cancel the pending enable.
	tests := []struct {
		name    string
		next    byte
		wantIME bool
	}{
		{name: "NOP enables", next: 0x00, wantIME: true},
		{name: "DI cancels", next: 0xF3},
		{name: "second EI preserves first delay", next: 0xFB, wantIME: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu, _ := newStepCPU(t, 0xFB, tt.next, 0x00)
			assert.NoError(t, cpu.Step())
			assert.False(t, cpu.ime)
			assert.NoError(t, cpu.Step())
			assert.Equal(t, tt.wantIME, cpu.ime)
			assert.NoError(t, cpu.Step())
			assert.Equal(t, tt.wantIME, cpu.ime)
		})
	}
}

func TestHALTBugFetch(t *testing.T) {
	// HALT with IME clear and an already pending interrupt suppresses one fetch
	// increment, repeating the opcode byte as an operand rather than rerunning HALT.
	tests := []struct {
		name    string
		program []byte
		wantPC  uint16
		wantA   byte
		wantBC  uint16
		wantE   byte
	}{
		{name: "one byte", program: []byte{0x3C}, wantPC: 0x101, wantA: 1},
		{name: "immediate byte", program: []byte{0x3E, 0x42}, wantPC: 0x102, wantA: 0x3E},
		{name: "immediate word", program: []byte{0x01, 0x12, 0x34}, wantPC: 0x103, wantBC: 0x1201},
		{name: "CB prefix", program: []byte{0xCB, 0x00}, wantPC: 0x102, wantE: 2},
		{name: "relative jump", program: []byte{0x18, 0x00}, wantPC: 0x11A},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu, mem := newStepCPU(t, append([]byte{0x76}, tt.program...)...)
			mem.Write(AddrIE, IntTimer)
			mem.Write(AddrIF, IntTimer)
			assert.NoError(t, cpu.Step())
			assert.False(t, cpu.Halted())
			assert.Equal(t, uint16(0x101), cpu.PC)
			assert.NoError(t, cpu.Step())
			assert.Equal(t, tt.wantPC, cpu.PC)
			assert.Equal(t, tt.wantA, cpu.A)
			assert.Equal(t, tt.wantBC, cpu.BC())
			assert.Equal(t, tt.wantE, cpu.E)
			assert.False(t, cpu.haltBug)
		})
	}
}

func TestHALTWakesWithoutFetchBug(t *testing.T) {
	// An interrupt arriving after HALT began must resume normal operand fetches.
	cpu, mem := newStepCPU(t, 0x76, 0x3E, 0x42)
	mem.Write(AddrIE, IntTimer)
	assert.NoError(t, cpu.Step())
	assert.True(t, cpu.Halted())
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x101), cpu.PC)
	mem.Write(AddrIF, IntTimer)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x103), cpu.PC)
	assert.Equal(t, uint8(0x42), cpu.A)
	assert.False(t, cpu.Halted())
}

func TestEIHaltInterruptReturn(t *testing.T) {
	// EI followed by a bugged HALT must return from the interrupt to HALT.
	cpu, mem := newStepCPU(t, 0xFB, 0x76, 0x00)
	mem.Write(AddrIE, IntTimer)
	mem.Write(AddrIF, IntTimer)
	mem.Write(VectorTimer, 0xD9) // RETI.
	assert.NoError(t, cpu.Step())
	assert.NoError(t, cpu.Step())
	assert.NoError(t, cpu.Step())
	assert.Equal(t, VectorTimer, cpu.PC)
	assert.Equal(t, uint16(0x101), mem.ReadWord(cpu.SP))
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x101), cpu.PC)
	assert.NoError(t, cpu.Step())
	assert.True(t, cpu.Halted())
}

func newStepCPU(t *testing.T, program ...byte) (*CPU, *BasicMemory) {
	t.Helper()
	mem := NewBasicMemory()
	mem.LoadProgram(0x100, program)
	cpu, err := New(mem, WithInitialPC(0x100), WithInitialSP(0xFFFE))
	assert.NoError(t, err)
	return cpu, mem
}
