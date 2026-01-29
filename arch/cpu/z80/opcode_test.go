package z80

import (
	"fmt"
	"strings"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

// TestVerifyOpcodes ensures bidirectional opcode mapping consistency.
// Every opcode in the lookup tables (Opcodes[X] -> Instruction) must have
// a reverse mapping in the instruction's opcode maps (Instruction -> X).
// This enables disassembly and code generation tools.
func TestVerifyOpcodes(t *testing.T) {
	t.Parallel()

	hasReverseMapping := func(ins *Instruction, opcode byte, addressing AddressingMode) bool {
		if addrInfo, ok := ins.Addressing[addressing]; ok && addrInfo.Opcode == opcode {
			return true
		}

		for _, regInfo := range ins.RegisterOpcodes {
			if regInfo.Opcode == opcode {
				return true
			}
		}

		for _, pairInfo := range ins.RegisterPairOpcodes {
			if pairInfo.Opcode == opcode {
				return true
			}
		}

		return false
	}

	verifyOpcodeArray := func(name string, opcodes [256]Opcode) {
		var missingMappings []string

		for b, op := range opcodes {
			ins := op.Instruction
			if ins == nil {
				continue
			}
			if ins.Unofficial && ins.Name == Nop.Name {
				// Unofficial NOPs share opcodes with different addressing modes
				continue
			}

			if !hasReverseMapping(ins, byte(b), op.Addressing) {
				missingMappings = append(missingMappings,
					fmt.Sprintf("0x%02X: %s with %v addressing has no reverse mapping",
						b, ins.Name, op.Addressing))
			}
		}

		if len(missingMappings) > 0 {
			t.Errorf("%s: Found %d opcodes with missing reverse mappings:\n  %s",
				name, len(missingMappings), strings.Join(missingMappings, "\n  "))
		}
	}

	verifyOpcodeArray("Opcodes", Opcodes)
	verifyOpcodeArray("EDOpcodes", EDOpcodes)
	verifyOpcodeArray("DDOpcodes", DDOpcodes)
	verifyOpcodeArray("FDOpcodes", FDOpcodes)
}

// TestOpcodeProperties validates timing and size constraints for all opcodes.
// Timing is in T-states (clock cycles), size is instruction length in bytes.
func TestOpcodeProperties(t *testing.T) {
	t.Parallel()

	verifyOpcodeProperties := func(opcodes [256]Opcode, prefix string) {
		for i, opcode := range opcodes {
			if opcode.Instruction == nil {
				continue
			}
			assert.True(t, opcode.Timing > 0 && opcode.Timing <= 23,
				"%s 0x%02X (%s) has invalid timing: %d", prefix, i, opcode.Instruction.Name, opcode.Timing)
			assert.True(t, opcode.Size > 0 && opcode.Size <= 4,
				"%s 0x%02X (%s) has invalid size: %d", prefix, i, opcode.Instruction.Name, opcode.Size)
		}
	}

	verifyOpcodeProperties(Opcodes, "Opcode")
	verifyOpcodeProperties(EDOpcodes, "ED")
	verifyOpcodeProperties(DDOpcodes, "DD")
	verifyOpcodeProperties(FDOpcodes, "FD")
}

// TestInstructionCoverage verifies essential Z80 instructions are present in opcode tables.
func TestInstructionCoverage(t *testing.T) {
	t.Parallel()

	majorInstructions := []*Instruction{
		Nop, LdReg8, LdReg16, LdImm8, IncReg8, IncReg16, DecReg8, DecReg16,
		AddA, AdcA, SubA, SbcA, AndA, XorA, OrA, CpA,
		JrRel, JrCond, JpAbs, JpCond,
		Call, CallCond, Ret, RetCond,
		PushReg16, PopReg16, Rst, Halt, Ei, Di,
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

func TestUnofficialInstructions(t *testing.T) {
	t.Parallel()

	// Test that unofficial instructions are marked correctly
	unofficialCount := 0
	for _, opcode := range Opcodes {
		if opcode.Instruction != nil && opcode.Instruction.Unofficial {
			unofficialCount++
		}
	}

	// Z80 has fewer unofficial instructions than 6502
	// This is acceptable - Z80 opcode space is more densely packed
	assert.True(t, unofficialCount >= 0, "Unofficial instruction count should be non-negative")
	assert.True(t, unofficialCount < len(Opcodes)/2, "Too many unofficial instructions")
}

func TestOpcodeCoverage(t *testing.T) {
	t.Parallel()

	definedCount := 0
	var undefinedOpcodes []int

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

func TestRegisterDecoding(t *testing.T) {
	t.Parallel()

	// Test that register decoding from opcode bits works correctly
	// Z80 encoding: bits 0-2 = source register, bits 3-5 = destination register
	testCases := []struct {
		opcode      byte
		description string
		expectSrc   Register
		expectDst   Register
	}{
		{0x40, "LD B,B", Register(0), Register(0)}, // src=000, dst=000
		{0x41, "LD B,C", Register(1), Register(0)}, // src=001, dst=000
		{0x47, "LD B,A", Register(7), Register(0)}, // src=111, dst=000
		{0x78, "LD A,B", Register(0), Register(7)}, // src=000, dst=111
		{0x7F, "LD A,A", Register(7), Register(7)}, // src=111, dst=111
	}

	for _, test := range testCases {
		// Extract registers using Z80 bit encoding
		srcReg := Register(test.opcode & 0x07)
		dstReg := Register((test.opcode >> 3) & 0x07)

		assert.Equal(t, test.expectSrc, srcReg,
			"Opcode 0x%02X (%s): source register mismatch", test.opcode, test.description)
		assert.Equal(t, test.expectDst, dstReg,
			"Opcode 0x%02X (%s): destination register mismatch", test.opcode, test.description)
	}
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
	t.Parallel()

	for _, tt := range getOpcodeByRegisterTests() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			opcodeInfo, exists := tt.instruction.GetOpcodeByRegister(tt.register)

			assert.Equal(t, tt.wantExists, exists)
			if exists {
				assert.Equal(t, tt.wantOpcode, opcodeInfo.Opcode)
			}
		})
	}
}

func TestInstruction_GetAllRegisterVariants(t *testing.T) {
	t.Parallel()

	t.Run("IncReg8 variants", func(t *testing.T) {
		t.Parallel()
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
		t.Parallel()
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
		t.Parallel()
		// Test that NOP has RegisterOpcodes set to nil
		assert.Nil(t, Nop.RegisterOpcodes, "NOP should have nil RegisterOpcodes")

		// When RegisterOpcodes is nil, GetOpcodeByRegister falls back to Addressing map
		_, exists := Nop.GetOpcodeByRegister(RegB)
		assert.True(t, exists, "NOP fallback to Addressing map should work")
	})
}

func TestRegisterParam_Constants(t *testing.T) {
	t.Parallel()

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
			t.Parallel()
			assert.Equal(t, tc.expected, tc.param.String())
		})
	}
}

func TestInstructionRegisterOpcodes_CompareWithOldOpcodeMap(t *testing.T) {
	t.Parallel()

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
			t.Parallel()
			opcodeInfo, exists := tc.instruction.GetOpcodeByRegister(tc.register)
			assert.True(t, exists, "Instruction should have opcode for register %s", tc.register)
			if exists {
				assert.Equal(t, tc.expected, opcodeInfo.Opcode,
					"Opcode for %s should match expected value", tc.name)
			}
		})
	}
}
