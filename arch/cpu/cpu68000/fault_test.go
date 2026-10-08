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
			writeLong(cpu.bus, VectorAddressErr*4, 0x2000)
			cpu.USP = 0x8000
			if tt.user {
				cpu.SetSR(0)
			}
			sr := cpu.GetSR()
			assert.NoError(t, cpu.Step())
			assert.Equal(t, uint32(0x2000), cpu.PC)
			assert.Equal(t, uint32(0xFFF2), cpu.A7())
			assert.Equal(t, uint32(0x3001), readLong(cpu.bus, cpu.A7()+2))
			assert.Equal(t, tt.opcode, cpu.bus.ReadWord(cpu.A7()+6))
			assert.Equal(t, sr, cpu.bus.ReadWord(cpu.A7()+8))
			assert.Equal(t, uint32(0x1000), readLong(cpu.bus, cpu.A7()+10))
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
	writeLong(cpu.bus, VectorAddressErr*4, 0x2000)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x2000), cpu.PC)
	assert.Equal(t, uint32(0x1001), readLong(cpu.bus, cpu.A7()+2))
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
			bus := &rejectingBus{
				BasicBus: cpu.bus.(*BasicBus),
				address:  tt.address,
				write:    tt.write,
				reject:   true,
				cause:    errors.New("unmapped device"),
			}
			cpu.bus = bus
			cpu.A[0], cpu.D[0] = 0x3000, 0x12345678
			bus.WriteWord(cpu.PC, tt.opcode)
			writeLong(bus, VectorBusError*4, 0x2000)
			assert.NoError(t, cpu.Step())
			assert.Equal(t, uint32(0x2000), cpu.PC)
			assert.False(t, cpu.Halted())
			assert.Equal(t, tt.address, readLong(bus, cpu.A7()+2))
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
	bus := &rejectingBus{
		BasicBus:    cpu.bus.(*BasicBus),
		address:     0x3000,
		reject:      true,
		rejectStack: true,
		cause:       errors.New("bus unavailable"),
	}
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
	writeLong(bus, 0, 0x8000)
	writeLong(bus, 4, 0x2000)
	assert.NoError(t, cpu.Reset())
	assert.False(t, cpu.Halted())
	assert.Equal(t, uint32(0x8000), cpu.A7())
	assert.Equal(t, uint32(0x2000), cpu.PC)
}

func TestResetBusError(t *testing.T) {
	cause := errors.New("reset ROM missing")
	bus := &rejectingBus{
		BasicBus: NewBasicBus(NewBasicMemory()),
		address:  4,
		reject:   true,
		cause:    cause,
	}
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
	writeLong(cpu.bus, VectorAddressErr*4, 0x2000)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint64(54), cpu.Cycles())
	assert.Equal(t, uint32(0x1002), readLong(cpu.bus, cpu.A7()+10))
}

func TestAddressErrorJSRBeforeStackWrite(t *testing.T) {
	// JSR checks its target before pushing a return PC; BSR pushes first.
	cpu := newTestCPU(t)
	cpu.A[0] = 0x3001
	cpu.bus.WriteWord(cpu.PC, 0x4E90)
	writeLong(cpu.bus, VectorAddressErr*4, 0x2000)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0xFFF2), cpu.A7())
	assert.Equal(t, uint32(0x3001), readLong(cpu.bus, cpu.A7()+2))
}

func TestMOVEWriteFaultPreservesPostincrement(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.A[0] = 0x3001
	cpu.bus.WriteWord(cpu.PC, 0x30C0) // MOVE.W D0,(A0)+.
	writeLong(cpu.bus, VectorAddressErr*4, 0x2000)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x3001), cpu.A[0])
}

func TestStatusPrivilegeBeforeExtensionFetch(t *testing.T) {
	// A forbidden SR write must trap before an inaccessible extension is read.
	for _, opcode := range []uint16{0x007C, 0x027C, 0x0A7C} {
		cpu := newTestCPU(t)
		bus := &rejectingBus{
			BasicBus: cpu.bus.(*BasicBus),
			address:  0x1002,
			reject:   true,
			cause:    errors.New("unmapped extension"),
		}
		cpu.bus = bus
		cpu.USP = 0x8000
		cpu.SetSR(0)
		bus.WriteWord(cpu.PC, opcode)
		writeLong(bus, VectorPrivilege*4, 0x2000)
		assert.NoError(t, cpu.Step())
		assert.Equal(t, uint32(0x2000), cpu.PC)
		assert.Equal(t, uint32(0xFFFA), cpu.A7())
		assert.Equal(t, uint32(0x1000), readLong(bus, cpu.A7()+2))
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
	bus := &rejectingBus{
		BasicBus: cpu.bus.(*BasicBus),
		address:  VectorTrap0 * 4,
		reject:   true,
		cause:    errors.New("unmapped vector"),
	}
	cpu.bus = bus
	bus.WriteWord(cpu.PC, 0x4E40) // TRAP #0.
	writeLong(bus, VectorBusError*4, 0x2000)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x2000), cpu.PC)
	assert.Equal(t, uint32(0xFFEC), cpu.A7()) // Ordinary frame plus bus-error frame.
	assert.Equal(t, uint32(VectorTrap0*4), readLong(bus, cpu.A7()+10))
}

func TestDoubleFaultAtHandlerAddress(t *testing.T) {
	cpu := newTestCPU(t)
	cpu.A[0] = 0x3001
	cpu.bus.WriteWord(cpu.PC, 0x3010)
	writeLong(cpu.bus, VectorAddressErr*4, 0x2001)
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
		writeLong(cpu.bus, uint32(tt.vector)*4, 0x2000)
		assert.NoError(t, cpu.Step())
		assert.Equal(t, uint32(0x2000), cpu.PC)
		assert.Equal(t, uint32(0x1000), readLong(cpu.bus, cpu.A7()+2))
	}
}

