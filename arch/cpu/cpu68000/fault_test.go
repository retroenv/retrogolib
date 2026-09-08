package cpu68000

import (
	"errors"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

type rejectingBus struct {
	*BasicBus
	address     uint32
	write       bool
	reject      bool
	rejectStack bool
	cause       error
	accesses    []uint32
}

func (bus *rejectingBus) OnBusAccess(address uint32, _ OperandSize, write bool) error {
	bus.accesses = append(bus.accesses, address)
	if bus.reject && (address == bus.address && write == bus.write || bus.rejectStack && write) {
		return bus.cause
	}
	return nil
}

func TestAddressErrorFrame(t *testing.T) {
	// Odd operands previously reached memory without taking vector 3.
	tests := []struct {
		name   string
		opcode uint16
		write  bool
		user   bool
	}{
		{name: "word read", opcode: 0x3010},               // MOVE.W (A0),D0.
		{name: "long read", opcode: 0x2010},               // MOVE.L (A0),D0.
		{name: "word write", opcode: 0x3080, write: true}, // MOVE.W D0,(A0).
		{name: "long write", opcode: 0x2080, write: true}, // MOVE.L D0,(A0).
		{name: "user read", opcode: 0x3010, user: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newTestCPU(t)
			cpu.A[0], cpu.D[0] = 0x3001, 0x1234
			cpu.bus.WriteWord(cpu.PC, tt.opcode)
			cpu.bus.WriteLong(VectorAddressErr*4, 0x2000)
			cpu.USP = 0x8000
			if tt.user {
				cpu.SetSR(0)
			}
			sr := cpu.GetSR()
			assert.NoError(t, cpu.Step())
			assert.Equal(t, uint32(0x2000), cpu.PC)
			assert.Equal(t, uint32(0xFFF2), cpu.A7())
			assert.Equal(t, uint32(0x3001), cpu.bus.ReadLong(cpu.A7()+2))
			assert.Equal(t, tt.opcode, cpu.bus.ReadWord(cpu.A7()+6))
			assert.Equal(t, sr, cpu.bus.ReadWord(cpu.A7()+8))
			assert.Equal(t, uint32(0x1000), cpu.bus.ReadLong(cpu.A7()+10))
			status := uint16(5)
			if tt.user {
				status = 1
			}
			if !tt.write {
				status |= 0x10
			}
			assert.Equal(t, status, cpu.bus.ReadWord(cpu.A7())&0x1F)
			assert.Equal(t, uint64(50), cpu.Cycles())
			assert.Equal(t, uint16(0), cpu.bus.ReadWord(0x3001))
		})
	}
}

func TestAddressErrorInstructionFetch(t *testing.T) {
	// Instruction fetches must use the same alignment checks as operands.
	cpu := newTestCPU(t)
	cpu.PC = 0x1001
	cpu.bus.WriteLong(VectorAddressErr*4, 0x2000)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x2000), cpu.PC)
	assert.Equal(t, uint32(0x1001), cpu.bus.ReadLong(cpu.A7()+2))
	assert.Equal(t, uint16(0x16), cpu.bus.ReadWord(cpu.A7())&0x1F)
}

func TestStackPointersAfterStep(t *testing.T) {
	// The active SP changed while the exported SSP and State snapshot stayed stale.
	cpu := newTestCPU(t)
	cpu.bus.WriteWord(cpu.PC, 0x2F00) // MOVE.L D0,-(SP).
	assert.NoError(t, cpu.Step())
	assert.Equal(t, cpu.A7(), cpu.SSP)
	assert.Equal(t, cpu.A7(), cpu.State().SSP)
}

