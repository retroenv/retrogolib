package cp1610

import (
	"errors"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

type testBus [0x10000]uint16

func (b *testBus) Read(address uint16) uint16  { return b[address] }
func (b *testBus) Write(address, value uint16) { b[address] = value }

func newTestCPU(t *testing.T) (*CPU, *testBus) {
	t.Helper()
	bus := &testBus{}
	memory, err := NewMemory(bus)
	assert.NoError(t, err)
	cpu, err := New(memory)
	assert.NoError(t, err)
	return cpu, bus
}

func TestNewAndReset(t *testing.T) {
	_, err := NewMemory(nil)
	assert.ErrorIs(t, err, ErrNilMemory)
	var nilBus *testBus
	_, err = NewMemory(nilBus)
	assert.ErrorIs(t, err, ErrNilMemory)
	_, err = New(nil)
	assert.ErrorIs(t, err, ErrNilMemory)

	cpu, _ := newTestCPU(t)
	assert.NoError(t, cpu.ValidateState())
	assert.Equal(t, uint16(ResetAddress), cpu.R[ProgramCounter])
	cpu.R[0] = 0x1234
	cpu.R[ProgramCounter] = 0x5000
	cpu.Flags.C = true
	cpu.Reset()
	assert.Equal(t, uint16(0), cpu.R[0])
	assert.Equal(t, uint16(ResetAddress), cpu.R[ProgramCounter])
	assert.False(t, cpu.Flags.C)
	assert.ErrorIs(t, (&CPU{}).ValidateState(), ErrNilMemory)
}

func TestRegisterArithmeticAndFlags(t *testing.T) {
	cpu, bus := newTestCPU(t)
	bus[0x1000] = 0x00C1 // ADDR R0, R1
	bus[0x1001] = 0x0141 // CMPR R0, R1
	bus[0x1002] = 0x0101 // SUBR R0, R1
	cpu.R[0] = 1
	cpu.R[1] = 0x7FFF
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x8000), cpu.R[1])
	assert.True(t, cpu.Flags.O)
	assert.True(t, cpu.Flags.S)
	assert.False(t, cpu.Flags.C)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x8000), cpu.R[1])
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x7FFF), cpu.R[1])
	assert.True(t, cpu.Flags.C)
}

func TestBranchAndJump(t *testing.T) {
	cpu, bus := newTestCPU(t)
	bus[0x1000] = 0x0204 // BEQ forward
	bus[0x1001] = 2
	bus[0x1004] = 0x0004 // JSR R5, $5010
	bus[0x1005] = 0x0100 | 0x0050
	bus[0x1006] = 0x0010
	cpu.Flags.Z = true
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x1004), cpu.R[ProgramCounter])
	assert.Equal(t, uint64(9), cpu.Cycles())
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x1007), cpu.R[5])
	assert.Equal(t, uint16(0x5010), cpu.R[ProgramCounter])

	cpu.R[ProgramCounter] = 0x1000
	cpu.Flags.Z = false
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x1002), cpu.R[ProgramCounter])
}

func TestMemoryModesAndDoubleByteData(t *testing.T) {
	cpu, bus := newTestCPU(t)
	bus[0x1000] = 0x02B8 // MVII #$1234, R0
	bus[0x1001] = 0x1234
	bus[0x1002] = 0x0260 // MVO@ R0, R4
	bus[0x1003] = 0x0001 // SDBD
	bus[0x1004] = 0x02B9 // MVII #$5678, R1 in double-byte form
	bus[0x1005] = 0x0078
	bus[0x1006] = 0x0056
	bus[0x1007] = 0x0282 // MVI $0200, R2
	bus[0x1008] = 0x0200
	cpu.R[4] = 0x0200
	for range 5 {
		assert.NoError(t, cpu.Step())
	}
	assert.Equal(t, uint16(0x1234), cpu.R[0])
	assert.Equal(t, uint16(0x1234), bus[0x0200])
	assert.Equal(t, uint16(0x0201), cpu.R[4])
	assert.Equal(t, uint16(0x5678), cpu.R[1])
	assert.Equal(t, uint16(0x1234), cpu.R[2])
	assert.False(t, cpu.Flags.D)
}