func TestUnassignedOpcodeRaisesIllegalInstruction(t *testing.T) {
	// Undecoded words returned a host error instead of taking vector 4.
	for _, opcode := range []uint16{0x4E7A, 0x42C0, 0x4E88} {
		cpu := newTestCPU(t)
		cpu.bus.WriteWord(cpu.PC, opcode)
		writeLong(cpu.bus, VectorIllegal*4, 0x2000)
		assert.NoError(t, cpu.Step())
		assert.Equal(t, uint32(0x2000), cpu.PC)
		assert.Equal(t, uint32(0x1000), readLong(cpu.bus, cpu.A7()+2))
		assert.Equal(t, uint64(34), cpu.Cycles())
	}
}

func TestTraceFollowsTrappingInstruction(t *testing.T) {
	// TRAP, TRAPV, CHK, and zero divide execute, so the trace exception
	// follows their own exception with the handler address as saved PC.
	tests := []struct {
		name   string
		words  []uint16
		vector int
		flags  Flags
	}{
		{name: "TRAP", words: []uint16{0x4E40}, vector: VectorTrap0},
		{name: "TRAPV", words: []uint16{0x4E76}, vector: VectorTRAPV, flags: Flags{V: 1}},
		{name: "CHK", words: []uint16{0x4582}, vector: VectorCHK},      // CHK D2,D2.
		{name: "DIVU", words: []uint16{0x82C0}, vector: VectorDivZero}, // DIVU D0,D1.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newTestCPU(t)
			cpu.sr |= MaskTrace
			cpu.Flags = tt.flags
			cpu.D[2] = 0xFFFF // Negative CHK operand; D0 stays a zero divisor.
			for i, word := range tt.words {
				cpu.bus.WriteWord(cpu.PC+uint32(i)*2, word)
			}
			writeLong(cpu.bus, uint32(tt.vector)*4, 0x2000)
			writeLong(cpu.bus, VectorTrace*4, 0x3000)
			assert.NoError(t, cpu.Step())
			assert.Equal(t, uint32(0x3000), cpu.PC)
			assert.Equal(t, uint32(0x2000), readLong(cpu.bus, cpu.A7()+2))
			assert.Equal(t, uint16(0), cpu.bus.ReadWord(cpu.A7())&MaskTrace)
			assert.Equal(t, uint32(0x1002), readLong(cpu.bus, cpu.A7()+8))
		})
	}
}

func TestTraceDroppedForUnexecutedInstruction(t *testing.T) {
	// Illegal, privileged, and line A/F words do not execute, so no trace follows.
	tests := []struct {
		name   string
		opcode uint16
		vector int
		user   bool
	}{
		{name: "illegal", opcode: 0x4AFC, vector: VectorIllegal},
		{name: "privilege", opcode: 0x4E70, vector: VectorPrivilege, user: true},
		{name: "line A", opcode: 0xA000, vector: VectorLineA},
		{name: "line F", opcode: 0xF000, vector: VectorLineF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newTestCPU(t)
			cpu.USP = 0x8000
			if tt.user {
				cpu.SetSR(0)
			}
			cpu.sr |= MaskTrace
			cpu.bus.WriteWord(cpu.PC, tt.opcode)
			writeLong(cpu.bus, uint32(tt.vector)*4, 0x2000)
			writeLong(cpu.bus, VectorTrace*4, 0x3000)
			assert.NoError(t, cpu.Step())
			assert.Equal(t, uint32(0x2000), cpu.PC)
			assert.Equal(t, uint32(0x1000), readLong(cpu.bus, cpu.A7()+2))
		})
	}
}

func TestSccReadsDestinationBeforeWrite(t *testing.T) {
	// The 68000 performs a read cycle at the Scc destination before the write.
	cpu := newTestCPU(t)
	bus := &rejectingBus{BasicBus: cpu.bus.(*BasicBus)}
	cpu.bus = bus
	cpu.A[0] = 0x3000
	bus.WriteWord(cpu.PC, 0x50D0) // ST (A0).
	assert.NoError(t, cpu.Step())
	assert.Equal(t, []uint32{0x1000, 0x3000, 0x3000}, bus.accesses)
	assert.Equal(t, uint8(0xFF), bus.Read(0x3000))
}

func TestTraceStartsAfterStatusWrite(t *testing.T) {
	// Setting T must trace the following instruction, not the SR write itself.
	cpu := newTestCPU(t)
	cpu.bus.WriteWord(cpu.PC, 0x007C) // ORI #T,SR.
	cpu.bus.WriteWord(cpu.PC+2, MaskTrace)
	cpu.bus.WriteWord(cpu.PC+4, 0x4E71) // NOP.
	writeLong(cpu.bus, VectorTrace*4, 0x2000)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x1004), cpu.PC)
	assert.NoError(t, cpu.Step())
	assert.Equal(t, uint32(0x2000), cpu.PC)
	assert.Equal(t, uint32(0x1006), readLong(cpu.bus, cpu.A7()+2))
	assert.Equal(t, uint64(58), cpu.Cycles())
}

func TestStatusRegisterReservedBits(t *testing.T) {
	// Bits unimplemented on the original 68000 must read back as zero.
	cpu := newTestCPU(t)
	cpu.SetSR(0xFFFF)
	assert.Equal(t, uint16(0xA71F), cpu.GetSR())
}
