package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestOpcodeMap_DecReg16_Problem_Solved(t *testing.T) {
	opcodeMap := NewOpcodeMap()
	
	// Test the core problem: multiple DecReg16 opcodes with RegisterAddressing
	// should now be differentiated by their parameters
	
	testCases := []struct {
		name           string
		params         []string
		expectedOpcode byte
	}{
		{"DEC BC", []string{"bc"}, 0x0B},
		{"DEC DE", []string{"de"}, 0x1B},
		{"DEC HL", []string{"hl"}, 0x2B},
		{"DEC SP", []string{"sp"}, 0x3B},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			detail := opcodeMap.GetOpcodeByInstructionAndParams("dec", RegisterAddressing, tc.params)
			assert.NotNil(t, detail, "Expected to find opcode for %s", tc.name)
			assert.Equal(t, tc.expectedOpcode, detail.Opcode)
			assert.Equal(t, tc.params, detail.Params)
		})
	}
}

func TestOpcodeMap_IncReg16_Problem_Solved(t *testing.T) {
	opcodeMap := NewOpcodeMap()
	
	testCases := []struct {
		name           string
		params         []string  
		expectedOpcode byte
	}{
		{"INC BC", []string{"bc"}, 0x03},
		{"INC DE", []string{"de"}, 0x13},
		{"INC HL", []string{"hl"}, 0x23},
		{"INC SP", []string{"sp"}, 0x33},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			detail := opcodeMap.GetOpcodeByInstructionAndParams("inc", RegisterAddressing, tc.params)
			assert.NotNil(t, detail, "Expected to find opcode for %s", tc.name)
			assert.Equal(t, tc.expectedOpcode, detail.Opcode)
			assert.Equal(t, tc.params, detail.Params)
		})
	}
}

func TestOpcodeMap_8BitRegisters(t *testing.T) {
	opcodeMap := NewOpcodeMap()
	
	// Test 8-bit register operations
	testCases := []struct {
		instruction    string
		params         []string
		expectedOpcode byte
	}{
		{"inc", []string{"b"}, 0x04},
		{"inc", []string{"c"}, 0x0C},
		{"inc", []string{"a"}, 0x3C},
		{"dec", []string{"b"}, 0x05},
		{"dec", []string{"c"}, 0x0D},
		{"dec", []string{"a"}, 0x3D},
	}
	
	for _, tc := range testCases {
		t.Run(tc.instruction+" "+tc.params[0], func(t *testing.T) {
			detail := opcodeMap.GetOpcodeByInstructionAndParams(tc.instruction, RegisterAddressing, tc.params)
			assert.NotNil(t, detail, "Expected to find opcode for %s %s", tc.instruction, tc.params[0])
			assert.Equal(t, tc.expectedOpcode, detail.Opcode)
		})
	}
}

func TestOpcodeMap_LoadOperations(t *testing.T) {
	opcodeMap := NewOpcodeMap()
	
	testCases := []struct {
		name           string
		addressing     AddressingMode
		params         []string
		expectedOpcode byte
	}{
		{"LD BC,nn", ImmediateAddressing, []string{"bc", "nn"}, 0x01},
		{"LD A,n", ImmediateAddressing, []string{"a", "n"}, 0x3E},
		{"LD B,C", RegisterAddressing, []string{"b", "c"}, 0x41},
		{"LD (HL),A", RegisterIndirectAddressing, []string{"(hl)", "a"}, 0x77},
		{"LD A,(BC)", RegisterIndirectAddressing, []string{"a", "(bc)"}, 0x0A},
		{"LD (nn),HL", ExtendedAddressing, []string{"(nn)", "hl"}, 0x22},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			detail := opcodeMap.GetOpcodeByInstructionAndParams("ld", tc.addressing, tc.params)
			assert.NotNil(t, detail, "Expected to find opcode for %s", tc.name)
			assert.Equal(t, tc.expectedOpcode, detail.Opcode)
			assert.Equal(t, tc.params, detail.Params)
		})
	}
}

func TestOpcodeMap_GetOpcodeByBytes(t *testing.T) {
	opcodeMap := NewOpcodeMap()
	
	// Test reverse lookup by opcode byte
	testCases := []struct {
		opcode         byte
		expectedInstr  string
		expectedParams []string
	}{
		{0x0B, "dec", []string{"bc"}},
		{0x1B, "dec", []string{"de"}},
		{0x2B, "dec", []string{"hl"}},
		{0x3B, "dec", []string{"sp"}},
		{0x03, "inc", []string{"bc"}},
		{0x3E, "ld", []string{"a", "n"}},
		{0x77, "ld", []string{"(hl)", "a"}},
	}
	
	for _, tc := range testCases {
		t.Run("Opcode 0x"+string(rune('0'+(tc.opcode>>4)))+string(rune('0'+(tc.opcode&0xF))), func(t *testing.T) {
			detail := opcodeMap.GetOpcodeByBytes(tc.opcode)
			assert.NotNil(t, detail, "Expected to find opcode detail for 0x%02X", tc.opcode)
			assert.Equal(t, tc.expectedInstr, detail.Instruction.Name)
			assert.Equal(t, tc.expectedParams, detail.Params)
		})
	}
}

