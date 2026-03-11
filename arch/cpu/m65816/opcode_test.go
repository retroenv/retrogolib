package m65816

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
	"github.com/retroenv/retrogolib/set"
)

// TestOpcodeTableComplete verifies that all 256 opcode slots are defined.
func TestOpcodeTableComplete(t *testing.T) {
	for i := range 256 {
		op := Opcodes[i]
		assert.NotNil(t, op.Instruction)
	}
}

// TestGetOpcodeInfo verifies the lookup function.
func TestGetOpcodeInfo(t *testing.T) {
	op, ok := GetOpcodeInfo(0xEA) // NOP
	assert.True(t, ok)
	assert.Equal(t, NopName, op.Instruction.Name)
	assert.Equal(t, ImpliedAddressing, op.Addressing)
}

// TestOpcodeTimings verifies that all opcodes have non-zero timing.
func TestOpcodeTimings(t *testing.T) {
	for i := range 256 {
		op := Opcodes[i]
		if op.Instruction == nil {
			continue
		}
		assert.NotEqual(t, byte(0), op.Timing)
	}
}

// TestOpcodeConsistency verifies each opcode references an instruction that has
// that addressing mode registered.
func TestOpcodeConsistency(t *testing.T) {
	for i := range 256 {
		op := Opcodes[i]
		if op.Instruction == nil {
			continue
		}
		_, ok := op.Instruction.Addressing[op.Addressing]
		assert.True(t, ok)
	}
}

// TestWidthFlagOnlyForVariableWidthInstructions verifies that WidthM/WidthX
// are only set on instructions that actually vary.
func TestWidthFlagCorrect(t *testing.T) {
	// These instructions should have WidthM
	wantWidthM := set.New[uint8]()
	wantWidthM.Add(0x69) // ADC #
	wantWidthM.Add(0x29) // AND #
	wantWidthM.Add(0x89) // BIT #
	wantWidthM.Add(0xC9) // CMP #
	wantWidthM.Add(0x49) // EOR #
	wantWidthM.Add(0xA9) // LDA #
	wantWidthM.Add(0x09) // ORA #
	wantWidthM.Add(0xE9) // SBC #

	// These should have WidthX
	wantWidthX := set.New[uint8]()
	wantWidthX.Add(0xC0) // CPY #
	wantWidthX.Add(0xE0) // CPX #
	wantWidthX.Add(0xA0) // LDY #
	wantWidthX.Add(0xA2) // LDX #

	for i := range 256 {
		op := Opcodes[i]
		if op.Instruction == nil {
			continue
		}
		b := uint8(i)
		if wantWidthM.Contains(b) && op.WidthFlag != WidthM {
			t.Errorf("opcode 0x%02X should have WidthM, has %v", i, op.WidthFlag)
		}
		if wantWidthX.Contains(b) && op.WidthFlag != WidthX {
			t.Errorf("opcode 0x%02X should have WidthX, has %v", i, op.WidthFlag)
		}
	}
}
