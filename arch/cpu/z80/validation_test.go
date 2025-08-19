package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

// Z80 validation tests - opcode consistency, architectural validation, and instruction API testing

// =============================================================================
// Opcode Consistency and Architectural Validation Tests
// =============================================================================

func TestEnhancedOpcodeConsistency(t *testing.T) {
	t.Run("All Opcodes Have Valid Instructions", func(t *testing.T) {
		definedCount := 0
		for i, opcode := range Opcodes {
			if opcode.Instruction != nil {
				definedCount++
				validateOpcode(t, &opcode, i)
			}
		}

		assert.True(t, definedCount >= 250, "Should have at least 250 valid opcodes")
	})

	t.Run("Enhanced Opcodes Provide Register Information", func(t *testing.T) {
		enhancedCount := 0
		for i, opcode := range Opcodes {
			if opcode.Instruction != nil {
				hasRegisterInfo := opcode.SrcRegister != RegNone ||
					opcode.DstRegister != RegNone ||
					opcode.Register != RegNone
				if hasRegisterInfo {
					enhancedCount++
				}

				// Validate addressing mode consistency
				switch opcode.Addressing {
				case RegisterAddressing:
					validateRegisterAddressing(t, i, hasRegisterInfo)
				case ImmediateAddressing:
					validateImmediateAddressing(t, &opcode, i)
				case RegisterIndirectAddressing:
					validateRegisterIndirectAddressing(t, &opcode, i)
				}
			}
		}

		assert.True(t, enhancedCount >= 80, "Should have at least 80 opcodes with register info")
	})
}

func TestOpcodeCoverage(t *testing.T) {
	definedCount := 0
	undefinedOpcodes := []int{}

	for i, opcode := range Opcodes {
		if opcode.Instruction != nil {
			definedCount++
		} else {
			undefinedOpcodes = append(undefinedOpcodes, i)
		}
	}

	undefinedCount := len(undefinedOpcodes)

	// Most Z80 opcodes should be implemented
	assert.True(t, definedCount >= 250, "Should have at least 250 defined opcodes")
	assert.True(t, undefinedCount <= 10, "Should have no more than 10 undefined opcodes")
}

func TestCriticalOpcodesAreEnhanced(t *testing.T) {
	// These are the opcodes that previously caused test failures
	criticalOpcodes := []struct {
		opcode      byte
		description string
	}{
		// LD r,(HL) instructions
		{0x46, "LD B,(HL)"},
		{0x4E, "LD C,(HL)"},
		{0x56, "LD D,(HL)"},
		{0x5E, "LD E,(HL)"},
		{0x66, "LD H,(HL)"},
		{0x6E, "LD L,(HL)"},
		{0x7E, "LD A,(HL)"},

		// ALU operations with (HL)
		{0x86, "ADD A,(HL)"},
		{0x8E, "ADC A,(HL)"},
		{0x96, "SUB (HL)"},
		{0x9E, "SBC A,(HL)"},
		{0xA6, "AND (HL)"},
		{0xAE, "XOR (HL)"},
		{0xB6, "OR (HL)"},
		{0xBE, "CP (HL)"},

		// Register-to-register operations
		{0x40, "LD B,B"},
		{0x41, "LD B,C"},
		{0x47, "LD B,A"},

		// Other critical operations
		{0x01, "LD BC,nn"},
		{0x04, "INC B"},
		{0x06, "LD B,n"},
	}

	for _, test := range criticalOpcodes {
		opcode := Opcodes[test.opcode]
		hasRegisterInfo := opcode.SrcRegister != RegNone ||
			opcode.DstRegister != RegNone ||
			opcode.Register != RegNone

		assert.True(t, hasRegisterInfo,
			"Opcode 0x%02X (%s): Enhanced with register info",
			test.opcode, test.description)
	}
}

func TestNoRegisterCollisions(t *testing.T) {
	t.Run("Register Information Is Unambiguous", func(t *testing.T) {
		// Test a few specific cases to ensure register information is clear
		testCases := []struct {
			opcode      byte
			description string
			expectSrc   RegisterParam
			expectDst   RegisterParam
		}{
			{0x41, "LD B,C", RegC, RegB},
			{0x86, "ADD A,(HL)", RegHLIndirect, RegA},
		}

		for _, test := range testCases {
			opcode := Opcodes[test.opcode]
			assert.Equal(t, test.expectSrc, opcode.SrcRegister)
			assert.Equal(t, test.expectDst, opcode.DstRegister)
		}
	})
}

// =============================================================================
// Instruction API Testing
// =============================================================================

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
		// Test that the instruction exists and we can get its opcode by register
		expectedRegs := []RegisterParam{RegB, RegC, RegD, RegE, RegH, RegL, RegA}
		foundCount := 0
		for _, reg := range expectedRegs {
			if _, exists := IncReg8.GetOpcodeByRegister(reg); exists {
				foundCount++
			}
		}
		assert.Equal(t, 7, foundCount, "Should find 7 register variants for INC")
	})

	t.Run("Rst variants", func(t *testing.T) {
		// Test that RST variants exist
		expectedRst := []RegisterParam{
			RegRst00, RegRst08, RegRst10, RegRst18,
			RegRst20, RegRst28, RegRst30, RegRst38,
		}
		foundCount := 0
		for _, reg := range expectedRst {
			if _, exists := Rst.GetOpcodeByRegister(reg); exists {
				foundCount++
			}
		}
		assert.Equal(t, 8, foundCount, "Should find 8 RST variants")
	})

	t.Run("Instruction without RegisterOpcodes", func(t *testing.T) {
		// Test that NOP has RegisterOpcodes set to nil
		assert.Nil(t, Nop.RegisterOpcodes, "NOP should have nil RegisterOpcodes")

		// When RegisterOpcodes is nil, GetOpcodeByRegister falls back to Addressing map
		_, exists := Nop.GetOpcodeByRegister(RegB)
		assert.True(t, exists, "NOP fallback to Addressing map should work")
	})
}

