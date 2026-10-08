package cpu65816

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/set"
)

// TestOpcodeTableComplete verifies that all 256 opcode slots are defined.
func TestOpcodeTableComplete(t *testing.T) {
	for i := range 256 {
		op := Opcodes[i]
		assert.NotNil(t, op.Instruction)
	}
}

// TestGetOpcodeInfo verifies the lookup function.
func TestGetOpcodeInfo(t *testing.T) {
	op, ok := GetOpcodeInfo(0xEA) // NOP
	assert.True(t, ok)
	assert.Equal(t, NopName, op.Instruction.Name)
	assert.Equal(t, ImpliedAddressing, op.Addressing)
}

// TestOpcodeTimings verifies that all opcodes have non-zero timing.
func TestOpcodeTimings(t *testing.T) {
	for i := range 256 {
		op := Opcodes[i]
		if op.Instruction == nil {
			continue
		}
		assert.NotEqual(t, byte(0), op.Timing)
	}
}

// TestVerifyOpcodes ensures bidirectional opcode mapping consistency.
// Every opcode in the lookup table (Opcodes[X] -> Instruction) must map back
// to the same opcode byte in the instruction's Addressing map.
func TestVerifyOpcodes(t *testing.T) {
	for i := range 256 {
		op := Opcodes[i]
		info, ok := op.Instruction.Addressing[op.Addressing]
		assert.True(t, ok, "opcode 0x%02X %s: addressing mode %d missing", i, op.Instruction.Name, op.Addressing)
		assert.Equal(t, uint8(i), info.Opcode, "opcode 0x%02X %s: reverse mapping", i, op.Instruction.Name)
	}
}

// TestInstructionRegistryCoversOpcodeTable verifies that the name registry
// and the OpcodeID tables cover every instruction.
func TestInstructionRegistryCoversOpcodeTable(t *testing.T) {
	seen := set.New[string]()
	for i := range 256 {
		name := Opcodes[i].Instruction.Name
		assert.Equal(t, Opcodes[i].Instruction, Instructions[name], "opcode 0x%02X", i)
		seen.Add(name)
	}
	assert.Len(t, Instructions, seen.Size())

	for id := OpcodeID(1); id <= OpcodeIDMax; id++ {
		name := OpcodeIDToName[id]
		assert.NotNil(t, Instructions[name], "opcode ID %d (%s)", id, name)
		assert.Equal(t, id, NameToOpcodeID[name])
	}
	assert.Len(t, Instructions, int(OpcodeIDMax))
}

// TestWidthFlagCorrect verifies that WidthM/WidthX are only set on
// instructions whose immediate operand varies with M or X.
func TestWidthFlagCorrect(t *testing.T) {
	wantWidthM := set.NewFromSlice([]uint8{0x69, 0x29, 0x89, 0xC9, 0x49, 0xA9, 0x09, 0xE9})
	wantWidthX := set.NewFromSlice([]uint8{0xC0, 0xE0, 0xA0, 0xA2})

	for i := range 256 {
		op := Opcodes[i]
		b := uint8(i)
		switch {
		case wantWidthM.Contains(b):
			assert.Equal(t, WidthM, op.WidthFlag, "opcode 0x%02X", i)
		case wantWidthX.Contains(b):
			assert.Equal(t, WidthX, op.WidthFlag, "opcode 0x%02X", i)
		default:
			assert.Equal(t, WidthNone, op.WidthFlag, "opcode 0x%02X", i)
		}
	}
}

// TestCycleRuleFlags verifies that the cycle rule flags follow the datasheet
// notes for each instruction family and addressing mode.
func TestCycleRuleFlags(t *testing.T) {
	accWidth := set.NewFromSlice([]string{AdcName, AndName, BitName, CmpName, EorName, LdaName, OraName, SbcName,
		StaName, StzName, PhaName, PlaName})
	readModifyWrite := set.NewFromSlice([]string{AslName, DecName, IncName, LsrName, RolName, RorName, TsbName, TrbName})
	idxWidth := set.NewFromSlice([]string{CpxName, CpyName, LdxName, LdyName, StxName, StyName, PhxName, PhyName,
		PlxName, PlyName})
	native := set.NewFromSlice([]string{BrkName, CopName, RtiName})
	directPage := set.NewFromSlice([]AddressingMode{DirectPageAddressing, DirectPageIndexedXAddressing,
		DirectPageIndexedYAddressing, DirectPageIndirectAddressing, DirectPageIndexedXIndirectAddressing,
		DirectPageIndirectIndexedYAddressing, DirectPageIndirectLongAddressing,
		DirectPageIndirectLongIndexedYAddressing})
	indexed := set.NewFromSlice([]AddressingMode{AbsoluteIndexedXAddressing, AbsoluteIndexedYAddressing,
		DirectPageIndirectIndexedYAddressing})
	noIndexPenalty := set.NewFromSlice([]string{StaName, StzName})

	for i := range 256 {
		op := Opcodes[i]
		name := op.Instruction.Name
		var want CycleRule
		if accWidth.Contains(name) {
			want |= CycleM
		}
		if readModifyWrite.Contains(name) && op.Addressing != AccumulatorAddressing {
			want |= CycleRMW
		}
		if idxWidth.Contains(name) {
			want |= CycleX
		}
		if directPage.Contains(op.Addressing) {
			want |= CycleDL
		}
		if indexed.Contains(op.Addressing) && !noIndexPenalty.Contains(name) && !readModifyWrite.Contains(name) {
			want |= CycleIndex
		}
		if native.Contains(name) {
			want |= CycleNative
		}
		assert.Equal(t, want, op.Cycles, "opcode 0x%02X %s", i, name)
	}
}

func TestOpcodeMemoryCategories(t *testing.T) {
	assert.True(t, Opcodes[0xAD].ReadsMemory(MemoryReadInstructions))            // LDA abs
	assert.False(t, Opcodes[0xA9].ReadsMemory(MemoryReadInstructions))           // LDA #
	assert.True(t, Opcodes[0x8D].WritesMemory(MemoryWriteInstructions))          // STA abs
	assert.False(t, Opcodes[0x8D].ReadWritesMemory(MemoryReadWriteInstructions)) // STA abs
	assert.True(t, Opcodes[0x0E].ReadWritesMemory(MemoryReadWriteInstructions))  // ASL abs
	assert.False(t, Opcodes[0x0A].ReadWritesMemory(MemoryReadWriteInstructions)) // ASL A
	assert.True(t, Opcodes[0x04].ReadWritesMemory(MemoryReadWriteInstructions))  // TSB dp
	assert.True(t, Opcodes[0xD4].ReadsMemory(MemoryReadInstructions))            // PEI (dp)
	assert.True(t, Opcodes[0x80].IsBranching(BranchingInstructions))             // BRA
	assert.False(t, (Opcode{Instruction: NopInst}).ReadsMemory(MemoryReadInstructions))
}
