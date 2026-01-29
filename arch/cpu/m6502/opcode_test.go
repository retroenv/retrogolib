package m6502

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
		if ins.Unofficial && ins.Name == Nop.Name {
			// Unofficial NOPs share opcodes with different addressing modes
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
		Adc, And, Asl, Bcc, Bcs, Beq, Bit, Bmi, Bne, Bpl, Brk, Bvc, Bvs,
		Clc, Cld, Cli, Clv, Cmp, Cpx, Cpy, Dec, Dex, Dey, Eor, Inc, Inx,
		Iny, Jmp, Jsr, Lda, Ldx, Ldy, Lsr, Nop, Ora, Pha, Php, Pla, Plp,
		Rol, Ror, Rti, Rts, Sbc, Sec, Sed, Sei, Sta, Stx, Sty, Tax, Tay,
		Tsx, Txa, Txs, Tya,
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
