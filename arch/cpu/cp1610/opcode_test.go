package cp1610

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestVerifyOpcodes(t *testing.T) {
	assert.Len(t, Opcodes, 1024)
	for word, op := range Opcodes {
		assert.True(t, op.Instruction != nil, "missing opcode %03x", word)
		found, ok := GetOpcodeInfo(uint16(word))
		assert.True(t, ok, "opcode %03x has no forward mapping", word)
		assert.Equal(t, op, found)
		encoded, ok := GetEncoding(op.Instruction.Name, op.OperandKey)
		assert.True(t, ok, "opcode %03x has no reverse mapping", word)
		assert.Equal(t, uint16(word), encoded.Word)
		assert.Equal(t, op.Size, encoded.Size)
		assert.Equal(t, op.Cycles, encoded.Cycles)
	}
	_, ok := GetOpcodeInfo(0x400)
	assert.False(t, ok)
	_, ok = GetEncoding("unknown", OperandKey{})
	assert.False(t, ok)
}

func TestOpcodeFamilies(t *testing.T) {
	tests := []struct {
		word uint16
		name string
		mode AddressingMode
		size uint8
	}{
		{0x001, SdbdName, ImpliedAddressing, 1},
		{0x004, JumpName, SpecialAddressing, 3},
		{0x00B, IncrName, RegisterAddressing, 1},
		{0x04D, SllName, RegisterAddressing, 1},
		{0x0D3, AddrName, RegisterAddressing, 1},
		{0x224, BeqName, RelativeAddressing, 2},
		{0x246, MvoName, DirectAddressing, 2},
		{0x26E, MvoName, IndirectAddressing, 1},
		{0x27E, MvoName, ImmediateAddressing, 2},
		{0x287, MviName, DirectAddressing, 2},
		{0x2B8, MviName, ImmediateAddressing, 2},
		{0x3FF, XorName, ImmediateAddressing, 2},
	}
	for _, tt := range tests {
		op, ok := GetOpcodeInfo(tt.word)
		assert.True(t, ok)
		assert.Equal(t, tt.name, op.Instruction.Name)
		assert.Equal(t, tt.mode, op.Addressing)
		assert.Equal(t, tt.size, op.Size)
	}
}

func TestOpcodeMemoryCategories(t *testing.T) {
	directRead, _ := GetOpcodeInfo(0x280)
	assert.True(t, directRead.ReadsMemory(MemoryReadInstructions))
	assert.False(t, directRead.WritesMemory(MemoryWriteInstructions))

	indirectWrite, _ := GetOpcodeInfo(0x260)
	assert.True(t, indirectWrite.WritesMemory(MemoryWriteInstructions))
	assert.False(t, indirectWrite.ReadWritesMemory(MemoryReadWriteInstructions))

	immediate, _ := GetOpcodeInfo(0x2B8)
	assert.False(t, immediate.ReadsMemory(MemoryReadInstructions))
	assert.False(t, Opcode{}.WritesMemory(MemoryWriteInstructions))
	assert.False(t, Opcode{}.ReadWritesMemory(MemoryReadWriteInstructions))
}

func TestInstructionCategories(t *testing.T) {
	assert.True(t, BranchingInstructions.Contains(BName))
	assert.True(t, BranchingInstructions.Contains(BextName))
	assert.True(t, NotExecutingFollowingOpcodeInstructions.Contains(JumpName))
	assert.False(t, BranchingInstructions.Contains(NoppName))
	assert.True(t, MemoryReadInstructions.Contains(MviName))
	assert.True(t, MemoryWriteInstructions.Contains(MvoName))
	assert.True(t, MemoryReadWriteInstructions.IsEmpty())
}