func TestStackAndInterrupt(t *testing.T) {
	cpu, bus := newTestCPU(t)
	bus[0x1000] = 0x0002 // EIS
	bus[0x1001] = 0x0034 // NOP
	bus[0x1002] = 0x0268 // MVO@ R0, R5
	bus[0x1004] = 0x02B0 // MVI@ R6, R0
	cpu.R[StackPointer] = 0x0300
	assert.NoError(t, cpu.Step())
	cpu.TriggerIRQ()
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x1002), cpu.R[ProgramCounter])
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x1004), cpu.R[ProgramCounter])
	assert.Equal(t, uint16(0x1002), bus[0x0300])
	assert.Equal(t, uint16(0x0301), cpu.R[StackPointer])
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x1002), cpu.R[0])
	assert.Equal(t, uint16(0x0300), cpu.R[StackPointer])
}

func TestShiftHaltAndInvalidJump(t *testing.T) {
	cpu, bus := newTestCPU(t)
	bus[0x1000] = 0x0050 // RLC R0
	bus[0x1001] = 0x0000 // HLT
	cpu.R[0] = 0x8000
	cpu.Flags.C = true
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(1), cpu.R[0])
	assert.True(t, cpu.Flags.C)
	assert.NoError(t, cpu.Step())
	assert.True(t, cpu.Halted())
	state := cpu.State()
	assert.Equal(t, uint16(0x1002), state.R[ProgramCounter])
	assert.NoError(t, cpu.Step())
	assert.Equal(t, state, cpu.State())

	cpu.Reset()
	bus[0x1000] = 0x0004
	bus[0x1001] = 0x0003
	assert.True(t, errors.Is(cpu.Step(), ErrInvalidOpcode))
}

func TestStatusAndUnaryInstructions(t *testing.T) {
	cpu, bus := newTestCPU(t)
	bus[0x1000] = 0x0007 // SETC
	bus[0x1001] = 0x0028 // ADCR R0
	bus[0x1002] = 0x0030 // GSWD R0
	bus[0x1003] = 0x0039 // RSWD R1
	bus[0x1004] = 0x0006 // CLRC
	cpu.R[0] = 0xFFFF
	cpu.R[1] = 0x0090
	for range 2 {
		assert.NoError(t, cpu.Step())
	}
	assert.Equal(t, uint16(0), cpu.R[0])
	assert.True(t, cpu.Flags.Z)
	assert.True(t, cpu.Flags.C)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x5050), cpu.R[0])
	assert.NoError(t, cpu.Step())
	assert.True(t, cpu.Flags.S)
	assert.True(t, cpu.Flags.C)
	assert.False(t, cpu.Flags.Z)
	assert.NoError(t, cpu.Step())
	assert.False(t, cpu.Flags.C)
}

func TestIndirectModes(t *testing.T) {
	cpu, bus := newTestCPU(t)
	bus[0x1000] = 0x0298 // MVI@ R3, R0
	bus[0x1001] = 0x02A1 // MVI@ R4, R1
	bus[0x1002] = 0x026A // MVO@ R2, R5
	bus[0x1003] = 0x0272 // PSHR R2
	bus[0x1004] = 0x02B3 // PULR R3
	cpu.R[3] = 0x0200
	cpu.R[4] = 0x0201
	cpu.R[5] = 0x0202
	cpu.R[6] = 0x0300
	cpu.R[2] = 0xCAFE
	bus[0x0200] = 0x1111
	bus[0x0201] = 0x2222
	for range 5 {
		assert.NoError(t, cpu.Step())
	}
	assert.Equal(t, uint16(0x1111), cpu.R[0])
	assert.Equal(t, uint16(0x2222), cpu.R[1])
	assert.Equal(t, uint16(0x0202), cpu.R[4])
	assert.Equal(t, uint16(0x0203), cpu.R[5])
	assert.Equal(t, uint16(0x0300), cpu.R[6])
	assert.Equal(t, uint16(0xCAFE), bus[0x0202])
	assert.Equal(t, uint16(0xCAFE), bus[0x0300])
	assert.Equal(t, uint16(0xCAFE), cpu.R[3])
}

