package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/arch"
	"github.com/retroenv/retrogolib/assert"
)

// Z80 instruction tests - arithmetic, logical, data movement, control flow, and rotate/shift operations

// =============================================================================
// Arithmetic Operations
// =============================================================================

func TestInc8(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test normal increment
	result := cpu.inc8(0x10)
	assert.Equal(t, uint8(0x11), result, "Increment should work correctly")
	assert.Equal(t, uint8(0), cpu.Flags.Z, "Zero flag should not be set")
	assert.Equal(t, uint8(0), cpu.Flags.S, "Sign flag should not be set")
	assert.Equal(t, uint8(0), cpu.Flags.H, "Half carry should not be set")
	assert.Equal(t, uint8(0), cpu.Flags.N, "N flag should be clear for increment")

	// Test zero result
	result = cpu.inc8(0xFF)
	assert.Equal(t, uint8(0x00), result, "Increment of 0xFF should wrap to 0x00")
	assert.Equal(t, uint8(1), cpu.Flags.Z, "Zero flag should be set")
	assert.Equal(t, uint8(1), cpu.Flags.H, "Half carry should be set (0xF + 1)")

	// Test overflow (0x7F -> 0x80)
	result = cpu.inc8(0x7F)
	assert.Equal(t, uint8(0x80), result, "Increment of 0x7F should be 0x80")
	assert.Equal(t, uint8(1), cpu.Flags.S, "Sign flag should be set")
	assert.Equal(t, uint8(1), cpu.Flags.P, "Parity/overflow flag should be set")
}

func TestDec8(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test normal decrement
	result := cpu.dec8(0x10)
	assert.Equal(t, uint8(0x0F), result, "Decrement should work correctly")
	assert.Equal(t, uint8(0), cpu.Flags.Z, "Zero flag should not be set")
	assert.Equal(t, uint8(0), cpu.Flags.S, "Sign flag should not be set")
	assert.Equal(t, uint8(1), cpu.Flags.N, "N flag should be set for decrement")

	// Test zero result
	result = cpu.dec8(0x01)
	assert.Equal(t, uint8(0x00), result, "Decrement of 0x01 should be 0x00")
	assert.Equal(t, uint8(1), cpu.Flags.Z, "Zero flag should be set")

	// Test underflow (0x80 -> 0x7F)
	result = cpu.dec8(0x80)
	assert.Equal(t, uint8(0x7F), result, "Decrement of 0x80 should be 0x7F")
	assert.Equal(t, uint8(0), cpu.Flags.S, "Sign flag should not be set")
	assert.Equal(t, uint8(1), cpu.Flags.P, "Parity/overflow flag should be set")

	// Test half carry
	cpu.dec8(0x10)
	assert.Equal(t, uint8(1), cpu.Flags.H, "Half carry should be set (0x0 - 1)")
}

func TestAdd8(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test normal addition
	result := cpu.add8(0x10, 0x20)
	assert.Equal(t, uint8(0x30), result, "Addition should work correctly")
	assert.Equal(t, uint8(0), cpu.Flags.C, "Carry should not be set")
	assert.Equal(t, uint8(0), cpu.Flags.Z, "Zero flag should not be set")
	assert.Equal(t, uint8(0), cpu.Flags.N, "N flag should be clear for addition")

	// Test carry
	result = cpu.add8(0xFF, 0x01)
	assert.Equal(t, uint8(0x00), result, "Addition with carry should wrap")
	assert.Equal(t, uint8(1), cpu.Flags.C, "Carry should be set")
	assert.Equal(t, uint8(1), cpu.Flags.Z, "Zero flag should be set")

	// Test half carry
	result = cpu.add8(0x0F, 0x01)
	assert.Equal(t, uint8(0x10), result, "Addition should work correctly")
	assert.Equal(t, uint8(1), cpu.Flags.H, "Half carry should be set")

	// Test overflow (0x7F + 0x01 = 0x80)
	result = cpu.add8(0x7F, 0x01)
	assert.Equal(t, uint8(0x80), result, "Addition should work correctly")
	assert.Equal(t, uint8(1), cpu.Flags.S, "Sign flag should be set")
	assert.Equal(t, uint8(1), cpu.Flags.P, "Overflow flag should be set")
}

