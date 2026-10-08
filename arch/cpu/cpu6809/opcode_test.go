package cpu6809

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/set"
)

// opcodeTables lists the three opcode pages with their prefix byte.
var opcodeTables = []struct {
	name   string
	prefix byte
	table  *[256]Opcode
}{
	{name: "base", prefix: 0x00, table: &Opcodes},
	{name: "page 2", prefix: Prefix10, table: &OpcodesPage2},
	{name: "page 3", prefix: Prefix11, table: &OpcodesPage3},
}

// TestVerifyOpcodes ensures bidirectional opcode mapping consistency over all
// three opcode pages. Every table entry must map back to an OpcodeInfo with the
// same prefix, opcode byte, size and cycle count, every OpcodeInfo must map
// forward to the same instruction pointer, and every reachable instruction
// must be registered under its own name.
func TestVerifyOpcodes(t *testing.T) {
	reachable := set.Set[*Instruction]{}

	for _, page := range opcodeTables {
		for b, op := range page.table {
			ins := op.Instruction
			if ins == nil {
				continue
			}
			reachable[ins] = struct{}{}

			info, ok := ins.Addressing[op.Addressing]
			assert.True(t, ok, "%s opcode 0x%02X: %s has no addressing entry for mode %d",
				page.name, b, ins.Name, op.Addressing)
			assert.Equal(t, page.prefix, info.Prefix, "%s opcode 0x%02X: %s prefix", page.name, b, ins.Name)
			assert.Equal(t, byte(b), info.Opcode, "%s opcode 0x%02X: %s opcode byte", page.name, b, ins.Name)
			assert.Equal(t, op.Size, info.Size, "%s opcode 0x%02X: %s size", page.name, b, ins.Name)
			assert.Equal(t, op.Timing, info.Cycles, "%s opcode 0x%02X: %s cycles", page.name, b, ins.Name)
		}
	}

	for ins := range reachable {
		assert.True(t, Instructions[ins.Name] == ins, "instruction %s is not registered under its name", ins.Name)
		assert.Equal(t, ins.ID, NameToOpcodeID[ins.Name], "instruction %s opcode ID", ins.Name)

		for mode, info := range ins.Addressing {
			table := tableForPrefix(info.Prefix)
			assert.NotNil(t, table, "instruction %s has unknown prefix 0x%02X", ins.Name, info.Prefix)
			if table == nil {
				continue
			}
			op := table[info.Opcode]
			assert.True(t, op.Instruction == ins, "instruction %s opcode 0x%02X 0x%02X maps to a different instruction",
				ins.Name, info.Prefix, info.Opcode)
			assert.Equal(t, mode, op.Addressing, "instruction %s opcode 0x%02X 0x%02X addressing",
				ins.Name, info.Prefix, info.Opcode)
		}
	}

	assert.Len(t, Instructions, len(reachable), "registry size must equal the reachable instruction count")
	for name, ins := range Instructions {
		assert.True(t, reachable.Contains(ins), "registered instruction %s is not in any opcode table", name)
	}
}

func tableForPrefix(prefix byte) *[256]Opcode {
	for _, page := range opcodeTables {
		if page.prefix == prefix {
			return page.table
		}
	}
	return nil
}

// TestGetOpcodeInfo verifies the lookup function.
func TestGetOpcodeInfo(t *testing.T) {
	op, ok := GetOpcodeInfo(0x12) // NOP
	assert.True(t, ok)
	assert.Equal(t, NopName, op.Instruction.Name)
	assert.Equal(t, ImpliedAddressing, op.Addressing)
}

// TestOpcodeTimings verifies that all defined opcodes have non-zero timing.
func TestOpcodeTimings(t *testing.T) {
	for _, page := range opcodeTables {
		for i, op := range page.table {
			if op.Instruction == nil {
				continue
			}
			assert.NotEqual(t, byte(0), op.Timing, "%s opcode 0x%02X has zero timing", page.name, i)
		}
	}
}

