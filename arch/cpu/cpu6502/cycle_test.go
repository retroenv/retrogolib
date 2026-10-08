package cpu6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestBusCycleAccesses(t *testing.T) {
	// Addressing modes must issue their dummy reads and both RMW writes.
	tests := []struct {
		name    string
		program []byte
		want    []BusCycle
	}{
		{"store", []byte{0x8d, 0x11, 0x40}, []BusCycle{
			{Address: 0x8000}, {Address: 0x8001}, {Address: 0x8002}, {Address: 0x4011, Value: 0x42, Write: true},
		}},
		{"indexed store", []byte{0x9d, 0xff, 0x40}, []BusCycle{
			{Address: 0x8000}, {Address: 0x8001}, {Address: 0x8002}, {Address: 0x4000},
			{Address: 0x4100, Value: 0x42, Write: true},
		}},
		{"indexed load", []byte{0xbd, 0xff, 0x40}, []BusCycle{
			{Address: 0x8000}, {Address: 0x8001}, {Address: 0x8002}, {Address: 0x4000}, {Address: 0x4100},
		}},
		{"increment", []byte{0xee, 0x11, 0x40}, []BusCycle{
			{Address: 0x8000}, {Address: 0x8001}, {Address: 0x8002}, {Address: 0x4011},
			{Address: 0x4011, Value: 0x23, Write: true}, {Address: 0x4011, Value: 0x24, Write: true},
		}},
		{"unofficial shift and OR", []byte{0x0f, 0x11, 0x40}, []BusCycle{
			{Address: 0x8000}, {Address: 0x8001}, {Address: 0x8002}, {Address: 0x4011},
			{Address: 0x4011, Value: 0x23, Write: true}, {Address: 0x4011, Value: 0x46, Write: true},
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mem := &ssSparseMemory{data: map[uint16]byte{0x4011: 0x23}}
			for index, value := range test.program {
				mem.Write(0x8000+uint16(index), value)
			}
			memory, err := NewMemory(mem)
			assert.NoError(t, err)
			var accesses []BusCycle
			cpu := New(memory, WithCycleHook(func(cycle BusCycle) bool {
				accesses = append(accesses, cycle)
				return false
			}))
			cpu.PC, cpu.A, cpu.X = 0x8000, 0x42, 1

			assert.NoError(t, cpu.Step())
			assert.Equal(t, test.want, accesses)
			assert.Equal(t, uint64(initialCycles+len(test.want)), cpu.Cycles())
		})
	}
}

// shFamilyCycleCases lists SH-family stores with registers A=0x33, X=0x15, Y=0x01.
// The stored value is register & (high+1).
var shFamilyCycleCases = []struct {
	name    string
	program []byte
	want    []BusCycle
}{
	{"SHX abs,Y", []byte{0x9e, 0x11, 0x40}, []BusCycle{
		{Address: 0x8000}, {Address: 0x8001}, {Address: 0x8002}, {Address: 0x4012},
		{Address: 0x4012, Value: 0x01, Write: true},
	}},
	{"SHX abs,Y page cross", []byte{0x9e, 0xff, 0x40}, []BusCycle{
		{Address: 0x8000}, {Address: 0x8001}, {Address: 0x8002}, {Address: 0x4000},
		{Address: 0x0100, Value: 0x01, Write: true},
	}},
	{"SHY abs,X", []byte{0x9c, 0x11, 0x40}, []BusCycle{
		{Address: 0x8000}, {Address: 0x8001}, {Address: 0x8002}, {Address: 0x4026},
		{Address: 0x4026, Value: 0x01, Write: true},
	}},
	{"SHY abs,X page cross", []byte{0x9c, 0xff, 0x40}, []BusCycle{
		{Address: 0x8000}, {Address: 0x8001}, {Address: 0x8002}, {Address: 0x4014},
		{Address: 0x0114, Value: 0x01, Write: true},
	}},
	{"SHA abs,Y", []byte{0x9f, 0x11, 0x40}, []BusCycle{
		{Address: 0x8000}, {Address: 0x8001}, {Address: 0x8002}, {Address: 0x4012},
		{Address: 0x4012, Value: 0x01, Write: true},
	}},
	{"SHA abs,Y page cross", []byte{0x9f, 0xff, 0x40}, []BusCycle{
		{Address: 0x8000}, {Address: 0x8001}, {Address: 0x8002}, {Address: 0x4000},
		{Address: 0x0100, Value: 0x01, Write: true},
	}},
	{"TAS abs,Y", []byte{0x9b, 0x11, 0x40}, []BusCycle{
		{Address: 0x8000}, {Address: 0x8001}, {Address: 0x8002}, {Address: 0x4012},
		{Address: 0x4012, Value: 0x01, Write: true},
	}},
	{"TAS abs,Y page cross", []byte{0x9b, 0xff, 0x40}, []BusCycle{
		{Address: 0x8000}, {Address: 0x8001}, {Address: 0x8002}, {Address: 0x4000},
		{Address: 0x0100, Value: 0x01, Write: true},
	}},
	{"SHA (ind),Y", []byte{0x93, 0x24}, []BusCycle{
		{Address: 0x8000}, {Address: 0x8001}, {Address: 0x0024}, {Address: 0x0025}, {Address: 0x4012},
		{Address: 0x4012, Value: 0x01, Write: true},
	}},
	{"SHA (ind),Y page cross", []byte{0x93, 0x26}, []BusCycle{
		{Address: 0x8000}, {Address: 0x8001}, {Address: 0x0026}, {Address: 0x0027}, {Address: 0x4000},
		{Address: 0x0100, Value: 0x01, Write: true},
	}},
}