func TestSub8(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test normal subtraction
	result := cpu.sub8(0x30, 0x10)
	assert.Equal(t, uint8(0x20), result, "Subtraction should work correctly")
	assert.Equal(t, uint8(0), cpu.Flags.C, "Carry should not be set")
	assert.Equal(t, uint8(1), cpu.Flags.N, "N flag should be set for subtraction")

	// Test borrow
	result = cpu.sub8(0x00, 0x01)
	assert.Equal(t, uint8(0xFF), result, "Subtraction with borrow should wrap")
	assert.Equal(t, uint8(1), cpu.Flags.C, "Carry should be set for borrow")

	// Test zero result
	result = cpu.sub8(0x42, 0x42)
	assert.Equal(t, uint8(0x00), result, "Subtraction should yield zero")
	assert.Equal(t, uint8(1), cpu.Flags.Z, "Zero flag should be set")
}

func TestArithmeticWithCarry(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test ADC (add with carry)
	cpu.Flags.C = 1
	result := cpu.adc(0x10, 0x20)
	assert.Equal(t, uint8(0x31), result, "ADC should include carry")

	cpu.Flags.C = 0
	result = cpu.adc(0x10, 0x20)
	assert.Equal(t, uint8(0x30), result, "ADC without carry should work normally")

	// Test SBC (subtract with carry)
	cpu.Flags.C = 1
	result = cpu.sbc(0x30, 0x10)
	assert.Equal(t, uint8(0x1F), result, "SBC should include carry")

	cpu.Flags.C = 0
	result = cpu.sbc(0x30, 0x10)
	assert.Equal(t, uint8(0x20), result, "SBC without carry should work normally")
}

func TestNegation(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test NEG with positive number
	result := cpu.neg(0x01)
	assert.Equal(t, uint8(0xFF), result, "NEG should work correctly")
	assert.Equal(t, uint8(1), cpu.Flags.C, "Carry should be set")
	assert.Equal(t, uint8(1), cpu.Flags.S, "Sign flag should be set")
	assert.Equal(t, uint8(1), cpu.Flags.N, "N flag should be set for negation")

	// Test NEG with zero
	result = cpu.neg(0x00)
	assert.Equal(t, uint8(0x00), result, "NEG of zero should be zero")
	assert.Equal(t, uint8(0), cpu.Flags.C, "Carry should not be set for zero")
	assert.Equal(t, uint8(1), cpu.Flags.Z, "Zero flag should be set")

	// Test NEG with 0x80 (overflow case)
	result = cpu.neg(0x80)
	assert.Equal(t, uint8(0x80), result, "NEG of 0x80 should be 0x80 (overflow)")
	assert.Equal(t, uint8(1), cpu.Flags.P, "Overflow flag should be set")
}

func TestCompare(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test equal values
	cpu.cp(0x42, 0x42)
	assert.Equal(t, uint8(1), cpu.Flags.Z, "Zero flag should be set for equal values")
	assert.Equal(t, uint8(0), cpu.Flags.C, "Carry should not be set")
	assert.Equal(t, uint8(1), cpu.Flags.N, "N flag should be set for compare")

	// Test first value greater
	cpu.cp(0x50, 0x30)
	assert.Equal(t, uint8(0), cpu.Flags.Z, "Zero flag should not be set")
	assert.Equal(t, uint8(0), cpu.Flags.C, "Carry should not be set")

	// Test first value less (borrow)
	cpu.cp(0x30, 0x50)
	assert.Equal(t, uint8(0), cpu.Flags.Z, "Zero flag should not be set")
	assert.Equal(t, uint8(1), cpu.Flags.C, "Carry should be set for borrow")
}

// =============================================================================
// Logical Operations
// =============================================================================

