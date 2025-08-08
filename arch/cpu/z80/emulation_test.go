package z80

import (
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

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