func TestBusCycleSHFamilyWrites(t *testing.T) {
	t.Parallel()

	// The bus-cycle path dispatches the SH-family stores directly to their
	// handlers, so each store must emit exactly one write cycle itself.
	// Registers: A=0x33, X=0x15, Y=0x01. The stored value is register & (high+1).
	for _, test := range shFamilyCycleCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mem := &ssSparseMemory{data: map[uint16]byte{0x24: 0x11, 0x25: 0x40, 0x26: 0xff, 0x27: 0x40}}
			for index, value := range test.program {
				mem.Write(0x8000+uint16(index), value)
			}
			memory, err := NewMemory(mem)
			assert.NoError(t, err)
			var accesses []BusCycle
			cpu := New(memory, WithCycleHook(func(cycle BusCycle) bool {
				accesses = append(accesses, cycle)
				return false
			}))
			cpu.PC, cpu.A, cpu.X, cpu.Y = 0x8000, 0x33, 0x15, 0x01

			assert.NoError(t, cpu.Step())
			assert.Equal(t, test.want, accesses)
			assert.Equal(t, uint64(initialCycles+len(test.want)), cpu.Cycles())
			write := test.want[len(test.want)-1]
			assert.Equal(t, write.Value, mem.Read(write.Address))
			assert.Equal(t, uint16(0x8000+len(test.program)), cpu.PC)
		})
	}
}

func TestBusCycleHookCanHoldReads(t *testing.T) {
	mem := &ssSparseMemory{data: map[uint16]byte{0x8000: 0x8d, 0x8001: 0x11, 0x8002: 0x40}}
	memory, err := NewMemory(mem)
	assert.NoError(t, err)
	var attempts []BusCycle
	cpu := New(memory, WithCycleHook(func(cycle BusCycle) bool {
		attempts = append(attempts, cycle)
		return len(attempts) <= 2 || cycle.Write
	}))
	cpu.PC, cpu.A = 0x8000, 0x42

	assert.NoError(t, cpu.Step())
	assert.Len(t, attempts, 6)
	assert.Equal(t, uint16(0x8000), attempts[0].Address)
	assert.Equal(t, attempts[0], attempts[1])
	assert.Equal(t, attempts[0], attempts[2])
	assert.Equal(t, byte(0x42), mem.Read(0x4011), "a write cannot be held")
}

func TestBusCycleResetPreservesRegisters(t *testing.T) {
	mem := &ssSparseMemory{data: map[uint16]byte{ResetAddress: 0x34, ResetAddress + 1: 0x12}}
	memory, err := NewMemory(mem)
	assert.NoError(t, err)
	var accesses []BusCycle
	cpu := New(memory, WithCycleHook(func(cycle BusCycle) bool {
		accesses = append(accesses, cycle)
		return false
	}))
	cpu.A, cpu.X, cpu.Y, cpu.SP = 1, 2, 3, 0
	cpu.Flags.C, cpu.Flags.I = 1, 0
	cpu.Reset()

	assert.Equal(t, byte(1), cpu.A)
	assert.Equal(t, byte(2), cpu.X)
	assert.Equal(t, byte(3), cpu.Y)
	assert.Equal(t, byte(0xfd), cpu.SP)
	assert.Equal(t, byte(1), cpu.Flags.C)
	assert.Equal(t, byte(1), cpu.Flags.I)
	assert.Equal(t, uint16(0x1234), cpu.PC)
	assert.Len(t, accesses, 7)
	for _, access := range accesses {
		assert.False(t, access.Write)
	}
}