func TestLogicalOperations(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test AND
	result := cpu.and8(0xF0, 0x0F)
	assert.Equal(t, uint8(0x00), result, "AND should work correctly")
	assert.Equal(t, uint8(1), cpu.Flags.Z, "Zero flag should be set")
	assert.Equal(t, uint8(1), cpu.Flags.H, "Half carry should be set for AND")
	assert.Equal(t, uint8(0), cpu.Flags.C, "Carry should be clear")
	assert.Equal(t, uint8(0), cpu.Flags.N, "N flag should be clear")

	result = cpu.and8(0xFF, 0xAA)
	assert.Equal(t, uint8(0xAA), result, "AND should work correctly")
	assert.Equal(t, uint8(1), cpu.Flags.S, "Sign flag should be set")

	// Test OR
	result = cpu.or8(0xF0, 0x0F)
	assert.Equal(t, uint8(0xFF), result, "OR should work correctly")
	assert.Equal(t, uint8(0), cpu.Flags.Z, "Zero flag should not be set")
	assert.Equal(t, uint8(1), cpu.Flags.S, "Sign flag should be set")
	assert.Equal(t, uint8(0), cpu.Flags.H, "Half carry should be clear")
	assert.Equal(t, uint8(0), cpu.Flags.C, "Carry should be clear")

	// Test XOR
	result = cpu.xor8(0xFF, 0xFF)
	assert.Equal(t, uint8(0x00), result, "XOR should work correctly")
	assert.Equal(t, uint8(1), cpu.Flags.Z, "Zero flag should be set")
	assert.Equal(t, uint8(0), cpu.Flags.S, "Sign flag should not be set")

	result = cpu.xor8(0xAA, 0x55)
	assert.Equal(t, uint8(0xFF), result, "XOR should work correctly")
	assert.Equal(t, uint8(1), cpu.Flags.S, "Sign flag should be set")
}

func TestBitOperations(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test BIT instruction
	cpu.bit(7, 0x80)
	assert.Equal(t, uint8(0), cpu.Flags.Z, "Zero flag should not be set when bit is 1")
	assert.Equal(t, uint8(1), cpu.Flags.H, "Half carry should be set for BIT")
	assert.Equal(t, uint8(0), cpu.Flags.N, "N flag should be clear")

	cpu.bit(7, 0x7F)
	assert.Equal(t, uint8(1), cpu.Flags.Z, "Zero flag should be set when bit is 0")

	// Test SET instruction
	result := cpu.setBit(3, 0x00)
	assert.Equal(t, uint8(0x08), result, "SET should set the specified bit")

	result = cpu.setBit(7, 0x7F)
	assert.Equal(t, uint8(0xFF), result, "SET should set the specified bit")

	// Test RES instruction
	result = cpu.res(3, 0xFF)
	assert.Equal(t, uint8(0xF7), result, "RES should clear the specified bit")

	result = cpu.res(7, 0x80)
	assert.Equal(t, uint8(0x00), result, "RES should clear the specified bit")
}

// =============================================================================
// Rotate and Shift Operations
// =============================================================================

func TestRotateOperations(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test RLCA (rotate left circular, accumulator)
	result := cpu.rlca(0x81)
	assert.Equal(t, uint8(0x03), result, "RLCA should work correctly")
	assert.Equal(t, uint8(1), cpu.Flags.C, "Carry should be set from bit 7")
	assert.Equal(t, uint8(0), cpu.Flags.H, "Half carry should be clear")
	assert.Equal(t, uint8(0), cpu.Flags.N, "N flag should be clear")

	// Test RRCA (rotate right circular, accumulator)
	result = cpu.rrca(0x81)
	assert.Equal(t, uint8(0xC0), result, "RRCA should work correctly")
	assert.Equal(t, uint8(1), cpu.Flags.C, "Carry should be set from bit 0")

	// Test RLC (rotate left circular with flags)
	result = cpu.rlc(0x00)
	assert.Equal(t, uint8(0x00), result, "RLC should work correctly")
	assert.Equal(t, uint8(1), cpu.Flags.Z, "Zero flag should be set")
	assert.Equal(t, uint8(0), cpu.Flags.C, "Carry should not be set")

	// Test with carry
	cpu.Flags.C = 1
	result = cpu.rl(0x80)
	assert.Equal(t, uint8(0x01), result, "RL should include old carry")
	assert.Equal(t, uint8(1), cpu.Flags.C, "Carry should be set from bit 7")
}

