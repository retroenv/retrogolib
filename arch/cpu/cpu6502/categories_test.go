package cpu6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

var memoryOperandCases = []memoryOperandCase{
	{
		name:   "adc",
		opcode: Opcodes[0x65],
		read:   true,
	},
	{
		name:   "sbc",
		opcode: Opcodes[0xe5],
		read:   true,
	},
	{
		name:   "or",
		opcode: Opcodes[0x05],
		read:   true,
	},
	{
		name:   "xor",
		opcode: Opcodes[0x45],
		read:   true,
	},
	{
		name:   "and",
		opcode: Opcodes[0x25],
		read:   true,
	},
	{
		name:   "load",
		opcode: Opcodes[0xa5],
		read:   true,
	},
	{
		name:   "unofficial load",
		opcode: Opcodes[0xbb],
		read:   true,
	},
	{
		name:   "unofficial nop read",
		opcode: Opcodes[0x04],
		read:   true,
	},
	{
		name:   "store",
		opcode: Opcodes[0x85],
		write:  true,
	},
	{
		name:   "unofficial sha",
		opcode: Opcodes[0x9f],
		write:  true,
	},
	{
		name:   "unofficial shx",
		opcode: Opcodes[0x9e],
		write:  true,
	},
	{
		name:   "unofficial shy",
		opcode: Opcodes[0x9c],
		write:  true,
	},
	{
		name:   "unofficial tas",
		opcode: Opcodes[0x9b],
		write:  true,
	},
	{
		name:   "memory shift",
		opcode: Opcodes[0x06],
		read:   true,
		write:  true,
	},
	{
		name:   "memory increment",
		opcode: Opcodes[0xe6],
		read:   true,
		write:  true,
	},
	{
		name:   "unofficial rmw",
		opcode: Opcodes[0x07],
		read:   true,
		write:  true,
	},
	{
		name:   "accumulator shift",
		opcode: Opcodes[0x0a],
	},
	{
		name:   "immediate",
		opcode: Opcodes[0x69],
	},
	{
		name:   "implied",
		opcode: Opcodes[0x1a],
	},
	{
		name:   "branch",
		opcode: Opcodes[0xd0],
	},
	{
		name:   "direct jump",
		opcode: Opcodes[0x4c],
	},
	{
		name:   "indirect jump",
		opcode: Opcodes[0x6c],
		read:   true,
	},
	{
		name:   "unknown",
		opcode: Opcode{},
	},
}

type memoryOperandCase struct {
	name        string
	opcode      Opcode
	read, write bool
}

type categoryMemory struct {
	testMemory
	writes int
}

func (mem *categoryMemory) Write(address uint16, value uint8) {
	mem.testMemory.Write(address, value)
	if address == 0x20 {
		mem.writes++
	}
}

func TestMemoryOperandCategories(t *testing.T) {
	t.Parallel()

	for _, test := range memoryOperandCases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rmw := test.opcode.ReadWritesMemory(MemoryReadWriteInstructions)
			assert.Equal(t, test.read, test.opcode.ReadsMemory(MemoryReadInstructions) || rmw)
			assert.Equal(t, test.write, test.opcode.WritesMemory(MemoryWriteInstructions) || rmw)
		})
	}
}

func TestExtendedInstructionCategories(t *testing.T) {
	t.Parallel()

	for _, instruction := range []*Instruction{BraInst, Bbr0, Bbr7, Bbs0, Bbs7} {
		assert.True(t, BranchingInstructions.Contains(instruction.Name), instruction.Name)
	}
	for _, instruction := range []*Instruction{Bbr0, Bbr7, Bbs0, Bbs7} {
		assert.True(t, MemoryReadInstructions.Contains(instruction.Name), instruction.Name)
		assert.False(t, MemoryReadWriteInstructions.Contains(instruction.Name), instruction.Name)
	}
	for _, instruction := range []*Instruction{Rmb0, Rmb7, Smb0, Smb7} {
		assert.True(t, MemoryReadWriteInstructions.Contains(instruction.Name), instruction.Name)
	}
	for _, instruction := range []*Instruction{AslInst, LsrInst, RolInst, RorInst, Inc65C02Inst, Dec65C02Inst} {
		opcode := Opcode{
			Instruction: instruction,
			Addressing:  AccumulatorAddressing,
		}
		assert.False(t, opcode.ReadWritesMemory(MemoryReadWriteInstructions), instruction.Name)
	}
}

func TestALUMemoryCategoriesMatchExecution(t *testing.T) {
	t.Parallel()

	for _, opcode := range []byte{0x65, 0xe5, 0x05, 0x45} {
		bus := &categoryMemory{}
		memory, err := NewMemory(bus)
		assert.NoError(t, err)
		memory.Write(0x8000, opcode)
		memory.Write(0x8001, 0x20)
		memory.Write(0x20, 2)
		bus.writes = 0
		cpu := newTestCPU(t, memory, WithVariant(VariantNES6502))
		cpu.PC = 0x8000
		cpu.A = 3
		cpu.Flags.C = 1
		assert.NoError(t, cpu.Step())
		assert.Equal(t, byte(2), bus.Read(0x20))
		assert.Equal(t, 0, bus.writes)
		assert.True(t, Opcodes[opcode].ReadsMemory(MemoryReadInstructions))
		assert.False(t, Opcodes[opcode].ReadWritesMemory(MemoryReadWriteInstructions))
	}
}
