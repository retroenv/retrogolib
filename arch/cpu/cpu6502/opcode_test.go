package cpu6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

// TestVerifyOpcodes ensures bidirectional opcode mapping consistency.
// Every opcode in the lookup table (Opcodes[X] -> Instruction) must have
// a reverse mapping in the instruction's Addressing map (Instruction -> X).
// This enables disassembly and code generation tools.
func TestVerifyOpcodes(t *testing.T) {
	t.Parallel()

	for b, op := range Opcodes {
		ins := op.Instruction
		if ins == nil {
			continue
		}
		if ins.Unofficial && ins.Name == NopInst.Name {
			// Unofficial NOPs share opcodes with different addressing modes
			continue
		}
		if ins.Name == KilName {
			// KIL/JAM instructions share the same addressing mode across multiple opcodes
			continue
		}

		info := ins.Addressing[op.Addressing]
		assert.Equal(t, b, info.Opcode, "Opcode mismatch for instruction %s with addressing %d", ins.Name, op.Addressing)
	}
}

// TestOpcodeProperties validates timing constraints for all opcodes.
// Timing is in CPU cycles, typically 2-7 cycles for most 6502 instructions.
func TestOpcodeProperties(t *testing.T) {
	t.Parallel()

	for i, opcode := range Opcodes {
		if opcode.Instruction == nil {
			continue
		}
		assert.True(t, opcode.Timing > 0 && opcode.Timing <= 8,
			"Opcode 0x%02X (%s) has invalid timing: %d", i, opcode.Instruction.Name, opcode.Timing)
	}
}

// TestInstructionCoverage verifies essential 6502 instructions are present in opcode table.
func TestInstructionCoverage(t *testing.T) {
	t.Parallel()

	majorInstructions := []*Instruction{
		AdcInst, AndInst, AslInst, BccInst, BcsInst, BeqInst, BitInst, BmiInst, BneInst, BplInst, BrkInst, BvcInst, BvsInst,
		ClcInst, CldInst, CliInst, ClvInst, CmpInst, CpxInst, CpyInst, DecInst, DexInst, DeyInst, EorInst, IncInst, InxInst,
		InyInst, JmpInst, JsrInst, LdaInst, LdxInst, LdyInst, LsrInst, NopInst, OraInst, PhaInst, PhpInst, PlaInst, PlpInst,
		RolInst, RorInst, RtiInst, RtsInst, SbcInst, SecInst, SedInst, SeiInst, StaInst, StxInst, StyInst, TaxInst, TayInst,
		TsxInst, TxaInst, TxsInst, TyaInst,
	}

	for _, ins := range majorInstructions {
		found := false
		for _, opcode := range Opcodes {
			if opcode.Instruction == ins {
				found = true
				break
			}
		}
		assert.True(t, found, "Instruction %s not found in opcodes", ins.Name)
	}
}

// TestUnofficialInstructions validates undocumented opcodes are marked correctly.
func TestUnofficialInstructions(t *testing.T) {
	t.Parallel()

	unofficialCount := 0
	for _, opcode := range Opcodes {
		if opcode.Instruction != nil && opcode.Instruction.Unofficial {
			unofficialCount++
		}
	}

	assert.True(t, unofficialCount > 0, "Expected some unofficial instructions")
	assert.True(t, unofficialCount < len(Opcodes)/2, "Too many unofficial instructions")
}

func TestInstructionRegistriesCoverOpcodeTables(t *testing.T) {
	t.Parallel()

	tables := map[string][256]Opcode{
		"NMOS":     Opcodes,
		"65C02":    Opcodes65C02,
		"Synertek": OpcodesSynertek65C02,
	}
	for name, table := range tables {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			for opcode, info := range table {
				assert.NotNil(t, info.Instruction, "opcode 0x%02x", opcode)
				assert.NotNil(t, Instructions[info.Instruction.Name], "instruction %s", info.Instruction.Name)
			}
		})
	}

	for id := OpcodeID(1); id <= OpcodeIDMax; id++ {
		assert.NotNil(t, InstructionsByID[id], "opcode ID %d (%s)", id, OpcodeIDToName[id])
	}
}

func TestOpcodeMemoryCategories(t *testing.T) {
	t.Parallel()

	assert.False(t, (Opcode{}).ReadsMemory(MemoryReadInstructions))
	assert.False(t, (Opcode{}).WritesMemory(MemoryWriteInstructions))
	assert.False(t, (Opcode{}).ReadWritesMemory(MemoryReadWriteInstructions))

	assert.True(t, Opcodes[0x6d].ReadsMemory(MemoryReadInstructions))
	assert.False(t, Opcodes[0x6d].ReadWritesMemory(MemoryReadWriteInstructions))
	assert.True(t, Opcodes[0x0e].ReadWritesMemory(MemoryReadWriteInstructions))
	assert.False(t, Opcodes[0x0a].ReadWritesMemory(MemoryReadWriteInstructions))
	assert.True(t, Opcodes65C02[0x0f].ReadsMemory(MemoryReadInstructions))
	assert.Equal(t, Bbr0.Name, Opcodes65C02[0x0f].Instruction.Name)
}
