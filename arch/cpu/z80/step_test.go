package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/arch"
	"github.com/retroenv/retrogolib/assert"
)

func TestStepNOP(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory, WithSystemType(arch.GameBoy)) // Game Boy starts at 0x0100

	// Set up NOP instruction at PC
	memory.Write(0x0100, 0x00) // NOP

	initialCycles := cpu.cycles
	initialPC := cpu.PC

	err := cpu.Step()
	assert.NoError(t, err, "Step should not return error for NOP")
	assert.Equal(t, initialPC+1, cpu.PC, "PC should increment by 1 for NOP")
	assert.Equal(t, initialCycles+4, cpu.cycles, "Cycles should increment by 4 for NOP")
}

func TestStepHalt(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory, WithSystemType(arch.GameBoy)) // Game Boy starts at 0x0100

	// Set up HALT instruction at PC
	memory.Write(0x0100, 0x76) // HALT

	assert.False(t, cpu.halted, "CPU should not be halted initially")

	err := cpu.Step()
	assert.NoError(t, err, "Step should not return error for HALT")
	assert.True(t, cpu.halted, "CPU should be halted after HALT instruction")

	// Test that halted CPU just advances cycles
	initialCycles := cpu.cycles
	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error when halted")
	assert.Equal(t, initialCycles+4, cpu.cycles, "Cycles should advance when halted")
}

func TestStepLoadInstructions(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory, WithSystemType(arch.GameBoy))

	// Test LD BC,nn (0x01)
	memory.Write(0x0100, 0x01) // LD BC,nn
	memory.Write(0x0101, 0x34) // Low byte
	memory.Write(0x0102, 0x12) // High byte

	err := cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint16(0x1234), cpu.BC(), "BC should be loaded with 0x1234")
	assert.Equal(t, uint16(0x0103), cpu.PC, "PC should advance by 3")

	// Test LD B,n (0x06)
	memory.Write(0x0103, 0x06) // LD B,n
	memory.Write(0x0104, 0x42) // Value

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint8(0x42), cpu.B, "B should be loaded with 0x42")
	assert.Equal(t, uint16(0x0105), cpu.PC, "PC should advance by 2")
}

func TestStepIncrementDecrement(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory, WithSystemType(arch.GameBoy))

	// Test INC BC (0x03)
	cpu.setBC(0x1234)
	memory.Write(0x0100, 0x03) // INC BC

	err := cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint16(0x1235), cpu.BC(), "BC should be incremented")
	assert.Equal(t, uint16(0x0101), cpu.PC, "PC should advance by 1")

	// Test INC B (0x04)
	cpu.B = 0x10
	memory.Write(0x0101, 0x04) // INC B

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint8(0x11), cpu.B, "B should be incremented")
	assert.Equal(t, uint8(0), cpu.Flags.Z, "Zero flag should not be set")
	assert.Equal(t, uint8(0), cpu.Flags.N, "N flag should be clear for increment")

	// Test DEC B (0x05)
	memory.Write(0x0102, 0x05) // DEC B

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint8(0x10), cpu.B, "B should be decremented")
	assert.Equal(t, uint8(1), cpu.Flags.N, "N flag should be set for decrement")
}

func TestStepMemoryOperations(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory, WithSystemType(arch.GameBoy))

	// Test LD (BC),A (0x02)
	cpu.A = 0x42
	cpu.setBC(0x2000)
	memory.Write(0x0100, 0x02) // LD (BC),A

	err := cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint8(0x42), memory.Read(0x2000), "Memory at BC should contain A")

	// Test LD A,(BC) (0x0A)
	memory.Write(0x2000, 0x55)
	memory.Write(0x0101, 0x0A) // LD A,(BC)

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint8(0x55), cpu.A, "A should be loaded from memory at BC")
}

func TestStepRotateInstructions(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory, WithSystemType(arch.GameBoy))

	// Test RLCA (0x07)
	cpu.A = 0x81
	memory.Write(0x0100, 0x07) // RLCA

	err := cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint8(0x03), cpu.A, "A should be rotated left")
	assert.Equal(t, uint8(1), cpu.Flags.C, "Carry should be set from bit 7")

	// Test RRCA (0x0F)
	cpu.A = 0x81
	memory.Write(0x0101, 0x0F) // RRCA

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint8(0xC0), cpu.A, "A should be rotated right")
	assert.Equal(t, uint8(1), cpu.Flags.C, "Carry should be set from bit 0")
}

