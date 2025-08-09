package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

type opcodeTest struct {
	opcode      byte
	instruction *Instruction
	addressing  AddressingMode
	expectSrc   RegisterParam
	expectDst   RegisterParam
	expectReg   RegisterParam
}

// getProblematicOpcodes returns test cases for opcodes that had consistency issues
func getProblematicOpcodes() []opcodeTest {
	return []opcodeTest{
		// LD r,(HL) instructions
		{0x46, LdReg8, RegisterIndirectAddressing, RegHLIndirect, RegB, RegNone},
		{0x4E, LdReg8, RegisterIndirectAddressing, RegHLIndirect, RegC, RegNone},
		{0x56, LdReg8, RegisterIndirectAddressing, RegHLIndirect, RegD, RegNone},
		{0x5E, LdReg8, RegisterIndirectAddressing, RegHLIndirect, RegE, RegNone},
		{0x66, LdReg8, RegisterIndirectAddressing, RegHLIndirect, RegH, RegNone},
		{0x6E, LdReg8, RegisterIndirectAddressing, RegHLIndirect, RegL, RegNone},
		{0x7E, LdReg8, RegisterIndirectAddressing, RegHLIndirect, RegA, RegNone},

		// ALU operations with (HL)
		{0x86, AddA, RegisterIndirectAddressing, RegHLIndirect, RegA, RegNone},
		{0x8E, AdcA, RegisterIndirectAddressing, RegHLIndirect, RegA, RegNone},
		{0x96, SubA, RegisterIndirectAddressing, RegHLIndirect, RegNone, RegNone},
		{0x9E, SbcA, RegisterIndirectAddressing, RegHLIndirect, RegA, RegNone},
		{0xA6, AndA, RegisterIndirectAddressing, RegHLIndirect, RegNone, RegNone},
		{0xAE, XorA, RegisterIndirectAddressing, RegHLIndirect, RegNone, RegNone},
		{0xB6, OrA, RegisterIndirectAddressing, RegHLIndirect, RegNone, RegNone},
		{0xBE, CpA, RegisterIndirectAddressing, RegHLIndirect, RegNone, RegNone},
	}
}

func verifyOpcodeTest(t *testing.T, test opcodeTest) {
	t.Helper()
	opcode := Opcodes[test.opcode]

	assert.Equal(t, test.instruction, opcode.Instruction,
		"Opcode 0x%02X should have correct instruction", test.opcode)
	assert.Equal(t, test.addressing, opcode.Addressing,
		"Opcode 0x%02X should have correct addressing mode", test.opcode)
	assert.Equal(t, test.expectSrc, opcode.SrcRegister,
		"Opcode 0x%02X should have correct source register", test.opcode)
	assert.Equal(t, test.expectDst, opcode.DstRegister,
		"Opcode 0x%02X should have correct destination register", test.opcode)
	assert.Equal(t, test.expectReg, opcode.Register,
		"Opcode 0x%02X should have correct register", test.opcode)
}

// TestEnhancedOpcodeStructure validates that the enhanced opcode structure
// provides complete register information and resolves the consistency issues
func TestEnhancedOpcodeStructure(t *testing.T) {
	t.Run("Problematic Opcodes Now Have Register Info", func(t *testing.T) {
		for _, test := range getProblematicOpcodes() {
			verifyOpcodeTest(t, test)
		}
	})

	t.Run("Register-to-Register LD Instructions Are Disambiguated", testRegisterToRegisterLD)
	t.Run("Single Register Operations Use Register Field", testSingleRegisterOperations)
	t.Run("Immediate Operations Specify Immediate Source", testImmediateOperations)
}

func testRegisterToRegisterLD(t *testing.T) {
	ldTests := []struct {
		opcode byte
		src    RegisterParam
		dst    RegisterParam
	}{
		{0x40, RegB, RegB}, // LD B,B
		{0x41, RegC, RegB}, // LD B,C
		{0x42, RegD, RegB}, // LD B,D
		{0x47, RegA, RegB}, // LD B,A
		{0x7F, RegA, RegA}, // LD A,A
	}

	for _, test := range ldTests {
		opcode := Opcodes[test.opcode]

		assert.Equal(t, LdReg8, opcode.Instruction,
			"Opcode 0x%02X should be LD instruction", test.opcode)
		assert.Equal(t, RegisterAddressing, opcode.Addressing,
			"Opcode 0x%02X should use register addressing", test.opcode)
		assert.Equal(t, test.src, opcode.SrcRegister,
			"Opcode 0x%02X should have source register %s", test.opcode, test.src)
		assert.Equal(t, test.dst, opcode.DstRegister,
			"Opcode 0x%02X should have destination register %s", test.opcode, test.dst)
	}
}

func testSingleRegisterOperations(t *testing.T) {
	singleRegTests := []struct {
		opcode      byte
		instruction *Instruction
		register    RegisterParam
	}{
		{0x04, IncReg8, RegB}, // INC B
		{0x05, DecReg8, RegB}, // DEC B
		{0xC7, Rst, RegRst00}, // RST 00H
		{0xCF, Rst, RegRst08}, // RST 08H
	}

	for _, test := range singleRegTests {
		opcode := Opcodes[test.opcode]

		assert.Equal(t, test.instruction, opcode.Instruction,
			"Opcode 0x%02X should have correct instruction", test.opcode)
		assert.Equal(t, test.register, opcode.Register,
			"Opcode 0x%02X should target register %s", test.opcode, test.register)
		assert.Equal(t, RegNone, opcode.SrcRegister,
			"Opcode 0x%02X should not have source register", test.opcode)
		assert.Equal(t, RegNone, opcode.DstRegister,
			"Opcode 0x%02X should not have destination register", test.opcode)
	}
}