func TestShiftOperations(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test SLA (shift left arithmetic)
	result := cpu.sla(0x81)
	assert.Equal(t, uint8(0x02), result, "SLA should work correctly")
	assert.Equal(t, uint8(1), cpu.Flags.C, "Carry should be set from bit 7")

	// Test SRA (shift right arithmetic, keep sign)
	result = cpu.sra(0x81)
	assert.Equal(t, uint8(0xC0), result, "SRA should keep sign bit")
	assert.Equal(t, uint8(1), cpu.Flags.C, "Carry should be set from bit 0")

	// Test SRL (shift right logical)
	result = cpu.srl(0x81)
	assert.Equal(t, uint8(0x40), result, "SRL should clear sign bit")
	assert.Equal(t, uint8(1), cpu.Flags.C, "Carry should be set from bit 0")
}

// =============================================================================
// Data Movement Operations
// =============================================================================

func TestRegisterPairs(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test BC register pair
	cpu.B = 0x12
	cpu.C = 0x34
	assert.Equal(t, uint16(0x1234), cpu.BC(), "BC register pair should return correct value")

	cpu.setBC(0x5678)
	assert.Equal(t, uint8(0x56), cpu.B, "B should be set correctly")
	assert.Equal(t, uint8(0x78), cpu.C, "C should be set correctly")

	// Test DE register pair
	cpu.D = 0xAB
	cpu.E = 0xCD
	assert.Equal(t, uint16(0xABCD), cpu.DE(), "DE register pair should return correct value")

	cpu.setDE(0xEF01)
	assert.Equal(t, uint8(0xEF), cpu.D, "D should be set correctly")
	assert.Equal(t, uint8(0x01), cpu.E, "E should be set correctly")

	// Test HL register pair
	cpu.H = 0x23
	cpu.L = 0x45
	assert.Equal(t, uint16(0x2345), cpu.HL(), "HL register pair should return correct value")

	cpu.setHL(0x6789)
	assert.Equal(t, uint8(0x67), cpu.H, "H should be set correctly")
	assert.Equal(t, uint8(0x89), cpu.L, "L should be set correctly")

	// Test AF register pair
	cpu.A = 0xF0
	cpu.setFlags(0x0F)
	assert.Equal(t, uint16(0xF00F), cpu.AF(), "AF register pair should return correct value")

	cpu.setAF(0x1E2D)
	assert.Equal(t, uint8(0x1E), cpu.A, "A should be set correctly")
	assert.Equal(t, uint8(0x2D), cpu.GetFlags(), "Flags should be set correctly")
}

func TestStackOperations(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Test push and pop byte
	originalSP := cpu.SP
	cpu.push(0x42)
	assert.Equal(t, originalSP-1, cpu.SP, "SP should decrement after push")
	assert.Equal(t, uint8(0x42), memory.Read(cpu.SP), "Value should be stored at SP")

	value := cpu.pop()
	assert.Equal(t, uint8(0x42), value, "Popped value should match pushed value")
	assert.Equal(t, originalSP, cpu.SP, "SP should return to original value")

	// Test push and pop 16-bit word
	cpu.push16(0x1234)
	assert.Equal(t, originalSP-2, cpu.SP, "SP should decrement by 2 after push16")

	word := cpu.pop16()
	assert.Equal(t, uint16(0x1234), word, "Popped word should match pushed word")
	assert.Equal(t, originalSP, cpu.SP, "SP should return to original value")
}

func TestExchange(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Set main registers
	cpu.A = 0x11
	cpu.B = 0x22
	cpu.C = 0x33
	cpu.D = 0x44
	cpu.E = 0x55
	cpu.H = 0x66
	cpu.L = 0x77
	cpu.setFlags(0x88)

	// Set alternate registers
	cpu.AltA = 0xAA
	cpu.AltB = 0xBB
	cpu.AltC = 0xCC
	cpu.AltD = 0xDD
	cpu.AltE = 0xEE
	cpu.AltH = 0xFF
	cpu.AltL = 0x00
	cpu.AltFlags.C = 1
	cpu.AltFlags.Z = 1

	// Test exchange
	cpu.exchange()

	// Check that main and alternate registers are swapped
	assert.Equal(t, uint8(0xAA), cpu.A, "A should be swapped")
	assert.Equal(t, uint8(0xBB), cpu.B, "B should be swapped")
	assert.Equal(t, uint8(0xCC), cpu.C, "C should be swapped")
	assert.Equal(t, uint8(0xDD), cpu.D, "D should be swapped")
	assert.Equal(t, uint8(0xEE), cpu.E, "E should be swapped")
	assert.Equal(t, uint8(0xFF), cpu.H, "H should be swapped")
	assert.Equal(t, uint8(0x00), cpu.L, "L should be swapped")

	assert.Equal(t, uint8(0x11), cpu.AltA, "A' should be swapped")
	assert.Equal(t, uint8(0x22), cpu.AltB, "B' should be swapped")
	assert.Equal(t, uint8(0x33), cpu.AltC, "C' should be swapped")
	assert.Equal(t, uint8(0x44), cpu.AltD, "D' should be swapped")
	assert.Equal(t, uint8(0x55), cpu.AltE, "E' should be swapped")
	assert.Equal(t, uint8(0x66), cpu.AltH, "H' should be swapped")
	assert.Equal(t, uint8(0x77), cpu.AltL, "L' should be swapped")
}