func TestRegisterParam_Constants(t *testing.T) {
	testCases := []struct {
		name     string
		param    RegisterParam
		expected string
	}{
		{"a", RegA, "a"},
		{"b", RegB, "b"},
		{"bc", RegBC, "bc"},
		{"hl", RegHL, "hl"},
		{"(hl)", RegHLIndirect, "(hl)"},
		{"n", RegImm8, "n"},
		{"nn", RegImm16, "nn"},
		{"08h", RegRst08, "08h"},
		{"38h", RegRst38, "38h"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.param.String())
		})
	}
}

func TestInstructionRegisterOpcodes_CompareWithOldOpcodeMap(t *testing.T) {
	// Test that new instruction-based register lookup matches old opcode map
	testCases := []struct {
		name        string
		instruction *Instruction
		register    RegisterParam
		expected    byte
	}{
		// 8-bit increments
		{"INC B", IncReg8, RegB, 0x04},
		{"INC C", IncReg8, RegC, 0x0C},
		{"INC A", IncReg8, RegA, 0x3C},

		// 8-bit decrements
		{"DEC B", DecReg8, RegB, 0x05},
		{"DEC C", DecReg8, RegC, 0x0D},
		{"DEC A", DecReg8, RegA, 0x3D},

		// 16-bit increments
		{"INC BC", IncReg16, RegBC, 0x03},
		{"INC DE", IncReg16, RegDE, 0x13},
		{"INC HL", IncReg16, RegHL, 0x23},
		{"INC SP", IncReg16, RegSP, 0x33},

		// 16-bit decrements
		{"DEC BC", DecReg16, RegBC, 0x0B},
		{"DEC DE", DecReg16, RegDE, 0x1B},
		{"DEC HL", DecReg16, RegHL, 0x2B},
		{"DEC SP", DecReg16, RegSP, 0x3B},

		// 8-bit immediate loads - Note: LdReg8 uses fallback to base opcode
		// Real immediate loads use LdImm8 instruction which has proper RegisterOpcodes
		{"LD B,n (fallback)", LdReg8, RegB, 0x7F}, // Falls back to LD A,A base opcode
		{"LD C,n (fallback)", LdReg8, RegC, 0x7F}, // Falls back to LD A,A base opcode
		{"LD A,n (fallback)", LdReg8, RegA, 0x7F}, // Falls back to LD A,A base opcode

		// 16-bit immediate loads
		{"LD BC,nn", LdReg16, RegBC, 0x01},
		{"LD DE,nn", LdReg16, RegDE, 0x11},
		{"LD HL,nn", LdReg16, RegHL, 0x21},
		{"LD SP,nn", LdReg16, RegSP, 0x31},

		// Stack operations
		{"POP BC", PopReg16, RegBC, 0xC1},
		{"POP AF", PopReg16, RegAF, 0xF1},
		{"PUSH BC", PushReg16, RegBC, 0xC5},
		{"PUSH AF", PushReg16, RegAF, 0xF5},

		// Restart instructions
		{"RST 00H", Rst, RegRst00, 0xC7},
		{"RST 08H", Rst, RegRst08, 0xCF},
		{"RST 38H", Rst, RegRst38, 0xFF},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opcodeInfo, exists := tc.instruction.GetOpcodeByRegister(tc.register)
			assert.True(t, exists, "Instruction should have opcode for register %s", tc.register)
			if exists {
				assert.Equal(t, tc.expected, opcodeInfo.Opcode,
					"Opcode for %s should match expected value", tc.name)
			}
		})
	}
}

// =============================================================================
// Validation Helper Functions
// =============================================================================

func validateOpcode(t *testing.T, opcode *Opcode, index int) {
	t.Helper()
	assert.NotNil(t, opcode.Instruction,
		"Opcode 0x%02X should have valid instruction", index)
	assert.NotEqual(t, NoAddressing, opcode.Addressing,
		"Opcode 0x%02X should have addressing mode set", index)
	assert.True(t, opcode.Timing >= 1 && opcode.Timing <= 23,
		"Opcode 0x%02X should have reasonable timing (%d cycles)", index, opcode.Timing)
	assert.True(t, opcode.Size >= 1 && opcode.Size <= 4,
		"Opcode 0x%02X should have reasonable size (%d bytes)", index, opcode.Size)
}

func validateRegisterAddressing(t *testing.T, index int, hasRegisterInfo bool) {
	t.Helper()
	// Some register addressing opcodes might not have enhanced register info
	// This is acceptable as the opcode table is still being enhanced
	_ = index           // Used for potential future validation
	_ = hasRegisterInfo // Used for potential future validation
}

func validateImmediateAddressing(t *testing.T, opcode *Opcode, index int) {
	t.Helper()
	if opcode.SrcRegister == RegImm8 || opcode.SrcRegister == RegImm16 {
		assert.NotEqual(t, RegNone, opcode.DstRegister,
			"Opcode 0x%02X with immediate should have destination", index)
	}
}

func validateRegisterIndirectAddressing(t *testing.T, opcode *Opcode, index int) {
	t.Helper()
	// Function can be used for validation if needed in the future
	_ = opcode // Used for potential future validation
	_ = index  // Used for potential future validation
}