// TestOpcodeIDMappingComplete verifies all instruction names have OpcodeIDs.
func TestOpcodeIDMappingComplete(t *testing.T) {
	for name := range Instructions {
		id, ok := NameToOpcodeID[name]
		assert.True(t, ok, "instruction %s missing from NameToOpcodeID", name)
		assert.Equal(t, name, OpcodeIDToName[id])
	}
	for id := OpcodeID(1); id <= OpcodeIDMax; id++ {
		name := OpcodeIDToName[id]
		assert.NotEqual(t, "", name, "opcode ID %d has no name", id)
		assert.Equal(t, id, NameToOpcodeID[name], "opcode ID %d round trip", id)
	}
}

// TestOpcodeIDsAlphabetical verifies that opcode IDs follow the alphabetical mnemonic order.
func TestOpcodeIDsAlphabetical(t *testing.T) {
	for id := OpcodeID(2); id <= OpcodeIDMax; id++ {
		assert.True(t, OpcodeIDToName[id-1] < OpcodeIDToName[id],
			"opcode ID order: %s must precede %s", OpcodeIDToName[id-1], OpcodeIDToName[id])
	}
}

func TestOpcodeInstructionMetadata(t *testing.T) {
	seen := set.Set[*Instruction]{}
	for _, page := range opcodeTables {
		for i, op := range page.table {
			if op.Instruction == nil {
				continue
			}

			ins := op.Instruction
			seen[ins] = struct{}{}
			assert.NotEqual(t, InvalidOpcodeID, ins.ID,
				"%s opcode 0x%02X has no instruction ID", page.name, i)
			assert.Equal(t, ins.Name, OpcodeIDToName[ins.ID],
				"%s opcode 0x%02X has mismatched instruction metadata", page.name, i)
		}
	}

	for ins := range seen {
		hasNoParamHandler := ins.noParamFunc != nil
		hasParamHandler := ins.paramFunc != nil
		assert.True(t, hasNoParamHandler != hasParamHandler,
			"instruction %s must have exactly one handler", ins.Name)
	}
}

// TestInherentAccumulatorFormsAreDistinct verifies that the inherent accumulator
// forms carry their own mnemonic and do not count as memory operations.
func TestInherentAccumulatorFormsAreDistinct(t *testing.T) {
	tests := []struct {
		opcode byte
		name   string
	}{
		{opcode: 0x40, name: NegaName},
		{opcode: 0x4F, name: ClraName},
		{opcode: 0x53, name: CombName},
		{opcode: 0x5C, name: IncbName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, ok := GetOpcodeInfo(tt.opcode)
			assert.True(t, ok)
			assert.Equal(t, tt.name, op.Instruction.Name)
			assert.False(t, MemoryReadWriteInstructions.Contains(tt.name))
			assert.False(t, op.ReadsMemory(MemoryReadInstructions))
			assert.False(t, op.WritesMemory(MemoryWriteInstructions))
		})
	}
}

func TestOpcodeCategories(t *testing.T) {
	undefined := Opcode{}
	assert.False(t, undefined.IsBranching(BranchingInstructions))
	assert.False(t, undefined.ReadsMemory(MemoryReadInstructions))
	assert.False(t, undefined.WritesMemory(MemoryWriteInstructions))

	jmp, ok := GetOpcodeInfo(0x0E)
	assert.True(t, ok)
	assert.True(t, jmp.ReadsMemory(MemoryReadInstructions))
	assert.False(t, jmp.WritesMemory(MemoryWriteInstructions))

	clr, ok := GetOpcodeInfo(0x0F)
	assert.True(t, ok)
	assert.True(t, MemoryReadWriteInstructions.Contains(clr.Instruction.Name))

	inc, ok := GetOpcodeInfo(0x0C)
	assert.True(t, ok)
	assert.True(t, MemoryReadWriteInstructions.Contains(inc.Instruction.Name))
}