func TestBusCycleInterruptPolling(t *testing.T) {
	tests := []struct {
		name    string
		program []byte
		pc      uint16
		mask    byte
		at      int
		want    bool
	}{
		{"NOP first cycle", []byte{0xea}, 0x8000, 0, 1, true},
		{"NOP last cycle", []byte{0xea}, 0x8000, 0, 2, false},
		{"STA before write", []byte{0x8d, 0, 2}, 0x8000, 0, 3, true},
		{"STA write cycle", []byte{0x8d, 0, 2}, 0x8000, 0, 4, false},
		{"CLI delayed mask", []byte{0x58}, 0x8000, 1, 1, false},
		{"SEI delayed mask", []byte{0x78}, 0x8000, 0, 1, true},
		{"PLP delayed mask", []byte{0x28}, 0x8000, 1, 1, false},
		{"RTI restored mask", []byte{0x40}, 0x8000, 1, 1, true},
		{"branch first cycle", []byte{0xd0, 2}, 0x8000, 0, 1, true},
		{"branch operand cycle", []byte{0xd0, 2}, 0x8000, 0, 2, false},
		{"branch page fixup", []byte{0xd0, 2}, 0x80fc, 0, 3, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mem := &ssSparseMemory{data: map[uint16]byte{0x1fe: 0x20}}
			for i, value := range test.program {
				mem.Write(test.pc+uint16(i), value)
			}
			memory, err := NewMemory(mem)
			assert.NoError(t, err)
			var cpu *CPU
			cycles := 0
			cpu = New(memory, WithCycleHook(func(BusCycle) bool {
				cycles++
				if cycles == test.at {
					cpu.SetIRQ(true)
				}
				return false
			}))
			cpu.PC, cpu.Flags.I = test.pc, test.mask
			assert.NoError(t, cpu.Step())
			assert.Equal(t, test.want, cpu.CheckInterrupts())
		})
	}
}

func TestBusCycleKILRequiresReset(t *testing.T) {
	mem := &ssSparseMemory{data: map[uint16]byte{0x8000: 0x02, ResetAddress + 1: 0x90}}
	memory, err := NewMemory(mem)
	assert.NoError(t, err)
	cpu := New(memory, WithCycleHook(func(BusCycle) bool { return false }))
	cpu.PC = 0x8000
	assert.NoError(t, cpu.Step())
	cpu.TriggerNMI()
	assert.NoError(t, cpu.Step())
	assert.False(t, cpu.CheckInterrupts())
	assert.Equal(t, uint16(0x8001), cpu.PC)
	cpu.Reset()
	assert.False(t, cpu.jammed)
	assert.Equal(t, uint16(0x9000), cpu.PC)
}

func TestBusCycleNMIReplacesBRKVector(t *testing.T) {
	for _, at := range []int{4, 5} {
		mem := &ssSparseMemory{data: map[uint16]byte{NMIAddress + 1: 0x90, IrqAddress + 1: 0xa0}}
		memory, err := NewMemory(mem)
		assert.NoError(t, err)
		var cpu *CPU
		cycles := 0
		cpu = New(memory, WithCycleHook(func(BusCycle) bool {
			cycles++
			if cycles == at {
				cpu.TriggerNMI()
			}
			return false
		}))
		cpu.PC = 0x8000
		assert.NoError(t, cpu.Step())
		expected := uint16(0xa000)
		if at == 4 {
			expected = 0x9000
		}
		assert.Equal(t, expected, cpu.PC)
		assert.Equal(t, byte(0x80), mem.Read(0x1fd))
		assert.Equal(t, byte(0x02), mem.Read(0x1fc))
		assert.Equal(t, byte(0x10), mem.Read(0x1fb)&0x10, "BRK still pushes the B bit")
		assert.False(t, cpu.CheckInterrupts(), "interrupt entry must not poll for a second interrupt")
	}
}