func TestExchangeAF(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory)
	assert.NoError(t, err)

	// Set main AF
	cpu.A = 0x12
	cpu.setFlags(0x34)

	// Set alternate AF
	cpu.AltA = 0x56
	cpu.AltFlags.C = 1
	cpu.AltFlags.Z = 1
	cpu.AltFlags.S = 1

	// Test AF exchange
	cpu.exchangeAF()

	// Check that AF is swapped
	assert.Equal(t, uint8(0x56), cpu.A, "A should be swapped")
	assert.Equal(t, uint8(0x12), cpu.AltA, "A' should be swapped")

	// Flags should be swapped
	assert.Equal(t, uint8(1), cpu.Flags.C, "C flag should be from alternate")
	assert.Equal(t, uint8(1), cpu.Flags.Z, "Z flag should be from alternate")
	assert.Equal(t, uint8(1), cpu.Flags.S, "S flag should be from alternate")
}

func TestLoadInstructionExecution(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory, WithSystemType(arch.GameBoy))
	assert.NoError(t, err)

	// Test LD BC,nn (0x01)
	memory.Write(0x0100, 0x01) // LD BC,nn
	memory.Write(0x0101, 0x34) // Low byte
	memory.Write(0x0102, 0x12) // High byte

	err = cpu.Step()
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

func TestMemoryOperations(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory, WithSystemType(arch.GameBoy))
	assert.NoError(t, err)

	// Test LD (BC),A (0x02)
	cpu.A = 0x42
	cpu.setBC(0x2000)
	memory.Write(0x0100, 0x02) // LD (BC),A

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint8(0x42), memory.Read(0x2000), "Memory at BC should contain A")

	// Test LD A,(BC) (0x0A)
	memory.Write(0x2000, 0x55)
	memory.Write(0x0101, 0x0A) // LD A,(BC)

	err = cpu.Step()
	assert.NoError(t, err, "Step should not return error")
	assert.Equal(t, uint8(0x55), cpu.A, "A should be loaded from memory at BC")
}

// =============================================================================
// Control Flow Operations
// =============================================================================

func TestJumpInstructions(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory, WithSystemType(arch.GameBoy))
	assert.NoError(t, err)

	// Test JP nn (0xC3)
	memory.Write(0x0100, 0xC3) // JP nn
	memory.Write(0x0101, 0x00) // Low byte
	memory.Write(0x0102, 0x20) // High byte

	err = cpu.Step()
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

func TestDJNZ(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory, WithSystemType(arch.GameBoy))
	assert.NoError(t, err)

	// Test DJNZ with branch taken
	cpu.B = 0x02
	memory.Write(0x0100, 0x10) // DJNZ
	memory.Write(0x0101, 0x05) // Offset +5

	err = cpu.Step()
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

func TestEndlessLoopDetection(t *testing.T) {
	memory := NewMemory()
	cpu, err := New(memory, WithSystemType(arch.GameBoy))
	assert.NoError(t, err)

	// Create an endless loop: JR -2 (0x18 0xFE)
	memory.Write(0x0100, 0x18) // JR relative
	memory.Write(0x0101, 0xFE) // -2 (jumps to 0x0100 + 2 + (-2) = 0x0100)

	initialCycles := cpu.cycles
	initialPC := cpu.PC

	// Execute the jump instruction once
	err = cpu.Step()
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
