package m6809

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

// TestOpcodeTableConsistency verifies each opcode references an instruction that has
// that addressing mode registered.
func TestOpcodeTableConsistency(t *testing.T) {
	for i := range 256 {
		op := Opcodes[i]
		if op.Instruction == nil {
			continue
		}
		_, ok := op.Instruction.Addressing[op.Addressing]
		assert.True(t, ok, "base opcode 0x%02X: instruction %s missing addressing mode %d",
			i, op.Instruction.Name, op.Addressing)
	}
}

// TestPage2OpcodeConsistency verifies page 2 opcodes.
func TestPage2OpcodeConsistency(t *testing.T) {
	for i := range 256 {
		op := OpcodesPage2[i]
		if op.Instruction == nil {
			continue
		}
		_, ok := op.Instruction.Addressing[op.Addressing]
		assert.True(t, ok, "page 2 opcode 0x10 0x%02X: instruction %s missing addressing mode %d",
			i, op.Instruction.Name, op.Addressing)
	}
}

// TestPage3OpcodeConsistency verifies page 3 opcodes.
func TestPage3OpcodeConsistency(t *testing.T) {
	for i := range 256 {
		op := OpcodesPage3[i]
		if op.Instruction == nil {
			continue
		}
		_, ok := op.Instruction.Addressing[op.Addressing]
		assert.True(t, ok, "page 3 opcode 0x11 0x%02X: instruction %s missing addressing mode %d",
			i, op.Instruction.Name, op.Addressing)
	}
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
	for i := range 256 {
		op := Opcodes[i]
		if op.Instruction == nil {
			continue
		}
		assert.NotEqual(t, byte(0), op.Timing, "base opcode 0x%02X has zero timing", i)
	}
}

// TestBidirectionalOpcodeMapping verifies that instruction addressing maps match the opcode tables.
func TestBidirectionalOpcodeMapping(t *testing.T) {
	// Check base page: for each instruction's addressing entry with no prefix,
	// the opcode table should point back to that instruction.
	for name, inst := range Instructions {
		for mode, info := range inst.Addressing {
			if info.Prefix != 0 {
				continue // skip prefixed opcodes
			}
			op := Opcodes[info.Opcode]
			assert.NotNil(t, op.Instruction,
				"instruction %s opcode 0x%02X not in base table", name, info.Opcode)
			if op.Instruction != nil {
				assert.Equal(t, mode, op.Addressing,
					"instruction %s opcode 0x%02X: addressing mismatch", name, info.Opcode)
			}
		}
	}
}

// TestBidirectionalPage2OpcodeMapping verifies page 2 opcode mappings.
func TestBidirectionalPage2OpcodeMapping(t *testing.T) {
	for name, inst := range Instructions {
		for mode, info := range inst.Addressing {
			if info.Prefix != 0x10 {
				continue
			}
			op := OpcodesPage2[info.Opcode]
			assert.NotNil(t, op.Instruction,
				"instruction %s opcode 0x10 0x%02X not in page 2 table", name, info.Opcode)
			if op.Instruction != nil {
				assert.Equal(t, mode, op.Addressing,
					"instruction %s opcode 0x10 0x%02X: addressing mismatch", name, info.Opcode)
			}
		}
	}
}

// TestBidirectionalPage3OpcodeMapping verifies page 3 opcode mappings.
func TestBidirectionalPage3OpcodeMapping(t *testing.T) {
	for name, inst := range Instructions {
		for mode, info := range inst.Addressing {
			if info.Prefix != 0x11 {
				continue
			}
			op := OpcodesPage3[info.Opcode]
			assert.NotNil(t, op.Instruction,
				"instruction %s opcode 0x11 0x%02X not in page 3 table", name, info.Opcode)
			if op.Instruction != nil {
				assert.Equal(t, mode, op.Addressing,
					"instruction %s opcode 0x11 0x%02X: addressing mismatch", name, info.Opcode)
			}
		}
	}
}

// TestOpcodeIDMappingComplete verifies all instruction names have OpcodeIDs.
func TestOpcodeIDMappingComplete(t *testing.T) {
	for name := range Instructions {
		_, ok := NameToOpcodeID[name]
		assert.True(t, ok, "instruction %s missing from NameToOpcodeID", name)
	}
}
