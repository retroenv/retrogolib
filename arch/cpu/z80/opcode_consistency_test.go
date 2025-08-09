package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

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
	assert.True(t, hasRegisterInfo,
		"Opcode 0x%02X with RegisterAddressing should have register info", index)
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
	hasIndirectRegister := opcode.SrcRegister == RegHLIndirect ||
		opcode.SrcRegister == RegBCIndirect ||
		opcode.SrcRegister == RegDEIndirect ||
		opcode.DstRegister == RegHLIndirect ||
		opcode.DstRegister == RegBCIndirect ||
		opcode.DstRegister == RegDEIndirect

	assert.True(t, hasIndirectRegister,
		"Opcode 0x%02X should specify indirect register", index)
}

func testAllOpcodesHaveValidInstructions(t *testing.T) {
	definedCount := 0
	for i, opcode := range Opcodes {
		if opcode.Instruction != nil {
			definedCount++
			validateOpcode(t, &opcode, i)
		}
	}

	t.Logf("Found %d valid opcodes out of 256 total", definedCount)
	assert.True(t, definedCount >= 240, "Should have most opcodes defined")
}

func testEnhancedOpcodesProvideRegisterInfo(t *testing.T) {
	enhancedCount := 0

	for i, opcode := range Opcodes {
		if opcode.Instruction == nil {
			continue
		}

		hasRegisterInfo := opcode.SrcRegister != RegNone ||
			opcode.DstRegister != RegNone ||
			opcode.Register != RegNone

		if hasRegisterInfo {
			enhancedCount++

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

	t.Logf("Found %d opcodes with enhanced register information", enhancedCount)
	assert.True(t, enhancedCount >= 25, "Should have enhanced register info for key opcodes")
}

// TestEnhancedOpcodeConsistency verifies that the enhanced opcode structure
// provides complete and consistent information for Z80 CPU emulation.
// This replaces the old consistency tests that checked the deprecated
// Addressing and RegisterOpcodes maps.
func TestEnhancedOpcodeConsistency(t *testing.T) {
	t.Run("All Opcodes Have Valid Instructions", testAllOpcodesHaveValidInstructions)
	t.Run("Enhanced Opcodes Provide Register Information", testEnhancedOpcodesProvideRegisterInfo)
}

// TestOpcodeCoverage provides statistics about opcode definition coverage
func TestOpcodeCoverage(t *testing.T) {
	definedCount := 0
	undefinedOpcodes := []byte{}

	for i, opcode := range Opcodes {
		if opcode.Instruction != nil {
			definedCount++
		} else {
			undefinedOpcodes = append(undefinedOpcodes, byte(i))
		}
	}

	undefinedCount := 256 - definedCount

	t.Logf("Opcode coverage summary:")
	t.Logf("  Defined opcodes: %d/256 (%.1f%%)", definedCount, float64(definedCount)/256*100)
	t.Logf("  Undefined opcodes: %d/256 (%.1f%%)", undefinedCount, float64(undefinedCount)/256*100)

	if len(undefinedOpcodes) > 0 && len(undefinedOpcodes) <= 10 {
		t.Logf("  Undefined opcode bytes: %v", undefinedOpcodes)
	}

	// Most Z80 opcodes should be defined (>90%)
	assert.True(t, definedCount >= 230, "Should have most Z80 opcodes defined")
}

// TestCriticalOpcodesAreEnhanced verifies that the opcodes we specifically
// enhanced (the 15 problematic ones) now have complete register information
func TestCriticalOpcodesAreEnhanced(t *testing.T) {
	criticalOpcodes := []struct {
		opcode      byte
		description string
	}{
		// The 15 problematic opcodes from our original test
		{0x46, "LD B,(HL)"},
		{0x4E, "LD C,(HL)"},
		{0x56, "LD D,(HL)"},
		{0x5E, "LD E,(HL)"},
		{0x66, "LD H,(HL)"},
		{0x6E, "LD L,(HL)"},
		{0x7E, "LD A,(HL)"},
		{0x86, "ADD A,(HL)"},
		{0x8E, "ADC A,(HL)"},
		{0x96, "SUB (HL)"},
		{0x9E, "SBC A,(HL)"},
		{0xA6, "AND (HL)"},
		{0xAE, "XOR (HL)"},
		{0xB6, "OR (HL)"},
		{0xBE, "CP (HL)"},

		// Some sample register-to-register loads
		{0x40, "LD B,B"},
		{0x41, "LD B,C"},
		{0x47, "LD B,A"},

		// Some other enhanced opcodes
		{0x01, "LD BC,nn"},
		{0x04, "INC B"},
		{0x06, "LD B,n"},
	}

	for _, test := range criticalOpcodes {
		opcode := Opcodes[test.opcode]

		// Verify the opcode exists and is enhanced
		assert.NotNil(t, opcode.Instruction,
			"Critical opcode 0x%02X (%s) should have instruction", test.opcode, test.description)

		// Verify it has register information
		hasRegisterInfo := opcode.SrcRegister != RegNone ||
			opcode.DstRegister != RegNone ||
			opcode.Register != RegNone
		assert.True(t, hasRegisterInfo,
			"Critical opcode 0x%02X (%s) should have register information", test.opcode, test.description)

		t.Logf("✅ Opcode 0x%02X (%s): Enhanced with register info", test.opcode, test.description)
	}
}

// TestNoRegisterCollisions verifies that the enhanced structure eliminates
// the register collision issues from the old architecture
func TestNoRegisterCollisions(t *testing.T) {
	t.Run("Register Information Is Unambiguous", func(t *testing.T) {
		// Test that register-to-register operations are clearly defined
		ambiguousTests := []struct {
			opcode byte
			desc   string
		}{
			{0x41, "LD B,C"},     // Previously ambiguous which was src/dst
			{0x86, "ADD A,(HL)"}, // Previously unclear about indirect register
		}

		for _, test := range ambiguousTests {
			opcode := Opcodes[test.opcode]

			// These should now have clear register definitions
			if test.opcode == 0x41 {
				assert.Equal(t, RegC, opcode.SrcRegister, "Should clearly identify C as source")
				assert.Equal(t, RegB, opcode.DstRegister, "Should clearly identify B as destination")
			}

			if test.opcode == 0x86 {
				assert.Equal(t, RegHLIndirect, opcode.SrcRegister, "Should clearly identify (HL) as source")
				assert.Equal(t, RegA, opcode.DstRegister, "Should clearly identify A as destination")
			}

			t.Logf("✅ %s (0x%02X): No ambiguity - src=%s, dst=%s, reg=%s",
				test.desc, test.opcode, opcode.SrcRegister, opcode.DstRegister, opcode.Register)
		}
	})
}
