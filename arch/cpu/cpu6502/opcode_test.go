package cpu6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestOpcodeSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		opcode byte
		size   int
	}{
		{
			name:   "official implied",
			opcode: 0xea,
			size:   1,
		},
		{
			name:   "official immediate",
			opcode: 0xa9,
			size:   2,
		},
		{
			name:   "brk has no encoded operand",
			opcode: 0x00,
			size:   1,
		},
		{
			name:   "official absolute",
			opcode: 0x4c,
			size:   3,
		},
		{
			name:   "unofficial implied kil",
			opcode: 0x02,
			size:   1,
		},
		{
			name:   "unofficial zero page dcp",
			opcode: 0xc7,
			size:   2,
		},
		{
			name:   "unofficial absolute lax",
			opcode: 0xaf,
			size:   3,
		},
		{
			name:   "unofficial absolute x isc",
			opcode: 0xff,
			size:   3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.size, Opcodes[tt.opcode].Size())
		})
	}
}

func TestOpcodeSizeMatchesOpcodeInfo(t *testing.T) {
	t.Parallel()

	for b, op := range Opcodes {
		size := op.Size()
		assert.True(t, size >= 1 && size <= MaxOpcodeSize, "opcode $%02X", b)

		// An unofficial opcode can have no size in OpcodeInfo. The size of BRK
		// in OpcodeInfo includes the byte that the CPU skips.
		info := op.Instruction.Addressing[op.Addressing]
		if info.Size != 0 && op.Instruction.Name != BrkName {
			assert.Equal(t, int(info.Size), size, "opcode $%02X", b)
		}
	}
}

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