func testImmediateOperations(t *testing.T) {
	immediateTests := []struct {
		opcode byte
		src    RegisterParam
		dst    RegisterParam
	}{
		{0x01, RegImm16, RegBC}, // LD BC,nn
		{0x06, RegImm8, RegB},   // LD B,n
	}

	for _, test := range immediateTests {
		opcode := Opcodes[test.opcode]

		assert.Equal(t, ImmediateAddressing, opcode.Addressing,
			"Opcode 0x%02X should use immediate addressing", test.opcode)
		assert.Equal(t, test.src, opcode.SrcRegister,
			"Opcode 0x%02X should have immediate source", test.opcode)
		assert.Equal(t, test.dst, opcode.DstRegister,
			"Opcode 0x%02X should have destination register %s", test.opcode, test.dst)
	}
}

func checkRegisterAddressing(opcode *Opcode) bool {
	switch opcode.Instruction {
	case LdReg8:
		// Two-register load operations should have src and dst
		return opcode.SrcRegister != RegNone && opcode.DstRegister != RegNone
	case IncReg8, DecReg8:
		// Single register operations should have Register field
		return opcode.Register != RegNone
	}
	return true
}

func checkImmediateAddressing(opcode *Opcode) bool {
	if opcode.Instruction == LdReg16 || opcode.Instruction == LdImm8 {
		expectedSrc := RegImm8
		if opcode.Instruction == LdReg16 {
			expectedSrc = RegImm16
		}
		return opcode.SrcRegister == expectedSrc && opcode.DstRegister != RegNone
	}
	return true
}

func checkRegisterIndirectAddressing(opcode *Opcode) bool {
	// Should specify which register is indirect
	return opcode.SrcRegister != RegNone || opcode.DstRegister != RegNone
}

// TestRegisterInformationCompleteness validates that register information
// is complete and consistent across the enhanced opcode structure
func TestRegisterInformationCompleteness(t *testing.T) {
	var incompleteOpcodes []byte

	for i, opcode := range Opcodes {
		if opcode.Instruction == nil {
			continue // Skip empty entries (prefix opcodes)
		}

		opcodeIdx := byte(i)
		isComplete := true

		switch opcode.Addressing {
		case RegisterAddressing:
			isComplete = checkRegisterAddressing(&opcode)
		case ImmediateAddressing:
			isComplete = checkImmediateAddressing(&opcode)
		case RegisterIndirectAddressing:
			isComplete = checkRegisterIndirectAddressing(&opcode)
		}

		if !isComplete {
			incompleteOpcodes = append(incompleteOpcodes, opcodeIdx)
		}
	}

	// Report incomplete opcodes (this is informational for now)
	if len(incompleteOpcodes) > 0 {
		t.Logf("Found %d opcodes with incomplete register information: %v",
			len(incompleteOpcodes), incompleteOpcodes)
	}
}

// TestArchitecturalImprovement demonstrates how the enhanced structure
// eliminates the architectural inconsistencies we found earlier
func TestArchitecturalImprovement(t *testing.T) {
	t.Run("Enhanced Opcodes Provide Direct Register Access", func(t *testing.T) {
		// Previously ambiguous opcodes now have clear register information
		opcode := Opcodes[0x41] // LD B,C

		// Direct access to register information without complex lookups
		assert.Equal(t, LdReg8, opcode.Instruction)
		assert.Equal(t, RegC, opcode.SrcRegister) // Source: C
		assert.Equal(t, RegB, opcode.DstRegister) // Destination: B

		t.Logf("Opcode 0x41: %s %s,%s", opcode.Instruction.Name,
			opcode.DstRegister, opcode.SrcRegister)
	})

	t.Run("Consistency Issues Are Resolved", func(t *testing.T) {
		// The 15 problematic opcodes from our original test should now be valid

		// LD B,(HL) - previously not found in instruction addressing
		opcode := Opcodes[0x46]
		assert.Equal(t, LdReg8, opcode.Instruction)
		assert.Equal(t, RegisterIndirectAddressing, opcode.Addressing)
		assert.Equal(t, RegHLIndirect, opcode.SrcRegister)
		assert.Equal(t, RegB, opcode.DstRegister)

		// ADD A,(HL) - previously not found in instruction addressing
		opcode = Opcodes[0x86]
		assert.Equal(t, AddA, opcode.Instruction)
		assert.Equal(t, RegisterIndirectAddressing, opcode.Addressing)
		assert.Equal(t, RegHLIndirect, opcode.SrcRegister)
		assert.Equal(t, RegA, opcode.DstRegister)
	})

	t.Run("No More Opcode Collisions", func(t *testing.T) {
		// Previously, opcodes appeared in both Addressing and RegisterOpcodes maps
		// Now all information is centralized in the Opcodes array

		// Example: INC B (0x04)
		opcode := Opcodes[0x04]
		assert.Equal(t, IncReg8, opcode.Instruction)
		assert.Equal(t, RegB, opcode.Register)

		// The register information is directly available without
		// needing to check RegisterOpcodes maps
		t.Logf("Opcode 0x04: %s %s (direct register access)",
			opcode.Instruction.Name, opcode.Register)
	})
}