func TestSDBDIndirectAndImmediateStore(t *testing.T) {
	cpu, bus := newTestCPU(t)
	bus[0x1000] = 0x0001 // SDBD
	bus[0x1001] = 0x02A0 // MVI@ R4, R0
	bus[0x1002] = 0x0279 // MVOI R1
	bus[0x0200] = 0x0034
	bus[0x0201] = 0x0012
	cpu.R[4] = 0x0200
	cpu.R[1] = 0xBEEF
	assert.NoError(t, cpu.Step())
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x1234), cpu.R[0])
	assert.Equal(t, uint16(0x0202), cpu.R[4])
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0xBEEF), bus[0x1003])
	assert.Equal(t, uint16(0x1004), cpu.R[ProgramCounter])
}

func TestReverseBranchAndProgramCounterRegister(t *testing.T) {
	cpu, bus := newTestCPU(t)
	bus[0x1000] = 0x0220 // B backward
	bus[0x1001] = 0
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x1001), cpu.R[ProgramCounter])

	cpu.R[ProgramCounter] = 0x1002
	bus[0x1002] = 0x0087 // MOVR R0, R7
	cpu.R[0] = 0x5000
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint16(0x5000), cpu.R[ProgramCounter])
}

func TestDirectALUInstructions(t *testing.T) {
	tests := []struct {
		name   string
		opcode uint16
		want   uint16
		carry  bool
	}{
		{"add", 0x02C1, 8, false},
		{"subtract", 0x0301, 2, true},
		{"compare", 0x0341, 5, true},
		{"and", 0x0381, 1, false},
		{"xor", 0x03C1, 6, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu, bus := newTestCPU(t)
			bus[0x1000] = tt.opcode
			bus[0x1001] = 0x0200
			bus[0x0200] = 3
			cpu.R[1] = 5
			assert.NoError(t, cpu.Step())
			assert.Equal(t, tt.want, cpu.R[1])
			assert.Equal(t, tt.carry, cpu.Flags.C)
			assert.Equal(t, uint16(0x1002), cpu.R[ProgramCounter])
		})
	}
}

func TestShiftInstructions(t *testing.T) {
	tests := []struct {
		name   string
		opcode uint16
		input  uint16
		carry  bool
		want   uint16
		wantC  bool
	}{
		{"swap", 0x0040, 0x1234, false, 0x3412, false},
		{"swap twice", 0x0044, 0x1234, false, 0x3434, false},
		{"logical left", 0x0048, 0x4001, false, 0x8002, false},
		{"rotate left", 0x0050, 0x8000, true, 0x0001, true},
		{"left with carry", 0x0058, 0x8001, false, 0x0002, true},
		{"logical right", 0x0060, 0x0101, false, 0x0080, false},
		{"arithmetic right", 0x0068, 0x8001, false, 0xC000, false},
		{"rotate right", 0x0070, 0x0001, true, 0x8000, true},
		{"right with carry", 0x0078, 0x8001, false, 0xC000, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu, bus := newTestCPU(t)
			bus[0x1000] = tt.opcode
			cpu.R[0] = tt.input
			cpu.Flags.C = tt.carry
			assert.NoError(t, cpu.Step())
			assert.Equal(t, tt.want, cpu.R[0])
			assert.Equal(t, tt.wantC, cpu.Flags.C)
		})
	}
}