func TestStepJumpInstructions(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory, WithSystemType(arch.GameBoy))

	// Test JP nn (0xC3)
	memory.Write(0x0100, 0xC3) // JP nn
	memory.Write(0x0101, 0x00) // Low byte
	memory.Write(0x0102, 0x20) // High byte

	err := cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint16(0x2000), cpu.PC, "PC should jump to 0x2000")

	// Test CALL nn (0xCD)
	cpu.PC = 0x1000
	cpu.SP = 0xFFFE
	memory.Write(0x1000, 0xCD) // CALL nn
	memory.Write(0x1001, 0x00) // Low byte
	memory.Write(0x1002, 0x30) // High byte

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint16(0x3000), cpu.PC, "PC should jump to 0x3000")
	assert.Equal(t, uint16(0xFFFC), cpu.SP, "SP should be decremented by 2")
	assert.Equal(t, uint16(0x1003), memory.ReadWord(cpu.SP), "Return address should be on stack")

	// Test RET (0xC9)
	memory.Write(0x3000, 0xC9) // RET

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint16(0x1003), cpu.PC, "PC should return to saved address")
	assert.Equal(t, uint16(0xFFFE), cpu.SP, "SP should be restored")
}

func TestStepInterruptInstructions(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory, WithSystemType(arch.GameBoy))

	// Test DI (0xF3)
	cpu.iff1 = true
	cpu.iff2 = true
	memory.Write(0x0100, 0xF3) // DI

	err := cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.False(t, cpu.iff1, "IFF1 should be disabled")
	assert.False(t, cpu.iff2, "IFF2 should be disabled")

	// Test EI (0xFB)
	memory.Write(0x0101, 0xFB) // EI

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.True(t, cpu.iff1, "IFF1 should be enabled")
	assert.True(t, cpu.iff2, "IFF2 should be enabled")
}

func TestStepDJNZ(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory, WithSystemType(arch.GameBoy))

	// Test DJNZ with branch taken
	cpu.B = 0x02
	memory.Write(0x0100, 0x10) // DJNZ
	memory.Write(0x0101, 0x05) // Offset +5

	err := cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint8(0x01), cpu.B, "B should be decremented")
	assert.Equal(t, uint16(0x0107), cpu.PC, "PC should branch (0x0102 + 5)")

	// Test DJNZ with branch not taken
	cpu.B = 0x01
	cpu.PC = 0x0200
	memory.Write(0x0200, 0x10) // DJNZ
	memory.Write(0x0201, 0x05) // Offset +5

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint8(0x00), cpu.B, "B should be decremented to 0")
	assert.Equal(t, uint16(0x0202), cpu.PC, "PC should not branch")
}

func TestStepAdd16(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory, WithSystemType(arch.GameBoy))

	// Test ADD HL,BC (0x09)
	cpu.setHL(0x1000)
	cpu.setBC(0x0234)
	memory.Write(0x0100, 0x09) // ADD HL,BC

	err := cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint16(0x1234), cpu.HL(), "HL should be HL + BC")
	assert.Equal(t, uint8(0), cpu.Flags.N, "N flag should be clear for addition")
}

func TestStepExtendedInstructions(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory, WithSystemType(arch.GameBoy))

	// Test CB prefix instruction (CB 00 - RLC B)
	cpu.B = 0x81
	memory.Write(0x0100, 0xCB) // CB prefix
	memory.Write(0x0101, 0x00) // RLC B

	err := cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint8(0x03), cpu.B, "B should be rotated left")
	assert.Equal(t, uint8(1), cpu.Flags.C, "Carry should be set")
	assert.Equal(t, uint16(0x0102), cpu.PC, "PC should advance by 2")

	// Test ED prefix instruction (ED 44 - NEG)
	cpu.A = 0x01
	memory.Write(0x0102, 0xED) // ED prefix
	memory.Write(0x0103, 0x44) // NEG

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint8(0xFF), cpu.A, "A should be negated")
	assert.Equal(t, uint8(1), cpu.Flags.N, "N flag should be set for negation")

	// Test DD prefix instruction (DD 21 - LD IX,nn)
	memory.Write(0x0104, 0xDD) // DD prefix
	memory.Write(0x0105, 0x21) // LD IX,nn
	memory.Write(0x0106, 0x34) // Low byte
	memory.Write(0x0107, 0x12) // High byte

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint16(0x1234), cpu.IX, "IX should be loaded")
	assert.Equal(t, uint16(0x0108), cpu.PC, "PC should advance by 4")
}

