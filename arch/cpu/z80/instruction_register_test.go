package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

type opcodeByRegisterTest struct {
	name        string
	instruction *Instruction
	register    RegisterParam
	wantOpcode  byte
	wantExists  bool
}

func getOpcodeByRegisterTests() []opcodeByRegisterTest {
	return []opcodeByRegisterTest{
		{"IncReg8 - INC B", IncReg8, RegB, 0x04, true},
		{"IncReg8 - INC A", IncReg8, RegA, 0x3C, true},
		{"DecReg8 - DEC C", DecReg8, RegC, 0x0D, true},
		{"LdReg16 - LD HL,nn", LdReg16, RegHL, 0x21, true},
		{"IncReg16 - INC SP", IncReg16, RegSP, 0x33, true},
		{"Rst - RST 08H", Rst, RegRst08, 0xCF, true},
		{"PopReg16 - POP AF", PopReg16, RegAF, 0xF1, true},
		{"PushReg16 - PUSH DE", PushReg16, RegDE, 0xD5, true},
		{"Non-existent register", IncReg8, RegIX, 0x00, false},
	}
}

func TestInstruction_GetOpcodeByRegister(t *testing.T) {
	for _, tt := range getOpcodeByRegisterTests() {
		t.Run(tt.name, func(t *testing.T) {
			opcodeInfo, exists := tt.instruction.GetOpcodeByRegister(tt.register)

			assert.Equal(t, tt.wantExists, exists)
			if exists {
				assert.Equal(t, tt.wantOpcode, opcodeInfo.Opcode)
			}
		})
	}
}

func TestInstruction_GetAllRegisterVariants(t *testing.T) {
	t.Run("IncReg8 variants", func(t *testing.T) {
		variants := IncReg8.GetAllRegisterVariants()

		// Should have 7 variants (B, C, D, E, H, L, A)
		assert.Equal(t, 7, len(variants))

		// Check specific variants
		assert.Equal(t, byte(0x04), variants[RegB].Opcode) // INC B
		assert.Equal(t, byte(0x0C), variants[RegC].Opcode) // INC C
		assert.Equal(t, byte(0x3C), variants[RegA].Opcode) // INC A
	})

	t.Run("Rst variants", func(t *testing.T) {
		variants := Rst.GetAllRegisterVariants()

		// Should have 8 RST variants
		assert.Equal(t, 8, len(variants))

		// Check specific variants
		assert.Equal(t, byte(0xC7), variants[RegRst00].Opcode) // RST 00H
		assert.Equal(t, byte(0xCF), variants[RegRst08].Opcode) // RST 08H
		assert.Equal(t, byte(0xFF), variants[RegRst38].Opcode) // RST 38H
	})

	t.Run("Instruction without RegisterOpcodes", func(t *testing.T) {
		variants := Nop.GetAllRegisterVariants()
		assert.Nil(t, variants)
	})
}

func TestRegisterParam_Constants(t *testing.T) {
	// Test that our register constants have the expected string values
	tests := []struct {
		param    RegisterParam
		expected string
	}{
		{RegA, "a"},
		{RegB, "b"},
		{RegBC, "bc"},
		{RegHL, "hl"},
		{RegHLIndirect, "(hl)"},
		{RegImm8, "n"},
		{RegImm16, "nn"},
		{RegRst08, "08h"},
		{RegRst38, "38h"},
	}

	for _, tt := range tests {
		t.Run(tt.param.String(), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.param.String())
		})
	}
}

type opcodeMapComparisonTest struct {
	name         string
	instruction  *Instruction
	register     RegisterParam
	expectedByte byte
}

func getOpcodeMapComparisonTests() []opcodeMapComparisonTest {
	return []opcodeMapComparisonTest{
		// 8-bit increment/decrement instructions
		{"INC B", IncReg8, RegB, 0x04},
		{"INC C", IncReg8, RegC, 0x0C},
		{"INC A", IncReg8, RegA, 0x3C},
		{"DEC B", DecReg8, RegB, 0x05},
		{"DEC C", DecReg8, RegC, 0x0D},
		{"DEC A", DecReg8, RegA, 0x3D},

		// 16-bit increment/decrement instructions
		{"INC BC", IncReg16, RegBC, 0x03},
		{"INC DE", IncReg16, RegDE, 0x13},
		{"INC HL", IncReg16, RegHL, 0x23},
		{"INC SP", IncReg16, RegSP, 0x33},
		{"DEC BC", DecReg16, RegBC, 0x0B},
		{"DEC DE", DecReg16, RegDE, 0x1B},
		{"DEC HL", DecReg16, RegHL, 0x2B},
		{"DEC SP", DecReg16, RegSP, 0x3B},

		// Load instructions
		{"LD B,n", LdImm8, RegB, 0x06},
		{"LD C,n", LdImm8, RegC, 0x0E},
		{"LD A,n", LdImm8, RegA, 0x3E},
		{"LD BC,nn", LdReg16, RegBC, 0x01},
		{"LD DE,nn", LdReg16, RegDE, 0x11},
		{"LD HL,nn", LdReg16, RegHL, 0x21},
		{"LD SP,nn", LdReg16, RegSP, 0x31},

		// Stack and RST operations
		{"POP BC", PopReg16, RegBC, 0xC1},
		{"POP AF", PopReg16, RegAF, 0xF1},
		{"PUSH BC", PushReg16, RegBC, 0xC5},
		{"PUSH AF", PushReg16, RegAF, 0xF5},
		{"RST 00H", Rst, RegRst00, 0xC7},
		{"RST 08H", Rst, RegRst08, 0xCF},
		{"RST 38H", Rst, RegRst38, 0xFF},
	}
}

// TestInstructionRegisterOpcodes_CompareWithOldOpcodeMap demonstrates that the new
// RegisterOpcodes system provides the same functionality as the old OpcodeMap system.
func TestInstructionRegisterOpcodes_CompareWithOldOpcodeMap(t *testing.T) {
	for _, tt := range getOpcodeMapComparisonTests() {
		t.Run(tt.name, func(t *testing.T) {
			opcodeInfo, exists := tt.instruction.GetOpcodeByRegister(tt.register)

			assert.True(t, exists, "Register should exist for instruction")
			assert.Equal(t, tt.expectedByte, opcodeInfo.Opcode, "Opcode should match expected value")
		})
	}
}
