package m65816

import "testing"

// TestOpcodeTableComplete verifies that all 256 opcode slots are defined.
func TestOpcodeTableComplete(t *testing.T) {
	for i := range 256 {
		op := Opcodes[i]
		if op.Instruction == nil {
			t.Errorf("opcode 0x%02X has nil Instruction", i)
		}
	}
}

// TestGetOpcodeInfo verifies the lookup function.
func TestGetOpcodeInfo(t *testing.T) {
	op, ok := GetOpcodeInfo(0xEA) // NOP
	if !ok {
		t.Fatal("GetOpcodeInfo(0xEA) returned false")
	}
	if op.Instruction.Name != NopName {
		t.Errorf("expected NOP, got %s", op.Instruction.Name)
	}
	if op.Addressing != ImpliedAddressing {
		t.Errorf("expected ImpliedAddressing, got %v", op.Addressing)
	}
}

// TestOpcodeTimings verifies that all opcodes have non-zero timing.
func TestOpcodeTimings(t *testing.T) {
	for i := range 256 {
		op := Opcodes[i]
		if op.Instruction == nil {
			continue
		}
		if op.Timing == 0 {
			t.Errorf("opcode 0x%02X (%s) has zero timing", i, op.Instruction.Name)
		}
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
		if _, ok := op.Instruction.Addressing[op.Addressing]; !ok {
			t.Errorf("opcode 0x%02X (%s): addressing mode %v not in instruction.Addressing map",
				i, op.Instruction.Name, op.Addressing)
		}
	}
}

// TestWidthFlagOnlyForVariableWidthInstructions verifies that WidthM/WidthX
// are only set on instructions that actually vary.
func TestWidthFlagCorrect(t *testing.T) {
	// These instructions should have WidthM
	wantWidthM := map[uint8]bool{
		0x69: true, // ADC #
		0x29: true, // AND #
		0x89: true, // BIT #
		0xC9: true, // CMP #
		0x49: true, // EOR #
		0xA9: true, // LDA #
		0x09: true, // ORA #
		0xE9: true, // SBC #
	}
	// These should have WidthX
	wantWidthX := map[uint8]bool{
		0xC0: true, // CPY #
		0xE0: true, // CPX #
		0xA0: true, // LDY #
		0xA2: true, // LDX #
	}

	for i := range 256 {
		op := Opcodes[i]
		if op.Instruction == nil {
			continue
		}
		b := uint8(i)
		if wantWidthM[b] && op.WidthFlag != WidthM {
			t.Errorf("opcode 0x%02X should have WidthM, has %v", i, op.WidthFlag)
		}
		if wantWidthX[b] && op.WidthFlag != WidthX {
			t.Errorf("opcode 0x%02X should have WidthX, has %v", i, op.WidthFlag)
		}
	}
}