func TestStepWithTracing(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory, WithSystemType(arch.GameBoy), WithTracing())

	// Set up instruction
	memory.Write(0x0100, 0x00) // NOP

	// Set some register values for tracing
	cpu.A = 0x42
	cpu.setFlags(0xFF)

	err := cpu.Step()
	assert.NoError(t, err, "Step should not return error")

	// Check trace information
	assert.Equal(t, uint16(0x0100), cpu.TraceStep.PC, "Trace should record PC")
	assert.Equal(t, len(cpu.TraceStep.OpcodeOperands), 1, "Trace should record opcode operands")
	assert.Equal(t, uint8(0x00), cpu.TraceStep.OpcodeOperands[0], "Trace should record opcode")

	// Verify CPU state after execution
	assert.Equal(t, uint8(0x42), cpu.A, "A register should be unchanged")
	assert.Equal(t, uint8(0xFF), cpu.GetFlags(), "Flags should be unchanged")
	assert.Equal(t, uint16(0x0101), cpu.PC, "PC should advance to next instruction")
}

func TestStepErrorHandling(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory, WithSystemType(arch.GameBoy))

	// Test unimplemented opcode (ED prefix with invalid instruction)
	memory.Write(0x0100, 0xED) // ED prefix
	memory.Write(0x0101, 0xFF) // Invalid ED instruction

	err := cpu.Step()
	assert.NotNil(t, err, "Step should return error for unimplemented opcode")
	assert.Contains(t, err.Error(), "unimplemented", "Error should mention unimplemented instruction")
}

func TestRefreshRegister(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory, WithSystemType(arch.GameBoy))

	// R register should increment on each instruction
	initialR := cpu.R
	memory.Write(0x0100, 0x00) // NOP

	err := cpu.Step()
	assert.NoError(t, err, "Step should not return error")

	expectedR := (initialR & 0x80) | ((initialR + 1) & 0x7F)
	assert.Equal(t, expectedR, cpu.R, "R register should increment correctly")
}

func TestEndlessLoopDetection(t *testing.T) {
	memory := NewMemory()
	cpu := New(memory, WithSystemType(arch.GameBoy))

	// Create an endless loop: JR -2 (0x18 0xFE)
	// This instruction jumps back to itself, creating an infinite loop
	memory.Write(0x0100, 0x18) // JR relative
	memory.Write(0x0101, 0xFE) // -2 (0xFE as signed int8 = -2, jumps to 0x0100 + 2 + (-2) = 0x0100)

	initialCycles := cpu.cycles
	initialPC := cpu.PC

	// Execute the jump instruction once
	err := cpu.Step()
	assert.NoError(t, err, "Step should not return error for JR")

	// Verify that PC jumped back to the same address (endless loop)
	assert.Equal(t, initialPC, cpu.PC, "PC should jump back to itself (endless loop)")
	assert.Greater(t, cpu.cycles, initialCycles, "Cycles should have advanced")

	// Execute the same instruction again to confirm it's truly an endless loop
	secondCycles := cpu.cycles
	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error for second iteration")
	assert.Equal(t, initialPC, cpu.PC, "PC should still be at the same address")
	assert.Greater(t, cpu.cycles, secondCycles, "Cycles should have advanced again")

	// Test another endless loop pattern: JP 0x0200 to itself
	cpu.PC = 0x0200
	memory.Write(0x0200, 0xC3) // JP absolute
	memory.Write(0x0201, 0x00) // Low byte of address (0x0200)
	memory.Write(0x0202, 0x02) // High byte of address

	thirdCycles := cpu.cycles
	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error for JP")
	assert.Equal(t, uint16(0x0200), cpu.PC, "PC should jump to same address (endless loop)")
	assert.Greater(t, cpu.cycles, thirdCycles, "Cycles should have advanced for JP")
}
