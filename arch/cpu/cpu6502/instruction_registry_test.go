package cpu6502

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestInstructionsForVariant(t *testing.T) {
	t.Parallel()

	nmos := InstructionsForVariant(VariantNMOS6502)
	assert.Equal(t, LdaInst, nmos[LdaName])
	assert.Equal(t, NopInst, nmos[NopName])
	_, ok := nmos[Bbr0.Name]
	assert.False(t, ok)

	wdc := InstructionsForVariant(Variant65C02)
	assert.Equal(t, Lda65C02Inst, wdc[LdaName])
	assert.Equal(t, Bbr0, wdc[Bbr0.Name])

	synertek := InstructionsForVariant(VariantSynertek65C02)
	assert.Equal(t, Lda65C02Inst, synertek[LdaName])
	_, ok = synertek[Bbr0.Name]
	assert.False(t, ok)

	assert.Nil(t, InstructionsForVariant(CPUVariant(255)))
}

func TestInstructionIdentityRegistryIsComplete(t *testing.T) {
	t.Parallel()

	for name, instruction := range Instructions {
		id, ok := NameToOpcodeID[name]
		assert.True(t, ok, "instruction %s", name)
		assert.NotEqual(t, InvalidOpcodeID, id, "instruction %s", name)
		assert.Equal(t, name, OpcodeIDToName[id])
		assert.Equal(t, instruction, InstructionsByID[id])
	}
}

func TestInstructionAddressingMetadataHasSizes(t *testing.T) {
	t.Parallel()

	seen := make(map[*Instruction]struct{})
	for _, table := range [][256]Opcode{Opcodes, Opcodes65C02, OpcodesSynertek65C02} {
		for _, opcode := range table {
			instruction := opcode.Instruction
			if instruction == nil {
				continue
			}
			if _, ok := seen[instruction]; ok {
				continue
			}
			seen[instruction] = struct{}{}
			for addressing, info := range instruction.Addressing {
				assert.True(t, info.Size > 0, "instruction %s addressing %d", instruction.Name, addressing)
				assert.True(
					t,
					info.Size <= MaxOpcodeSize,
					"instruction %s addressing %d size %d",
					instruction.Name,
					addressing,
					info.Size,
				)
			}
		}
	}
}