func TestBusErrorAbortsTransfer(t *testing.T) {
	// A host can reject a fetch, operand, or individual half of a long write.
	tests := []struct {
		name    string
		opcode  uint16
		address uint32
		write   bool
	}{
		{name: "opcode fetch", opcode: 0x3010, address: 0x1000},
		{name: "operand read", opcode: 0x3010, address: 0x3000},
		{name: "operand write", opcode: 0x2080, address: 0x3000, write: true},
		{name: "second word write", opcode: 0x2080, address: 0x3002, write: true},
		{name: "extension fetch", opcode: 0x203C, address: 0x1002},
		{name: "untaken branch extension", opcode: 0x6700, address: 0x1002},
		{name: "stack write", opcode: 0x2F00, address: 0xFFFC, write: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newTestCPU(t)
			bus := &rejectingBus{BasicBus: cpu.bus.(*BasicBus), address: tt.address,
				write: tt.write, reject: true, cause: errors.New("unmapped device")}
			cpu.bus = bus
			cpu.A[0], cpu.D[0] = 0x3000, 0x12345678
			bus.WriteWord(cpu.PC, tt.opcode)
			bus.WriteLong(VectorBusError*4, 0x2000)
			assert.NoError(t, cpu.Step())
			assert.Equal(t, uint32(0x2000), cpu.PC)
			assert.False(t, cpu.Halted())
			assert.Equal(t, tt.address, bus.ReadLong(cpu.A7()+2))
			assert.Equal(t, uint16(0), bus.ReadWord(0x3002))
			if tt.address == 0x3002 {
				assert.Equal(t, uint16(0x1234), bus.ReadWord(0x3000))
			} else {
				assert.Equal(t, uint16(0), bus.ReadWord(0x3000))
			}
		})
	}
}

func TestDoubleFaultRequiresReset(t *testing.T) {
	// A second fault while building the fault frame must halt without recursion.
	cpu := newTestCPU(t)
	bus := &rejectingBus{BasicBus: cpu.bus.(*BasicBus), address: 0x3000,
		reject: true, rejectStack: true, cause: errors.New("bus unavailable")}
	cpu.bus = bus
	cpu.A[0] = 0x3000
	bus.WriteWord(cpu.PC, 0x3010)
	assert.NoError(t, cpu.Step())
	assert.True(t, cpu.Halted())
	cpu.Resume()
	assert.True(t, cpu.Halted())
	accesses := len(bus.accesses)
	assert.NoError(t, cpu.Step())
	assert.Len(t, bus.accesses, accesses)
	bus.reject = false
	bus.WriteLong(0, 0x8000)
	bus.WriteLong(4, 0x2000)
	assert.NoError(t, cpu.Reset())
	assert.False(t, cpu.Halted())
	assert.Equal(t, uint32(0x8000), cpu.A7())
	assert.Equal(t, uint32(0x2000), cpu.PC)
}

func TestResetBusError(t *testing.T) {
	cause := errors.New("reset ROM missing")
	bus := &rejectingBus{BasicBus: NewBasicBus(NewBasicMemory()), address: 4, reject: true, cause: cause}
	cpu, err := New(bus)
	assert.Nil(t, cpu)
	assert.ErrorIs(t, err, ErrBusError)
	assert.ErrorIs(t, err, cause)
}

func TestFaultBoundaryPreservesHostPanics(t *testing.T) {
	assert.Panics(t, func() {
		_ = catchAccessFault(func() error { panic("host bug") })
	})
}

func TestAddressErrorAfterExtension(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.A[0] = 0x3000
	cpu.bus.WriteWord(cpu.PC, 0x3028) // MOVE.W 1(A0),D0.
	cpu.bus.WriteWord(cpu.PC+2, 1)
	cpu.bus.WriteLong(VectorAddressErr*4, 0x2000)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint64(54), cpu.Cycles())
	assert.Equal(t, uint32(0x1002), cpu.bus.ReadLong(cpu.A7()+10))
}

func TestAddressErrorJSRBeforeStackWrite(t *testing.T) {
	// JSR checks its target before pushing a return PC; BSR pushes first.
	cpu := newTestCPU(t)
	cpu.A[0] = 0x3001
	cpu.bus.WriteWord(cpu.PC, 0x4E90)
	cpu.bus.WriteLong(VectorAddressErr*4, 0x2000)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0xFFF2), cpu.A7())
	assert.Equal(t, uint32(0x3001), cpu.bus.ReadLong(cpu.A7()+2))
}

func TestMOVEWriteFaultPreservesPostincrement(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.A[0] = 0x3001
	cpu.bus.WriteWord(cpu.PC, 0x30C0) // MOVE.W D0,(A0)+.
	cpu.bus.WriteLong(VectorAddressErr*4, 0x2000)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x3001), cpu.A[0])
}