func TestOpcodeMap_GetInstructionVariants(t *testing.T) {
	opcodeMap := NewOpcodeMap()
	
	// Test getting all variants of an instruction
	decVariants := opcodeMap.GetInstructionVariants("dec")
	assert.NotEmpty(t, decVariants, "Expected DEC instruction to have variants")
	
	// Should have both 8-bit and 16-bit DEC variants
	has8BitDec := false
	has16BitDec := false
	
	for _, variant := range decVariants {
		if len(variant.Params) == 1 {
			param := variant.Params[0]
			if param == "bc" || param == "de" || param == "hl" || param == "sp" {
				has16BitDec = true
			} else if param == "a" || param == "b" || param == "c" {
				has8BitDec = true
			}
		}
	}
	
	assert.True(t, has8BitDec, "Expected to find 8-bit DEC variants")
	assert.True(t, has16BitDec, "Expected to find 16-bit DEC variants")
}

func TestOpcodeMap_AssemblerUseCases(t *testing.T) {
	opcodeMap := NewOpcodeMap()
	
	// Test real-world assembler use cases
	testCases := []struct {
		assembly       string
		instruction    string
		addressing     AddressingMode
		params         []string
		expectedOpcode byte
	}{
		{"DEC BC", "dec", RegisterAddressing, []string{"bc"}, 0x0B},
		{"DEC DE", "dec", RegisterAddressing, []string{"de"}, 0x1B},
		{"INC HL", "inc", RegisterAddressing, []string{"hl"}, 0x23},
		{"LD A,n", "ld", ImmediateAddressing, []string{"a", "n"}, 0x3E},
		{"ADD HL,DE", "add", RegisterAddressing, []string{"hl", "de"}, 0x19},
		{"PUSH AF", "push", RegisterAddressing, []string{"af"}, 0xF5},
		{"RST 08H", "rst", ImpliedAddressing, []string{"08h"}, 0xCF},
	}
	
	for _, tc := range testCases {
		t.Run(tc.assembly, func(t *testing.T) {
			detail := opcodeMap.GetOpcodeByInstructionAndParams(tc.instruction, tc.addressing, tc.params)
			assert.NotNil(t, detail, "Expected to find opcode for assembly: %s", tc.assembly)
			assert.Equal(t, tc.expectedOpcode, detail.Opcode, "Wrong opcode for assembly: %s", tc.assembly)
		})
	}
}

func TestOpcodeMap_NoMoreDuplicateKeys(t *testing.T) {
	opcodeMap := NewOpcodeMap()
	
	// Verify that we can now distinguish between all the problematic cases
	// that would have been duplicate keys in the old system
	
	// All DEC with RegisterAddressing should be unique
	decBC := opcodeMap.GetOpcodeByInstructionAndParams("dec", RegisterAddressing, []string{"bc"})
	decDE := opcodeMap.GetOpcodeByInstructionAndParams("dec", RegisterAddressing, []string{"de"})
	decHL := opcodeMap.GetOpcodeByInstructionAndParams("dec", RegisterAddressing, []string{"hl"})
	decSP := opcodeMap.GetOpcodeByInstructionAndParams("dec", RegisterAddressing, []string{"sp"})
	
	assert.NotNil(t, decBC)
	assert.NotNil(t, decDE)
	assert.NotNil(t, decHL)
	assert.NotNil(t, decSP)
	
	// All should have different opcodes
	opcodes := []byte{decBC.Opcode, decDE.Opcode, decHL.Opcode, decSP.Opcode}
	for i := 0; i < len(opcodes); i++ {
		for j := i + 1; j < len(opcodes); j++ {
			assert.NotEqual(t, opcodes[i], opcodes[j], "Opcodes should be unique")
		}
	}
	
	// Verify the actual opcode values
	assert.Equal(t, byte(0x0B), decBC.Opcode)
	assert.Equal(t, byte(0x1B), decDE.Opcode)
	assert.Equal(t, byte(0x2B), decHL.Opcode)
	assert.Equal(t, byte(0x3B), decSP.Opcode)
}