func TestStatusPrivilegeBeforeExtensionFetch(t *testing.T) {
	// A forbidden SR write must trap before an inaccessible extension is read.
	for _, opcode := range []uint16{0x007C, 0x027C, 0x0A7C} {
		cpu := newTestCPU(t)
		bus := &rejectingBus{BasicBus: cpu.bus.(*BasicBus), address: 0x1002,
			reject: true, cause: errors.New("unmapped extension")}
		cpu.bus = bus
		cpu.USP = 0x8000
		cpu.SetSR(0)
		bus.WriteWord(cpu.PC, opcode)
		bus.WriteLong(VectorPrivilege*4, 0x2000)
		assert.NoError(t, cpu.Step())
		assert.Equal(t, uint32(0x2000), cpu.PC)
		assert.Equal(t, uint32(0xFFFA), cpu.A7())
		assert.Equal(t, uint32(0x1000), bus.ReadLong(cpu.A7()+2))
	}
}

func TestByteAccessAtOddAddress(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.A[0] = 0x3001
	cpu.bus.WriteWord(cpu.PC, 0x1010) // MOVE.B (A0),D0.
	cpu.bus.Write(0x3001, 0xAB)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0xAB), cpu.D[0])
	assert.Equal(t, uint32(0x1002), cpu.PC)
}

func TestFaultWhileFetchingExceptionVector(t *testing.T) {
	cpu := newTestCPU(t)
	bus := &rejectingBus{BasicBus: cpu.bus.(*BasicBus), address: VectorTrap0 * 4,
		reject: true, cause: errors.New("unmapped vector")}
	cpu.bus = bus
	bus.WriteWord(cpu.PC, 0x4E40) // TRAP #0.
	bus.WriteLong(VectorBusError*4, 0x2000)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x2000), cpu.PC)
	assert.Equal(t, uint32(0xFFEC), cpu.A7()) // Ordinary frame plus bus-error frame.
	assert.Equal(t, uint32(VectorTrap0*4), bus.ReadLong(cpu.A7()+10))
}

func TestDoubleFaultAtHandlerAddress(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.A[0] = 0x3001
	cpu.bus.WriteWord(cpu.PC, 0x3010)
	cpu.bus.WriteLong(VectorAddressErr*4, 0x2001)
	assert.NoError(t, cpu.Step())
	assert.True(t, cpu.Halted())
}

func TestLineExceptionsSaveInstructionAddress(t *testing.T) {
	// Keeping the whole opcode distinguishes line A/F and saves the trapping PC.
	for _, tt := range []struct {
		opcode uint16
		vector int
	}{
		{opcode: 0xA123, vector: VectorLineA},
		{opcode: 0xF123, vector: VectorLineF},
	} {
		cpu := newTestCPU(t)
		cpu.bus.WriteWord(cpu.PC, tt.opcode)
		cpu.bus.WriteLong(uint32(tt.vector)*4, 0x2000)
		assert.NoError(t, cpu.Step())
		assert.Equal(t, uint32(0x2000), cpu.PC)
		assert.Equal(t, uint32(0x1000), cpu.bus.ReadLong(cpu.A7()+2))
	}
}

func TestTraceStartsAfterStatusWrite(t *testing.T) {
	// Setting T must trace the following instruction, not the SR write itself.
	cpu := newTestCPU(t)
	cpu.bus.WriteWord(cpu.PC, 0x007C) // ORI #T,SR.
	cpu.bus.WriteWord(cpu.PC+2, MaskTrace)
	cpu.bus.WriteWord(cpu.PC+4, 0x4E71) // NOP.
	cpu.bus.WriteLong(VectorTrace*4, 0x2000)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x1004), cpu.PC)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x2000), cpu.PC)
	assert.Equal(t, uint32(0x1006), cpu.bus.ReadLong(cpu.A7()+2))
	assert.Equal(t, uint64(58), cpu.Cycles())
}

func TestStatusRegisterReservedBits(t *testing.T) {
	// Bits unimplemented on the original 68000 must read back as zero.
	cpu := newTestCPU(t)
	cpu.SetSR(0xFFFF)
	assert.Equal(t, uint16(0xA71F), cpu.GetSR())
}
